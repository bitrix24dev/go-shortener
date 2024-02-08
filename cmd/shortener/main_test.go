package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ищем ключ по значению в MAP
func getKeyByValue(m map[string]string, value string) (string, bool) {
	for key, val := range m {
		if val == value {
			return key, true
		}
	}
	return "", false
}

func TestHandleURLShortening(t *testing.T) {
	// Создаем фейковый запрос с телом
	url := "https://practicum.yandex.ru/"
	reqBody := []byte(url)
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	// Создаем фейковый ResponseWriter
	w := httptest.NewRecorder()

	// Вызываем функцию-обработчик
	handleURLShortening(w, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusCreated, w.Code)

	// Получаем тело ответа
	respBody, err := io.ReadAll(w.Body)
	require.NoError(t, err)

	// Проверяем тело ответа
	expectedBody, ok := getKeyByValue(urlMap, url)
	assert.Equal(t, ok, true)
	assert.Equal(t, expectedBody, string(respBody))
}

func TestHandleRedirect(t *testing.T) {
	// Генерируем случайный ключ для карты urlMap
	key := generateRandomString(8)

	// Добавляем запись в карту urlMap
	urlMap["http://localhost:8080/"+key] = "https://practicum.yandex.ru/"

	// Создаем фейковый запрос с переменной в URL
	req, err := http.NewRequest("GET", "/"+key, nil)
	require.NoError(t, err)

	// Создаем фейковый ResponseWriter
	w := httptest.NewRecorder()

	// Вызываем функцию-обработчик
	handleRedirect(w, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)

	// Проверяем заголовок Location
	expectedLocation := "https://practicum.yandex.ru/"
	assert.Equal(t, expectedLocation, w.Header().Get("Location"))
}
