// Package docs provides Swagger documentation for USERFC API.
// Run `swag init -g main.go` in the USERFC directory to regenerate from annotations.
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "title": "USERFC API",
        "description": "User registration, authentication (JWT), and user info for Go Commerce.",
        "version": "1.0"
    },
    "host": "localhost:28080",
    "basePath": "/",
    "paths": {
        "/ping": {
            "get": {
                "description": "Health check",
                "produces": ["application/json"],
                "summary": "Ping",
                "responses": {
                    "200": {"description": "pong"}
                }
            }
        },
        "/v1/register": {
            "post": {
                "description": "Register a new user",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "summary": "Register",
                "parameters": [{
                    "in": "body",
                    "name": "body",
                    "required": true,
                    "schema": {
                        "type": "object",
                        "required": ["name","email","password","confirm_password"],
                        "properties": {
                            "name": {"type": "string"},
                            "email": {"type": "string"},
                            "password": {"type": "string"},
                            "confirm_password": {"type": "string"}
                        }
                    }
                }],
                "responses": {
                    "201": {"description": "User registered successfully"},
                    "400": {"description": "Bad request"},
                    "500": {"description": "Internal server error"}
                }
            }
        },
        "/v1/login": {
            "post": {
                "description": "Login and get JWT token",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "summary": "Login",
                "parameters": [{
                    "in": "body",
                    "name": "body",
                    "required": true,
                    "schema": {
                        "type": "object",
                        "required": ["email","password"],
                        "properties": {
                            "email": {"type": "string"},
                            "password": {"type": "string"}
                        }
                    }
                }],
                "responses": {
                    "200": {"description": "Returns token", "schema": {"type": "object", "properties": {"token": {"type": "string"}}}},
                    "400": {"description": "Bad request"},
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/api/v1/user-info": {
            "get": {
                "security": [{"BearerAuth": []}],
                "description": "Get current user info (requires JWT)",
                "produces": ["application/json"],
                "summary": "Get user info",
                "responses": {
                    "200": {"description": "name, email"},
                    "401": {"description": "Unauthorized"}
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

func init() {
	swag.Register(swag.Name, &s{})
}

type s struct{}

func (s *s) ReadDoc() string {
	return docTemplate
}
