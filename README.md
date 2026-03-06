# go-url-shortener

# Go URL Shortener

A simple and scalable URL shortener API built with **Go** and **PostgreSQL**.  
Inspired by **Bitly**.

---

## Features

- Shorten long URLs with custom codes
- Redirect short links to original URLs
- Track click counts
- Support link expiration

---

## Tech Stack

- **Backend:** Go — fast, concurrent, production-ready  
- **Database:** PostgreSQL — strong indexing and production-ready  
- **Cache (optional):** Redis — reduces DB load and latency  
- **Containerization (optional):** Docker  
- **Orchestration (optional):** Kubernetes for scaling  

---

## API Endpoints

- **POST /shorten** — Create a short URL  

```json
{
  "url": "https://example.com",
  "custom_code": "myCode",
  "expires_in": 60
}

GET /:short_code — Redirect to original URL

GET /stats/:short_code — Get click count and link info

Run Locally

Setup PostgreSQL database:

CREATE DATABASE urlshortener;
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    click_count BIGINT DEFAULT 0,
    expires_at TIMESTAMP
);

Run API:

go run cmd/main.go

Test API with Postman or curl.

Author

Ummugulsum Zengin
