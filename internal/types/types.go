package types

type Link struct {
	Short     string `json:"short"`
	Long      string `json:"long"`
	Created   int64  `json:"created"`
	CreatedBy string `json:"created_by"`
	Clicks    int    `json:"clicks"`
}
