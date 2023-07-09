package handler

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tunedev/GoShortLink/pkg/model"
	"github.com/tunedev/GoShortLink/pkg/store"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	store *store.Store
}

func NewHandler(s *store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/shorten":
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}
		h.Shorten(w, r)
	default:
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}
		h.Redirect(w, r)
	}
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.URL.Path[len("/"):]
	fmt.Printf("db query for url :=== %v\n", shortUrl)
	if shortUrl == "" {
		if r.Method != http.MethodPost {
			http.Error(w, "Short url required", http.StatusBadRequest)
			return
		}
	}
	url, err := h.store.GetUrlByShortURL(r.Context(), shortUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.LongUrl, http.StatusMovedPermanently)
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var reqUrl model.Url
  err := json.NewDecoder(r.Body).Decode(&reqUrl)
	if err != nil || reqUrl.LongUrl == "" {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Validate URL
	_, err = url.ParseRequestURI(reqUrl.LongUrl)
	if err != nil {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	
	existingUrl, err := h.store.GetUrlByLongURL(r.Context(), reqUrl.LongUrl)
	
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(existingUrl)
		return
	} else if err != mongo.ErrNoDocuments {
			http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
			return
	}
		
	reqUrl.CreatedAt = time.Now()
	err = h.store.SaveUrl(r.Context(), &reqUrl)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	reqUrl.ShortUrl = base62Encode(strings.TrimPrefix(reqUrl.ID, "id_"))
	fmt.Printf("after base64Conversion : %#v\n",reqUrl)


	err = h.store.UpdateUrl(r.Context(), &reqUrl)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	urls, err := h.store.GetAllUrls(r.Context())
	if err != nil {
		fmt.Printf("Error occured: %v",err)
		http.Error(w, "Failed to fetch all short URLs", http.StatusInternalServerError)
		return
	}
	urlsVal := make([]model.Url, len(urls))
	for i, v := range urls {
		urlsVal[i] = *v
	}
	fmt.Printf("Entries in the DB: %#v", urlsVal)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reqUrl)
}

func base62Encode(str string) string {
	id := new(big.Int)
	id.SetString(str, 10)
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if id.Cmp(big.NewInt(0)) == 0 {
		return string(chars[0])
	}

	s := ""
	n := new(big.Int).Set(id) // make a copy of id
	zero := big.NewInt(0)
	mod := new(big.Int)
	base := big.NewInt(62)

	for n.Cmp(zero) > 0 {
		n.DivMod(n, base, mod)
		s = string(chars[mod.Int64()]) + s
	}

	return s
}
