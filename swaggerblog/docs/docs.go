// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-10-05 17:43:17.911062698 +0500 +05 m=+0.245508546

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
            "name": "Rishat Ishbulatov",
            "email": "progjb@gmail.com"
        },
        "license": {
            "name": "MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/docs/swagger.json": {
            "get": {
                "description": "Returns swagger's json docs",
                "produces": [
                    "application/json"
                ],
                "summary": "returns json docs",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/posts": {
            "post": {
                "description": "Create new post and return json",
                "produces": [
                    "application/json"
                ],
                "summary": "create new post",
                "parameters": [
                    {
                        "description": "title, author, content",
                        "name": "post",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.Post"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/posts/{id}": {
            "put": {
                "description": "renew post by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "renew post by id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of post",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "title, author, content",
                        "name": "post",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.Post"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "deletes a post by id",
                "summary": "deletes a post by id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of post",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/{template}": {
            "get": {
                "description": "return templated page if exist",
                "produces": [
                    "text/html"
                ],
                "summary": "serve templated page",
                "parameters": [
                    {
                        "type": "string",
                        "description": "in form like index.html",
                        "name": "template",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ErrorModel": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "desc": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                },
                "internal": {
                    "type": "object"
                }
            }
        },
        "models.Post": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "content": {
                    "type": "string"
                },
                "hex": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
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
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Mongoblog server",
	Description: "Blog server with swagger and tests for handlers.",
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