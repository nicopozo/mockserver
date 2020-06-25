// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2020-06-25 16:47:25.243629 -0300 -03 m=+0.048626854

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Nicolas Pozo",
            "email": "nicopozo@gmail.com"
        },
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/rules": {
            "get": {
                "description": "Search Rule by key, name, application, method or status",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rules"
                ],
                "summary": "Search Rule",
                "operationId": "search-rule",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Rule key generated by service",
                        "name": "key",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Name of the key",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Application",
                        "name": "application",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Method",
                        "name": "method",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Enabled/Disabled",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "number",
                        "default": 30,
                        "description": "Max expected number of results",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "0",
                        "description": "number of results to be skipped",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Result",
                        "schema": {
                            "$ref": "#/definitions/model.RuleList"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a Rule for serving a mock response",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rules"
                ],
                "summary": "Create a Rule",
                "operationId": "create-rule",
                "parameters": [
                    {
                        "description": "The rule to be created",
                        "name": "rule",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Rule"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Rule successfully created",
                        "schema": {
                            "$ref": "#/definitions/model.Rule"
                        }
                    },
                    "400": {
                        "description": "Validation of the rule failed",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    }
                }
            }
        },
        "/rules/{key}": {
            "get": {
                "description": "Get a Rule, if not found return 404",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rules"
                ],
                "summary": "Get Rule by Key",
                "operationId": "get-rule",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Key generated by service",
                        "name": "key",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Result",
                        "schema": {
                            "$ref": "#/definitions/model.Rule"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    }
                }
            },
            "put": {
                "description": "Update an existing Rule for serving a mock response",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rules"
                ],
                "summary": "Update a Rule",
                "operationId": "update-rule",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Key of rule to update",
                        "name": "key",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "The rule to be updated",
                        "name": "rule",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Rule"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Rule successfully updated",
                        "schema": {
                            "$ref": "#/definitions/model.Rule"
                        }
                    },
                    "400": {
                        "description": "Validation of the rule failed",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete Rule by Key",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rules"
                ],
                "summary": "Delete Rule by key",
                "operationId": "delete-rule",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Key generated by service",
                        "name": "key",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {},
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    }
                }
            }
        },
        "/rules/{key}/status": {
            "put": {
                "description": "Update an existing Rule for serving a mock response",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rules"
                ],
                "summary": "Update a Rule Status",
                "operationId": "update-rule",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Key of rule to update",
                        "name": "key",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "The rule to be updated",
                        "name": "rule",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.RuleStatus"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Rule successfully updated",
                        "schema": {
                            "$ref": "#/definitions/model.Rule"
                        }
                    },
                    "400": {
                        "description": "Validation of the rule failed",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Error": {
            "type": "object",
            "properties": {
                "cause": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.ErrorCause"
                    }
                },
                "error": {
                    "type": "string",
                    "example": "Not Found"
                },
                "message": {
                    "type": "string",
                    "example": "no rule found with key: banks_get_55603295"
                },
                "status": {
                    "type": "integer",
                    "example": 404
                }
            }
        },
        "model.ErrorCause": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 1030
                },
                "description": {
                    "type": "string",
                    "example": "Resource Not Found"
                }
            }
        },
        "model.Paging": {
            "type": "object",
            "properties": {
                "limit": {
                    "type": "integer",
                    "default": 30,
                    "maximum": 1000,
                    "minimum": 0
                },
                "offset": {
                    "type": "integer",
                    "default": 0,
                    "maximum": 1000,
                    "minimum": 0
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "model.Response": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "string",
                    "example": "{\"id\":5804214224, \"payer_id\": 548390723, \"external_reference\": \"X281924481\"}"
                },
                "content_type": {
                    "type": "string",
                    "example": "application/json"
                },
                "delay": {
                    "type": "integer",
                    "example": 0
                },
                "http_status": {
                    "type": "integer",
                    "example": 200
                }
            }
        },
        "model.Rule": {
            "type": "object",
            "properties": {
                "application": {
                    "type": "string",
                    "example": "payments"
                },
                "key": {
                    "type": "string",
                    "example": "payments_get_556032950"
                },
                "method": {
                    "type": "string",
                    "example": "GET"
                },
                "name": {
                    "type": "string",
                    "example": "get payment"
                },
                "path": {
                    "type": "string",
                    "example": "/v1/payments/{payment_id}"
                },
                "responses": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Response"
                    }
                },
                "status": {
                    "type": "string",
                    "example": "enabled"
                },
                "strategy": {
                    "type": "string",
                    "example": "normal"
                },
                "variables": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Variable"
                    }
                }
            }
        },
        "model.RuleList": {
            "type": "object",
            "properties": {
                "paging": {
                    "type": "object",
                    "$ref": "#/definitions/model.Paging"
                },
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Rule"
                    }
                }
            }
        },
        "model.RuleStatus": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "enabled"
                }
            }
        },
        "model.Variable": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string",
                    "example": "$.nickname"
                },
                "name": {
                    "type": "string",
                    "example": "nickname"
                },
                "type": {
                    "type": "string",
                    "example": "body"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:8080",
	BasePath:    "/mock-server",
	Schemes:     []string{"http"},
	Title:       "Mock Server",
	Description: "Mock Server is intended to serve mocks for different APIs during development process.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
