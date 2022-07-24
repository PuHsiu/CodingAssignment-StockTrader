package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

type RestService interface {
	RestRouter(*Router)
}

type Router struct {
	muxRouter *mux.Router
}

func (r *Router) Get(path string, fn func(Context)) {
	r.muxRouter.HandleFunc(path, wrapContext(fn)).Methods(http.MethodGet)
}

func (r *Router) Post(path string, fn func(Context)) {
	r.muxRouter.HandleFunc(path, wrapContext(fn)).Methods(http.MethodPost)
}

func (r *Router) Put(path string, fn func(Context)) {
	r.muxRouter.HandleFunc(path, wrapContext(fn)).Methods(http.MethodPut)
}

func (r *Router) Patch(path string, fn func(Context)) {
	r.muxRouter.HandleFunc(path, wrapContext(fn)).Methods(http.MethodPatch)
}

func (r *Router) Delete(path string, fn func(Context)) {
	r.muxRouter.HandleFunc(path, wrapContext(fn)).Methods(http.MethodDelete)
}

func wrapContext(fn func(Context)) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		fn(Context{Req: req, Res: res})
	}
}
