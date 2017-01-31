package main

import (
	"encoding/json"
	//	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func GlobalServiceStatus(w http.ResponseWriter, req *http.Request) {

	for _, value := range State {
		if value != "OK" {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}
	bs, err := json.Marshal(State)
	if err != nil {
		//		TODO..do not panic; use a recovery handler
		panic(err)
	}
	w.Write(bs)
}

func SingleServiceStatus(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	ts := vars["testService"]
	if val, ok := State[ts]; !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("☄ hey!, the requested test service could not be found."))

	} else {
		bs, err := json.Marshal(val)
		if err != nil {
			//		TODO..do not panic; use a recovery handler
			panic(err)
		}
		if val != "OK" {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		w.Write(bs)
	}
}

type supportCORS struct {
	router *mux.Router
}

func (server *supportCORS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", `POST, GET, OPTIONS,
        	PUT, DELETE`)
		w.Header().Set("Access-Control-Allow-Headers",
			`Accept, Content-Type, Content-Length, Accept-Encoding,
            X-CSRF-Token, Authorization`)
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	server.router.ServeHTTP(w, r)
}

func servHome(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "index/home", nil)
}

//AddHandlers creates a router and adds handlers
func AddHandlers() *mux.Router {
	router := mux.NewRouter()
	http.Handle("/", &supportCORS{router})
	router.HandleFunc("/services", GlobalServiceStatus).Methods("GET")
	router.HandleFunc("/service/{testService}", SingleServiceStatus).Methods("GET")
	router.HandleFunc("/", servHome).Methods("GET")
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	return router
}
