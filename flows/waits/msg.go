package waits

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/triggers"
)

func init() {
	RegisterType(TypeMsg, func() flows.Wait { return &MsgWait{} })
}

// TypeMsg is the type of our message wait
const TypeMsg string = "msg"

// MsgWait is a wait which waits for an incoming message (i.e. a msg_received event)
type MsgWait struct {
	baseWait
}

// NewMsgWait creates a new message wait
func NewMsgWait(timeout *int) *MsgWait {
	return &MsgWait{baseWait{Timeout_: timeout}}
}

// Type returns the type of this wait
func (w *MsgWait) Type() string { return TypeMsg }

// Begin beings waiting at this wait
func (w *MsgWait) Begin(run flows.FlowRun, step flows.Step) bool {
	if !w.baseWait.Begin(run) {
		return false
	}

	// if we have a msg trigger and we're the first thing to happen... then we skip ourselves
	triggerHasMsg := run.Session().Trigger().Type() == triggers.TypeMsg

	if triggerHasMsg && len(run.Session().Runs()) == 1 && len(run.Path()) == 1 {
		return false
	}

	run.LogEvent(step, events.NewMsgWait(w.TimeoutOn_))
	return true
}

// End ends this wait or returns an error
func (w *MsgWait) End(resume flows.Resume) error {
	if resume.Type() == resumes.TypeMsg {
		return nil
	}

	return w.baseWait.End(resume)
}

var _ flows.Wait = (*MsgWait)(nil)
