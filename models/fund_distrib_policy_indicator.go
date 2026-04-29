package models

// -----------------------------------------------------------------------------

type FundDistributionPolicyIndicator string

const (
	FundDistributionPolicyIndicatorNone             FundDistributionPolicyIndicator = ""
	FundDistributionPolicyIndicatorAccumulationFund FundDistributionPolicyIndicator = "N"
	FundDistributionPolicyIndicatorIncomeFund       FundDistributionPolicyIndicator = "Y"
)
