package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpRequest struct {
	Request *http.Request          `json:"-"`
	Method  string                 `json:"method"       binding:"required"`
	Session string                 `json:"session"`
	Params  map[string]interface{} `json:"params"`
}

func NewParseRequest(c *gin.Context) *HttpRequest {
	var req HttpRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(HTTP_200, NewErrorResonse(ERR_CODE_INVALID_PARAMS, err.Error()))
		return nil
	}
	req.Request = c.Request
	if len(req.Session) <= 0 {
		if s, err := c.Cookie("session"); err == nil {
			req.Session = s
		}
	}
	return &req
}
