package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/scythe504/tiny-rl/internal"
	"github.com/scythe504/tiny-rl/internal/database"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	// Apply CORS middleware
	r.Use(s.corsMiddleware)

	r.HandleFunc("/", s.HelloWorldHandler)

	r.HandleFunc("/health", s.healthHandler)

	r.HandleFunc("/api/shorten", s.shortenURL)

	r.HandleFunc("/{shortCode}", s.getFullUrl)

	return r
}

// CORS middleware
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS Headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Wildcard allows all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Credentials not allowed with wildcard origins

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) shortenURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	// TODO: LOAD IN FROM ENV VAR
	const backendUrl = "http://localhost:8080"

	if err != nil {
		log.Println("[ShortenURL] error while reading body: ", err)
		http.Error(w, "error in reading the request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var link struct {
		URL string `json:"url"`
	}

	if err = json.Unmarshal(body, &link); err != nil {
		log.Println("[ShortenURL] error while unmarshaling json ", err)
		http.Error(w, "error in parsing request body", http.StatusBadRequest)
		return
	}

	if !internal.ValidURL(link.URL) {
		http.Error(w, "invalid url format, we do not support these format", http.StatusBadRequest)
		return
	}

	count := 0
	shortCode := internal.ShortCode()
	link_map := database.LinkMap{
		Id:  shortCode,
		Url: link.URL,
	}
	for {
		if count > 5 {
			http.Error(w, "error while generating short code", http.StatusInternalServerError)
			return
		}
		err = s.db.InsertShortenedLink(link_map)
		if err != nil {
			var pgErr *pgconn.PgError

			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				link_map.Id = internal.ShortCode()
				log.Println("[ShortenURL] duplicate key, retrying with new code")
				count++
				continue
			}

			// It's some other error - don't retry
			log.Println("[ShortenURL] database error: ", err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		} else {
			break
		}
	}

	var resp = struct {
		Data string `json:"data"`
	}{
		Data: fmt.Sprintf("%s/%s", backendUrl, shortCode),
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println("[ShortenURL] error while marshaling resp into json ", err)
		http.Error(w, "error while sending shortened url", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func (s *Server) getFullUrl(w http.ResponseWriter, r *http.Request) {
	shCode := mux.Vars(r)["shortCode"]

	linkMap, err := s.db.GetLink(shCode)
	if err != nil {
		switch err.Error() {
		case pgx.ErrNoRows.Error():
			http.Error(w, "short url is invalid", http.StatusNotFound)
			return
		default:
			log.Println("[GetFullUrl] error occured while getting link", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	var resp map[string]any = make(map[string]any)

	resp["data"] = struct {
		Id  string `json:"id"`
		Url string `json:"url"`
	}{
		Id:  linkMap.Id,
		Url: linkMap.Url,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println("[GetFullUrl] error while marshaling resp into json ", err)
		http.Error(w, "Failed to redirect to url", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
