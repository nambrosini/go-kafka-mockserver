package routes

import (
	"encoding/json"
	m "github.com/nambrosini/go-kafka-mockserver/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"log"
	"net/http"
	"strings"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	route := path[len(path)-1]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	m.KafkaMessages.With(prometheus.Labels{"topic": route}).Inc()

	err = validateJsonRequest(string(body))

	if err != nil {
		log.Println("Invalid JSON", err)
		w.WriteHeader(http.StatusBadRequest)
		m.KafkaInvalidMessages.With(prometheus.Labels{"topic": route}).Inc()
		return
	}

	log.Println("Valid JSON", string(body))
	m.KafkaValidMessages.With(prometheus.Labels{"topic": route}).Inc()
	w.WriteHeader(http.StatusCreated)
}

func validateJsonRequest(body string) error {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(body), &data)
	return err
}
