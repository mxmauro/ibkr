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

// -----------------------------------------------------------------------------

func NewFundAssetFromString(fundAsset string) FundAsset {
	switch fundAsset {
	case string(FundAssetOthers):
		return FundAssetOthers
	case string(FundAssetMoneyMarket):
		return FundAssetMoneyMarket
	case string(FundAssetFixedIncome):
		return FundAssetFixedIncome
	case string(FundAssetMultiAsset):
		return FundAssetMultiAsset
	case string(FundAssetEquity):
		return FundAssetEquity
	case string(FundAssetSector):
		return FundAssetSector
	case string(FundAssetGuaranteed):
		return FundAssetGuaranteed
	case string(FundAssetAlternative):
		return FundAssetAlternative
	}
	return FundAssetNone
}
