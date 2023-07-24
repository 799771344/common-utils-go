package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Request struct {
	Name string `json:"name"`
}

type Response struct {
	Code    int32                  `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type Router struct {
	server *http.Server
	router *mux.Router
}

func NewRouter(addr string) *Router {
	r := &Router{
		server: &http.Server{
			Addr: fmt.Sprintf(":%s", addr),
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			}),
		},
		router: mux.NewRouter(),
	}

	r.server.Handler = r.router

	return r
}

func (r *Router) Start() error {
	log.Printf("Server started at %s", r.server.Addr)
	return r.server.ListenAndServe()
}

func (r *Router) HandleFunc(ops ...interface{}) {
	for _, routerNew := range ops {
		RouterNew := routerNew.(router)
		methodGet := "GET"
		methodPost := "POST"
		switch RouterNew.Method {
		case methodGet:
			r.router.HandleFunc(RouterNew.Addr, RouterNew.Handler).Methods("GET")
		case methodPost:
			r.router.HandleFunc(RouterNew.Addr, RouterNew.Handler).Methods("GET")
		default:
			r.router.HandleFunc(RouterNew.Addr, RouterNew.Handler).Methods("GET", "POST")
		}
	}
}

type router struct {
	Addr    string
	Handler func(w http.ResponseWriter, req *http.Request)
	Method  string
}

type subrouter struct {
	Addr    string
	Handler func(w http.ResponseWriter, req *http.Request)
	Method  string
}

func (r *Router) PathPrefix(path string, ops ...interface{}) *mux.Router {
	rr := r.router.PathPrefix(path).Subrouter()
	for _, subRouter := range ops {
		subRouterNew := subRouter.(subrouter)
		methodGet := "GET"
		methodPost := "POST"
		switch subRouterNew.Method {
		case methodGet:
			rr.HandleFunc(subRouterNew.Addr, subRouterNew.Handler).Methods("GET")
		case methodPost:
			rr.HandleFunc(subRouterNew.Addr, subRouterNew.Handler).Methods("GET")
		default:
			rr.HandleFunc(subRouterNew.Addr, subRouterNew.Handler).Methods("GET", "POST")
		}
	}
	return rr
}

// 处理请求函数，在项目里定义
func handleHello(w http.ResponseWriter, req *http.Request) {
	var request Request
	//req.URL.Query()
	//req.URL.Parse("aa")
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := Response{
		Code:    0,
		Message: fmt.Sprintf("Hello, %s!", request.Name),
		Data: map[string]interface{}{
			"1": "!213",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := NewRouter("8085")

	r.PathPrefix("/api",
		subrouter{Addr: "/hello", Handler: handleHello},
		subrouter{Addr: "/helloOne", Handler: handleHello},
		subrouter{Addr: "/helloTwo", Handler: handleHello},
	)

	r.HandleFunc(
		router{Addr: "/hello", Handler: handleHello},
		router{Addr: "/helloOne", Handler: handleHello},
		router{Addr: "/helloTwo", Handler: handleHello},
	)

	if err := r.Start(); err != nil {
		log.Fatal(err)
	}
}
