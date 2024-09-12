package model

// 搜索参数
type SearchParams struct {
	Index     string  `json:"Index"`
	StartTime string  `json:"startTime"`
	EndTime   string  `json:"endTime"`
	Keyword   string  `json:"keyword"`
	Search    string  `json:"search"`
	// 其他搜索字段
}
