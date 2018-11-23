package main
import (
		"fmt"
		"net/http"
		"google.golang.org/appengine"
)
func askHandler(w http.ResponseWriter, r *http.Request){
	keys, ok := r.URL.Query()["qu"]
	if !ok || len(keys[0]) < 1 {
		return
	}
	qu := keys[0]
	fmt.Fprintln(w, "You asked:", qu)
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}
func main() {
	http.HandleFunc("/ask", askHandler)
	http.HandleFunc("/", indexHandler)
	appengine.Main()
}