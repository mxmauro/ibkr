package models

// -----------------------------------------------------------------------------

type SecurityType string

const (
	SecurityTypeNone                      SecurityType = ""
	SecurityTypeStock                     SecurityType = "STK"
	SecurityTypeOption                    SecurityType = "OPT"
	SecurityTypeFuture                    SecurityType = "FUT"
	SecurityTypeIndex                     SecurityType = "IND"
	SecurityTypeFuturesOption             SecurityType = "FOP"
	SecurityTypeForex                     SecurityType = "CASH"
	SecurityTypePair                      SecurityType = "BAG"
	SecurityTypeWarrant                   SecurityType = "WAR"
	SecurityTypeBond                      SecurityType = "BOND"
	SecurityTypeCommodity                 SecurityType = "CMDTY"
	SecurityTypeNews                      SecurityType = "NEWS"
	SecurityTypeMutualFund                SecurityType = "FUND"
	SecurityTypeContinuousFuturesContract SecurityType = "CONTFUT"
	SecurityTypeContractForDifference     SecurityType = "CFD"
	SecurityTypeIOPT                      SecurityType = "IOPT"
	SecurityTypeForwardContract           SecurityType = "FWD"
	SecurityTypeTreasuryBill              SecurityType = "BILL"
	SecurityTypeFixedIncome               SecurityType = "FIXED"
	SecurityTypeStockLoanBorrow           SecurityType = "SLB"
	SecurityTypeBasket                    SecurityType = "BSK"
	SecurityTypeInterCommodityUnit        SecurityType = "ICU"
	SecurityTypeInterCommoditySpread      SecurityType = "ICS"
	SecurityTypeCrypto                    SecurityType = "CRYPTO"
)

// -----------------------------------------------------------------------------

func NewSecurityTypeFromString(securityType string) SecurityType {
	switch securityType {
	case string(SecurityTypeStock):
		return SecurityTypeStock
	case string(SecurityTypeOption):
		return SecurityTypeOption
	case string(SecurityTypeFuture):
		return SecurityTypeFuture
	case string(SecurityTypeIndex):
		return SecurityTypeIndex
	case string(SecurityTypeFuturesOption):
		return SecurityTypeFuturesOption
	case string(SecurityTypeForex):
		return SecurityTypeForex
	case string(SecurityTypePair):
		return SecurityTypePair
	case string(SecurityTypeWarrant):
		return SecurityTypeWarrant
	case string(SecurityTypeBond):
		return SecurityTypeBond
	case string(SecurityTypeCommodity):
		return SecurityTypeCommodity
	case string(SecurityTypeNews):
		return SecurityTypeNews
	case string(SecurityTypeMutualFund):
		return SecurityTypeMutualFund
	case string(SecurityTypeContinuousFuturesContract):
		return SecurityTypeContinuousFuturesContract
	case string(SecurityTypeContractForDifference):
		return SecurityTypeContractForDifference
	case string(SecurityTypeIOPT):
		return SecurityTypeIOPT
	case string(SecurityTypeForwardContract):
		return SecurityTypeForwardContract
	case string(SecurityTypeTreasuryBill):
		return SecurityTypeTreasuryBill
	case string(SecurityTypeFixedIncome):
		return SecurityTypeFixedIncome
	case string(SecurityTypeStockLoanBorrow):
		return SecurityTypeStockLoanBorrow
	case string(SecurityTypeBasket):
		return SecurityTypeBasket
	case string(SecurityTypeInterCommodityUnit):
		return SecurityTypeInterCommodityUnit
	case string(SecurityTypeInterCommoditySpread):
		return SecurityTypeInterCommoditySpread
	case string(SecurityTypeCrypto):
		return SecurityTypeCrypto
	}
	return SecurityTypeNone
}
