package CommenDb

import (
	"errors"
	validation2 "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"regexp"
	"strconv"
	"strings"
)

func ConvertAnyToInt(param any) int {
	switch p := param.(type) {
	case int:
		return p
	case string:
		i, err := strconv.Atoi(p)
		if err != nil {
			return int(0)
		}
		return i
	case float64:
		return int(p)
	case int64:
		return int(p)
	default:
		return int(0)
	}
}
func ConvertAnyToFloat64(param any) float64 {
	switch p := param.(type) {
	case float64:
		return p
	case string:
		f, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return float64(0)
		}
		return f
	case int:
		return float64(p)
	case int64:
		return float64(p)
	default:
		return float64(0)
	}
}

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
	err := validation2.Validate(value, validation2.Length(ConvertAnyToInt(param), 10000))
	if err != nil {
		return err
	}
	return nil
}
func maxLen(value interface{}, param any) error {
	err := validation2.Validate(value, validation2.Length(0, ConvertAnyToInt(param)))
	if err != nil {
		return err
	}
	return nil
}
func minValue(value interface{}, param any) error {
	return validation2.Validate(value, validation2.Min(ConvertAnyToFloat64(param)))
}
func maxValue(value interface{}, param any) error {
	return validation2.Validate(value, validation2.Max(ConvertAnyToFloat64(param)))
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
func ValidationArrayOfObjInput(values []map[string]any, dbName string, validate string, param any,
	account string, lang string, title string) *ResponseErrors {
	funcs := setValidate()
	titleParam := make(map[string]string)
	titleParam["$information$"] = title
	meta := make(map[string]any)
	meta["dbName"] = dbName
	field := strings.Split(dbName, ".")
	if len(field) > 1 {
		for _, value := range values {
			err := funcs[validate](value[field[1]], param)
			if err != nil {
				NewErr := GetErrors(validate, account, lang, titleParam)
				NewErr.MetaData = meta
				return NewErr
			}
		}
	}
	return nil
}
func ValidationObjInput(values map[string]any, dbName string, validate string, param any,
	account string, lang string, title string) *ResponseErrors {
	funcs := setValidate()
	titleParam := make(map[string]string)
	titleParam["$information$"] = title
	meta := make(map[string]any)
	meta["dbName"] = dbName
	field := strings.Split(dbName, ".")
	if len(field) > 1 {
		err := funcs[validate](values[field[1]], param)
		if err != nil {
			NewErr := GetErrors(validate, account, lang, titleParam)
			NewErr.MetaData = meta
			return NewErr
		}

	}
	return nil
}
func ValidationArray(values []any, validate string, param any,
	account string, lang string, title string) *ResponseErrors {
	funcs := setValidate()
	titleParam := make(map[string]string)
	titleParam["$information$"] = title
	for _, value := range values {
		err := funcs[validate](value, param)
		if err != nil {
			NewErr := GetErrors(validate, account, lang, titleParam)
			return NewErr
		}
	}
	return nil
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
