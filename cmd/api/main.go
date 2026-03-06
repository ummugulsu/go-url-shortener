package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

type ShortenRequest struct {
	URL       string `json:"url"`
	Custom    string `json:"custom_code"`
	ExpiresIn int    `json:"expires_in"` // dakika
}

func generateCode(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}


func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.URL == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var code string

	// custom code varsa kullan
	if req.Custom != "" {
		code = req.Custom
	} else {
		code = generateCode(6)
	}

	// expire hesapla
	var expiresAt interface{}
	if req.ExpiresIn > 0 {
		expiresAt = time.Now().Add(time.Duration(req.ExpiresIn) * time.Minute)
	} else {
		expiresAt = nil
	}

	_, err = db.Exec(
		"INSERT INTO urls (original_url, short_code, expires_at) VALUES ($1, $2, $3)",
		req.URL, code, expiresAt,
	)

	if err != nil {
		http.Error(w, "Short code already exists or DB error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"short_url": "http://localhost:8080/" + code,
	}

	json.NewEncoder(w).Encode(response)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:]

	var original string
	var expiresAt sql.NullTime

	err := db.QueryRow(
		"SELECT original_url, expires_at FROM urls WHERE short_code=$1",
		code,
	).Scan(&original, &expiresAt)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	// expire kontrolü
	if expiresAt.Valid && time.Now().After(expiresAt.Time) {
		http.Error(w, "Link expired", http.StatusGone)
		return
	}

	db.Exec("UPDATE urls SET click_count = click_count + 1 WHERE short_code=$1", code)

	http.Redirect(w, r, original, http.StatusFound)
}
func statsHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/stats/"):]

	row := db.QueryRow(
		"SELECT original_url, click_count, created_at, expires_at FROM urls WHERE short_code=$1",
		code,
	)

	var original string
	var clicks int64
	var created time.Time
	var expires sql.NullTime

	err := row.Scan(&original, &clicks, &created, &expires)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	response := map[string]interface{}{
		"original_url": original,
		"click_count":  clicks,
		"created_at":   created,
		"expires_at":   expires,
	}

	json.NewEncoder(w).Encode(response)
}


func main() {
	rand.Seed(time.Now().UnixNano())

	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=urlshortener sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database bağlantısı başarısız:", err)
	}

	http.HandleFunc("/stats/", statsHandler)
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}