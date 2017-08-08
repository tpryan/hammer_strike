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
	http.HandleFunc("/log", logHandler)
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
	writeLog(results, token)
	sendMessage(w, "success - handled ab for token:"+token+" on ip:"+r.RemoteAddr)
	return

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	sendMessage(w, "success - system is healthy")
	return
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	podname := os.Getenv("HOSTNAME")
	content := ""

	if len(token) == 0 {
		handleError(w, errors.New("You did not pass a token"))
		return
	}

	dat, err := ioutil.ReadFile("/go/src/abrunner/logs/" + token + ".log")

	if err != nil {
		content = "No logs have been generated yet on pod:  " + podname + " " + err.Error()
	} else {
		content = podname + "\n" + string(dat)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain;  charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, content)

	return
}

func ab(n, c, u string) ([]byte, error) {
	args := []string{"-l", "-n", n, "-c", c, "-v", "2", "-q", u}
	cmd := "ab"
	return exec.Command(cmd, args...).Output()
}

func writeLog(data []byte, token string) error {
	name := "/go/src/abrunner/logs/" + token + ".log"

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
