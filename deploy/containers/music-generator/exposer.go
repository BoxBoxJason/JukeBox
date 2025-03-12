package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

var lock sync.Mutex

type RequestData struct {
	Request  string `json:"request"`
	Duration int    `json:"duration,omitempty"`
}

func triggerScript(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var requestData RequestData
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if requestData.Request == "" {
		http.Error(w, "Missing 'request' in the request body", http.StatusBadRequest)
		return
	}

	// Determine the duration
	duration := 10 // default value
	if requestData.Duration > 0 {
		duration = requestData.Duration
	}

	// Log the request
	fmt.Printf(time.Now().Format("2006-01-02_15-04-05"), "Processing request")

	// Acquire the lock
	lock.Lock()
	defer lock.Unlock()

	// Generated music filename
	filename := time.Now().Format("2006-01-02_15-04-05") + ".wav"

	// Execute the script with the provided content and duration
	cmd := exec.Command("/usr/bin/musicgpt-wrapper", requestData.Request, strconv.Itoa(duration), filename)
	err = cmd.Run()
	if err != nil {
		http.Error(w, fmt.Sprintf("Script execution failed: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"filename": "%s"}`, filename)))
}

func main() {
	http.HandleFunc("/", triggerScript)
	fmt.Println("Starting server on port 5556")
	log.Fatal(http.ListenAndServe(":5556", nil))
}
