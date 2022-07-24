package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Res http.ResponseWriter
	Req *http.Request
}

type BadRequestResp struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

func (ctx *Context) Ok() {
	ctx.Res.WriteHeader(http.StatusOK)
}

func (ctx *Context) BadRequest(typ string, key ...string) {
	resp := BadRequestResp{
		Type: typ,
	}

	if len(key) != 0 {
		resp.Key = key[0]
	}

	ctx.JSON(http.StatusBadRequest, resp)
}

func (ctx *Context) Forbidden() {
	ctx.Res.WriteHeader(http.StatusForbidden)
}

func (ctx *Context) InternalError(err error) {
	fmt.Printf("internal error: %v\n", err.Error())
	ctx.Res.WriteHeader(http.StatusInternalServerError)
}

func (ctx *Context) JSON(statusCode int, payload interface{}) {
	str, _ := json.Marshal(payload)

	ctx.Res.WriteHeader(statusCode)
	_, _ = ctx.Res.Write(str)
}

func (ctx *Context) BindBody(out interface{}) error {
	return json.NewDecoder(ctx.Req.Body).Decode(out)
}
