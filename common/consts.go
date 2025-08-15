package common

// -----------------------------------------------------------------------------

// const (
// UnsetInt   = math.MaxInt32
// UnsetLong  = math.MaxInt64
// UnsetFloat = math.MaxFloat64
// INFINITY_STRING string  = "Infinity"
// )

const (
	MessageDelimiter = '\x00'
)

const (
	ServerVersionProtobuf                  = 201
	ServerVersionZeroStrike                = 202
	ServerVersionProtobufPlaceOrder        = 203
	ServerVersionProtobufCompletedOrder    = 204
	ServerVersionProtobufContractData      = 205
	ServerVersionProtobufMarketData        = 206
	ServerVersionProtobufAccountsPositions = 207
	ServerVersionProtobufHistoricalData    = 208

	MinClientVersion = ServerVersionProtobuf
	MaxClientVersion = ServerVersionProtobufAccountsPositions

	MinServerVersion = ServerVersionProtobuf
)

const (
	TICK_PRICE                               uint32 = 1
	TICK_SIZE                                uint32 = 2
	ORDER_STATUS                             uint32 = 3
	ERR_MSG                                  uint32 = 4
	OPEN_ORDER                               uint32 = 5
	ACCT_VALUE                               uint32 = 6
	PORTFOLIO_VALUE                          uint32 = 7
	ACCT_UPDATE_TIME                         uint32 = 8
	NEXT_VALID_ID                            uint32 = 9
	CONTRACT_DATA                            uint32 = 10
	EXECUTION_DATA                           uint32 = 11
	MARKET_DEPTH                             uint32 = 12
	MARKET_DEPTH_L2                          uint32 = 13
	NEWS_BULLETINS                           uint32 = 14
	MANAGED_ACCTS                            uint32 = 15
	RECEIVE_FA                               uint32 = 16
	HISTORICAL_DATA                          uint32 = 17
	BOND_CONTRACT_DATA                       uint32 = 18
	SCANNER_PARAMETERS                       uint32 = 19
	SCANNER_DATA                             uint32 = 20
	TICK_OPTION_COMPUTATION                  uint32 = 21
	TICK_GENERIC                             uint32 = 45
	TICK_STRING                              uint32 = 46
	TICK_EFP                                 uint32 = 47
	CURRENT_TIME                             uint32 = 49
	REAL_TIME_BARS                           uint32 = 50
	FUNDAMENTAL_DATA                         uint32 = 51
	CONTRACT_DATA_END                        uint32 = 52
	OPEN_ORDER_END                           uint32 = 53
	ACCT_DOWNLOAD_END                        uint32 = 54
	EXECUTION_DATA_END                       uint32 = 55
	DELTA_NEUTRAL_VALIDATION                 uint32 = 56
	TICK_SNAPSHOT_END                        uint32 = 57
	MARKET_DATA_TYPE                         uint32 = 58
	COMMISSION_AND_FEES_REPORT               uint32 = 59
	POSITION_DATA                            uint32 = 61
	POSITION_END                             uint32 = 62
	ACCOUNT_SUMMARY                          uint32 = 63
	ACCOUNT_SUMMARY_END                      uint32 = 64
	VERIFY_MESSAGE_API                       uint32 = 65
	VERIFY_COMPLETED                         uint32 = 66
	DISPLAY_GROUP_LIST                       uint32 = 67
	DISPLAY_GROUP_UPDATED                    uint32 = 68
	VERIFY_AND_AUTH_MESSAGE_API              uint32 = 69
	VERIFY_AND_AUTH_COMPLETED                uint32 = 70
	POSITION_MULTI                           uint32 = 71
	POSITION_MULTI_END                       uint32 = 72
	ACCOUNT_UPDATE_MULTI                     uint32 = 73
	ACCOUNT_UPDATE_MULTI_END                 uint32 = 74
	SECURITY_DEFINITION_OPTION_PARAMETER     uint32 = 75
	SECURITY_DEFINITION_OPTION_PARAMETER_END uint32 = 76
	SOFT_DOLLAR_TIERS                        uint32 = 77
	FAMILY_CODES                             uint32 = 78
	SYMBOL_SAMPLES                           uint32 = 79
	MKT_DEPTH_EXCHANGES                      uint32 = 80
	TICK_REQ_PARAMS                          uint32 = 81
	SMART_COMPONENTS                         uint32 = 82
	NEWS_ARTICLE                             uint32 = 83
	TICK_NEWS                                uint32 = 84
	NEWS_PROVIDERS                           uint32 = 85
	HISTORICAL_NEWS                          uint32 = 86
	HISTORICAL_NEWS_END                      uint32 = 87
	HEAD_TIMESTAMP                           uint32 = 88
	HISTOGRAM_DATA                           uint32 = 89
	HISTORICAL_DATA_UPDATE                   uint32 = 90
	REROUTE_MKT_DATA_REQ                     uint32 = 91
	REROUTE_MKT_DEPTH_REQ                    uint32 = 92
	MARKET_RULE                              uint32 = 93
	PNL                                      uint32 = 94
	PNL_SINGLE                               uint32 = 95
	HISTORICAL_TICKS                         uint32 = 96
	HISTORICAL_TICKS_BID_ASK                 uint32 = 97
	HISTORICAL_TICKS_LAST                    uint32 = 98
	TICK_BY_TICK                             uint32 = 99
	ORDER_BOUND                              uint32 = 100
	COMPLETED_ORDER                          uint32 = 101
	COMPLETED_ORDERS_END                     uint32 = 102
	REPLACE_FA_END                           uint32 = 103
	WSH_META_DATA                            uint32 = 104
	WSH_EVENT_DATA                           uint32 = 105
	HISTORICAL_SCHEDULE                      uint32 = 106
	USER_INFO                                uint32 = 107
	HISTORICAL_DATA_END                      uint32 = 108
	CURRENT_TIME_IN_MILLIS                   uint32 = 109
)

