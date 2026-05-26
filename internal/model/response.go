package model

// WebhookConfig represents the configuration for sending a webhook when a response is returned.
type WebhookConfig struct {
	URL     string            `json:"url" example:"https://hooks.example.com/callback"`
	Method  string            `json:"method" example:"POST"`
	Headers map[string]string `json:"headers" example:"{\"Authorization\":\"Bearer token\"}"`
	Body    string            `json:"body" example:"{\"event\":\"payment_created\",\"id\":\"{payment_id}\"}"`
	Enabled bool              `json:"enabled" example:"true"`
	Timeout *int              `json:"timeout,omitempty" example:"5000"`
}

type Response struct {
	Body        string         `json:"body" example:"{\"id\":5804214224, \"payer_id\": 548390723, \"external_reference\": \"X281924481\"}"` //nolint:lll
	ContentType string         `json:"content_type" example:"application/json"`
	HTTPStatus  int            `json:"http_status" example:"200"`
	Delay       int            `json:"delay" example:"0"`
	Scene       string         `json:"scene" example:"normal"`
	Description string         `json:"description" example:"success response"`
	Webhook     *WebhookConfig `json:"webhook,omitempty"`
}
