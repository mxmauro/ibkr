package models

// -----------------------------------------------------------------------------

type FaData int32

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
	}
	return ""
}
