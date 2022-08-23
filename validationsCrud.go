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
		Response := GetErrors("REF.INVALIDATION_ERROR", api.Account, api.Lang, nil)
		Err := make(map[string]any)
		Err["validation"] = err.Error()
		Response.MetaData = Err
		errors = append(errors, Response)
		return nil, errors
	}
	validation, err := getEndPointFileds(api.Route, api.Method, api.Service)
	if err != nil {
		Response := GetErrors("REF.CANNOT_INSERT", api.Account, api.Lang, nil)
		errors = append(errors, Response)
		return nil, errors
	}
	var userRole []string
	if validation.Meta.SecurityLevel != "1" {
		userAccount := getCurrentAccount(api.Account, "user")
		userRole, err = getUserRole(api.User, userAccount)
		if err != nil {
			Response := GetErrors("REF.CANNOT_ACCESS", api.Account, api.Lang, nil)
			errors = append(errors, Response)
			return nil, errors
		}
	}
	for _, field := range validation.Validator {
		isReqBody := false
		for _, feature := range field.Features {
			if feature == "10" {
				isReqBody = true
			}
		}
		if !isReqBody {
			continue
		}
		if _, ok := value[field.DbName]; ok {
			if len(field.DenyRoleKeys) != 0 && validation.Meta.SecurityLevel == "1" {
				Response := GetErrors("REF.SERVICE_UNKNOWN_ERROR", api.Account, api.Lang, nil)
				errors = append(errors, Response)
				return nil, errors
			}
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
								if field.MultiLang {
									for _, v := range value[field.DbName].(map[string]any) {
										validateErr := ValidationInput(v, validator.Rule, validator.Param, api.Account,
											api.Lang, field.Title[api.Lang])
										if validateErr != nil {
											errors = append(errors, validateErr)
										}
									}
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
							if field.MultiLang {
								for _, v := range value[field.DbName].(map[string]any) {
									validateErr := ValidationInput(v, validator.Rule, validator.Param, api.Account,
										api.Lang, field.Title[api.Lang])
									if validateErr != nil {
										errors = append(errors, validateErr)
									}
								}
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
			var parent []map[string]any
			arrayMap, err := json.Marshal(value[element.Parent])
			if err != nil {
				Response := GetErrors("REF.INVALIDATION_ERROR", api.Account, api.Lang, nil)
				Err := make(map[string]any)
				Err["validation"] = err.Error()
				Response.MetaData = Err
				errors = append(errors, Response)
				return nil, errors
			}
			err = json.Unmarshal(arrayMap, &parent)
			if err != nil {
				Response := GetErrors("REF.INVALIDATION_ERROR", api.Account, api.Lang, nil)
				Err := make(map[string]any)
				Err["validation"] = err.Error()
				Response.MetaData = Err
				errors = append(errors, Response)
				return nil, errors
			}
			for _, validator := range element.Validators {
				validateErr := ValidationArrayOfObjInput(parent, element.DbName, validator.Rule, validator.Param, api.Account,
					api.Lang, element.Title[api.Lang])
				if validateErr != nil {
					errors = append(errors, validateErr)
				}
			}
		} else if _, ok = ObjFields[element.Parent]; ok {
			parent := make(map[string]any)
			arrayMap, err := json.Marshal(value[element.Parent])
			if err != nil {
				Response := GetErrors("REF.INVALIDATION_ERROR", api.Account, api.Lang, nil)
				Err := make(map[string]any)
				Err["validation"] = err.Error()
				Response.MetaData = Err
				errors = append(errors, Response)
				return nil, errors
			}
			err = json.Unmarshal(arrayMap, &parent)
			if err != nil {
				Response := GetErrors("REF.INVALIDATION_ERROR", api.Account, api.Lang, nil)
				Err := make(map[string]any)
				Err["validation"] = err.Error()
				Response.MetaData = Err
				errors = append(errors, Response)
				return nil, errors
			}
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
		Response := GetErrors("REF.INVALIDATION_ERROR", api.Account, api.Lang, nil)
		Err := make(map[string]any)
		Err["validation"] = err.Error()
		Response.MetaData = Err
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
	for _, field := range validation.Validator {
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
	for _, field := range validation.Validator {
		if len(field.DenyRoleKeys) != 0 {
			canSee := CheckRoles(field.DenyRoleKeys, userRole)
			if canSee && field.Parent == "" {
				fields[field.DbName] = 1
			}
		} else {
			if field.Parent == "" {
				fields[field.DbName] = 1
			}
		}

	}
	return fields
}
