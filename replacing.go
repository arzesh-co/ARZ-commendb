package CommenDb

import (
	"fmt"
	"regexp"
	"strings"
)

func ConvertMap(values map[string]any, params map[string]string) map[string]any {
	NewMap := make(map[string]any)
	for key, element := range params {
		switch ConvertorType(element) {
		case "func":
			NewMap[key] = findFunc(element)
		case "string":
			NewMap[key] = values[element]
		case "map":
			NewMap[key] = FindValueOfMap(values, element)
		case "replace":
			NewMap[key] = Replacer(values, element)
		}
	}
	return NewMap
}
func Replacer(values map[string]any, param string) string {
	re := regexp.MustCompile("{{(.*?)}}")
	submatchall := re.FindAllString(param, -1)
	NewString := ""
	for _, element := range submatchall {
		repVal := element
		element = strings.Trim(element, "}")
		element = strings.Trim(element, "{")
		switch ConvertorType(element) {
		case "func":
			value := findFunc(element)
			NewString = strings.Replace(param, repVal, fmt.Sprintf("%v", value), -1)
		case "string":
			NewString = strings.Replace(param, repVal, fmt.Sprintf("%v", values[element]), -1)
		case "map":
			value := FindValueOfMap(values, element)
			NewString = strings.Replace(param, repVal, fmt.Sprintf("%v", value), -1)
		}
	}
	return NewString
}
func findFunc(param string) any {
	funcs := SetCommenFunc()
	re := regexp.MustCompile(`\$(.*?)\$`)
	Findfunc := re.FindStringSubmatch(param)
	if len(Findfunc) > 0 {
		re = regexp.MustCompile(`(.*?)\(`)
		match := re.FindStringSubmatch(Findfunc[1])
		funcName := match[1]
		var params string
		re = regexp.MustCompile(`\((.*?)\)`)
		// Text between parentheses:
		submatchall := re.FindAllString(Findfunc[1], -1)
		for _, element := range submatchall {
			element = strings.Trim(element, "(")
			element = strings.Trim(element, ")")
			params = element
		}
		if params != "" {
			return funcs[funcName](params)
		} else {
			return funcs[funcName](0)
		}
	}
	return nil
}
func FindValueOfMap(value map[string]any, key string) any {
	partOfMap := strings.Split(key, ".")
	if len(partOfMap) > 0 {
		part := value[partOfMap[0]]
		if len(partOfMap) > 1 {
			for _, s := range partOfMap {
				part = part.(map[string]any)[s]
			}
		}
		return part
	}
	return nil
}
func ConvertorType(param string) string {
	re := regexp.MustCompile(`\$(.*?)\$`)
	isFunc := re.MatchString(param)
	if isFunc {
		return "func"
	}
	partOfMap := strings.Split(param, ".")
	if len(partOfMap) > 1 {
		return "map"
	}
	re = regexp.MustCompile("{{(.*?)}}")
	isReplacer := re.MatchString(param)
	if isReplacer {
		return "replace"
	}
	return "string"
}
