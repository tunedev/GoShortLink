package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"math/big"

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
	case "/":
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}
		h.Redirect(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.URL.Path[len("/"):]
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
	println("is it getting here =========>>>>>>>>>")
	if err != nil || reqUrl.LongUrl == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	reqUrl.ShortUrl = base62Encode(reqUrl.ID)

	err = h.store.UpdateUrl(r.Context(), &reqUrl)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reqUrl)
}

func base62Encode(uuid string) string {
	uuid = strings.ReplaceAll(uuid, "-", "")

	n := new(big.Int)
	n.SetString(uuid, 16)

	const chars = "0123456789abcdefgjijklmnoqprstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	s := ""
	for n.Cmp(big.NewInt(0)) > 0 {
		s = string(chars[int(new(big.Int).Mod(n, big.NewInt(62)).Int64())]) + s
		n.Div(n, big.NewInt(62))
	}
	return s
}