const (
	PROTOBUF_MSG_ID uint32 = 200

	REQ_MKT_DATA                  = 1
	CANCEL_MKT_DATA               = 2
	PLACE_ORDER                   = 3
	CANCEL_ORDER                  = 4
	REQ_OPEN_ORDERS               = 5
	REQ_ACCT_DATA                 = 6
	REQ_EXECUTIONS                = 7
	REQ_IDS                       = 8
	REQ_CONTRACT_DATA             = 9
	REQ_MKT_DEPTH                 = 10
	CANCEL_MKT_DEPTH              = 11
	REQ_NEWS_BULLETINS            = 12
	CANCEL_NEWS_BULLETINS         = 13
	SET_SERVER_LOGLEVEL           = 14
	REQ_AUTO_OPEN_ORDERS          = 15
	REQ_ALL_OPEN_ORDERS           = 16
	REQ_MANAGED_ACCTS             = 17
	REQ_FA                        = 18
	REPLACE_FA                    = 19
	REQ_HISTORICAL_DATA           = 20
	EXERCISE_OPTIONS              = 21
	REQ_SCANNER_SUBSCRIPTION      = 22
	CANCEL_SCANNER_SUBSCRIPTION   = 23
	REQ_SCANNER_PARAMETERS        = 24
	CANCEL_HISTORICAL_DATA        = 25
	REQ_CURRENT_TIME              = 49
	REQ_REAL_TIME_BARS            = 50
	CANCEL_REAL_TIME_BARS         = 51
	REQ_FUNDAMENTAL_DATA          = 52
	CANCEL_FUNDAMENTAL_DATA       = 53
	REQ_CALC_IMPLIED_VOLAT        = 54
	REQ_CALC_OPTION_PRICE         = 55
	CANCEL_CALC_IMPLIED_VOLAT     = 56
	CANCEL_CALC_OPTION_PRICE      = 57
	REQ_GLOBAL_CANCEL             = 58
	REQ_MARKET_DATA_TYPE          = 59
	REQ_POSITIONS                 = 61
	REQ_ACCOUNT_SUMMARY           = 62
	CANCEL_ACCOUNT_SUMMARY        = 63
	CANCEL_POSITIONS              = 64
	VERIFY_REQUEST                = 65
	VERIFY_MESSAGE                = 66
	QUERY_DISPLAY_GROUPS          = 67
	SUBSCRIBE_TO_GROUP_EVENTS     = 68
	UPDATE_DISPLAY_GROUP          = 69
	UNSUBSCRIBE_FROM_GROUP_EVENTS = 70
	START_API                     = 71
	VERIFY_AND_AUTH_REQUEST       = 72
	VERIFY_AND_AUTH_MESSAGE       = 73
	REQ_POSITIONS_MULTI           = 74
	CANCEL_POSITIONS_MULTI        = 75
	REQ_ACCOUNT_UPDATES_MULTI     = 76
	CANCEL_ACCOUNT_UPDATES_MULTI  = 77
	REQ_SEC_DEF_OPT_PARAMS        = 78
	REQ_SOFT_DOLLAR_TIERS         = 79
	REQ_FAMILY_CODES              = 80
	REQ_MATCHING_SYMBOLS          = 81
	REQ_MKT_DEPTH_EXCHANGES       = 82
	REQ_SMART_COMPONENTS          = 83
	REQ_NEWS_ARTICLE              = 84
	REQ_NEWS_PROVIDERS            = 85
	REQ_HISTORICAL_NEWS           = 86
	REQ_HEAD_TIMESTAMP            = 87
	REQ_HISTOGRAM_DATA            = 88
	CANCEL_HISTOGRAM_DATA         = 89
	CANCEL_HEAD_TIMESTAMP         = 90
	REQ_MARKET_RULE               = 91
	REQ_PNL                       = 92
	CANCEL_PNL                    = 93
	REQ_PNL_SINGLE                = 94
	CANCEL_PNL_SINGLE             = 95
	REQ_HISTORICAL_TICKS          = 96
	REQ_TICK_BY_TICK_DATA         = 97
	CANCEL_TICK_BY_TICK_DATA      = 98
	REQ_COMPLETED_ORDERS          = 99
	REQ_WSH_META_DATA             = 100
	CANCEL_WSH_META_DATA          = 101
	REQ_WSH_EVENT_DATA            = 102
	CANCEL_WSH_EVENT_DATA         = 103
	REQ_USER_INFO                 = 104
	REQ_CURRENT_TIME_IN_MILLIS    = 105
)

