package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const baseURL = "http://192.168.0.50/Develop_EDT_Donskov/hs/aapi/"

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	originalURL := r.URL.String()

	// Изменяем URL
	newURLStr := fmt.Sprintf("%s%s", baseURL, r.URL.Path)
	newURL, err := url.Parse(newURLStr)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Создаем новый запрос
	req, err := http.NewRequest(r.Method, newURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Копируем заголовки
	for key, values := range r.Header {
		req.Header[key] = values
	}

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to forward request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	var bodyBytes []byte
	if resp.Body != nil {
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response body", http.StatusInternalServerError)
			return
		}
	}

	// Восстанавливаем оригинальный URL
	if location := resp.Header.Get("Location"); location != "" {
		parsedLocation, _ := url.Parse(location)
		if parsedLocation != nil && strings.HasPrefix(parsedLocation.Host, "example.com") {
			originalLocation := fmt.Sprintf("%s%s", originalURL, parsedLocation.Path)
			resp.Header.Set("Location", originalLocation)
		}
	}

	// Отправляем ответ клиенту
	w.WriteHeader(resp.StatusCode)
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	if len(bodyBytes) > 0 {
		w.Write(bodyBytes)
	}
}
