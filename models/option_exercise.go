package models

// -----------------------------------------------------------------------------

type OptionExercise int

const (
	OptionExerciseNone                 OptionExercise = -1
	OptionExerciseExercise             OptionExercise = 1
	OptionExerciseLapse                OptionExercise = 2
	OptionExerciseDoNothing            OptionExercise = 3
	OptionExerciseAssigned             OptionExercise = 100
	OptionExerciseAutoexerciseClearing OptionExercise = 101
	OptionExerciseExpired              OptionExercise = 102
	OptionExerciseNetting              OptionExercise = 103
	OptionExerciseAutoexerciseTrading  OptionExercise = 104
)

// -----------------------------------------------------------------------------

func (e OptionExercise) String() string {
	switch e {
	case OptionExerciseNone:
		return "None"
	case OptionExerciseExercise:
		return "Exercise"
	case OptionExerciseLapse:
		return "Lapse"
	case OptionExerciseDoNothing:
		return "DoNothing"
	case OptionExerciseAssigned:
		return "Assigned"
	case OptionExerciseAutoexerciseClearing:
		return "AutoexerciseClearing"
	case OptionExerciseExpired:
		return "Expired"
	case OptionExerciseNetting:
		return "Netting"
	case OptionExerciseAutoexerciseTrading:
		return "AutoexerciseTrading"
	default:
		return "Unknown"
	}
}
