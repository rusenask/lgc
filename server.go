package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

// StublistHandler gets stubs, e.g.: stubo/api/get/stublist?scenario=first
func StublistHandler(w http.ResponseWriter, r *http.Request) {
	scenario := r.URL.Query().Get("scenario")
	fmt.Println("got:", r.URL.Query())
	fmt.Println("Scenario name:", scenario)
}

func main() {
	mux := bone.New()
	mux.Get("/stubo/api/get/stublist", http.HandlerFunc(StublistHandler))
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}
