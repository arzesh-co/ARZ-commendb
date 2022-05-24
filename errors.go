package CommenDb

func GetErrors(key string, account string, lang string, params map[string]string) *ResponseErrors {
	res := getError(key, account, lang, params)
	if res == nil {
		return nil
	}
	return res
}
