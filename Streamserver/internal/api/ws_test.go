package api_test

import (
	"Streamserver/internal/api"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestWSMessageRouting(t *testing.T) {
	// --- 1) Настраиваем in-memory БД и Gin с вашим WSHandler
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	router := gin.New()
	// монтируем WSHandler из вашего пакета api
	router.GET("/ws", api.WSHandler(db))

	// --- 2) Запускаем временный HTTP-сервер
	srv := httptest.NewServer(router)
	defer srv.Close()

	// Переключаем http:// -> ws:// для WebSocket
	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws?roomId=testRoom"

	// --- 3) Подключаемся как WebSocket-клиент
	conn, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err, "Dial should succeed")
	assert.Equal(t, 101, resp.StatusCode) // 101 Switching Protocols
	defer conn.Close()

	// --- 4) Проверяем эхо для разных types
	testCases := []struct {
		msgType string
	}{
		{"join"},
		{"offer"},
		{"answer"},
		{"candidate"},
		{"leave"},
		{"invalid"},
	}

	for _, tc := range testCases {
		// Формируем входящее сообщение
		in := map[string]interface{}{
			"type": tc.msgType,
			"data": map[string]string{"foo": "bar"},
		}
		raw, _ := json.Marshal(in)

		// Отправляем
		err := conn.WriteMessage(websocket.TextMessage, raw)
		assert.NoError(t, err, "WriteMessage should not error for type "+tc.msgType)

		// Получаем эхо
		_, outRaw, err := conn.ReadMessage()
		assert.NoError(t, err, "ReadMessage should not error for type "+tc.msgType)

		// Для «invalid» ожидаем объект {type:"error",message:"unknown message type"}
		if tc.msgType == "invalid" {
			var resp map[string]string
			err := json.Unmarshal(outRaw, &resp)
			assert.NoError(t, err)
			assert.Equal(t, "error", resp["type"])
			assert.Equal(t, "unknown message type", resp["message"])
		} else {
			// Для остальных — ровно то же, что отправили (эко)
			assert.JSONEq(t, string(raw), string(outRaw), "echo should match for "+tc.msgType)
		}
	}
}
