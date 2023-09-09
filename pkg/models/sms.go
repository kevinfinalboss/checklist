package models

type SmsRequest struct {
	Key    string `json:"key"`
	Type   int    `json:"type"`
	Number string `json:"number"`
	Msg    string `json:"msg"`
}

type SmsResponse struct {
	Situacao  string `json:"situacao"`
	Codigo    string `json:"codigo"`
	ID        string `json:"id"`
	Descricao string `json:"descricao"`
}
