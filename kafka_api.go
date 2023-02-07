package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var LogFile = os.Getenv("LOG_FILE")
var kafkaMessages = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "kafka_messages",
		Help: "Number of messages from the topic",
	},
	[]string{"topic"},
)

var kafkaValidMessages = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "kafka_valid_messages",
		Help: "Number of valid JSON messages from the topic",
	},
	[]string{"topic"},
)

var kafkaInvalidMessages = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "kafka_invalid_messages",
		Help: "Number of invalid JSON messages from the topic",
	},
	[]string{"topic"},
)

func main() {
	setupLogger()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/logs", logsHandler)
	http.Handle("/metrics", promhttp.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func setupLogger() {
	logFile, err := os.OpenFile(LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	route := path[len(path)-1]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	kafkaMessages.With(prometheus.Labels{"topic": route}).Inc()

	err = validateJsonRequest(string(body))

	if err != nil {
		log.Println("Invalid JSON", err)
		w.WriteHeader(http.StatusBadRequest)
		kafkaInvalidMessages.With(prometheus.Labels{"topic": route}).Inc()
		return
	}

	log.Println("Valid JSON", string(body))
	kafkaValidMessages.With(prometheus.Labels{"topic": route}).Inc()
	w.WriteHeader(http.StatusCreated)
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	data, err := os.ReadFile(LogFile)
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}

func validateJsonRequest(body string) error {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(body), &data)
	return err
}
