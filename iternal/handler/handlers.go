package handler

import (
	"bytes"
	"counter-api/iternal/counter"
	"encoding/json"
	"net/http"
)

type CounterHandler struct {
	Counter counter.CounterIntarface
}

func NewCounterHandler(counter counter.CounterIntarface) *CounterHandler {
	return &CounterHandler{
		Counter: counter,
	}
}

func (h *CounterHandler) GetCount(w http.ResponseWriter, r *http.Request) {
	value, err := h.Counter.GetCount(r.Context())
	if err != nil {
		status := http.StatusServiceUnavailable
		http.Error(w, http.StatusText(status), status)
		return
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(value)
	if err != nil {
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		return
	}
}
func (h *CounterHandler) IncrCount(w http.ResponseWriter, r *http.Request) {
	value, err := h.Counter.IncrCount(r.Context())
	if err != nil {
		status := http.StatusServiceUnavailable
		http.Error(w, http.StatusText(status), status)
		return
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(value)
	if err != nil {
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		return
	}
}
func (h *CounterHandler) Health(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(map[string]string{"status": "ok"})
	if err != nil {
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		return
	}
}