var PROTOBUF_MSG_IDS = map[uint32]int32{
	REQ_EXECUTIONS:               ServerVersionProtobuf,
	PLACE_ORDER:                  ServerVersionProtobufPlaceOrder,
	CANCEL_ORDER:                 ServerVersionProtobufPlaceOrder,
	REQ_GLOBAL_CANCEL:            ServerVersionProtobufPlaceOrder,
	REQ_ALL_OPEN_ORDERS:          ServerVersionProtobufCompletedOrder,
	REQ_AUTO_OPEN_ORDERS:         ServerVersionProtobufCompletedOrder,
	REQ_OPEN_ORDERS:              ServerVersionProtobufCompletedOrder,
	REQ_COMPLETED_ORDERS:         ServerVersionProtobufCompletedOrder,
	REQ_CONTRACT_DATA:            ServerVersionProtobufContractData,
	REQ_MKT_DATA:                 ServerVersionProtobufMarketData,
	CANCEL_MKT_DATA:              ServerVersionProtobufMarketData,
	REQ_MKT_DEPTH:                ServerVersionProtobufMarketData,
	CANCEL_MKT_DEPTH:             ServerVersionProtobufMarketData,
	REQ_MARKET_DATA_TYPE:         ServerVersionProtobufMarketData,
	REQ_ACCT_DATA:                ServerVersionProtobufAccountsPositions,
	REQ_MANAGED_ACCTS:            ServerVersionProtobufAccountsPositions,
	REQ_POSITIONS:                ServerVersionProtobufAccountsPositions,
	CANCEL_POSITIONS:             ServerVersionProtobufAccountsPositions,
	REQ_ACCOUNT_SUMMARY:          ServerVersionProtobufAccountsPositions,
	CANCEL_ACCOUNT_SUMMARY:       ServerVersionProtobufAccountsPositions,
	REQ_POSITIONS_MULTI:          ServerVersionProtobufAccountsPositions,
	CANCEL_POSITIONS_MULTI:       ServerVersionProtobufAccountsPositions,
	REQ_ACCOUNT_UPDATES_MULTI:    ServerVersionProtobufAccountsPositions,
	CANCEL_ACCOUNT_UPDATES_MULTI: ServerVersionProtobufAccountsPositions,
	REQ_HISTORICAL_DATA:          ServerVersionProtobufHistoricalData,
	CANCEL_HISTORICAL_DATA:       ServerVersionProtobufHistoricalData,
	REQ_REAL_TIME_BARS:           ServerVersionProtobufHistoricalData,
	CANCEL_REAL_TIME_BARS:        ServerVersionProtobufHistoricalData,
	REQ_HEAD_TIMESTAMP:           ServerVersionProtobufHistoricalData,
	CANCEL_HEAD_TIMESTAMP:        ServerVersionProtobufHistoricalData,
	REQ_HISTOGRAM_DATA:           ServerVersionProtobufHistoricalData,
	CANCEL_HISTOGRAM_DATA:        ServerVersionProtobufHistoricalData,
	REQ_HISTORICAL_TICKS:         ServerVersionProtobufHistoricalData,
	REQ_TICK_BY_TICK_DATA:        ServerVersionProtobufHistoricalData,
	CANCEL_TICK_BY_TICK_DATA:     ServerVersionProtobufHistoricalData,
}
