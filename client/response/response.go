package response

type Pagination struct {
	TotalRows   uint32 `json:"total_rows"`
	CurrentRows uint64 `json:"current_rows"`
	CurrentPage uint32 `json:"current_page"`
	TotalPages  uint32 `json:"total_pages"`
	Limit       uint32 `json:"limit"`
}
