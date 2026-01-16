package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"
)

type CepService struct {
	BaseURL    string
	HTTPClient *http.Client
}

type WeatherService struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

var (
	cepBaseURL     = "https://viacep.com.br/ws"
	weatherBaseURL = "http://api.weatherapi.com/v1"
)

var cepRegex = regexp.MustCompile(`^\d{8}$`)

var ErrMsgCepNotFound = "can not find zipcode"
var ErrMsgInvalidCep = "invalid zipcode"
var ErrMsgCepAPI = "cep status not ok"
var ErrMsgWeatherAPI = "weather status not ok"

var ErrCepNotFound = errors.New(ErrMsgCepNotFound)
var ErrCepAPI = errors.New(ErrMsgCepAPI)
var ErrWeatherAPI = errors.New(ErrMsgWeatherAPI)

func main() {
	godotenv.Load()
	weatherAPIKey := os.Getenv("WEATHERAPI_KEY")
	port := os.Getenv("PORT")

	httpClient := &http.Client{Timeout: 10 * time.Second}

	cepSvc := &CepService{BaseURL: cepBaseURL, HTTPClient: httpClient}
	weatherSvc := &WeatherService{BaseURL: weatherBaseURL, APIKey: weatherAPIKey, HTTPClient: httpClient}

	http.HandleFunc("/", weatherHandler(cepSvc, weatherSvc))

	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}

func (c *CepService) Lookup(ctx context.Context, cep string) (string, error) {

	var CepData struct {
		Localidade string `json:"localidade"`
		Erro       bool   `json:"erro"`
	}

	url := fmt.Sprintf("%s/%s/json/", c.BaseURL, cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", ErrCepAPI
	}

	if err := json.NewDecoder(resp.Body).Decode(&CepData); err != nil {
		return "", err
	}

	if CepData.Erro || CepData.Localidade == "" {
		return "", ErrCepNotFound
	}

	return CepData.Localidade, nil
}

func (wSvc *WeatherService) GetTempC(ctx context.Context, city string) (float64, error) {

	url := fmt.Sprintf("%s/current.json?key=%s&q=%s", wSvc.BaseURL, wSvc.APIKey, url.QueryEscape(city))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := wSvc.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, ErrWeatherAPI
	}

	var wData struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&wData); err != nil {
		return 0, err
	}
	return wData.Current.TempC, nil
}

func weatherHandler(cepSvc *CepService, weatherSvc *WeatherService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cep := r.URL.Query().Get("cep")
		if !cepRegex.MatchString(cep) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(ErrMsgInvalidCep))
			return
		}

		city, err := cepSvc.Lookup(r.Context(), cep)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(ErrMsgCepNotFound))
			return
		}

		// fmt.Printf("CEP Lookup | cep=%s city=%s\n", cep, city)

		tempC, err := weatherSvc.GetTempC(r.Context(), city)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := map[string]float64{
			"temp_C": tempC,
			"temp_F": tempC*1.8 + 32,
			"temp_K": tempC + 273,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
