package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang/gddo/httputil/header"
)

const (
	TimeUnit   = 1 * time.Millisecond
	BaseFactor = 50
	Extra      = 1
)

type TextRequest struct {
	Length int      `json:"length"`
	Ids    []string `json:"ids"`
	Texts  []string `json:"texts"`
}

type TextResponse struct {
	Length      int      `json:"length"`
	Ids         []string `json:"ids"`
	Texts       []string `json:"texts"`
	ProcessTime int      `json:"process_time"`
}

func processTexts(tReq TextRequest) TextResponse {
	length := tReq.Length
	processFactor := BaseFactor + (length-1)*Extra
	fmt.Printf("Request with length = %d, Process for %d ms.\n", length, processFactor)
	time.Sleep(time.Duration(processFactor) * TimeUnit)
	outputs := make([]string, length)
	copy(outputs, tReq.Texts)
	return TextResponse{Length: length, Ids: tReq.Ids, Texts: outputs, ProcessTime: processFactor}
}

func process(w http.ResponseWriter, r *http.Request) {
	var tReq TextRequest
	if r.Method == "POST" {
		if r.Header.Get("Content-Type") != "" {
			value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
			if value != "application/json" {
				msg := "Content-Type header is not application/json"
				http.Error(w, msg, http.StatusUnsupportedMediaType)
				return
			}
		}
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			msg := "Error when reading request body."
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &tReq)
		if err != nil {
			msg := "Error when unmarshal request body."
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		outputs := processTexts(tReq)

		respBody, err := json.Marshal(outputs)
		if err != nil {
			msg := "Error when marshal response body."
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(respBody)
	}
}

func main() {
	address := flag.String("address", "localhost:8089", "host:port")
	http.HandleFunc("/infer", process)

	err := http.ListenAndServe(*address, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed.")
	} else if err != nil {
		fmt.Printf("Server failed when starting: %s", err)
		os.Exit(1)
	}
}
