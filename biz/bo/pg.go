package bo

type PgBase struct {
	CallbackPathRequest
	TraceId               string `json:"trace_id"`
	PlayerName            string `json:"player_name"`
	OperatorPlayerSession string `json:"operator_player_session"`
	OperatorToken         string `json:"operator_token"`
}
