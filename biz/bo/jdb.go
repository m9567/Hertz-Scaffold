package bo

type JdbKey struct {
	IV  string `json:"iv"`
	KEY string `json:"key"`
	DC  string `json:"dc"`
}

type JdbCallbackRequest struct {
	CallbackPathRequest
	X string `json:"x"`
}

type JdbRequestUrl struct {
	ApiRequest string `json:"apiRequest"`
}
