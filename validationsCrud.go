package CommenDb

import (
	"encoding/json"
)

type Api struct {
	Route   string
	Method  string
	Service string
	User    string
	Account string
	Lang    string
}

func CheckRoles(endPointRoles, userRoles []string) bool {
	for _, role := range userRoles {
		for _, pointRole := range endPointRoles {
			if pointRole == role {
				return false
			}
		}
	}
	return true
}

func WriteValidations(data []byte, api *Api) ([]byte, []*ResponseErrors) {
	value := make(map[string]any)
	var errors []*ResponseErrors
	arrayOfObjFields := make(map[string]fieldsEntities)
	arrayFields := make(map[string]fieldsEntities)
	ObjFields := make(map[string]fieldsEntities)
	children := make(map[string]fieldsEntities)
	err := json.Unmarshal(data, &value)
	if err != nil {
		Response := GetErrors("ARZ-input", api.Account, api.Lang, nil)
		errors = append(errors, Response)
		return nil, errors
	}
	validation, err := getEndPointFileds(api.Route, api.Method, api.Service)
	if err != nil {
		Response := GetErrors("ARZ-write_in", api.Account, api.Lang, nil)
		errors = append(errors, Response)
		return nil, errors
	}
	userAccount := getCurrentAccount(api.Account, "user")
	userRole, err := getUserRole(api.User, userAccount)
	if err != nil {
		Response := GetErrors("ARZ-access", api.Account, api.Lang, nil)
		errors = append(errors, Response)
		return nil, errors
	}
	for _, field := range validation {
		if _, ok := value[field.DbName]; ok {
			if len(field.DenyRoleKeys) != 0 {
				canSet := CheckRoles(field.DenyRoleKeys, userRole)
				if !canSet {
					delete(value, field.DbName)
				} else {
					if len(field.Validators) != 0 {
						for _, validator := range field.Validators {
							if field.DataType == "array_of_object" {
								arrayOfObjFields[field.DbName] = field
							} else if field.DataType == "array" {
								arrayFields[field.DbName] = field
							} else if field.DataType == "object" {
								ObjFields[field.DbName] = field
							} else if field.Parent != "" {
								children[field.DbName] = field
							} else {
								validateErr := ValidationInput(value[field.DbName], validator.Rule, validator.Param, api.Account,
									api.Lang, field.Title[api.Lang])
								if validateErr != nil {
									errors = append(errors, validateErr)
								}
							}
						}
					}
				}
			} else {
				if len(field.Validators) != 0 {
					for _, validator := range field.Validators {
						if field.DataType == "array_of_object" {
							arrayOfObjFields[field.DbName] = field
						} else if field.DataType == "array" {
							arrayFields[field.DbName] = field
						} else if field.DataType == "object" {
							ObjFields[field.DbName] = field
						} else if field.Parent != "" {
							children[field.DbName] = field
						} else {
							validateErr := ValidationInput(value[field.DbName], validator.Rule, validator.Param, api.Account,
								api.Lang, field.Title[api.Lang])
							if validateErr != nil {
								errors = append(errors, validateErr)
							}
						}
					}
				}
			}
		} else {
			if len(field.Validators) != 0 {
				if field.Parent != "" {
					children[field.DbName] = field
				} else {
					canSet := CheckRoles(field.DenyRoleKeys, userRole)
					for _, validator := range field.Validators {
						if validator.Rule == "required" && canSet {
							validateErr := ValidationInput(nil, validator.Rule, validator.Param, api.Account,
								api.Lang, field.Title[api.Lang])
							if validateErr != nil {
								errors = append(errors, validateErr)
							}
						}
					}
				}
			}
		}
	}
	for _, element := range children {
		if _, ok := arrayOfObjFields[element.Parent]; ok {
			if _, ok = value[element.Parent]; !ok {
				continue
			}
			for _, validator := range element.Validators {
				validateErr := ValidationArrayOfObjInput(value[element.Parent].([]map[string]any), element.DbName, validator.Rule, validator.Param, api.Account,
					api.Lang, element.Title[api.Lang])
				if validateErr != nil {
					errors = append(errors, validateErr)
				}
			}
		} else if _, ok = ObjFields[element.Parent]; ok {
			for _, validator := range element.Validators {
				validateErr := ValidationObjInput(value[element.Parent].(map[string]any), element.DbName, validator.Rule, validator.Param, api.Account,
					api.Lang, element.Title[api.Lang])
				if validateErr != nil {
					errors = append(errors, validateErr)
				}
			}
		}
	}
	for _, element := range arrayFields {
		for _, validator := range element.Validators {
			validateErr := ValidationArray(value[element.DbName].([]any), validator.Rule, validator.Param, api.Account,
				api.Lang, element.Title[api.Lang])
			if validateErr != nil {
				errors = append(errors, validateErr)
			}
		}
	}
	body, err := json.Marshal(value)
	if err != nil {
		Response := GetErrors("ARZ-input", api.Account, api.Lang, nil)
		errors = append(errors, Response)
		return nil, errors
	}
	return body, errors
}
func GetOneValidations(value map[string]any, api *Api) map[string]any {
	validation, err := getEndPointFileds(api.Route, api.Method, api.Service)
	if err != nil {
		return nil
	}
	userAccount := getCurrentAccount(api.Account, "user")
	userRole, err := getUserRole(api.User, userAccount)
	if err != nil {
		return nil
	}
	for _, field := range validation {
		if _, ok := value[field.DbName]; ok {
			if len(field.DenyRoleKeys) != 0 {
				canSet := CheckRoles(field.DenyRoleKeys, userRole)
				if !canSet {
					delete(value, field.DbName)
				}
			}
		}
	}
	return value
}
func GetArrayValidations(api *Api) map[string]int8 {
	fields := make(map[string]int8)
	validation, err := getEndPointFileds(api.Route, api.Method, api.Service)
	if err != nil {
		return nil
	}
	userAccount := getCurrentAccount(api.Account, "user")
	userRole, err := getUserRole(api.User, userAccount)
	if err != nil {
		return nil
	}
	for _, field := range validation {
		if len(field.DenyRoleKeys) != 0 {
			canSee := CheckRoles(field.DenyRoleKeys, userRole)
			if canSee {
				fields[field.DbName] = 1
			}
		} else {
			fields[field.DbName] = 1
		}

	}
	return fields
}
