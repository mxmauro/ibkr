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
