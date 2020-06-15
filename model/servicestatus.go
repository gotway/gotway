package model

type serviceStatus int

const (
	// Healthy service is responding to healtchecks
	Healthy serviceStatus = 0
	// Idle service is not responding to healthchecks
	Idle serviceStatus = 1
)

func (status serviceStatus) String() string {
	return [...]string{"healty", "idle"}[status]
}
