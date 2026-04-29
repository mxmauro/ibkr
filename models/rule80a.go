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
