// package main

// import (
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"sync"
// )

// const (
// 	maxBufferSize = 1 * 1024 * 16 // 16KB
// 	outputDir     = "data"
// )

// type JSONData struct {
// 	Data string `json:"data"`
// }

// type JSONParseError struct {
// 	Msg string
// }

// func (e *JSONParseError) Error() string {
// 	return e.Msg
// }

// func isBlackholePathValid(path string) bool {
// 	// Remove leading and trailing slashes
// 	path = strings.Trim(path, "/")
// 	// Check if it has a valid name
// 	if !isValidFilename(path) {
// 		return false
// 	}
// 	// Check if the blackhole directory exists
// 	blackholeDir := filepath.Join(outputDir, "blackhole")
// 	info, err := os.Stat(blackholeDir)
// 	if err != nil || !info.IsDir() {
// 		return false
// 	}
// 	return true
// }
// func isValidFilename(filename string) bool {
// 	// Check for special characters
// 	if strings.ContainsAny(filename, `/\:*?"<>|`) {
// 		return false
// 	}
// 	return true
// }

// func bufferJSONs(jsonChan <-chan []byte, bufferChan chan<- []byte, wg *sync.WaitGroup) {
// 	buffer := []byte{}
// 	bufferSize := 0

// 	for json := range jsonChan {
// 		buffer = append(buffer, json...)
// 		bufferSize += len(json)

// 		if bufferSize >= maxBufferSize {
// 			bufferChan <- buffer
// 			buffer = []byte{}
// 			bufferSize = 0
// 		}
// 	}

// 	if len(buffer) > 0 {
// 		bufferChan <- buffer
// 	}

// 	close(bufferChan)
// 	wg.Done()
// }

// func storeBuffer(bufferChan <-chan []byte, wg *sync.WaitGroup) {
// 	file, err := os.OpenFile(getOutputFilename("/default"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Fatalf("Error opening output file: %v", err)
// 	}
// 	defer file.Close()

// 	for buffer := range bufferChan {
// 		_, err := file.Write(buffer)
// 		if err != nil {
// 			log.Printf("Error writing to output file: %v", err)
// 		}
// 	}

// 	wg.Done()
// }

// func getOutputFilename(path string) string {
// 	return outputDir + path + ".json"
// }

// func main() {
// 	jsonChan := make(chan []byte)
// 	bufferChan := make(chan []byte)

// 	var wg sync.WaitGroup
// 	wg.Add(2)

// 	go bufferJSONs(jsonChan, bufferChan, &wg)
// 	go storeBuffer(bufferChan, &wg)

// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == http.MethodPost {
// 			body, err := ioutil.ReadAll(r.Body)
// 			if err != nil {
// 				http.Error(w, "Error reading request body", http.StatusInternalServerError)
// 				return
// 			}
// 			defer r.Body.Close()

// 			path := r.URL.Path
// 			if path == "/" {
// 				path = "/default"
// 			}

// 			if strings.HasPrefix(path, "/blackhole") {
// 				if !isBlackholePathValid(path) {
// 					log.Printf("Invalid blackhole path: %s", path)
// 					http.Error(w, "Invalid blackhole path", http.StatusBadRequest)
// 					return
// 				}

// 				outputFile := getOutputFilename(path)

// 				jsonChan <- body

// 				w.Write([]byte("Request received and processed"))
// 				log.Printf("JSON data received and stored in file: %s", outputFile)
// 			} else {
// 				w.Write([]byte("Invalid request method"))
// 			}
// 		}
// 	})

// 	log.Fatal(http.ListenAndServe(":8080", nil))

// 	close(jsonChan)
// 	wg.Wait()
// }

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	maxBufferSize = 1 * 1024 * 16 // 16KB
	outputDir     = "data"
)

type JSONData struct {
	Data string `json:"data"`
}

type JSONParseError struct {
	Msg string
}

func (e *JSONParseError) Error() string {
	return e.Msg
}

func main() {
	// Create the output directory if it doesn't exist
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatal("Error creating output directory:", err)
	}

	jsonChan := make(chan []byte)
	bufferChan := make(chan []byte)

	var wg sync.WaitGroup
	wg.Add(2)

	go bufferJSONs(jsonChan, bufferChan, &wg)
	go storeBuffer(bufferChan, &wg)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()

			path := r.URL.Path
			if path == "/" {
				path = "/default"
			}

			if strings.HasPrefix(path, "/blackhole") {
				if !isBlackholePathValid(path) {
					log.Printf("Invalid blackhole path: %s", path)
					http.Error(w, "Invalid blackhole path", http.StatusBadRequest)
					return
				}

				outputFile := getOutputFilename(path)

				jsonChan <- body

				w.Write([]byte("Request received and processed"))
				log.Printf("JSON data received and stored in file: %s", outputFile)
			} else {
				jsonChan <- body

				w.Write([]byte("Request received"))
			}
		} else {
			w.Write([]byte("Invalid request method"))
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

	close(jsonChan)
	wg.Wait()
}

func bufferJSONs(jsonChan <-chan []byte, bufferChan chan<- []byte, wg *sync.WaitGroup) {
	buffer := []byte{}
	bufferSize := 0

	for json := range jsonChan {
		buffer = append(buffer, json...)
		bufferSize += len(json)

		if bufferSize >= maxBufferSize {
			bufferChan <- buffer
			buffer = []byte{}
			bufferSize = 0
		}
	}

	if len(buffer) > 0 {
		bufferChan <- buffer
	}

	close(bufferChan)
	wg.Done()
}

func storeBuffer(bufferChan <-chan []byte, wg *sync.WaitGroup) {
	for buffer := range bufferChan {
		outputFile := filepath.Join(outputDir, "received-jsons")

		err := ioutil.WriteFile(outputFile, buffer, 0644)
		if err != nil {
			log.Printf("Error writing to output file: %v", err)
		}
	}

	wg.Done()
}

func isBlackholePathValid(path string) bool {
	// Remove leading and trailing slashes
	path = strings.Trim(path, "/")

	// Check if it has a valid name
	if !isValidFilename(path) {
		return false
	}
	return true
}

func isValidFilename(filename string) bool {
	// Check for special characters
	if strings.ContainsAny(filename, `/\:*?"<>|`) {
		return false
	}

	return true
}

func getOutputFilename(path string) string {
	// Remove leading and trailing slashes
	path = strings.Trim(path, "/")

	return filepath.Join(outputDir, "blackhole", path)
}
