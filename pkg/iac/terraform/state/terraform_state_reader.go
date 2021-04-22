package state

import (
	"fmt"

	"github.com/cloudskiff/driftctl/pkg/iac"
	"github.com/cloudskiff/driftctl/pkg/iac/config"
	"github.com/cloudskiff/driftctl/pkg/iac/terraform/state/backend"
	"github.com/cloudskiff/driftctl/pkg/iac/terraform/state/enumerator"
	"github.com/cloudskiff/driftctl/pkg/remote/deserializer"
	"github.com/cloudskiff/driftctl/pkg/resource"

	"github.com/cloudskiff/driftctl/pkg/terraform"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/states"
	"github.com/hashicorp/terraform/states/statefile"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
	ctyconvert "github.com/zclconf/go-cty/cty/convert"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

const TerraformStateReaderSupplier = "tfstate"

type TerraformStateReader struct {
	library                  *terraform.ProviderLibrary
	config                   config.SupplierConfig
	backend                  backend.Backend
	enumerator               enumerator.StateEnumerator
	deserializers            []deserializer.CTYDeserializer
	backendOptions           *backend.Options
	resourceSchemaRepository resource.SchemaRepositoryInterface
}

func (r *TerraformStateReader) initReader() error {
	r.enumerator = enumerator.GetEnumerator(r.config)
	return nil
}

func NewReader(config config.SupplierConfig, library *terraform.ProviderLibrary, backendOpts *backend.Options, resourceSchemaRepository resource.SchemaRepositoryInterface) (*TerraformStateReader, error) {
	reader := TerraformStateReader{library: library, config: config, deserializers: iac.Deserializers(), backendOptions: backendOpts, resourceSchemaRepository: resourceSchemaRepository}
	err := reader.initReader()
	if err != nil {
		return nil, err
	}
	return &reader, nil
}

func (r *TerraformStateReader) retrieve() (map[string][]cty.Value, error) {
	b, err := backend.GetBackend(r.config, r.backendOptions)
	if err != nil {
		return nil, err
	}
	r.backend = b

	state, err := read(r.backend)
	defer r.backend.Close()
	if err != nil {
		return nil, err
	}

	resMap := make(map[string][]cty.Value)
	for moduleName, module := range state.Modules {
		logrus.WithFields(logrus.Fields{
			"module":        moduleName,
			"resourceCount": fmt.Sprintf("%d", len(module.Resources)),
		}).Debug("Found module in state")
		for _, stateRes := range module.Resources {
			resName := stateRes.Addr.Resource.Name
			resType := stateRes.Addr.Resource.Type
			if stateRes.Addr.Resource.Mode != addrs.ManagedResourceMode {
				logrus.WithFields(logrus.Fields{
					"mode": stateRes.Addr.Resource.Mode,
					"name": resName,
					"type": resType,
				}).Debug("Skipping state entry as it is not a managed resource")
				continue
			}
			providerType := stateRes.ProviderConfig.Provider.Type
			provider := r.library.Provider(providerType)
			if provider == nil {
				logrus.WithFields(logrus.Fields{
					"providerKey": providerType,
				}).Debug("Unsupported provider found in state")
				continue
			}
			schema := provider.Schema()[stateRes.Addr.Resource.Type]
			for _, instance := range stateRes.Instances {
				decodedVal, err := instance.Current.Decode(schema.Block.ImpliedType())
				if err != nil {
					// Try to do a manual type conversion if we got a path error
					// It will allow driftctl to read state generated with a superior version of provider
					// than the actually supported one
					// by ignoring new fields
					_, isPathError := err.(cty.PathError)
					if isPathError {
						logrus.WithFields(logrus.Fields{
							"name": resName,
							"type": resType,
							"err":  err.Error(),
						}).Debug("Got a cty path error when deserializing state")

						decodedVal, err = r.convertInstance(instance.Current, schema.Block.ImpliedType())
					}

					if err != nil {
						logrus.WithFields(logrus.Fields{
							"name": resName,
							"type": resType,
						}).Error("Unable to decode resource from state")
						return nil, err
					}
				}
				_, exists := resMap[stateRes.Addr.Resource.Type]
				if !exists {
					resMap[stateRes.Addr.Resource.Type] = []cty.Value{
						decodedVal.Value,
					}
				} else {
					resMap[stateRes.Addr.Resource.Type] = append(resMap[stateRes.Addr.Resource.Type], decodedVal.Value)
				}
			}
		}
	}

	return resMap, nil
}

func (r *TerraformStateReader) convertInstance(instance *states.ResourceInstanceObjectSrc, ty cty.Type) (*states.ResourceInstanceObject, error) {
	inputType, err := ctyjson.ImpliedType(instance.AttrsJSON)
	if err != nil {
		return nil, err
	}
	input, err := ctyjson.Unmarshal(instance.AttrsJSON, inputType)
	if err != nil {
		return nil, err
	}

	convertedVal, err := ctyconvert.Convert(input, ty)
	if err != nil {
		return nil, err
	}

	instanceObj := &states.ResourceInstanceObject{
		Value:               convertedVal,
		Status:              instance.Status,
		Dependencies:        instance.Dependencies,
		Private:             instance.Private,
		CreateBeforeDestroy: instance.CreateBeforeDestroy,
	}

	logrus.Debug("Successfully converted resource")

	return instanceObj, nil
}

func (r *TerraformStateReader) decode(values map[string][]cty.Value) ([]resource.Resource, error) {
	results := make([]resource.Resource, 0)
	for _, deserializer := range r.deserializers {

		typ := deserializer.HandledType().String()
		vals, exists := values[typ]
		if !exists {
			continue
		}
		decodedResources, err := deserializer.Deserialize(vals)
		if err != nil {
			logrus.Warnf("Could not read from decoder for %s: %+v", typ, err)
			continue
		}
		for _, res := range decodedResources {
			logrus.WithFields(logrus.Fields{
				"path":    r.config.Path,
				"backend": r.config.Backend,
				"id":      res.TerraformId(),
				"type":    res.TerraformType(),
			}).Debug("Found IAC resource")
			if resource.IsRefactoredResource(res.TerraformType()) {
				schema, exist := r.resourceSchemaRepository.GetSchema(res.TerraformType())
				ctyAttr := resource.ToResourceAttributes(res.CtyValue())
				if exist && schema.NormalizeFunc != nil {
					schema.NormalizeFunc(ctyAttr)
				}

				newRes := &resource.AbstractResource{
					Id:    res.TerraformId(),
					Type:  res.TerraformType(),
					Attrs: ctyAttr,
				}
				results = append(results, newRes)
				continue
			}
			normalisable, ok := res.(resource.NormalizedResource)
			if ok {
				normalizedRes, err := normalisable.NormalizeForState()
				if err != nil {
					logrus.Errorf("Could not normalize state for res %s: %+v", res.TerraformId(), err)
					results = append(results, res)
				}

				if err == nil {
					results = append(results, normalizedRes)
				}
			}
			if !ok {
				results = append(results, res)
			}
		}
	}

	return results, nil
}

func (r *TerraformStateReader) Resources() ([]resource.Resource, error) {

	if r.enumerator == nil {
		return r.retrieveForState(r.config.Path)
	}

	return r.retrieveMultiplesStates()
}

func (r *TerraformStateReader) retrieveForState(path string) ([]resource.Resource, error) {
	r.config.Path = path
	logrus.WithFields(logrus.Fields{
		"path":    r.config.Path,
		"backend": r.config.Backend,
	}).Debug("Reading resources from state")
	values, err := r.retrieve()
	if err != nil {
		return nil, err
	}
	return r.decode(values)
}

func (r *TerraformStateReader) retrieveMultiplesStates() ([]resource.Resource, error) {
	keys, err := r.enumerator.Enumerate()
	if err != nil {
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"keys": keys,
	}).Debug("Enumerated keys")
	results := make([]resource.Resource, 0)

	for _, key := range keys {
		resources, err := r.retrieveForState(key)
		if err != nil {
			return nil, err
		}
		results = append(results, resources...)
	}

	return results, nil
}

func read(reader backend.Backend) (*states.State, error) {
	state, err := readState(reader)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func readState(reader backend.Backend) (*states.State, error) {
	state, err := statefile.Read(reader)
	if err != nil {
		return nil, err
	}
	return state.State, nil
}
