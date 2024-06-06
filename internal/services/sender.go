package services

import (
	"log"
	"net/http"
	"net/url"
	"strings"
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

type Request struct {
	PeriodStart         string `json:"period_start"`
	PeriodEnd           string `json:"period_end"`
	PeriodKey           string `json:"period_key"`
	IndicatorToMoId     string `json:"indicator_to_mo_id"`
	IndicatorToMoFactId string `json:"indicator_to_mo_fact_id"`
	Value               string `json:"value"`
	FactTime            string `json:"fact_time"`
	IsPlane             string `json:"is_plane"`
	AuthUserId          string `json:"auth_user_id"`
	Comment             string `json:"comment"`
}

type Sender struct {
	client *http.Client
}

func NewSender() *Sender {
	return &Sender{
		client: &http.Client{},
	}
}

func (s *Sender) SendRequest(req *Request) {
	posturl := "https://development.kpi-drive.ru/_api/facts/save_fact"

	body := url.Values{
		"period_start":            {req.PeriodStart},
		"period_end":              {req.PeriodEnd},
		"period_key":              {req.PeriodKey},
		"indicator_to_mo_id":      {req.IndicatorToMoId},
		"indicator_to_mo_fact_id": {req.IndicatorToMoFactId},
		"value":                   {req.Value},
		"fact_time":               {req.FactTime},
		"is_plane":                {req.IsPlane},
		"auth_user_id":            {req.AuthUserId},
		"comment":                 {req.Comment},
	}

	request, err := http.NewRequest("POST", posturl, strings.NewReader(body.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Authorization", "Bearer 48ab34464a5573519725deb5865cc74c")

	resp, err := s.client.Do(request)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	log.Println("Request is sended.")

	// res := &Response{}
	// derr := json.NewDecoder(resp.Body).Decode(res)
	// if derr != nil {
	// 	panic(derr)
	// }
	// log.Println(res)
}
