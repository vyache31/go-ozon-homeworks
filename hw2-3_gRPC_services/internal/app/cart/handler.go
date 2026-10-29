package cart

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/vyache31/go-ozon-homeworks/homework_2/internal/domain"
)

type CartHandler struct {
	service *CartService
}

func NewCartHandler(service *CartService) *CartHandler {
	return &CartHandler{service: service}
}

func (ch *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("user_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": invalid user id}`)
		return
	}

	resp, err := ch.service.GetProduct(r.Context(), domain.GetProductRequest{UserID: int64(userID)})
	if errors.Is(err, domain.ErrorCartNoExistsOrEmpty) {
		ch.writeError(w, http.StatusNotFound, err.Error())
		return
	} else if errors.Is(err, domain.ErrorUserInvalidID) {
		ch.writeError(w, http.StatusBadRequest, err.Error())
		return
	} else if err != nil {
		ch.writeError(w, http.StatusInternalServerError, "неизвестная ошибка")
		log.Printf("ERROR: CartHandler: %v", err)
		return
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		ch.writeError(w, http.StatusInternalServerError, "неизвестная ошибка")
		log.Printf("ERROR: CartHandler: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func (ch *CartHandler) PostProduct(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("user_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": invalid user id}`)
		return
	}
	sku, err := strconv.Atoi(r.PathValue("sku_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": invalid sku}`)
		return
	}

	var body struct {
		Count uint32 `json:"count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		ch.writeError(w, http.StatusInternalServerError, "неизвестная ошибка")
		log.Printf("ERROR: CartHandler: %v", err)
		return
	}

	count := body.Count

	postReq := domain.PostProductRequest{
		UserID:     int64(userID),
		ProductSKU: int64(sku),
		Count:      count,
	}

	err = ch.service.PostProduct(r.Context(), postReq)
	if errors.Is(err, domain.ErrorUserInvalidID) {
		ch.writeError(w, http.StatusBadRequest, err.Error())
		return
	} else if errors.Is(err, domain.ErrorCartInvalidSKU) {
		ch.writeError(w, http.StatusBadRequest, err.Error())
		return
	} else if errors.Is(err, domain.ErrorPositionInvalidCount) {
		ch.writeError(w, http.StatusBadRequest, err.Error())
		return
	} else if errors.Is(err, domain.ErrorProductNotFound) {
		ch.writeError(w, http.StatusPreconditionFailed, err.Error())
		return
	} else if err != nil {
		ch.writeError(w, http.StatusInternalServerError, "неизвестная ошибка")
		log.Printf("ERROR: CartHandler: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ch *CartHandler) writeError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}
