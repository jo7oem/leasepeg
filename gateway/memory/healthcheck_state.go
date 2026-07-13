package memory

import (
	"fmt"
	"sync"
)

const (
	StateStarting ServerState = iota
	StateReady
	StateDegraded
	StateError
	StateShuttingDown
)

type ServerState int

func (s ServerState) String() string {
	switch s {
	case StateStarting:
		return "Starting"
	case StateReady:
		return "Ready"
	case StateDegraded:
		return "Degraded"
	case StateError:
		return "Error"
	case StateShuttingDown:
		return "ShuttingDown"
	default:
		return fmt.Sprintf("UnknownStatus: %d", s)
	}
}

type ServerStatus struct {
	State ServerState
	Err   error
}

func (s ServerStatus) Copy() ServerStatus {
	return ServerStatus{
		State: s.State,
		Err:   s.Err,
	}
}

var healthCheckStatusMutex sync.RWMutex
var healthCheckStatus = ServerStatus{
	State: StateStarting,
}

func SetHealthCheckState(state ServerState, err error) {
	healthCheckStatusMutex.Lock()
	defer healthCheckStatusMutex.Unlock()
	healthCheckStatus.State = state
	healthCheckStatus.Err = err
}
func GetHealthCheckState() ServerStatus {
	healthCheckStatusMutex.RLock()
	defer healthCheckStatusMutex.RUnlock()
	return healthCheckStatus.Copy()
}
