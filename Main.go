package main
import (
		"fmt"
		"net/http"
		"google.golang.org/appengine"
		"time"
		"google.golang.org/appengine/datastore"
		"google.golang.org/appengine/log"
		"strings"
		"regexp"
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
	// Convert question to lowercase and remove all punctuation and other
	// non-word characters
	reg := regexp.MustCompile(`[^\w ]+`)
	qu := strings.TrimSpace(
		strings.ToLower(
			reg.ReplaceAllString(keys[0], "")))
	key := datastore.NewKey(ctx, "Fact", qu, 0, nil)
	var fact Fact
	err := datastore.Get(ctx, key, &fact)
	if err == datastore.ErrNoSuchEntity{
		fact = Fact{
			Question: qu,
			Answer: "",
			LastTime: time.Now(),
			Count: 1,
		}
	} else if (err != nil) {
		log.Errorf(ctx, "datastore.Get: %v", err)
		return
	} else {
		fact.LastTime = time.Now()
		fact.Count ++
	}
	if _, err := datastore.Put(ctx, key, &fact); err != nil {
    	log.Errorf(ctx, "datastore.Put: %v", err)
    }
    if fact.Answer == ""{
    	fact.Answer = "I don't know yet. Come back later."
    }
	fmt.Fprintln(w, fact.Answer)
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