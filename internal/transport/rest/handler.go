package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/osiaeg/buffered_request/internal/services"
)

type Producer interface {
	Produce(context.Context, []byte, []byte) error
}

type Reporter interface {
	Report() (interface{}, error)
}

func NewHandler(producer Producer) *Handler {
	return &Handler{
		p: producer,
	}
}

type Handler struct {
	p Producer
}

func (h *Handler) SaveAct(w http.ResponseWriter, r *http.Request) {
	log.Println(fmt.Sprintf("%s %s", r.Method, r.URL.Path))

	decoder := json.NewDecoder(r.Body)

	if err := checkBracket(decoder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Incorrect start bracket.")
		return
	}

	for decoder.More() {
		var request services.Request
		if err := decoder.Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Invalid request data.")
			return
		}

		key := []byte(fmt.Sprintf("address-%s", r.RemoteAddr))
		value, _ := json.Marshal(request)

		if err := h.p.Produce(r.Context(), key, value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Can't write msg to kafka.")
			return
		}

		w.WriteHeader(http.StatusAccepted)

		log.Println("Body from request succesfully written to kafka.")
	}
	if err := checkBracket(decoder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Incorrect start bracket.")
		return
	}

}

func checkBracket(decoder *json.Decoder) error {
	_, err := decoder.Token()
	return err
}

func Mux(prod Producer) http.Handler {
	mux := http.NewServeMux()
	handler := NewHandler(prod)

	mux.HandleFunc("POST /save_fact", handler.SaveAct)
	return mux
}
