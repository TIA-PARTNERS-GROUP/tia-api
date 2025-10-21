package docs

import (
	"gorm.io/datatypes"
)

// @Description The structure contains a raw JSON payload, typically validated against a specific schema.
type SwaggerIgnoreGormJSON struct {
	Config datatypes.JSON `json:"config" swaggerignore:"true"`
}
