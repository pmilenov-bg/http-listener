package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

var buffer []interface{}
var mutex sync.Mutex

func main() {
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetRequest(w, r)
	case http.MethodPost:
		handlePostRequest(w, r)
	case http.MethodPut:
		handlePutRequest(w, r)
	case http.MethodDelete:
		handleDeleteRequest(w, r)
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "GET request received")
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	buffer = append(buffer, data)
	bufferSize := calculateBufferSize()
	mutex.Unlock()

	if bufferSize >= 0.5*1024*16 {
		go writeBufferedDataToFile()
	}

	fmt.Fprint(w, "POST request received")
}

func handlePutRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "PUT request received")
}

func handleDeleteRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "DELETE request received")
}

func calculateBufferSize() int {
	bufferSize := 0
	for _, data := range buffer {
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Println("Error marshaling JSON data:", err)
			continue
		}
		bufferSize += len(jsonData)
	}
	return bufferSize
}

func writeBufferedDataToFile() {
	mutex.Lock()
	defer mutex.Unlock()
	// --- test of the writing to the file -----
	file, err := os.OpenFile("received-jsonl", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	for _, data := range buffer {
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Println("Error marshaling JSON data:", err)
			continue
		}

		_, err = file.Write(jsonData)
		if err != nil {
			log.Println("Error writing to file:", err)
		}
	}
	// --- test of the writing to the file -----
	buffer = nil
	fmt.Println("Buffered data written to file")
}

func dumpt()
