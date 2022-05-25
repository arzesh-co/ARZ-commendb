package CommenDb

import (
	validation2 "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"regexp"
)

func alphaNumeric(value interface{}) error {
	err := validation2.Validate(value, is.Alphanumeric)
	if err != nil {
		return err
	}
	return nil
}
func minLen(value interface{}, len int64) error {
	err := validation2.Validate(value, validation2.Min(len))
	if err != nil {
		return err
	}
	return nil
}
func maxLen(value interface{}, len int64) error {
	err := validation2.Validate(value, validation2.Max(len))
	if err != nil {
		return err
	}
	return nil
}
func minValue(value interface{}, min int64) error {
	return validation2.Validate(value, validation2.Min(min))
}
func maxValue(value interface{}, max int64) error {
	return validation2.Validate(value, validation2.Max(max))
}
func phone(value interface{}) error {
	err := validation2.Validate(value, validation2.Match(regexp.MustCompile("^\\+[1-9]{1}[0-9]{3,14}$")))
	if err != nil {
		return err
	}
	return nil
}
func email(value interface{}) error {
	err := validation2.Validate(value, is.Email)
	if err != nil {
		return err
	}
	return nil
}
func numericString(value interface{}) error {
	err := validation2.Validate(value, is.Digit)
	if err != nil {
		return err
	}
	return nil
}

func ValidationInput(value interface{}, validate string, param interface{}) *ResponseErrors {
	switch validate {
	case "alphaNumeric":
		err := alphaNumeric(value)
		if err != nil {

		}
	}
	return nil
}
