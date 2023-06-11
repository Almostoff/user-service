package utils

import (
	"UsersService/internal/cConstants"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

// Use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(ctx context.Context, s interface{}) error {
	return validate.StructCtx(ctx, s)
}

func ValidateStructSize(s interface{}) error {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String {
			tag := typ.Field(i).Tag.Get("size")
			if tag == "" {
				tag = cConstants.MAxSizeStringStruct
			}
			if err := validateSize(field, tag); err != nil {
				return err
			}
			tag = typ.Field(i).Tag.Get("min")
			if tag == "" {
				tag = cConstants.MinSizeStringStruct
			}
			if err := validateMinSize(field, tag); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateSize(field reflect.Value, tagValue string) error {
	if field.String() == "" {
		return nil
	}
	if int64(len(field.String())) > StringToInt(tagValue) {
		return fmt.Errorf("%s exceeds maximum length of %s", field.String(), tagValue)
	}
	return nil
}

func validateMinSize(field reflect.Value, tagValue string) error {
	if field.String() == "" {
		return nil
	}
	if int64(len(field.String())) < StringToInt(tagValue) {
		return fmt.Errorf("%s exceeds minimum length of %s", field.String(), tagValue)
	}
	return nil
}
