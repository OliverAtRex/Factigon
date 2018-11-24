package main
import (
		"fmt"
		"net/http"
		"google.golang.org/appengine"
		"time"
		"google.golang.org/appengine/datastore"
		"google.golang.org/appengine/log"
)
type Fact struct{
	Question string
	Answer string
	LastTime time.Time
	Count int32
}
func askHandler(w http.ResponseWriter, r *http.Request){
	ctx := appengine.NewContext(r)
	keys, ok := r.URL.Query()["qu"]
	if !ok || len(keys[0]) < 1 {
		return
	}
	qu := keys[0]
	key := datastore.NewKey(ctx, "Fact", qu, 0, nil)
	fact := Fact{
		Question: qu,
		//Answer:
		LastTime: time.Now(),
		Count: 1,
	}
	if _, err := datastore.Put(ctx, key, &fact); err != nil {
        log.Errorf(ctx, "datastore.Put: %v", err)
    }
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