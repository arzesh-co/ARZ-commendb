package CommenDb

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"
	"os"
)

type ResponseErrors struct {
	ErrorType struct {
		Title string `json:"title"`
		Desc  string `json:"desc"`
		Url   string `json:"url"`
	} `json:"error_type"`
	StatusKey string         `json:"status_key"`
	Detail    string         `json:"detail"`
	Title     string         `json:"title"`
	HelpUrl   string         `json:"help_url"`
	MetaData  map[string]any `json:"meta_data"`
}

func getError(key string, account string, lang string, params map[string]string) *ResponseErrors {
	req, err := http.NewRequest("GET", os.Getenv("coreApi")+"/api/core/errors/key/"+key, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("account_uuid", account)
	q := req.URL.Query()
	q.Add("lang", lang)
	paramsS, _ := json.Marshal(params)
	q.Add("params", string(paramsS))
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	Info := &ResponseErrors{}
	err = json.Unmarshal(body, Info)
	if err != nil {
		return nil
	}
	return Info
}

type roles struct {
	Data []string `json:"data"`
}
type Validator struct {
	Rule  string `json:"rule" bson:"rule"`
	Param int64  `json:"param" bson:"param"`
}
type fieldsEntities struct {
	DbName       string                 `json:"db_name" bson:"db_name"`
	Title        map[string]string      `json:"title" bson:"title"`
	Parent       string                 `json:"parent" bson:"parent"`
	MultiLang    bool                   `json:"multi_lang" bson:"multi_lang"`
	Validators   []Validator            `json:"validators" bson:"validators"`
	Required     bool                   `json:"required" bson:"required"`
	DataType     string                 `json:"data_type" bson:"data_type"`
	DisplayType  string                 `json:"display_type" bson:"display_type"`
	Conf         map[string]interface{} `json:"conf" bson:"conf"`
	Features     []string               `json:"features" bson:"features"`
	DenyRoleKeys []string               `json:"deny_role_keys" bson:"deny_role_keys"`
}
type validation struct {
	Meta struct {
		SecurityLevel string `json:"security_level"`
	} `json:"meta"`
	Validator []fieldsEntities `json:"data"`
}

func getUserRole(user string, account string) ([]string, error) {
	req, err := http.NewRequest("GET", os.Getenv("userApi")+"/api/user/users/uuid/"+user+"/roles", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("account_uuid", account)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	userInfo := &roles{}
	err = json.Unmarshal(body, userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo.Data, nil
}
func getEndPointFileds(route string, method string, service string) (*validation, error) {
	req, err := http.NewRequest("GET", os.Getenv("coreApi")+"/api/core/endpoint/validators", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("service", service)
	q.Add("route", route)
	q.Add("method", method)
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	userInfo := &validation{}
	err = json.Unmarshal(body, userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

type accountService struct {
	Account string `json:"data"`
}

func getCurrentAccount(account string, service string) string {
	req, err := http.NewRequest("GET", os.Getenv("coreApi")+"/api/core/account/"+account+"/"+service, nil)
	if err != nil {
		return ""
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return ""
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	accountService := &accountService{}
	err = json.Unmarshal(body, accountService)
	if err != nil {
		return ""
	}
	return accountService.Account
}

type ResponseDomain struct {
	Data []map[string]any `json:"data"`
	Err  map[string]any   `json:"errors"`
}

func (a *Api) getDomainValuesDataByRefId(service, route string, refId []any,
	MainField string, bodyFields []any) []map[string]any {
	ServiceAccount := getCurrentAccount(a.Account, service)
	servicePort := os.Getenv(service + "Api")
	req, err := http.NewRequest("GET", servicePort+route, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("account_uuid", ServiceAccount)
	req.Header.Set("user_uuid", a.User)
	var Apifilter []bson.M
	Apifilter = append(Apifilter, bson.M{
		"label":     MainField,
		"operation": "Equal",
		"condition": refId,
	})
	jsonF, _ := json.Marshal(Apifilter)
	q := req.URL.Query()
	q.Add("filter", string(jsonF))
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	domain := &ResponseDomain{}
	err = json.Unmarshal(body, domain)
	if err != nil {
		return nil
	}
	if domain.Err != nil {
		return nil
	}
	var Infos []map[string]any
	for _, s := range refId {
		for _, datum := range domain.Data {
			mainId := FindValueInMap(datum, MainField)
			if mainId == s {
				info := make(map[string]any)
				info[MainField] = mainId
				for _, field := range bodyFields {
					info[field.(string)] = datum[field.(string)]
				}
				Infos = append(Infos, info)
			}
		}
	}
	return Infos
}

type DomainValueResponse struct {
	Data struct {
		Value []map[string]any `json:"values"`
	} `json:"data"`
	Err map[string]any `json:"errors"`
}

func (a *Api) getDomainValuesDataByRefKey(key string) []map[string]any {
	servicePort := os.Getenv("coreApi")
	req, err := http.NewRequest("GET", servicePort+"/api/core/domains/key/"+key, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("account_uuid", a.Account)
	req.Header.Set("user_uuid", a.User)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	domain := &DomainValueResponse{}
	err = json.Unmarshal(body, domain)
	if err != nil {
		return nil
	}
	if domain.Err != nil {
		return nil
	}
	return domain.Data.Value
}
