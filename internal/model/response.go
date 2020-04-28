package model

type Response struct {
	Body        string `json:"body"`
	ContentType string `json:"content_type"`
	HTTPStatus  int    `json:"http_status"`
}
