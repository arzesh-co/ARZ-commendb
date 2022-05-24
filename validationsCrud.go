package CommenDb

type Api struct {
	Route   string
	Method  string
	Service string
	User    string
	Account string
}

func CheckRoles(endPointRoles, userRoles []string) bool {
	for _, role := range userRoles {
		for _, pointRole := range endPointRoles {
			if pointRole == role {
				return true
			}
		}
	}
	return false
}

func WriteValidations(value map[string]any, api *Api) map[string]any {
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
			} else {
				fields[field.DbName] = 0
			}
		}

	}
	return fields
}
