# API de Temperatura por CEP (Go) — Cloud Run

Sistema em Go que recebe um **CEP de 8 dígitos**, identifica a localização via **ViaCEP**, consulta a temperatura via **WeatherAPI** e retorna as temperaturas em **Celsius, Fahrenheit e Kelvin**.

Este repositório atende o desafio: **deploy no Google Cloud Run + testes automatizados + Docker/docker-compose**.

---

## ✅ Requisitos atendidos

- Recebe CEP válido (8 dígitos)
- Busca cidade/localização via ViaCEP
- Consulta temperatura via WeatherAPI
- Converte e retorna:
  - `temp_C`
  - `temp_F` (F = C * 1.8 + 32)
  - `temp_K` (K = C + 273)
- Cenários:
  - `200` → JSON com temperaturas
  - `422` → mensagem `invalid zipcode`
  - `404` → mensagem `can not find zipcode`
- Testes automatizados (`go test ./...`)
- Docker + docker-compose
- Deploy no Cloud Run

---

## 🔧 Variáveis de ambiente

Você precisa de uma chave da WeatherAPI:

- `WEATHERAPI_KEY` (obrigatória)

### Recomendações (pra não vazar segredo no Git)
Crie um arquivo `.env` local :

```env
WEATHERAPI_KEY=coloque_sua_chave_aqui
```

---

## ▶️ Como rodar

### 1) Local (Go)
```bash
go run main.go
```

### 2) Docker Compose
```bash
docker compose up --build
```

---

## 🧪 Testes automatizados

```bash
go test ./...
```

---

## ☁️ Serviço em Produção (Google Cloud Run)

A aplicação está publicada e acessível publicamente em:

https://fullcycle-desafio-738354502644.us-central1.run.app/

## 📡 API

### Endpoint
```
GET /?cep={CEP}
```

Exemplo:
```
GET /?cep=29902555
```

### Respostas

#### ✅ Sucesso
**HTTP 200**
```json
{"temp_C":28.5,"temp_F":83.3,"temp_K":301.5}
```

#### ❌ CEP inválido (formato)
**HTTP 422**
```
invalid zipcode
```

#### ❌ CEP não encontrado
**HTTP 404**
```
can not find zipcode
```

---

## ☁️ Deploy no Google Cloud Run

### Pré-requisitos
- `gcloud` instalado e autenticado
- Projeto GCP configurado (`gcloud config set project <PROJECT_ID>`)
- Billing habilitado (free tier costuma cobrir testes leves)

### Deploy (via source)
Substitua os valores:

```bash
gcloud run deploy cep-weather   --source .   --region southamerica-east1   --allow-unauthenticated   --set-env-vars WEATHERAPI_KEY=SEU_TOKEN_AQUI
```

Ao final, o Cloud Run vai te retornar a URL pública do serviço.

### Teste no Cloud Run
```bash
curl "https://fullcycle-desafio-738354502644.us-central1.run.app/?cep=20040001"
```

---

## 🧭 Dicas de avaliação

- Garanta que o serviço responda exatamente:
  - `422` com texto `invalid zipcode`
  - `404` com texto `can not find zipcode`
- Não comite chave de API no repositório (use env vars).

---

## 📦 Tecnologias

- Go
- ViaCEP
- WeatherAPI
- Docker / Docker Compose
- Google Cloud Run
