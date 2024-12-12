package models

// EnvironmentalImpact は環境影響の分析データを表します
type EnvironmentalImpact struct {
	TotalCO2Saved          float64
	AverageCO2SavedPerItem float64
	TotalEcoFriendlyItems  int
}
