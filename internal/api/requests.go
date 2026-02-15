package api

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"

    "github.com/go-playground/validator/v10"

    "segment-service/internal/entities"
)

var validate *validator.Validate

func init() {
    validate = validator.New()
}

type segmentRequestInput struct {
    Name       string `json:"name" validate:"required"`
    TTLSeconds *int   `json:"ttl_seconds" validate:"omitempty,gt=0"`
}

// Formatting and helper methods

func validateSegmentInput(r *http.Request, dst any) error {
    if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
        return entities.ErrInvalidJsonBody
    }
    if err := validate.Struct(dst); err != nil {
        return fmt.Errorf("%w: %w", entities.ErrInvalidSegment, formatValidationError(err))
    }
    return nil
}

func (c *segmentRequestInput) UnmarshalJSON(data []byte) error {
    type Alias segmentRequestInput
    aux := &struct {
        TTLSeconds interface{} `json:"ttl_seconds"`
        *Alias
    }{
        Alias: (*Alias)(c),
    }
    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }
    if aux.TTLSeconds != nil {
        switch v := aux.TTLSeconds.(type) {
        case string:
            i, err := strconv.Atoi(v)
            if err != nil {
                return err
            }
            c.TTLSeconds = &i
        case float64:
            i := int(v)
            c.TTLSeconds = &i
        }
    }
    return nil
}

func formatValidationError(err error) error {
    if validationErrs, ok := err.(validator.ValidationErrors); ok {
        if len(validationErrs) > 0 {
            e := validationErrs[0]
            field := e.Field()
            switch e.Tag() {
            case "required":
                return fmt.Errorf("%s is required", field)
            case "max":
                return fmt.Errorf("%s must not exceed %s characters", field, e.Param())
            case "gt":
                return fmt.Errorf("%s must be a greater than %s", field, e.Param())
            default:
                return fmt.Errorf("%s is invalid", field)
            }
        }
    }
    return err
}
