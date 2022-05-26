package CommenDb

import (
	"encoding/json"
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
type fieldsEntities struct {
	DbName     string            `json:"db_name" bson:"db_name"`
	Title      map[string]string `json:"title" bson:"title"`
	MultiLang  bool              `json:"multi_lang" bson:"multi_lang"`
	Validators []struct {
		Rule  string `json:"rule" bson:"rule"`
		Param int64  `json:"param" bson:"param"`
	} `json:"validators" bson:"validators"`
	Required     bool                   `json:"required" bson:"required"`
	DataType     string                 `json:"data_type" bson:"data_type"`
	DisplayType  string                 `json:"display_type" bson:"display_type"`
	Conf         map[string]interface{} `json:"conf" bson:"conf"`
	Features     []string               `json:"features" bson:"features"`
	DenyRoleKeys []string               `json:"deny_role_keys" bson:"deny_role_keys"`
}
type validation struct {
	Validator []fieldsEntities `json:"validator"`
}

func getUserRole(user string, account string) ([]string, error) {
	req, err := http.NewRequest("GET", os.Getenv("userApi")+"/api/user/roles/"+user, nil)
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
func getEndPointFileds(route string, method string, service string) ([]fieldsEntities, error) {
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
	return userInfo.Validator, nil
}
