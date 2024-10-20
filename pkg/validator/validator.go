package corevalidator

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

var errorsMap = map[string]error{}

var nullableMap = map[string]bool{}

type Validator[T any] struct {
	tagName  string
	validate *validator.Validate
}

func New[T any](validate *validator.Validate) *Validator[T] {
	return &Validator[T]{
		tagName:  "validate",
		validate: validate,
	}
}

func (v *Validator[T]) SetTagName(tagName string) {
	v.validate.SetTagName(tagName)
}

func (v *Validator[T]) RegisterCustomTag(tag string, validate validator.Func, nullable bool, errorMsg string) error {
	if err := v.validate.RegisterValidation(tag, validate, !nullable); err != nil {
		return fmt.Errorf("failed to add validator: %w", err)
	}
	nullableMap[tag] = nullable
	errorsMap[tag] = errors.New(errorMsg)
	return nil
}

func (v *Validator[T]) Validate(data T) (err error) {
	value := reflect.ValueOf(data)
	for idx := range value.NumField() {
		tagValue, exists := value.Type().Field(idx).Tag.Lookup(v.tagName)
		if !exists {
			continue
		}

		for idx := range strings.Split(tagValue, ",") {
			if validateErr := v.validate.Var(
				value.Field(idx).Interface(),
				tagValue,
			); validateErr != nil && !(nullableMap[tagValue] && value.Field(idx).IsNil()) {
				err = errors.Join(err, errorsMap[tagValue])
			}
		}

		if err != nil {
			return err
		}
	}

	return err
}
