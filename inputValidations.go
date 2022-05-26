package CommenDb

import (
	"errors"
	validation2 "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"regexp"
)

func required(value interface{}, param any) error {
	return validation2.Validate(value, validation2.Required)
}
func alphaNumeric(value interface{}, param any) error {
	err := validation2.Validate(value, is.Alphanumeric)
	if err != nil {
		return err
	}
	return nil
}
func minLen(value interface{}, param any) error {
	err := validation2.Validate(value, validation2.Min(param))
	if err != nil {
		return err
	}
	return nil
}
func maxLen(value interface{}, param any) error {
	err := validation2.Validate(value, validation2.Max(param))
	if err != nil {
		return err
	}
	return nil
}
func minValue(value interface{}, param any) error {
	return validation2.Validate(value, validation2.Min(param))
}
func maxValue(value interface{}, param any) error {
	return validation2.Validate(value, validation2.Max(param))
}
func phone(value interface{}, param any) error {
	err := validation2.Validate(value, validation2.Match(regexp.MustCompile("^\\+[1-9]{1}[0-9]{3,14}$")))
	if err != nil {
		return err
	}
	return nil
}
func email(value interface{}, param any) error {
	err := validation2.Validate(value, is.Email)
	if err != nil {
		return err
	}
	return nil
}
func numericString(value interface{}, param any) error {
	err := validation2.Validate(value, is.Digit)
	if err != nil {
		return err
	}
	return nil
}
func integer(value interface{}, param any) error {
	switch value.(type) {
	case int:
		return nil
	case float64:
		return errors.New("value is note integer")
	case string:
		return errors.New("value is note integer")
	default:
		return errors.New("value is note integer")
	}
}
func float(value interface{}, param any) error {
	switch value.(type) {
	case int:
		return errors.New("value is note float")
	case float64:
		return nil
	case string:
		return errors.New("value is note float")
	default:
		return errors.New("value is note float")
	}
}
func setValidate() map[string]func(value interface{}, param any) error {
	m := make(map[string]func(value interface{}, param any) error)
	m["required"] = required
	m["alphaNum"] = alphaNumeric
	m["minLength"] = minLen
	m["maxLength"] = maxLen
	m["minValue"] = minValue
	m["maxValue"] = maxValue
	m["phone"] = phone
	m["email"] = email
	m["numeric"] = numericString
	m["integer"] = integer
	m["float"] = float
	return m
}

func ValidationInput(value interface{}, validate string, param any,
	account string, lang string, title string) *ResponseErrors {
	funcs := setValidate()
	titleParam := make(map[string]string)
	titleParam["$information$"] = title
	err := funcs[validate](value, param)
	if err != nil {
		return GetErrors(validate, account, lang, titleParam)
	}
	return nil
}
