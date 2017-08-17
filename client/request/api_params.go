package request

type ApiParams struct {
	ApiHash   string `url:"hash,omitempty"`
	ApiFeedId string `url:"feed_id,omitempty"`
	Format    string `url:"format,omitempty"`
}

func (this ApiParams) IsValid() bool {
	return len(this.ApiHash) > 0 && len(this.ApiFeedId) > 0 && this.Format == "json"
}
