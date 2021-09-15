package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func parseJsonAndValidate(data []byte, s interface{}) error {
	if err := json.Unmarshal(data, s); err != nil {
		return fmt.Errorf("failed unmarshaling JSON: %s", err)
	}
	if err := validate.Struct(s); err != nil {
		return fmt.Errorf("failed validating unmarshaled JSON: %s", err)
	}

	return nil
}
