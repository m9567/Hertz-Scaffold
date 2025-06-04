package bo

type CallbackPathRequest struct {
	Currency string `path:"currency,required" json:"currency"`
}
