package resource

import (
	"encoding/json"
	"sort"

	"github.com/zclconf/go-cty/cty"
)

type Resource interface {
	TerraformId() string
	TerraformType() string
	CtyValue() *cty.Value
}

var refactoredResources = []string{
	"aws_s3_bucket",
}

func IsRefactoredResource(typ string) bool {
	for _, refactoredResource := range refactoredResources {
		if typ == refactoredResource {
			return true
		}
	}
	return false
}

type AbstractResource struct {
	Id    string
	Type  string
	Attrs map[string]interface{}
}

func (a *AbstractResource) TerraformId() string {
	return a.Id
}

func (a *AbstractResource) TerraformType() string {
	return a.Type
}

func (a *AbstractResource) CtyValue() *cty.Value {
	return nil
}

type ResourceFactory interface {
	CreateResource(data interface{}, ty string) (*cty.Value, error)
}

type SerializableResource struct {
	Resource
}

type SerializedResource struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

func (u SerializedResource) TerraformId() string {
	return u.Id
}

func (u SerializedResource) TerraformType() string {
	return u.Type
}

func (u SerializedResource) CtyValue() *cty.Value {
	return &cty.NilVal
}

func (s *SerializableResource) UnmarshalJSON(bytes []byte) error {
	var res SerializedResource

	if err := json.Unmarshal(bytes, &res); err != nil {
		return err
	}
	s.Resource = res
	return nil
}

func (s SerializableResource) MarshalJSON() ([]byte, error) {
	return json.Marshal(SerializedResource{Id: s.TerraformId(), Type: s.TerraformType()})
}

type NormalizedResource interface {
	NormalizeForState() (Resource, error)
	NormalizeForProvider() (Resource, error)
}

func IsSameResource(rRs, lRs Resource) bool {
	return rRs.TerraformType() == lRs.TerraformType() && rRs.TerraformId() == lRs.TerraformId()
}

func Sort(res []Resource) []Resource {
	sort.SliceStable(res, func(i, j int) bool {
		if res[i].TerraformType() != res[j].TerraformType() {
			return res[i].TerraformType() < res[j].TerraformType()
		}
		return res[i].TerraformId() < res[j].TerraformId()
	})
	return res
}
