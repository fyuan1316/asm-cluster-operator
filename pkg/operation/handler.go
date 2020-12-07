package operation

type Handler interface {
	Do() Resulter
	/*
		Deploy() Resulter
		Remove() Resulter
		Upgrade() Resulter
		HealthCheck()
		ProcessRunning()
		//ProcessFailedProvisioned()

	*/
}

//---------------------------------------

type Resulter interface {
	GetState() State
	GetReason() Reason
	GetError() error
	IsStatusChanged() bool
}

type Result struct {
	State   State
	Reason  Reason
	Err     error
	Changed bool
}

func (r Result) GetState() State {
	return r.State
}

func (r Result) GetReason() Reason {
	return r.Reason
}

func (r Result) GetError() error {
	return r.Err
}
func (r Result) IsStatusChanged() bool {
	return r.Changed
}

var _ Resulter = Result{}

type State string

var (
	Succeed   State = "Succeed"
	Failed    State = "Failed"
	Operating State = "Operating"
)

type Reason string

//------
