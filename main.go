package main

import (
	"labix.org/v2/mgo"
	"html/template"
	"net/http"
	"log"
)

var session *mgo.Session

var index = template.Must(template.ParseFiles(
	"templates/_base.html",
	"templates/index.html",
))


func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})

}

func hello(w http.ResponseWriter, req *http.Request) {
	s := session.Clone()
	defer s.Close()

	// set up collection and query
	coll := s.DB("gostbook").C("entries")
	query := coll.Find(nil).Sort("-timestamp")

	var entries []Entry

	// execute the query
	if err := query.All(&entries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	// execute the template
	if err := index.Execute(w, entries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	var err error
	session, err = mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", hello) 
	http.HandleFunc("/sign", sign)

	if err := http.ListenAndServe(":8080", Log(http.DefaultServeMux)); err != nil {
		panic(err)
	}
}
