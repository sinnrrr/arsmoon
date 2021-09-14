package main

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func parseJsonAndValidate(data []byte, s interface{}) error {
	if err := json.Unmarshal(data, s); err != nil {
		return err
	}
	if err := validate.Struct(s); err != nil {
		return err
	}

	return nil
}