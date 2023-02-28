package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeCallExternalService, func() flows.Action { return &CallExternalServiceAction{} })
}

var externalServiceCategories = []string{CategorySuccess, CategoryFailure}

// TypeCallExternalService is the type for the external service calls action
const TypeCallExternalService string = "call_external_service"

// CallExternalServiceAction is used to call a external service in the context of contact
//   {
//     "uuid": "3ed82e8e-409c-46b8-99ec-4ef9c7ec0270",
//     "type": "call_external_service",
//     "external_service": {
//         "uuid": "0ebd32fd-362b-4253-89a1-3796aa499b82",
//         "name": "service foo",
//     },
//		 "params": [{"data":{"value":"foo"}, "filter": {"value":{"name":"foo","type":"bar","verboseName":"barz"}}, "type": "foo", "verboseName": "bar"}],
//     "result_name": "external_service_call"
//   }
type CallExternalServiceAction struct {
	baseAction
	onlineAction

	ExternalService *assets.ExternalServiceReference `json:"external_service,omitempty"`
	CallAction      assets.ExternalServiceCallAction `json:"call"`
	Params          []assets.ExternalServiceParam    `json:"params,omitempty"`
	ResultName      string                           `json:"result_name,omitempty"`
}

func NewCallExternalService(uuid flows.ActionUUID, externalService *assets.ExternalServiceReference, callAction assets.ExternalServiceCallAction, params []assets.ExternalServiceParam, resultName string) *CallExternalServiceAction {
	return &CallExternalServiceAction{
		baseAction:      newBaseAction(TypeCallExternalService, uuid),
		ExternalService: externalService,
		CallAction:      callAction,
		Params:          params,
		ResultName:      resultName,
	}
}

func (a *CallExternalServiceAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	externalServices := run.Session().Assets().ExternalServices()
	externalService := externalServices.Get(a.ExternalService.UUID)

	return a.call(run, step, externalService, a.CallAction, a.Params, logEvent)
}

func (a *CallExternalServiceAction) call(run flows.FlowRun, step flows.Step, externalService *flows.ExternalService, callAction assets.ExternalServiceCallAction, params []assets.ExternalServiceParam, logEvent flows.EventCallback) error {
	if externalService == nil {
		logEvent(events.NewDependencyError(a.ExternalService))
		return nil
	}

	svc, err := run.Session().Engine().Services().ExternalService(run.Session(), externalService)
	if err != nil {
		logEvent(events.NewError(err))
		return nil
	}

	httpLogger := &flows.HTTPLogger{}

	call, err := svc.Call(run.Session(), callAction, params, httpLogger.Log)
	if err != nil {
		logEvent(events.NewError(err))
	}
	if len(httpLogger.Logs) > 0 {
		logEvent(events.NewExternalServiceCalled(externalService.Reference(), httpLogger.Logs))
	}

	if call != nil {
		if a.ResultName != "" {
			input := fmt.Sprintf("%s %s", call.Request.Method, call.Request.URL.String())
			a.saveResult(run, step, a.ResultName, string(call.ResponseJSON), CategorySuccess, "", input, call.ResponseJSON, logEvent)
		}
	}

	return nil
}

func (a *CallExternalServiceAction) Results(include func(*flows.ResultInfo)) {
	if a.ResultName != "" {
		include(flows.NewResultInfo(a.ResultName, externalServiceCategories))
	}
}
