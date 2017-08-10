// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package distributor

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/taskqueue"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	http.HandleFunc("/distributor", mainHandler)
	http.HandleFunc("/distributor/flush", flushHandler)
	http.HandleFunc("/distributor/list", listHandler)
	http.HandleFunc("/distributor/url", urlHandler)
	http.HandleFunc("/distributor/report", reportHandler)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	token := r.URL.Query().Get("token")
	param1 := r.URL.Query().Get("n")
	param2 := os.Getenv("TARGET_QPS")

	lns, err := List(c)
	if err != nil {
		handleError(c, w, errors.New("Could not retrieve a list of VMs in the project: "+err.Error()))
		return
	}

	if len(lns) < 1 {
		handleError(c, w, errors.New("There are no VMs currently running"))
		return
	}

	log.Debugf(c, "param1: %s param2: %s len(lns): %d", param1, param2, len(lns))
	n, cc, err := calcRates(param1, param2, len(lns))
	if err != nil {
		handleError(c, w, errors.New("Problem with parameters: "+err.Error()))
		return
	}

	//Using task queues because go on GAE does not support parallel processing.
	for _, ln := range lns {
		log.Infof(c, "host: %s n: %s c: %s token %s", ln.IP, n, cc, token)
		t := taskqueue.NewPOSTTask("/distributor/url", map[string][]string{
			"host":  {ln.IP},
			"n":     {n},
			"c":     {cc},
			"token": {token},
		})
		_, err := taskqueue.Add(c, t, "urlcaller")
		if err != nil {
			handleError(c, w, errors.New("Problem adding task to queue: "+err.Error()))
		}

	}

	sendMessage(w, "Success")
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	vms, err := List(c)
	if err != nil {
		handleError(c, w, errors.New("Could not retrieve a list of VMs in the project: "+err.Error()))
		return
	}

	b, err := json.Marshal(vms)
	if err != nil {
		handleError(c, w, errors.New("Could not marshal the json: "+err.Error()))
		return
	}
	sendJSON(w, string(b), http.StatusOK)
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	n := r.FormValue("n")
	if len(n) == 0 {
		handleError(c, w, errors.New("Invalid n: "+n))
	}
	cc := r.FormValue("c")
	if len(cc) == 0 {
		handleError(c, w, errors.New("Invalid c: "+cc))
	}
	host := r.FormValue("host")
	if len(host) == 0 {
		handleError(c, w, errors.New("Invalid host: "+host))
	}

	token := r.FormValue("token")

	u := fmt.Sprintf("http://%s:30000?n=%s&c=%s&token=%s", host, n, cc, token)

	cWithDeadline, _ := context.WithTimeout(c, 1*time.Minute)

	client := urlfetch.Client(cWithDeadline)
	_, err := client.Get(u)
	if err != nil {
		handleError(c, w, errors.New("Problem sending the load to "+u+" Error: "+err.Error()))
		return
	}
	sendMessage(w, "Success")
	return

}

func flushHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	memcache.Flush(c)
	sendMessage(w, "Cache flushed.")
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	token := r.URL.Query().Get("token")

	rep, err := GetReport(c, token)

	if err != nil {
		handleError(c, w, err)
		return
	}

	b, err := json.Marshal(rep)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	sendJSON(w, string(b), http.StatusOK)

}

func calcRates(n string, c string, count int) (string, string, error) {
	nInt, err := strconv.Atoi(n)
	if err != nil {
		return "", "", errors.New("Could not get valid value for `n`: " + n)
	}

	cInt, err := strconv.Atoi(c)
	if err != nil {
		return "", "", errors.New("Could not get valid value for env variable `TARGET_QPS: " + c)
	}

	nodeN := nInt / count
	nodeC := cInt / count

	// Ensures that C never exceeds N cause if that happens Apache Bench fails.
	if nodeC > nodeN {
		nodeC = nodeN
	}
	return strconv.Itoa(nodeN), strconv.Itoa(nodeC), nil
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

func sendMessage(w http.ResponseWriter, msg string) {
	content := "{ \"msg\" : \"" + msg + "\" }"
	sendJSON(w, content, http.StatusOK)
}

func handleError(c context.Context, w http.ResponseWriter, err error) {
	content := "{ \"error\" : \"" + err.Error() + "\" }"
	sendJSON(w, content, http.StatusInternalServerError)
	log.Errorf(c, err.Error())
}
