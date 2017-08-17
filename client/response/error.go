package response

type Error struct {
	ErrorMessage string `json:"error_message"`
	Error        int    `json:"errno"`
}
