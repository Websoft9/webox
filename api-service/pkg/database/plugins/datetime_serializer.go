package plugins

import (
	"context"
	"database/sql/driver"
	"fmt"
	"reflect"
	"time"

	"api-service/pkg/logger"

	"gorm.io/gorm/schema"
)

// DateTimeSerializer is a custom time serializer for GORM
// It formats time values to "2006-01-02 15:04:05" format when writing to database
type DateTimeSerializer struct{}

// Scan deserializes data from database
// When reading from database, it returns the original value without any processing
func (DateTimeSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue any) error {
	// No processing needed for deserialization, return original value
	return nil
}

// Value serializes data to database
// When writing to database, it formats time to "2006-01-02 15:04:05" format
func (DateTimeSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue any) (any, error) {
	logger.Debug("DateTimeSerializer.Value called",
		logger.String("field_name", field.Name),
		logger.String("field_value_type", fmt.Sprintf("%T", fieldValue)))

	switch v := fieldValue.(type) {
	case time.Time:
		// Handle time.Time type
		if v.IsZero() {
			logger.Debug("Time value is zero, returning nil",
				logger.String("field_name", field.Name))
			return nil, nil
		}
		formatted := v.Format(time.DateTime)
		logger.Debug("Formatted time.Time value",
			logger.String("field_name", field.Name),
			logger.String("formatted_value", formatted))
		return formatted, nil

	case *time.Time:
		// Handle *time.Time type
		if v == nil || v.IsZero() {
			logger.Debug("Time pointer is nil or zero, returning nil",
				logger.String("field_name", field.Name))
			return nil, nil
		}
		formatted := v.Format(time.DateTime)
		logger.Debug("Formatted *time.Time value",
			logger.String("field_name", field.Name),
			logger.String("original_value", v.String()),
			logger.String("formatted_value", formatted))
		return formatted, nil

	default:
		// Unsupported field type
		err := fmt.Errorf("unsupported field type: %T", fieldValue)
		logger.Debug("Unsupported field type in DateTimeSerializer",
			logger.String("field_name", field.Name),
			logger.String("field_type", fmt.Sprintf("%T", fieldValue)),
			logger.ErrorField(err))
		return nil, err
	}
}

// SerializerValuerInterface implements driver.Valuer interface
// This interface is used for custom value serialization in database operations
type SerializerValuerInterface interface {
	driver.Valuer
}
