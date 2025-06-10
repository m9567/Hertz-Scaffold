package constant

// 用于定义静态数据

const (
	InitTrade = iota
	Trading
	TradeOver
)

const (
	CallbackAPIModule = "callbackAPI"
	DefaultAPIModule  = "defaultAPI"
	DevOpsAPIModule   = "devOpsAPI"
	AdminApIModel     = "adminAPI"
	InnerApIModel     = "innerAPI"

	CallbackURLPrefix = "/callback"
	DefaultURLPrefix  = "/app-api"
	DevOpsURLPrefix   = "/devops-api"
	AdminURLPrefix    = "/admin-api"
	InnerURLPrefix    = "/inner-api"
)

const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH" // RFC 5789
	MethodDelete  = "DELETE"
	MethodOptions = "OPTIONS"
)
