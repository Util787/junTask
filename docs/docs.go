// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/users": {
            "get": {
                "description": "Get users using flexible query filters and pagination. You can provide partial values for ` + "`" + `name` + "`" + `, ` + "`" + `surname` + "`" + `, or ` + "`" + `patronymic` + "`" + ` — filtering will still work. Each of these parameters is optional and can be used independently or in combination.\n\nExample: ?page=5\u0026page_size=10\nResponse: 10 users with offset=40\n\nExample2: ?name=al\nResponse: Alex, Alina, etc.\n\nExample3: ?name=al\u0026surname=sh\nResponse: Alexandr Shprot, Alina Sham, etc.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "get all users with optionally filters and pagination",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name filter",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "surname filter",
                        "name": "surname",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "patronymic filter",
                        "name": "patronymic",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "gender filter can be only male or female",
                        "name": "gender",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "min:5",
                        "name": "page_size",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "min:1",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.User"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/internal_handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal_handlers.errorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "creating new user with provided name, surname, patronymic(optional)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "create user",
                "parameters": [
                    {
                        "description": "Users fullname: name, surname, patronymic(optional)",
                        "name": "fullname",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.FullName"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "message with created user's id",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/internal_handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal_handlers.errorResponse"
                        }
                    }
                }
            }
        },
        "/users/{user_id}": {
            "get": {
                "description": "recieve user info by providing id in path",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "get user by id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user_id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/internal_handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal_handlers.errorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "deleting user by id if exists",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "delete user by id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user_id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "successful deleting message",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/internal_handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal_handlers.errorResponse"
                        }
                    }
                }
            },
            "patch": {
                "description": "updating user info by id provided in path. In request body you can optionally provide: name, surname, patronymic, age, gender, nationality. Update_at will change automatically",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "update user info by id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user_id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "parameters for update",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.UpdateUserParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message about user update",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/internal_handlers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal_handlers.errorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entities.FullName": {
            "type": "object",
            "required": [
                "name",
                "surname"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "patronymic": {
                    "type": "string"
                },
                "surname": {
                    "type": "string"
                }
            }
        },
        "entities.UpdateUserParams": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "integer"
                },
                "gender": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "nationality": {
                    "type": "string"
                },
                "patronymic": {
                    "type": "string"
                },
                "surname": {
                    "type": "string"
                }
            }
        },
        "entities.User": {
            "type": "object",
            "required": [
                "name",
                "surname"
            ],
            "properties": {
                "age": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "nationality": {
                    "type": "string"
                },
                "patronymic": {
                    "type": "string"
                },
                "surname": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "internal_handlers.errorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "User manager api",
	Description:      "Rest api for managing users crud operations",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
