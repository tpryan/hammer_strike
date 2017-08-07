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
	if result == 1 {
		log.Debugf(c, "Instance Never seen before: "+instanceID)
		// Add the instance id to the lists of ID
		if err := recordInstance(c, token, instanceID); err != nil {
			return 0, errors.New("cannot add another instance to list: " + err.Error())
		}
		// Increment the total count of all instances
		if _, err = memcache.Increment(c, token+"_totalInstances", 1, 0); err != nil {
			return 0, errors.New("cannot increment the total instance count: " + err.Error())
		}
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
	log.Debugf(c, "Instance Being Recorded: "+instanceID)

	_, err := memcache.Get(c, iList)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			item1 := &memcache.Item{
				Key:   iList,
				Value: []byte(instanceID),
			}
			if err := memcache.Set(c, item1); err != nil {
				return errors.New("cannot create list of instances: " + err.Error())
			}
		} else {
			return errors.New("1st attempt - cannot get list of instances: " + err.Error())
		}
	}

	listItem := "|" + instanceID

	if err := appendInstanceList(c, iList, listItem, 0); err != nil {
		return errors.New("cannot append an instance list: " + err.Error())
	}

	return nil
}

func appendInstanceList(c context.Context, list string, listItem string, count int) error {
	log.Debugf(c, "Instance Being Appended to List: "+listItem+" list:"+list)
	item, err := memcache.Get(c, list)
	if err != nil {
		return errors.New("cannot get list of instances in edit instance list: " + err.Error())
	}
	item.Value = fastAppend(item.Value, listItem)

	if err = memcache.Set(c, item); err != nil {
		if count > 5 {
			return errors.New("5th attempt - cannot append list of instances: " + err.Error())
		}
		if err == memcache.ErrCASConflict {
			return appendInstanceList(c, list, listItem, count+1)
		}
		return err
	}
	log.Debugf(c, "Successfully appended to list: "+listItem+" list:"+list)
	return nil
}

func fastAppend(list []byte, newitem string) []byte {
	length := len(list) + len(newitem)
	result := make([]byte, length)
	c := 0

	for _, value := range list {
		result[c] = value
		c++
	}
	for _, value := range newitem {
		result[c] = byte(value)
		c++
	}

	return result

}

func whichList(i, token string) string {
	return token + "_instances_" + i[len(i)-1:]
}
