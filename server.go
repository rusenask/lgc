package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

func MyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %v\n", bone.GetValue(r, "id"))
}

func main() {
	mux := bone.New()
	mux.Get("/stubo/api/get/stublist/:id", http.HandlerFunc(MyHandler))
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}
