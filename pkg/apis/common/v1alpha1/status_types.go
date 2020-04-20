package v1alpha1

var (
	statusReady   = "Ready"
	statusPending = "Pending"
	statusError   = "Error"
)

// Status defines the observed state of a New Relic resource
type Status struct {
	// The value will be set to `Ready` once the policy has been created in New Relic
	Status string `json:"status"`
	// When a policy fails to be created, the value will be set to the error message received from New Relic
	Reason string `json:"reason,omitempty"`
	// The resource id in New Relic
	NewrelicId *int64 `json:"newrelicId,omitempty"`
}

func NewError(newrelicId *int64, err error) Status {
	return Status{
		Status:     statusError,
		Reason:     err.Error(),
		NewrelicId: newrelicId,
	}
}

func NewPending(newrelicId *int64) Status {
	return Status{
		Status:     statusPending,
		Reason:     "",
		NewrelicId: newrelicId,
	}
}

func NewReady(newrelicId *int64) Status {
	return Status{
		Status:     statusReady,
		Reason:     "",
		NewrelicId: newrelicId,
	}
}

func (s Status) IsReady() bool {
	return s.Status == statusReady
}

func (s Status) IsPending() bool {
	return s.Status == statusPending
}

func (s Status) IsError() bool {
	return s.Status == statusError
}
