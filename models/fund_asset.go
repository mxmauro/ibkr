package models

// -----------------------------------------------------------------------------

type FundAsset string

const (
	FundAssetNone        FundAsset = ""
	FundAssetOthers      FundAsset = "000"
	FundAssetMoneyMarket FundAsset = "001"
	FundAssetFixedIncome FundAsset = "002"
	FundAssetMultiAsset  FundAsset = "003"
	FundAssetEquity      FundAsset = "004"
	FundAssetSector      FundAsset = "005"
	FundAssetGuaranteed  FundAsset = "006"
	FundAssetAlternative FundAsset = "007"
)
