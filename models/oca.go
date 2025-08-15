package models

// -----------------------------------------------------------------------------

type Oca int32

const (
	OcaNone                  Oca = 0
	OcaCancelWithBlocking    Oca = 1
	OcaReduceWithBlocking    Oca = 2
	OcaReduceWithoutBlocking Oca = 3
)

// -----------------------------------------------------------------------------

func (oca Oca) String() string {
	switch oca {
	case OcaNone:
		return "None"
	case OcaCancelWithBlocking:
		return "Cancel with blocking"
	case OcaReduceWithBlocking:
		return "Reduce with blocking"
	case OcaReduceWithoutBlocking:
		return "Reduce without blocking"
	}
	return ""
}
