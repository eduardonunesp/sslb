package cfg

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// Generated at http://jsonschema.net/#/
const ConfigSchema = `
{
  "id": "/",
  "type": "object",
  "properties": {
    "general": {
      "id": "general",
      "type": "object",
      "properties": {
        "maxProcs": {
          "id": "maxProcs",
          "type": "integer"
        },
        "workerPoolSize": {
          "id": "workerPoolSize",
          "type": "integer"
        },
        "gracefulShutdown": {
          "id": "gracefulShutdown",
          "type": "boolean"
        },
        "logLevel": {
          "id": "logLevel",
          "type": "string"
        },
        "websocket": {
          "id": "websocket",
          "type": "boolean"
        },
        "host": {
          "id": "rpchost",
          "type": "string"
        },
        "port": {
          "id": "rpcport",
          "type": "integer"
        }
      }
    },
    "frontends": {
      "id": "frontends",
      "type": "array",
      "items": {
        "id": "0",
        "type": "object",
        "properties": {
          "name": {
            "id": "name",
            "type": "string"
          },
          "host": {
            "id": "host",
            "type": "string"
          },
          "port": {
            "id": "port",
            "type": "integer"
          },
          "route": {
            "id": "route",
            "type": "string"
          },
          "timeout": {
            "id": "timeout",
            "type": "integer"
          },
          "backends": {
            "id": "backends",
            "type": "array",
            "items": {
              "id": "0",
              "type": "object",
              "properties": {
                "name": {
                  "id": "name",
                  "type": "string"
                },
                "address": {
                  "id": "address",
                  "type": "string"
                },
                "heartbeat": {
                  "id": "heartbeat",
                  "type": "string"
                },
                "hbmethod": {
                  "id": "hbmethod",
                  "type": "string"
                },
                "weigth": {
                  "id": "weigth",
                  "type": "integer"
                },
                "inactiveAfter": {
                  "id": "inactiveAfter",
                  "type": "integer"
                },
                "activeAfter": {
                  "id": "activeAfter",
                  "type": "integer"
                },
                "heartbeatTime": {
                  "id": "heartbeatTime",
                  "type": "integer"
                },
                "retryTime": {
                  "id": "retryTime",
                  "type": "integer"
                }
              },
              "required": [
                "name",
                "address",
                "heartbeat"
              ]
            }
          }
        },
        "required": [
          "name",
          "host",
          "port",
          "route"
        ]
      }
    }
  },
  "required": [
    "general",
    "frontends"
  ]
}
`

func Validate(file []byte) error {
	schemaLoader := gojsonschema.NewStringLoader(ConfigSchema)
	schema, err := gojsonschema.NewSchema(schemaLoader)

	documentLoader := gojsonschema.NewStringLoader(string(file))

	result, err := schema.Validate(documentLoader)
	if err != nil {
		log.Println("Failed to validate", err.Error())
		return err
	}

	if !result.Valid() {
		errs := []string{}
		for _, desc := range result.Errors() {
			e := fmt.Sprintf("%s", desc)
			errs = append(errs, e)
		}
		res := strings.Join(errs, ",")
		return errors.New(res)
	}

	return nil
}
