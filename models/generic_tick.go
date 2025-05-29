package models

// -----------------------------------------------------------------------------

type GenericTick int

const (
	GenericTickOptionCallAndPutVolume                       GenericTick = 100
	GenericTickOptionCallAndPutOpenInterest                 GenericTick = 101
	GenericTickOptionHistoricalVolatility                   GenericTick = 104
	GenericTickAverageOptionVolume                          GenericTick = 105
	GenericTickOptionImpliedVolatility                      GenericTick = 106
	GenericTickIndexFuturePremium                           GenericTick = 162
	GenericTickLowHigh13Weeks26Weeks52WeeksAndAverageVolume GenericTick = 165
	GenericTickAuctionVolumePriceAndImbalance               GenericTick = 225
	GenericTickMarkPrice                                    GenericTick = 232
	GenericTickShortableAndShortableShares                  GenericTick = 236
	GenericTickTradeCount                                   GenericTick = 293
	GenericTickTradeRate                                    GenericTick = 294
	GenericTickVolumeRate                                   GenericTick = 295
	GenericTickLastRTHTrade                                 GenericTick = 318
	GenericTickRTTradeVolume                                GenericTick = 375
	GenericTickRTHistoricalVolatility                       GenericTick = 411
	GenericTickIBDividends                                  GenericTick = 456
	GenericTickBondFactorMultiplier                         GenericTick = 460
	GenericTickETFNavBidAndAsk                              GenericTick = 576
	GenericTickETFNavLast                                   GenericTick = 577
	GenericTickETFNavCloseAndPriorClose                     GenericTick = 578
	GenericTickEstimatedIPOMidpointAndFinalIPOPrice         GenericTick = 586
	GenericTickFuturesOpenInterest                          GenericTick = 588
	GenericTickETFNavHighAndLow                             GenericTick = 614
)
