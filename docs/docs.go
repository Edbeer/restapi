// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/auth/all": {
            "get": {
                "description": "Get the list of all users",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Get users",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "page",
                        "description": "page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "format": "size",
                        "description": "number of elements per page",
                        "name": "size",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "format": "orderBy",
                        "description": "filter name",
                        "name": "orderBy",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.UsersList"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpe.RestError"
                        }
                    }
                }
            }
        },
        "/auth/find": {
            "get": {
                "description": "Find user by name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Find by name",
                "parameters": [
                    {
                        "type": "string",
                        "format": "username",
                        "description": "username",
                        "name": "name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.UsersList"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpe.RestError"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "login user, returns user and set session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login new user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.User"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "description": "logout user removing session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Logout user",
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
        "/auth/me": {
            "get": {
                "description": "Get current user by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Get user by id",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.User"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpe.RestError"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "register new user, returns user and token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Register new user",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/entity.User"
                        }
                    }
                }
            }
        },
        "/auth/token": {
            "get": {
                "description": "Get CSRF token, required auth session cookie",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Get CSRF token",
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpe.RestError"
                        }
                    }
                }
            }
        },
        "/auth/{id}": {
            "get": {
                "description": "get string by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "get user by id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user_id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.User"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpe.RestError"
                        }
                    }
                }
            },
            "put": {
                "description": "update existing user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Update user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user_id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.User"
                        }
                    }
                }
            },
            "delete": {
                "description": "some description",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Delete user account",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user_id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpe.RestError"
                        }
                    }
                }
            }
        },
        "/comments": {
            "post": {
                "description": "create new comment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Comments"
                ],
                "summary": "Create new comment",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/entity.Comment"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/comments/byNewsId/{id}": {
            "get": {
                "description": "Get all comment by news id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Comments"
                ],
                "summary": "Get comments by news",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "news_id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "format": "page",
                        "description": "page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "format": "size",
                        "description": "number of elements per page",
                        "name": "size",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "format": "orderBy",
                        "description": "filter name",
                        "name": "orderBy",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.CommentsList"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/comments/{id}": {
            "get": {
                "description": "Get comment by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Comments"
                ],
                "summary": "Get comment",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "comment_id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.Comment"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "put": {
                "description": "update new comment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Comments"
                ],
                "summary": "Update comment",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "comment_id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.Comment"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "delete": {
                "description": "delete comment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Comments"
                ],
                "summary": "Delete comment",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "comment_id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/news": {
            "get": {
                "description": "Get all news with pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "News"
                ],
                "summary": "Get all news",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "page",
                        "description": "page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "format": "size",
                        "description": "number of elements per page",
                        "name": "size",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "format": "orderBy",
                        "description": "filter name",
                        "name": "orderBy",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.NewsList"
                        }
                    }
                }
            }
        },
        "/news/create": {
            "post": {
                "description": "Create news handler",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "News"
                ],
                "summary": "Create news",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/entity.News"
                        }
                    }
                }
            }
        },
        "/news/search": {
            "get": {
                "description": "Search news by title",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "News"
                ],
                "summary": "Search by title",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "page",
                        "description": "page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "format": "size",
                        "description": "number of elements per page",
                        "name": "size",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "format": "orderBy",
                        "description": "filter name",
                        "name": "orderBy",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.NewsList"
                        }
                    }
                }
            }
        },
        "/news/{id}": {
            "get": {
                "description": "Get by id news handler",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "News"
                ],
                "summary": "Get by id news",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "news_id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.News"
                        }
                    }
                }
            },
            "put": {
                "description": "Update news handler",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "News"
                ],
                "summary": "Update news",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "news_id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.News"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete by id news handler",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "News"
                ],
                "summary": "Delete news",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "news_id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.Comment": {
            "type": "object",
            "required": [
                "author_id",
                "message",
                "news_id"
            ],
            "properties": {
                "author_id": {
                    "type": "string"
                },
                "comment_id": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "likes": {
                    "type": "integer"
                },
                "message": {
                    "type": "string",
                    "minLength": 5
                },
                "news_id": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "entity.CommentBase": {
            "type": "object",
            "required": [
                "author",
                "author_id",
                "message"
            ],
            "properties": {
                "author": {
                    "type": "string"
                },
                "author_id": {
                    "type": "string"
                },
                "avatar_url": {
                    "type": "string"
                },
                "comment_id": {
                    "type": "string"
                },
                "likes": {
                    "type": "integer"
                },
                "message": {
                    "type": "string",
                    "minLength": 5
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "entity.CommentsList": {
            "type": "object",
            "properties": {
                "comments": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.CommentBase"
                    }
                },
                "has_more": {
                    "type": "boolean"
                },
                "page": {
                    "type": "integer"
                },
                "size": {
                    "type": "integer"
                },
                "total_count": {
                    "type": "integer"
                },
                "total_pages": {
                    "type": "integer"
                }
            }
        },
        "entity.News": {
            "type": "object",
            "required": [
                "author_id",
                "content",
                "title"
            ],
            "properties": {
                "author_id": {
                    "type": "string"
                },
                "category": {
                    "type": "string",
                    "maxLength": 10
                },
                "content": {
                    "type": "string",
                    "minLength": 20
                },
                "created_at": {
                    "type": "string"
                },
                "image_url": {
                    "type": "string",
                    "maxLength": 512
                },
                "news_id": {
                    "type": "string"
                },
                "title": {
                    "type": "string",
                    "minLength": 10
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "entity.NewsList": {
            "type": "object",
            "properties": {
                "has_more": {
                    "type": "boolean"
                },
                "news": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.News"
                    }
                },
                "page": {
                    "type": "integer"
                },
                "size": {
                    "type": "integer"
                },
                "total_count": {
                    "type": "integer"
                },
                "total_pages": {
                    "type": "integer"
                }
            }
        },
        "entity.User": {
            "type": "object",
            "required": [
                "password"
            ],
            "properties": {
                "address": {
                    "type": "string",
                    "maxLength": 250
                },
                "avatar": {
                    "type": "string"
                },
                "balance": {
                    "type": "number"
                },
                "city": {
                    "type": "string",
                    "maxLength": 24
                },
                "country": {
                    "type": "string",
                    "maxLength": 24
                },
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string",
                    "maxLength": 60
                },
                "first_name": {
                    "type": "string",
                    "maxLength": 30
                },
                "last_name": {
                    "type": "string",
                    "maxLength": 30
                },
                "password": {
                    "type": "string",
                    "minLength": 6
                },
                "phone_number": {
                    "type": "string",
                    "maxLength": 20
                },
                "postcode": {
                    "type": "integer",
                    "maximum": 10
                },
                "role": {
                    "type": "string",
                    "maxLength": 10
                },
                "updated_at": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "entity.UsersList": {
            "type": "object",
            "properties": {
                "has_more": {
                    "type": "boolean"
                },
                "page": {
                    "type": "integer"
                },
                "size": {
                    "type": "integer"
                },
                "total_count": {
                    "type": "integer"
                },
                "total_pages": {
                    "type": "integer"
                },
                "users": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.User"
                    }
                }
            }
        },
        "httpe.RestError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.1",
	Host:             "",
	BasePath:         "/api/",
	Schemes:          []string{},
	Title:            "restapi",
	Description:      "This is an example of an implementation RESTApi",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
