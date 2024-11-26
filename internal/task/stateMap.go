package task

type stateMap map[State][]State

var stateTransitionMap = map[State][]State{
	Pending:   {Scheduled, Failed},
	Scheduled: {Running, Failed},
	Running:   {Completed, Failed},
	Completed: {},
	Failed:    {},
}

func contains(states []State, state State) bool {
	for _, s := range states {
		if s == state {
			return true
		}
	}

	return false
}

func ValidStateTransition(src, dst State) bool {
	return src == dst || contains(stateTransitionMap[src], dst)
}
