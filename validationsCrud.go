package CommenDb

import (
	"encoding/json"
	"reflect"
	"strings"
)

type Api struct {
	Route          string
	Method         string
	Service        string
	User           string
	Account        string
	Lang           string
	ClientToken    string
	UserToken      string
	TraceId        string
	SpanId         string
	ServiceVersion string
	Context        map[string]any
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

func (a *Api) WriteValidations(data []byte, api *Api) ([]byte, []*ResponseErrors) {
	value := make(map[string]any)
	var errors []*ResponseErrors
	arrayOfObjFields := make(map[string]fieldsEntities)
	arrayFields := make(map[string]fieldsEntities)
	ObjFields := make(map[string]fieldsEntities)
	children := make(map[string]fieldsEntities)
	err := json.Unmarshal(data, &value)
	if err != nil {
		Response := a.GetErrors("REF.INVALIDATION_ERROR", api.Account, api.Lang, nil)
		Err := make(map[string]any)
		Err["validation"] = err.Error()
		Response.MetaData = Err
		errors = append(errors, Response)
		return nil, errors
	}
	validation, err := a.getEndPointFileds(api.Route, api.Method, api.Service)
	if err != nil {
		Response := a.GetErrors("REF.CANNOT_INSERT", api.Account, api.Lang, nil)
		errors = append(errors, Response)
		return nil, errors
	}
	var userRole []string
	if validation.Meta.SecurityLevel != "1" {
		userAccount := a.GetCurrentAccount(api.Account, "user")
		userRole, err = a.getUserRole(api.User, userAccount)
		if err != nil {
			Response := a.GetErrors("REF.CANNOT_ACCESS", api.Account, api.Lang, nil)
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
				Response := a.GetErrors("REF.SERVICE_UNKNOWN_ERROR", api.Account, api.Lang, nil)
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
										validateErr := a.ValidationInput(v, validator.Rule, validator.Param, api.Account,
											api.Lang, field.Title[api.Lang])
										if validateErr != nil {
											errors = append(errors, validateErr)
										}
									}
								} else {
									validateErr := a.ValidationInput(value[field.DbName], validator.Rule, validator.Param, api.Account,
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
									validateErr := a.ValidationInput(v, validator.Rule, validator.Param, api.Account,
										api.Lang, field.Title[api.Lang])
									if validateErr != nil {
										errors = append(errors, validateErr)
									}
								}
							} else {
								validateErr := a.ValidationInput(value[field.DbName], validator.Rule, validator.Param, api.Account,
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
							validateErr := a.ValidationInput(nil, validator.Rule, validator.Param, api.Account,
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
				Response := a.GetErrors("REF.INVALIDATION_ERROR", api.Account, api.Lang, nil)
				Err := make(map[string]any)
				Err["validation"] = err.Error()
				Response.MetaData = Err
				errors = append(errors, Response)
				return nil, errors
			}
			err = json.Unmarshal(arrayMap, &parent)
			if err != nil {
				Response := a.GetErrors("REF.INVALIDATION_ERROR", api.Account, api.Lang, nil)
				Err := make(map[string]any)
				Err["validation"] = err.Error()
				Response.MetaData = Err
				errors = append(errors, Response)
				return nil, errors
			}
			for _, validator := range element.Validators {
				validateErr := a.ValidationArrayOfObjInput(parent, element.DbName, validator.Rule, validator.Param, api.Account,
					api.Lang, element.Title[api.Lang])
				if validateErr != nil {
					errors = append(errors, validateErr)
				}
			}
		} else if _, ok = ObjFields[element.Parent]; ok {
			parent := make(map[string]any)
			arrayMap, err := json.Marshal(value[element.Parent])
			if err != nil {
				Response := a.GetErrors("REF.INVALIDATION_ERROR", api.Account, api.Lang, nil)
				Err := make(map[string]any)
				Err["validation"] = err.Error()
				Response.MetaData = Err
				errors = append(errors, Response)
				return nil, errors
			}
			err = json.Unmarshal(arrayMap, &parent)
			if err != nil {
				Response := a.GetErrors("REF.INVALIDATION_ERROR", api.Account, api.Lang, nil)
				Err := make(map[string]any)
				Err["validation"] = err.Error()
				Response.MetaData = Err
				errors = append(errors, Response)
				return nil, errors
			}
			for _, validator := range element.Validators {
				validateErr := a.ValidationObjInput(value[element.Parent].(map[string]any), element.DbName, validator.Rule, validator.Param, api.Account,
					api.Lang, element.Title[api.Lang])
				if validateErr != nil {
					errors = append(errors, validateErr)
				}
			}
		}
	}
	for _, element := range arrayFields {
		for _, validator := range element.Validators {
			validateErr := a.ValidationArray(value[element.DbName].([]any), validator.Rule, validator.Param, api.Account,
				api.Lang, element.Title[api.Lang])
			if validateErr != nil {
				errors = append(errors, validateErr)
			}
		}
	}
	body, err := json.Marshal(value)
	if err != nil {
		Response := a.GetErrors("REF.INVALIDATION_ERROR", api.Account, api.Lang, nil)
		Err := make(map[string]any)
		Err["validation"] = err.Error()
		Response.MetaData = Err
		errors = append(errors, Response)
		return nil, errors
	}
	return body, errors
}
func (a *Api) GetOneValidations(value map[string]any, api *Api) map[string]any {
	validation, err := a.getEndPointFileds(api.Route, api.Method, api.Service)
	if err != nil {
		return nil
	}
	userAccount := a.GetCurrentAccount(api.Account, "user")
	userRole, err := a.getUserRole(api.User, userAccount)
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
func (a *Api) GetArrayValidations(api *Api) map[string]int8 {
	fields := make(map[string]int8)
	validation, err := a.getEndPointFileds(api.Route, api.Method, api.Service)
	if err != nil {
		return nil
	}
	userAccount := a.GetCurrentAccount(api.Account, "user")
	userRole, err := a.getUserRole(api.User, userAccount)
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
func IsSlice(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}
func (a *Api) FindDomainValues(data any) any {
	if IsSlice(data) {
		return a.FindDomainValuesInArray(data)
	} else {
		return a.FindDomainValuesInMap(data)
	}
}
func FindValueInMap(mapInfo map[string]any, key string) any {
	splitKey := strings.Split(key, ".")
	if len(splitKey) > 1 {
		if value, found := mapInfo[splitKey[0]]; found {
			splitKey = splitKey[1:]
			if IsSlice(value) {
				var slice []any
				for _, a := range value.([]any) {
					values := FindValueInMap(a.(map[string]any), strings.Join(splitKey, "."))
					if IsSlice(values) {
						slice = append(slice, values.([]any)...)
					} else {
						slice = append(slice, values.(string))
					}
				}
				return slice
			}
			return FindValueInMap(value.(map[string]any), strings.Join(splitKey, "."))
		}
	} else {
		if value, found := mapInfo[key]; found {
			return value
		}
	}
	return nil
}
func SetDataToFieldOfMap(mapInfos map[string]any, key, MainField string, values []map[string]any) map[string]any {
	splitKey := strings.Split(key, ".")
	if len(splitKey) > 1 {
		if _, found := mapInfos[splitKey[0]]; found {
			parentKey := splitKey[0]
			splitKey = splitKey[1:]
			if IsSlice(mapInfos[parentKey]) {
				for i, _ := range mapInfos[parentKey].([]any) {
					mapInfos[parentKey].([]any)[i] = SetDataToFieldOfMap(mapInfos[parentKey].([]any)[i].(map[string]any), strings.Join(splitKey, "."), MainField, values)
				}
			} else {
				mapInfos[parentKey] = SetDataToFieldOfMap(mapInfos[parentKey].(map[string]any), strings.Join(splitKey, "."), MainField, values)
			}
		}
	} else {
		if _, found := mapInfos[key]; found {
			if IsSlice(mapInfos[key]) {
				var newInfo []map[string]any
				for _, s := range mapInfos[key].([]any) {
					for _, value := range values {
						if s == value[MainField] {
							newInfo = append(newInfo, value)
						}
					}
				}
				mapInfos[key] = newInfo
			} else {
				for _, value := range values {
					if mapInfos[key] == value[MainField] {
						mapInfos[key] = value
					}
				}
			}

		}
	}
	return mapInfos
}

type DomainVal struct {
	Ref   string
	Value any
}

func (a *Api) FindDomainValuesInArray(data any) any {
	var ArrData []map[string]any
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(jsonData, &ArrData)
	if err != nil {
		return nil
	}
	endpointInfo, _ := a.getEndPointFileds(a.Route, a.Method, a.Service)
	if endpointInfo == nil {
		return data
	}
	for _, entity := range endpointInfo.Validator {
		if entity.DisplayType == "combo" {
			if domainType, found := entity.Conf["domain_type"]; found {
				if domainType == "endpoint" {
					var ArrayIds []any
					service := FindValueInMap(entity.Conf, "service")
					route := FindValueInMap(entity.Conf, "route")
					fields := FindValueInMap(entity.Conf, "fields")
					MainField := FindValueInMap(entity.Conf, "main_field")
					for _, data := range ArrData {
						id := FindValueInMap(data, entity.DbName)
						if IsSlice(id) {
							ArrayIds = append(ArrayIds, id.([]any)...)
						} else {
							ArrayIds = append(ArrayIds, id.(string))
						}
					}
					values := a.getDomainValuesDataByRefId(service.(string), route.(string),
						ArrayIds, MainField.(string), fields.([]any))
					for i, _ := range ArrData {
						ArrData[i] = SetDataToFieldOfMap(ArrData[i], entity.DbName, MainField.(string), values)
					}
				} else if domainType == "domain_key" {
					var ArrayKeys []any
					for _, data := range ArrData {
						key := FindValueInMap(data, entity.DbName)
						if key != nil {
							if IsSlice(key) {
								ArrayKeys = append(ArrayKeys, key.([]any)...)
							} else {
								ArrayKeys = append(ArrayKeys, key.(string))
							}
						}
					}
					values := a.getDomainValuesDataByRefKey(entity.Conf["domain_key"].(string))
					for i, _ := range ArrData {
						ArrData[i] = SetDataToFieldOfMap(ArrData[i], entity.DbName, "code", values)
					}
				}
			}
		}
	}
	return ArrData
}
func (a *Api) FindDomainValuesInMap(data any) any {
	var MapData map[string]any
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(jsonData, &MapData)
	if err != nil {
		return nil
	}
	endpointInfo, _ := a.getEndPointFileds(a.Route, a.Method, a.Service)
	if endpointInfo == nil {
		return data
	}
	for _, entity := range endpointInfo.Validator {
		if entity.DisplayType == "combo" {
			if domainType, found := entity.Conf["domain_type"]; found {
				if domainType == "endpoint" {
					var ArrayIds []any
					service := FindValueInMap(entity.Conf, "service")
					route := FindValueInMap(entity.Conf, "route")
					fields := FindValueInMap(entity.Conf, "fields")
					MainField := FindValueInMap(entity.Conf, "main_field")
					id := FindValueInMap(MapData, entity.DbName)
					if IsSlice(id) {
						ArrayIds = append(ArrayIds, id.([]any)...)
					} else {
						ArrayIds = append(ArrayIds, id.(string))
					}
					values := a.getDomainValuesDataByRefId(service.(string), route.(string),
						ArrayIds, MainField.(string), fields.([]any))
					MapData = SetDataToFieldOfMap(MapData, entity.DbName, MainField.(string), values)
				} else if domainType == "domain_key" {
					var ArrayKeys []string
					key := FindValueInMap(MapData, entity.DbName)
					if IsSlice(key) {
						ArrayKeys = append(ArrayKeys, key.([]string)...)
					} else {
						ArrayKeys = append(ArrayKeys, key.(string))
					}
					values := a.getDomainValuesDataByRefKey(entity.Conf["domain_key"].(string))
					MapData = SetDataToFieldOfMap(MapData, entity.DbName, "key", values)
				}

			}
		}
	}
	return MapData
}
