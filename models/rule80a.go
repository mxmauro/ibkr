package models

// -----------------------------------------------------------------------------

type Rule80A string

const (
	Rule80aNone                 Rule80A = ""
	Rule80aIndividual           Rule80A = "I"
	Rule80aAgency               Rule80A = "A"
	Rule80aAgentOtherMember     Rule80A = "W"
	Rule80aIndividualPTIA       Rule80A = "J"
	Rule80aAgencyPTIA           Rule80A = "U"
	Rule80aAgentOtherMemberPTIA Rule80A = "M"
	Rule80aIndividualPT         Rule80A = "K"
	Rule80aAgencyPT             Rule80A = "Y"
	Rule80aAgentOtherMemberPT   Rule80A = "N"
)

// -----------------------------------------------------------------------------

func NewRule80aFromString(rule80a string) Rule80A {
	switch rule80a {
	case string(Rule80aIndividual):
		return Rule80aIndividual
	case string(Rule80aAgency):
		return Rule80aAgency
	case string(Rule80aAgentOtherMember):
		return Rule80aAgentOtherMember
	case string(Rule80aIndividualPTIA):
		return Rule80aIndividualPTIA
	case string(Rule80aAgencyPTIA):
		return Rule80aAgencyPTIA
	case string(Rule80aAgentOtherMemberPTIA):
		return Rule80aAgentOtherMemberPTIA
	case string(Rule80aIndividualPT):
		return Rule80aIndividualPT
	case string(Rule80aAgencyPT):
		return Rule80aAgencyPT
	case string(Rule80aAgentOtherMemberPT):
		return Rule80aAgentOtherMemberPT
	}
	return Rule80aNone
}
