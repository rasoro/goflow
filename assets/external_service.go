package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
)

// ExternalServiceUUID is the UUID of a external service
type ExternalServiceUUID uuids.UUID

// ExternalService is a third party service that can be called
//
//    {
//		  "uuid": "4e19fc3c-ae17-4f6b-acb5-7d915e29dc27",
//		  "name": "Third party service integration 1",
//		  "uuid": "generic third party service",
//    }
//
// @asset external_service
type ExternalService interface {
	UUID() ExternalServiceUUID
	Name() string
	Type() string
}

// ExternalServiceReference is used to reference a external service
type ExternalServiceReference struct {
	UUID ExternalServiceUUID `json:"uuid" validate:"required,uuid"`
	Name string              `json:"name"`
}

type ExternalServiceCallAction struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ExternalServiceParam struct {
	Data struct {
		Value string `json:"value,omitempty"`
	} `json:"data,omitempty"`
	Filter struct {
		Value *ExternalServiceFilterValue `json:"value"`
	} `json:"filter,omitempty"`
	Type        string `json:"type,omitempty"`
	VerboseName string `json:"verboseName,omitempty"`
}

type ExternalServiceFilterValue struct {
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	VerboseName string `json:"verboseName,omitempty"`
}

func NewExternalServiceParam(
	dataValue,
	filterName,
	filterType,
	filterVerboseName,
	pType,
	verboseName string,
) *ExternalServiceParam {
	p := &ExternalServiceParam{}
	p.Data.Value = dataValue
	if filterName != "" && filterType != "" && filterVerboseName != "" {
		p.Filter.Value = &ExternalServiceFilterValue{}
		p.Filter.Value.Name = filterName
		p.Filter.Value.Type = filterType
		p.Filter.Value.VerboseName = filterVerboseName
	} else {
		p.Filter.Value = nil
	}
	p.Type = pType
	p.VerboseName = verboseName
	return p
}

// NewExternalServiceReference creates a new external service reference with the given UUID and name
func NewExternalServiceReference(uuid ExternalServiceUUID, name string) *ExternalServiceReference {
	return &ExternalServiceReference{UUID: uuid, Name: name}
}

// Type returns the name of the asset type
func (r *ExternalServiceReference) Type() string {
	return "external_service"
}

// GenericUUID returns the untyped UUID
func (r *ExternalServiceReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity retusns the unique identity of the asset
func (r *ExternalServiceReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *ExternalServiceReference) Variable() bool {
	return false
}

// String returns a formated string for the external service referenced
func (r *ExternalServiceReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ UUIDReference = (*ExternalServiceReference)(nil)
