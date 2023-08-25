package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"tservice/common/logger"
	. "tservice/models"
)

// 解析request
func DecodeAuthRequest(c context.Context, r *http.Request) (interface{}, error) {
	var req HttpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	} else {
		return req, nil
	}
	/*query := r.URL.Query()
	name := query.Get("name")
	pwd := query.Get("pwd")
	logger.Debugln(pwd, name)*/
}

func EncodeAuthResponse(c context.Context, w http.ResponseWriter, response interface{}) error {
	logger.Debugf("%v", response)
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
