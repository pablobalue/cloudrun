package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherHandlerSuccess(t *testing.T) {

	viaCepAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"localidade":"Rio de Janeiro"}`))
	}))
	defer viaCepAPI.Close()

	weatherAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"current":{"temp_c":30.5}}`))
	}))
	defer weatherAPI.Close()

	handler := weatherHandler(
		&CepService{BaseURL: viaCepAPI.URL, HTTPClient: http.DefaultClient},
		&WeatherService{BaseURL: weatherAPI.URL, HTTPClient: http.DefaultClient},
	)

	req := httptest.NewRequest(http.MethodGet, "/?cep=12345678", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", rec.Code)
	}

	var data map[string]float64
	if err := json.Unmarshal(rec.Body.Bytes(), &data); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if data["temp_C"] != 30.5 {
		t.Fatalf("expected temp_C 30.5 got %v", data["temp_C"])
	}
	if data["temp_F"] != 30.5*1.8+32 {
		t.Fatalf("unexpected temp_F %v", data["temp_F"])
	}
	if data["temp_K"] != 30.5+273 {
		t.Fatalf("unexpected temp_K %v", data["temp_K"])
	}
}

func TestWeatherHandlerInvalidCEP(t *testing.T) {
	handler := weatherHandler(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/?cep=123", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 got %d", rec.Code)
	}
	if rec.Body.String() != ErrMsgInvalidCep {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestWeatherHandlerCEPNOTFOUND(t *testing.T) {
	viaCep := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"erro":true}`))
	}))
	defer viaCep.Close()

	handler := weatherHandler(
		&CepService{BaseURL: viaCep.URL, HTTPClient: http.DefaultClient},
		&WeatherService{BaseURL: "http://xpto", HTTPClient: http.DefaultClient},
	)

	req := httptest.NewRequest(http.MethodGet, "/?cep=12345678", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", rec.Code)
	}
	if rec.Body.String() != ErrMsgCepNotFound {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
