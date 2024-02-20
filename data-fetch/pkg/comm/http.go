package comm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func AddRequestHeader(req *http.Request) *http.Request {
	req.Header.Add("X-RapidAPI-Key", "8f33e913famsh7d8fd47144bb6b9p1bf9c9jsn817f5049196f")
	req.Header.Add("X-RapidAPI-Host", "api-football-v1.p.rapidapi.com")
	return req
}

func UnmarshalData(body []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetHttpBody(url string, args ...interface{}) (map[string]interface{}, error) {

	var err error
	var req *http.Request
	if args != nil {
		req, err = http.NewRequest("GET", fmt.Sprintf(url, args...), nil)
	} else {
		req, err = http.NewRequest("GET", url, nil)
	}

	if err != nil {
		return nil, err
	}

	req = AddRequestHeader(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	data, err := UnmarshalData(body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
