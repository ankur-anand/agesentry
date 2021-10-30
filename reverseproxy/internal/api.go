package internal

// StateDirector gives information about the State Machine to the Load Director
type StateDirector interface {
	// IsHealthy should return the health status of the provided named state machine
	IsHealthy(name string) bool
	// LoadCount should return the current active load count of the provided named state machine
	// Bool denotes the presence and absence of the provided named state machine.
	LoadCount(name string) (int, bool)
}

// LoadDirector wraps dispatches methods that helps in load balancing the cluster
// based on the configured policy of load balancer.
type LoadDirector interface {
	// DispatchTo returns the name of the state machine based on State and configured load policy.
	DispatchTo(sd StateDirector, sortedName []string) string
}
