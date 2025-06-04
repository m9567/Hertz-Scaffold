package bo

type JdbKey struct {
	IV  string `json:"iv"`
	KEY string `json:"key"`
	DC  string `json:"dc"`
}

type JdbRequest struct {
	CallbackPathRequest
	X string `json:"x"`
}
