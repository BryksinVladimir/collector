package request

type PageLimit struct {
	Limit uint32 `url:"limit,omitempty"`
	Page  uint32 `url:"page,omitempty"`
}

func (this PageLimit) IsValid() bool {
	return this.Limit > 0 && this.Page > 0
}
