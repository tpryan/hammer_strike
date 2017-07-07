package load

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

func init() {
	http.HandleFunc("/load/", indexHandler)
	http.HandleFunc("/load/flush", flushHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	token := r.URL.Query().Get("token")

	_, err := cache(c, token)
	if err != nil {
		sendError(c, w, "Memcache Error: "+err.Error())
		return
	}
	sendMessage(w, "Success, token="+token)
}

func flushHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	memcache.Flush(c)
	sendMessage(w, "Cache flushed.")
}

func sendError(c context.Context, w http.ResponseWriter, msg string) {
	log.Errorf(c, msg)
	content := "{ \"error\" : \"" + msg + "\" }"
	sendJSON(w, content, http.StatusInternalServerError)
}

func sendMessage(w http.ResponseWriter, msg string) {
	content := "{ \"msg\" : \"" + msg + "\" }"
	sendJSON(w, content, http.StatusOK)
}

func sendJSON(w http.ResponseWriter, content string, code int) {
	w.Header().Set("Content-Type", "application/json;  charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if content == "null" || content == "[]" {
		code = http.StatusNotFound
		content = "{ \"error\" : \"Not Found\" }"
	}
	w.WriteHeader(code)

	fmt.Fprint(w, content)
}

func cache(c context.Context, token string) (uint64, error) {
	instanceID := token + "_" + appengine.InstanceID()
	total, err := memcache.Increment(c, token+"_total", 1, 0)

	if err != nil {
		return 0, errors.New("cannot increment session total: " + err.Error())
	}

	// Record the begining and end of the token's usage
	if total == 1 {
		if err := start(c, token); err != nil {
			return 0, errors.New("Cannot record start: " + err.Error())
		}
	} else {
		if err := last(c, token); err != nil {
			return 0, errors.New("Cannot record end: " + err.Error())
		}
	}

	// Increment the number of requests handled by this instance
	result, err := memcache.Increment(c, instanceID, 1, 0)
	if err != nil {
		return 0, errors.New("cannot increment instance session count: " + err.Error())
	}

	// If the result is greater then 1 then we have seen this instance before,
	// so our work is done.
	if result > 1 {
		return total, nil
	}

	// Increment the total count of all instances
	if _, err = memcache.Increment(c, token+"_totalInstances", 1, 0); err != nil {
		return 0, errors.New("cannot increment instance count: " + err.Error())
	}

	// Add the instance id to the lists of ID
	if err := recordInstance(c, token, instanceID); err != nil {
		return 0, err
	}

	return total, nil

}

func start(c context.Context, token string) error {
	start := &memcache.Item{
		Key:   token + "_start",
		Value: []byte(strconv.Itoa(int(time.Now().UnixNano()))),
	}
	return memcache.Set(c, start)
}

func last(c context.Context, token string) error {
	last := &memcache.Item{
		Key:   token + "_end",
		Value: []byte(strconv.Itoa(int(time.Now().UnixNano()))),
	}
	return memcache.Set(c, last)
}

func recordInstance(c context.Context, token string, instanceID string) error {
	iList := whichList(instanceID, token)

	item0, err := memcache.Get(c, iList)
	if err != nil && err != memcache.ErrCacheMiss {
		return errors.New("cannot get list of instances: " + err.Error())
	}
	if err != nil && err == memcache.ErrCacheMiss {
		item1 := &memcache.Item{
			Key:   iList,
			Value: []byte(instanceID),
		}
		if err := memcache.Set(c, item1); err != nil {
			return errors.New("cannot create list of instances: " + err.Error())
		}
	}
	if err == nil {
		list := string(item0.Value)
		b := []byte(list + "|" + instanceID)
		item0.Value = b
		if err := memcache.Set(c, item0); err != nil {
			return errors.New("cannot append list of instances: " + err.Error())
		}
	}
	return nil
}

func whichList(i, token string) string {
	return token + "_instances-" + i[len(i)-1:]
}
