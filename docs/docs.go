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
        "license": {
            "name": "TMPS"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/dog": {
            "post": {
                "description": "Register new dog",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dog"
                ],
                "summary": "Register new dog",
                "parameters": [
                    {
                        "description": "dog",
                        "name": "dog",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/model.DogRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.DogResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "allOf": [
                                    {
                                        "type": "object"
                                    },
                                    {
                                        "type": "object",
                                        "properties": {
                                            " message": {
                                                "type": "string"
                                            },
                                            "error": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                ]
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "allOf": [
                                    {
                                        "type": "object"
                                    },
                                    {
                                        "type": "object",
                                        "properties": {
                                            " message": {
                                                "type": "string"
                                            },
                                            "error": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                ]
                            }
                        }
                    }
                }
            },
            "patch": {
                "description": "Updates existing dog",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dog"
                ],
                "summary": "Updates dog",
                "parameters": [
                    {
                        "description": "dog",
                        "name": "dog",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/model.DogRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.DogResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "allOf": [
                                    {
                                        "type": "string"
                                    },
                                    {
                                        "type": "object",
                                        "properties": {
                                            " message": {
                                                "type": "string"
                                            },
                                            "error": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                ]
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "allOf": [
                                    {
                                        "type": "string"
                                    },
                                    {
                                        "type": "object",
                                        "properties": {
                                            " message": {
                                                "type": "string"
                                            },
                                            "error": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                ]
                            }
                        }
                    }
                }
            }
        },
        "/dog/found": {
            "patch": {
                "description": "Reunite dog with owner. Making him its only host and removing other hosts",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dog"
                ],
                "summary": "Reunite dog with owner",
                "parameters": [
                    {
                        "type": "string",
                        "description": "dog ID",
                        "name": "dogID",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "dog owner ID",
                        "name": "ownerID",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "dog host ID",
                        "name": "hostID",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.DogResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "allOf": [
                                    {
                                        "type": "string"
                                    },
                                    {
                                        "type": "object",
                                        "properties": {
                                            " message": {
                                                "type": "string"
                                            },
                                            "error": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                ]
                            }
                        }
                    }
                }
            }
        },
        "/dog/missing": {
            "get": {
                "description": "If no argument is given it returns all missing dogs. If user location and a search radius is sent, then it returns all missing dogs within that radius.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dog"
                ],
                "summary": "Brings list of missing dogs",
                "parameters": [
                    {
                        "type": "number",
                        "description": "user latitude",
                        "name": "userLatitude",
                        "in": "query"
                    },
                    {
                        "type": "number",
                        "description": "user longitude",
                        "name": "userLongitude",
                        "in": "query"
                    },
                    {
                        "type": "number",
                        "description": "radio to look for dogs",
                        "name": "radius",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.DogResponse"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "allOf": [
                                    {
                                        "type": "string"
                                    },
                                    {
                                        "type": "object",
                                        "properties": {
                                            " message": {
                                                "type": "string"
                                            },
                                            "error": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                ]
                            }
                        }
                    }
                }
            }
        },
        "/img": {
            "post": {
                "description": "Add image to storage",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "imgs"
                ],
                "summary": "Add Image",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        },
                        "description": "dog img to save",
                        "name": "img",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "name\" saved correctly!",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "allOf": [
                                    {
                                        "type": "object"
                                    },
                                    {
                                        "type": "object",
                                        "properties": {
                                            " message": {
                                                "type": "string"
                                            },
                                            "error": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                ]
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "allOf": [
                                    {
                                        "type": "object"
                                    },
                                    {
                                        "type": "object",
                                        "properties": {
                                            " message": {
                                                "type": "string"
                                            },
                                            "error": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                ]
                            }
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "do ping",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ping"
                ],
                "summary": "Ping",
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
        "/user": {
            "post": {
                "description": "Register new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Register new user",
                "parameters": [
                    {
                        "description": "new user",
                        "name": "user",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "allOf": [
                                    {
                                        "type": "string"
                                    },
                                    {
                                        "type": "object",
                                        "properties": {
                                            " message": {
                                                "type": "string"
                                            },
                                            "error": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                ]
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "allOf": [
                                    {
                                        "type": "string"
                                    },
                                    {
                                        "type": "object",
                                        "properties": {
                                            " message": {
                                                "type": "string"
                                            },
                                            "error": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                ]
                            }
                        }
                    }
                }
            },
            "patch": {
                "description": "Updates existing user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Updates user",
                "parameters": [
                    {
                        "description": "user to update",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "allOf": [
                                    {
                                        "type": "string"
                                    },
                                    {
                                        "type": "object",
                                        "properties": {
                                            " message": {
                                                "type": "string"
                                            },
                                            "error": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                ]
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "allOf": [
                                    {
                                        "type": "string"
                                    },
                                    {
                                        "type": "object",
                                        "properties": {
                                            " message": {
                                                "type": "string"
                                            },
                                            "error": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                ]
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.DogRequest": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "string"
                },
                "breed": {
                    "type": "string"
                },
                "coatColor": {
                    "type": "string"
                },
                "coatLength": {
                    "type": "string"
                },
                "host": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "imgUrl": {
                    "type": "string"
                },
                "imgs": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        }
                    }
                },
                "isLost": {
                    "type": "boolean"
                },
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                },
                "size": {
                    "type": "string"
                }
            }
        },
        "model.DogResponse": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "string"
                },
                "breed": {
                    "type": "string"
                },
                "coatColor": {
                    "type": "string"
                },
                "coatLength": {
                    "type": "string"
                },
                "host": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "img": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "imgUrl": {
                    "type": "string"
                },
                "isLost": {
                    "type": "boolean"
                },
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                },
                "size": {
                    "type": "string"
                }
            }
        },
        "model.User": {
            "type": "object",
            "properties": {
                "city": {
                    "description": "City\nin: string",
                    "type": "string"
                },
                "email": {
                    "description": "Email\nin: string",
                    "type": "string"
                },
                "firstName": {
                    "description": "First name\nin: string",
                    "type": "string"
                },
                "id": {
                    "description": "User ID\nin: string",
                    "type": "string"
                },
                "lastName": {
                    "description": "Last name\nin: string",
                    "type": "string"
                },
                "phone": {
                    "description": "Phone number\nin: string",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.1",
	Host:             "https://tpms-fdwva.ondigitalocean.app/tpms-be2",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "TPMS-BE Api",
	Description:      "tpms back-end Api docs",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}