package model

type Response struct {
	Body        string `json:"body" example:"{\"id\":5804214224, \"payer_id\": 548390723, \"external_reference\": \"X281924481\"}"` //nolint:lll
	ContentType string `json:"content_type" example:"application/json"`
	HTTPStatus  int    `json:"http_status" example:"200"`
	Delay       int    `json:"delay" example:"0"`
	Scene       string `json:"scene" example:"normal"`
}
