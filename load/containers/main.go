package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/last", lastHandler)
	http.HandleFunc("/healthz", healthHandler)
	if err := http.ListenAndServe(":80", nil); err != nil {
		fmt.Println(err.Error())
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	n := r.URL.Query().Get("n")
	c := r.URL.Query().Get("c")

	urltohit := os.Getenv("LOAD_URL")

	if len(urltohit) == 0 {
		handleError(w, errors.New("Could not get LOAD_URL from env"))
		return
	}

	if len(n) == 0 {
		handleError(w, errors.New("n request variable not set"))
		return
	}

	if len(c) == 0 {
		handleError(w, errors.New("c request variable not set"))
		return
	}

	urltohit += "?token=" + token

	results, err := ab(n, c, urltohit)
	if err != nil {
		if err.Error() == "exit status 22" {
			handleError(w, errors.New("Error: Might be an issue with env variable `urltohit` value="+urltohit))
			return
		}
		handleError(w, err)
		return
	}
	writeLog(results)
	sendMessage(w, "success")
	return

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	sendMessage(w, "success")
	return
}

func lastHandler(w http.ResponseWriter, r *http.Request) {
	code := http.StatusOK
	content := ""

	dat, err := ioutil.ReadFile("/go/src/abrunner/logs/last.log")

	if err != nil {
		code = http.StatusNotFound
		content = "{ \"error\" : \"Not Found\" }"
	}

	content = string(dat)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain;  charset=UTF-8")
	w.WriteHeader(code)

	fmt.Fprint(w, content)

	return
}

func ab(n, c, u string) ([]byte, error) {
	args := []string{"-l", "-n", n, "-c", c, "-v", "2", "-q", u}
	cmd := "ab"
	return exec.Command(cmd, args...).Output()
}

func writeLog(data []byte) error {
	name := "/go/src/abrunner/logs/" + time.Now().Format("20060102150405.9999999") + ".log"
	last := "/go/src/abrunner/logs/last.log"

	if err := ioutil.WriteFile(last, data, 0644); err != nil {
		return errors.New("Error writing last log file: " + err.Error())
	}

	return ioutil.WriteFile(name, data, 0644)
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
	fmt.Println(time.Now().Format("2006/01/02 15:04:05") + " SUCCESS: " + msg)
	sendJSON(w, content, http.StatusOK)
}

func handleError(w http.ResponseWriter, err error) {
	content := "{ \"error\" : \"" + err.Error() + "\" }"
	fmt.Println(time.Now().Format("2006/01/02 15:04:05") + " ERROR: " + err.Error())
	sendJSON(w, content, http.StatusInternalServerError)
}
