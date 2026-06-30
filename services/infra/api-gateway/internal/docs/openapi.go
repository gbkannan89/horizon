package docs

import "github.com/gin-gonic/gin"

const spec = `{
  "openapi": "3.1.0",
  "info": {
    "title": "Horizon API",
    "description": "Goal-Centric Personal Financial Operating System — REST API",
    "version": "0.1.0",
    "contact": {"name": "Horizon Team"}
  },
  "servers": [{"url": "/api/v1", "description": "API v1"}],
  "components": {
    "securitySchemes": {
      "bearerAuth": {"type": "http", "scheme": "bearer", "bearerFormat": "JWT"}
    },
    "schemas": {
      "Error": {
        "type": "object",
        "properties": {
          "error": {
            "type": "object",
            "properties": {
              "type": {"type": "string"},
              "code": {"type": "string"},
              "message": {"type": "string"}
            }
          },
          "request_id": {"type": "string"}
        }
      },
      "CreateUserRequest": {
        "type": "object",
        "required": ["name","email"],
        "properties": {
          "name": {"type": "string", "example": "John Doe"},
          "email": {"type": "string", "format": "email", "example": "john@example.com"}
        }
      },
      "CreateGoalRequest": {
        "type": "object",
        "required": ["name","target_amount"],
        "properties": {
          "name": {"type": "string", "example": "Retirement"},
          "target_amount": {"type": "number", "example": 10000000}
        }
      },
      "CreateAccountRequest": {
        "type": "object",
        "required": ["name","type","currency"],
        "properties": {
          "name": {"type": "string", "example": "Savings Account"},
          "type": {"type": "string", "example": "Savings"},
          "currency": {"type": "string", "example": "INR"}
        }
      }
    }
  },
  "security": [{"bearerAuth": []}],
  "paths": {
    "/health": {"get": {"summary": "Health check","tags":["System"],"responses":{"200":{"description":"OK"}}}},
    "/auth/login": {"post": {"summary":"Login","tags":["Auth"],"requestBody":{"content":{"application/json":{"schema":{"type":"object","properties":{"email":{"type":"string"},"password":{"type":"string"}}}}}},"responses":{"200":{"description":"Login successful"},"401":{"description":"Invalid credentials"}}}},
    "/auth/me": {"get": {"summary":"Current user","tags":["Auth"],"security":[{"bearerAuth":[]}],"responses":{"200":{"description":"Current user profile"},"401":{"description":"Not authenticated"}}}},
    "/users": {
      "post": {"summary":"Create user","tags":["Users"],"requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/CreateUserRequest"}}}},"responses":{"201":{"description":"User created"},"400":{"description":"Validation error"}}},
      "get": {"summary":"List users","tags":["Users"],"responses":{"200":{"description":"User list"}}}
    },
    "/goals": {
      "post": {"summary":"Create goal","tags":["Goals"],"requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/CreateGoalRequest"}}}},"responses":{"201":{"description":"Goal created"}}},
      "get": {"summary":"List goals","tags":["Goals"],"responses":{"200":{"description":"Goal list"}}}
    },
    "/accounts": {
      "post": {"summary":"Create account","tags":["Accounts"],"requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/CreateAccountRequest"}}}},"responses":{"201":{"description":"Account created"}}},
      "get": {"summary":"List accounts","tags":["Accounts"],"responses":{"200":{"description":"Account list"}}}
    },
    "/events": {
      "post": {"summary":"Create event","tags":["Events"],"requestBody":{"content":{"application/json":{"schema":{"type":"object","properties":{"type":{"type":"string"},"amount":{"type":"number"},"currency":{"type":"string"}}}}}},"responses":{"201":{"description":"Event created"}}},
      "get": {"summary":"List events","tags":["Events"],"responses":{"200":{"description":"Event list"}}}
    },
    "/allocations": {
      "post": {"summary":"Create allocation","tags":["Allocations"]},
      "get": {"summary":"List allocations","tags":["Allocations"]}
    },
    "/institutions": {
      "post": {"summary":"Create institution","tags":["Institutions"]},
      "get": {"summary":"List institutions","tags":["Institutions"]}
    },
    "/assets": {
      "post": {"summary":"Create asset","tags":["Assets"]},
      "get": {"summary":"List assets","tags":["Assets"]}
    },
    "/liabilities": {
      "post": {"summary":"Create liability","tags":["Liabilities"]},
      "get": {"summary":"List liabilities","tags":["Liabilities"]}
    },
    "/portfolios": {
      "post": {"summary":"Create portfolio","tags":["Portfolios"]},
      "get": {"summary":"List portfolios","tags":["Portfolios"]}
    }
  }
}`

func RegisterOpenAPI(r *gin.RouterGroup) {
	r.GET("/openapi.json", func(c *gin.Context) {
		c.Data(200, "application/json", []byte(spec))
	})
}
