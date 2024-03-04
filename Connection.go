package CommenDb

import (
	"encoding/json"
	"fmt"
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

func (a *Api) getError(key string, account string, lang string, params map[string]string) *ResponseErrors {

	Error := FindError(key)

	ResErr := convertErrorToResponseErr(Error, lang)

	ResErr = setParamsToResponseErr(ResErr, Error.Params, params, lang)
	return ResErr
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

func (a *Api) getUserRole(user string, account string) ([]string, error) {

	userRoles := GetUserRolesFromToken(a.UserToken)

	return userRoles, nil
}
func (a *Api) getEndPointFileds(route string, method string, service string) (*validation, error) {
	req, err := http.NewRequest("GET", os.Getenv("meta_model")+"/api/meta-model/endpoints/validators", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("service", service)
	q.Add("route", route)
	q.Add("method", method)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("client", a.ClientToken)
	client := &http.Client{
		Transport: &http.Transport{},
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error is :", err.Error())
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error is :", err.Error())
		return nil, err
	}
	userInfo := &validation{}
	err = json.Unmarshal(body, userInfo)
	if err != nil {
		fmt.Println("error is :", err.Error())
		return nil, err
	}
	return userInfo, nil
}

type accountService struct {
	Account string `json:"data"`
}

func (a *Api) GetCurrentAccount(account string, service string) string {
	serviceAccount := GetStringValueFromToken(a.ClientToken, service)

	return serviceAccount
}

type ResponseDomain struct {
	Data []map[string]any `json:"data"`
	Err  map[string]any   `json:"errors"`
}

func (a *Api) getDomainValuesDataByRefId(service, route string, refId []any,
	MainField string, bodyFields []any) []map[string]any {
	ServiceAccount := a.GetCurrentAccount(a.Account, service)
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
		"operation": "=",
		"condition": refId,
	})
	jsonF, _ := json.Marshal(Apifilter)
	q := req.URL.Query()
	q.Add("filter", string(jsonF))
	q.Add("limit", "100")
	req.URL.RawQuery = q.Encode()
	client := &http.Client{
		Transport: &http.Transport{},
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error is :", err.Error())
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
	servicePort := os.Getenv("domain")
	req, err := http.NewRequest("GET", servicePort+"/api/domain/domains/key/"+key, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("client", a.ClientToken)

	client := &http.Client{
		Transport: &http.Transport{},
	}
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
