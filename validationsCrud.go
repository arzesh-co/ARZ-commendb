package CommenDb

import "encoding/json"

type Api struct {
	Route   string
	Method  string
	Service string
	User    string
	Account string
	lang    string
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

func WriteValidations(data []byte, api *Api) (map[string]any, []*ResponseErrors) {
	value := make(map[string]any)
	var errors []*ResponseErrors
	err := json.Unmarshal(data, &value)
	if err != nil {
		Response := GetErrors("ARZ-input", api.Account, api.lang, nil)
		errors = append(errors, Response)
		return nil, errors
	}
	validation, err := getEndPointFileds(api.Route, api.Method, api.Service)
	if err != nil {
		Response := GetErrors("ARZ-write_in", api.Account, api.lang, nil)
		errors = append(errors, Response)
		return nil, errors
	}
	userRole, err := getUserRole(api.User, api.Account)
	if err != nil {
		Response := GetErrors("ARZ-access", api.Account, api.lang, nil)
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
							validateErr := ValidationInput(value[field.DbName], validator.Rule, validator.Param, api.Account,
								api.lang, field.Title[api.lang])
							if validateErr != nil {
								errors = append(errors, validateErr)
							}
						}
					}
				}
			} else {
				if len(field.Validators) != 0 {
					for _, validator := range field.Validators {
						validateErr := ValidationInput(value[field.DbName], validator.Rule, validator.Param, api.Account,
							api.lang, field.Title[api.lang])
						if validateErr != nil {
							errors = append(errors, validateErr)
						}
					}
				}
			}
		}
	}
	return value, errors
}
func GetOneValidations(value map[string]any, api *Api) map[string]any {
	validation, err := getEndPointFileds(api.Route, api.Method, api.Service)
	if err != nil {
		return nil
	}
	userRole, err := getUserRole(api.User, api.Account)
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
	userRole, err := getUserRole(api.User, api.Account)
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
