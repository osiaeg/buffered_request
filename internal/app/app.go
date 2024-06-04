package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/osiaeg/buffered_request/internal/config"
)

type Response struct {
	Messages struct {
		Error   []string `json:"error"`
		Warning []string `json:"warning"`
		Info    []string `json:"info"`
	} `json:"MESSAGES"`

	Data struct {
		IndicatorToMoFactId int `json:"indicator_to_mo_fact_id"`
	} `json:"DATA"`

	Status string `json:"STATUS"`
}

func saveAct(w http.ResponseWriter, r *http.Request) {
	posturl := "https://development.kpi-drive.ru/_api/facts/save_fact"

	body := url.Values{
		"period_start":            {"2024-05-01"},
		"period_end":              {"2024-05-31"},
		"period_key":              {"month"},
		"indicator_to_mo_id":      {"227373"},
		"indicator_to_mo_fact_id": {"0"},
		"value":                   {"1"},
		"fact_time":               {"2024-05-31"},
		"is_plane":                {"0"},
		"auth_user_id":            {"40"},
		"comment":                 {"buffer Last_name"},
	}

	request, err := http.NewRequest("POST", posturl, strings.NewReader(body.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Authorization", "Bearer 48ab34464a5573519725deb5865cc74c")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	res := &Response{}
	derr := json.NewDecoder(resp.Body).Decode(res)
	if derr != nil {
		panic(derr)
	}
	log.Println(res)
}

func Run() {
	cfg := config.Parse("local")
	mux := http.NewServeMux()

	mux.HandleFunc("GET /save_acts", saveAct)

	fmt.Printf("Server launched.\nURL: http://%s:%s\n", cfg.Server.Host, cfg.Server.Port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
	}
}
