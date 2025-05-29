package models

// -----------------------------------------------------------------------------

type FundDistributionPolicyIndicator string

const (
	FundDistributionPolicyIndicatorNone             FundDistributionPolicyIndicator = ""
	FundDistributionPolicyIndicatorAccumulationFund FundDistributionPolicyIndicator = "N"
	FundDistributionPolicyIndicatorIncomeFund       FundDistributionPolicyIndicator = "Y"
)

// -----------------------------------------------------------------------------

func NewFundDistributionPolicyIndicatorFromString(fundIndicator string) FundDistributionPolicyIndicator {
	switch fundIndicator {
	case string(FundDistributionPolicyIndicatorAccumulationFund):
		return FundDistributionPolicyIndicatorAccumulationFund
	case string(FundDistributionPolicyIndicatorIncomeFund):
		return FundDistributionPolicyIndicatorIncomeFund
	default:
		return FundDistributionPolicyIndicatorNone
	}
}
