package model

type Response struct {
	Status  int         `json:"status" `
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseLogin struct {
	Status  int    `json:"status" `
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}
