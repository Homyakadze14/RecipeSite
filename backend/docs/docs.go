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
        "/auth/checktgtoken": {
            "post": {
                "description": "Check user telegram token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Check user telegram token",
                "operationId": "Check user telegram token",
                "parameters": [
                    {
                        "description": "token",
                        "name": "token",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.JWTToken"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.JWTData"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "description": "Logout user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Logout",
                "operationId": "Logout",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/auth/signin": {
            "post": {
                "description": "Sign in user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Sign in",
                "operationId": "signin",
                "parameters": [
                    {
                        "description": "User params",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.UserLogin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "login",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/auth/signup": {
            "post": {
                "description": "Sign up user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Sign up",
                "operationId": "signup",
                "parameters": [
                    {
                        "description": "User params",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.UserLogin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/auth/tgtoken": {
            "get": {
                "description": "Generate user telegram token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Generate user telegram token",
                "operationId": "generate user telegram token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.JWTToken"
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/recipe": {
            "get": {
                "description": "Get all recipe",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipe"
                ],
                "summary": "Get all recipe",
                "operationId": "get all recipe",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.Recipe"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "post": {
                "description": "Get filtered recipe",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipe"
                ],
                "summary": "Get filtered recipe",
                "operationId": "get filtered recipe",
                "parameters": [
                    {
                        "description": "filter",
                        "name": "filter",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/entities.RecipeFilter"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.Recipe"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/recipe/{id}": {
            "get": {
                "description": "Get recipe",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipe"
                ],
                "summary": "Get recipe",
                "operationId": "get recipe",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.RecipeInfo"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/recipe/{id}/comment": {
            "put": {
                "description": "Update comment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "comments"
                ],
                "summary": "Update comment",
                "operationId": "update comment",
                "parameters": [
                    {
                        "description": "Comment params",
                        "name": "comment",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.CommentUpdate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "post": {
                "description": "Create comment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "comments"
                ],
                "summary": "Create comment",
                "operationId": "create comment",
                "parameters": [
                    {
                        "description": "Comment params",
                        "name": "comment",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.CommentCreate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "delete": {
                "description": "Delete comment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "comments"
                ],
                "summary": "Delete comment",
                "operationId": "delete comment",
                "parameters": [
                    {
                        "description": "Comment params",
                        "name": "comment",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.CommentDelete"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/recipe/{id}/like": {
            "post": {
                "description": "Like recipe",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "likes"
                ],
                "summary": "Like",
                "operationId": "like",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/recipe/{id}/unlike": {
            "post": {
                "description": "Unlike recipe",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "likes"
                ],
                "summary": "Unlike",
                "operationId": "unlike",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/user/subscribe": {
            "post": {
                "description": "Subscribe to user",
                "produces": [
                    "application/json",
                    "application/json"
                ],
                "tags": [
                    "subscription"
                ],
                "summary": "Subscribe to user",
                "operationId": "subscribe to user",
                "parameters": [
                    {
                        "description": "User id to whom we subscribe",
                        "name": "creator",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.SubscribeCreator"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/user/unsubscribe": {
            "post": {
                "description": "Unsubscribe from user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscription"
                ],
                "summary": "Unsubscribe from user",
                "operationId": "unsubscribe from user",
                "parameters": [
                    {
                        "description": "User id to whom we unsubscribe",
                        "name": "creator",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.SubscribeCreator"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/user/{login}": {
            "get": {
                "description": "Get user info",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get user info",
                "operationId": "get user info",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.JSONUserInfo"
                        }
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "put": {
                "description": "Update user",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Update user",
                "operationId": "update user",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Icon",
                        "name": "icon",
                        "in": "formData"
                    },
                    {
                        "maxLength": 1500,
                        "type": "string",
                        "name": "about",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "name": "email",
                        "in": "formData"
                    },
                    {
                        "maxLength": 20,
                        "minLength": 3,
                        "type": "string",
                        "name": "login",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "login",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/user/{login}/password": {
            "put": {
                "description": "Update user password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Update user password",
                "operationId": "update user password",
                "parameters": [
                    {
                        "description": "User params",
                        "name": "user",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/entities.UserPasswordUpdate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/user/{login}/recipe": {
            "post": {
                "description": "Create recipe",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipe"
                ],
                "summary": "Create recipe",
                "operationId": "create recipe",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Photos",
                        "name": "photos",
                        "in": "formData"
                    },
                    {
                        "maxLength": 2500,
                        "type": "string",
                        "name": "about",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "maximum": 3,
                        "minimum": 1,
                        "enum": [
                            1,
                            2,
                            3
                        ],
                        "type": "integer",
                        "name": "complexitiy",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "maxLength": 1500,
                        "type": "string",
                        "name": "ingridients",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "need_time",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "maxLength": 50,
                        "minLength": 3,
                        "type": "string",
                        "name": "title",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/user/{login}/recipe/{id}": {
            "put": {
                "description": "Update recipe",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipe"
                ],
                "summary": "Update recipe",
                "operationId": "update recipe",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Photos",
                        "name": "photos",
                        "in": "formData"
                    },
                    {
                        "maxLength": 2500,
                        "type": "string",
                        "name": "about",
                        "in": "formData"
                    },
                    {
                        "maximum": 3,
                        "minimum": 1,
                        "enum": [
                            1,
                            2,
                            3
                        ],
                        "type": "integer",
                        "name": "complexitiy",
                        "in": "formData"
                    },
                    {
                        "maxLength": 1500,
                        "type": "string",
                        "name": "ingridients",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "name": "need_time",
                        "in": "formData"
                    },
                    {
                        "maxLength": 50,
                        "minLength": 3,
                        "type": "string",
                        "name": "title",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "delete": {
                "description": "Delete recipe",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipe"
                ],
                "summary": "Delete recipe",
                "operationId": "delete recipe",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "entities.Author": {
            "type": "object",
            "properties": {
                "icon_url": {
                    "type": "string"
                },
                "login": {
                    "type": "string"
                }
            }
        },
        "entities.Comment": {
            "type": "object",
            "required": [
                "text"
            ],
            "properties": {
                "author": {
                    "$ref": "#/definitions/entities.Author"
                },
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "text": {
                    "type": "string",
                    "maxLength": 250,
                    "minLength": 1
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "entities.CommentCreate": {
            "type": "object",
            "required": [
                "text"
            ],
            "properties": {
                "text": {
                    "type": "string",
                    "maxLength": 250,
                    "minLength": 1
                }
            }
        },
        "entities.CommentDelete": {
            "type": "object",
            "required": [
                "id"
            ],
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "entities.CommentUpdate": {
            "type": "object",
            "required": [
                "id",
                "text"
            ],
            "properties": {
                "id": {
                    "type": "integer"
                },
                "text": {
                    "type": "string",
                    "maxLength": 250,
                    "minLength": 1
                }
            }
        },
        "entities.FullRecipe": {
            "type": "object",
            "properties": {
                "author": {
                    "$ref": "#/definitions/entities.Author"
                },
                "comments": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Comment"
                    }
                },
                "is_liked": {
                    "type": "boolean"
                },
                "likes_count": {
                    "type": "integer"
                },
                "recipe": {
                    "$ref": "#/definitions/entities.Recipe"
                }
            }
        },
        "entities.JSONUserInfo": {
            "type": "object",
            "properties": {
                "user": {
                    "$ref": "#/definitions/entities.UserInfo"
                }
            }
        },
        "entities.JWTData": {
            "type": "object",
            "properties": {
                "user_id": {}
            }
        },
        "entities.JWTToken": {
            "type": "object",
            "required": [
                "token"
            ],
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "entities.Recipe": {
            "type": "object",
            "required": [
                "about",
                "complexitiy",
                "ingridients",
                "need_time",
                "title"
            ],
            "properties": {
                "about": {
                    "type": "string",
                    "maxLength": 2500
                },
                "complexitiy": {
                    "type": "integer",
                    "maximum": 3,
                    "minimum": 1,
                    "enum": [
                        1,
                        2,
                        3
                    ]
                },
                "created_at": {
                    "type": "string"
                },
                "creator_user_id": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "ingridients": {
                    "type": "string",
                    "maxLength": 1500
                },
                "need_time": {
                    "type": "string"
                },
                "photos_urls": {
                    "type": "string"
                },
                "title": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 3
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "entities.RecipeFilter": {
            "type": "object",
            "properties": {
                "limit": {
                    "type": "integer",
                    "example": 25
                },
                "offset": {
                    "type": "integer",
                    "example": 0
                },
                "order_by": {
                    "type": "integer",
                    "maximum": 1,
                    "minimum": -1,
                    "enum": [
                        -1,
                        0,
                        1
                    ]
                },
                "order_field": {
                    "type": "string",
                    "enum": [
                        "title",
                        "about",
                        "ingridients",
                        "emtpy"
                    ],
                    "example": "title"
                },
                "query": {
                    "type": "string",
                    "example": "tasty food"
                }
            }
        },
        "entities.RecipeInfo": {
            "type": "object",
            "properties": {
                "info": {
                    "$ref": "#/definitions/entities.FullRecipe"
                }
            }
        },
        "entities.SubscribeCreator": {
            "type": "object",
            "required": [
                "creator_id"
            ],
            "properties": {
                "creator_id": {
                    "type": "integer"
                }
            }
        },
        "entities.UserInfo": {
            "type": "object",
            "properties": {
                "about": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "icon_url": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "liked_recipies": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Recipe"
                    }
                },
                "login": {
                    "type": "string"
                },
                "recipies": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Recipe"
                    }
                }
            }
        },
        "entities.UserLogin": {
            "type": "object",
            "required": [
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "test@test.com"
                },
                "login": {
                    "type": "string",
                    "example": "testuser"
                },
                "password": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 8,
                    "example": "testpassword"
                }
            }
        },
        "entities.UserPasswordUpdate": {
            "type": "object",
            "required": [
                "password"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 8,
                    "example": "testpassword"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "RecipeSite",
	Description:      "RestAPI for recipe site",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
