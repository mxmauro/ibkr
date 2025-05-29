package models

// -----------------------------------------------------------------------------

type FaData int64

const (
	FaDataGroups  FaData = 1
	FaDataAliases FaData = 3
)

// -----------------------------------------------------------------------------

func (fa FaData) String() string {
	switch fa {
	case FaDataGroups:
		return "GROUPS"
	case FaDataAliases:
		return "ALIASES"
	default:
		return ""
	}
}
