package main
import (
		"context"
		"fmt"
		"time"
		"strings"
		"regexp"
		
		"net/http"
		"google.golang.org/appengine"
		"google.golang.org/appengine/datastore"
		"google.golang.org/appengine/log"
		"google.golang.org/appengine/memcache"
)
type Fact struct{
	Question string
	Answer string
	LastTime time.Time
	Count int32
}
func askHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	params, ok := r.URL.Query()["qu"]
	if !ok || len(params[0]) < 1 {
		return
	}
	rawQ := params[0]
	
	// Convert question to lowercase and remove all punctuation and other
	// non-word characters
	reg := regexp.MustCompile(`[^\w ]+`)
	qu := strings.TrimSpace(
		strings.ToLower(
			reg.ReplaceAllString(rawQ, "")))
	key := datastore.NewKey(ctx, "Fact", qu, 0, nil)
	var fact Fact
	err := datastore.Get(ctx, key, &fact)
	if err == datastore.ErrNoSuchEntity{
		fact = Fact{
			Question: rawQ,
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
    if fact.Answer == "" {
    	fact.Answer = "I don't know yet. Come back later."
    } else if strings.Contains(fact.Answer, "[RANDOM_FACT]") {
    	rf := getRandomFact(ctx)
    	fact.Answer = strings.Replace(fact.Answer,"[RANDOM_FACT]", rf, -1)
	}
	fmt.Fprintln(w, fact.Answer)
}

func getRandomFact(ctx context.Context) string {
	q := datastore.NewQuery("Fact")

	// If we stored a cursor during a previous request, use it.
	item, err := memcache.Get(ctx, "fact_cursor")
	if err == nil {
		cursor, err := datastore.DecodeCursor(string(item.Value))
		if err == nil {
			q = q.Start(cursor)
		}
	}

	// Iterate over the results looking for any answer
	t := q.Run(ctx)
	var fact Fact
	for {
		_, err := t.Next(&fact)
		if err == datastore.Done {
			memcache.Delete(ctx, "fact_cursor")
			return "Hmm, I'm not sure!"
        }
        if err != nil {
			log.Errorf(ctx, "fetching next Fact: %v", err)
			return "Internal error"
        }
        if fact.Answer != "" {
        	break
    	}
	}

	// Get updated cursor and store it for next time.
	if cursor, err := t.Cursor(); err == nil {
		memcache.Set(ctx, &memcache.Item{
			Key:   "fact_cursor",
			Value: []byte(cursor.String()),
		})
	}
	
	return fact.Question
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