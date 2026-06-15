// EXAMPLE FILE — clean server stub, no findings.
// Used in the demo to show Atheon only flags what actually matches.
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "myapp running")
	})

	fmt.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
