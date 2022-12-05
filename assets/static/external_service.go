package static

import (
	"github.com/nyaruka/goflow/assets"
)

// ExternalService is a JSON serializable implementation of a external service asset
type ExternalService struct {
	UUID_ assets.ExternalServiceUUID `json:"uuid" validate:"required,uuid`
	Name_ string                     `json:"name"`
	Type_ string                     `json:"type"`
}

// NewExternalService creates a new external service
func NewExternalService(uuid assets.ExternalServiceUUID, name string, type_ string) assets.ExternalService {
	return &ExternalService{
		UUID_: uuid,
		Name_: name,
		Type_: type_,
	}
}

// UUID returns the UUIUD of this external service
func (e *ExternalService) UUID() assets.ExternalServiceUUID { return e.UUID_ }

// Name returns the name of this external service
func (e *ExternalService) Name() string { return e.Name_ }

// Type returns the type of this external service
func (e *ExternalService) Type() string { return e.Type_ }
