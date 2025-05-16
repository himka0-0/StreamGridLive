package api

import (
	"Streamserver/internal/sfu"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// upgrader конфигурирует параметры апгрейда соединения до WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Разрешаем подключения со всех origin, при необходимости уточните
		return true
	},
}

// WSMessage представляет базовую структуру входящих WS-сообщений
type WSMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// WSHandler возвращает gin.HandlerFunc, содержащую доступ к базе данных
func WSHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Апгрейд HTTP-соединения до WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}
		defer conn.Close()

		// Читаем обязательный параметр roomId из query
		roomID := c.Query("roomId")
		if roomID == "" {
			closeMsg := websocket.FormatCloseMessage(
				websocket.ClosePolicyViolation,
				"missing roomId",
			)
			conn.WriteMessage(websocket.CloseMessage, closeMsg)
			return
		}

		log.Printf("WebSocket connected: %s (room=%s)", conn.RemoteAddr(), roomID)
		handleWS(conn, roomID, db)
	}
}

// handleWS читает сообщения, маршрутизирует их по типу и выполняет эхо-ответ
func handleWS(conn *websocket.Conn, roomID string, db *gorm.DB) {
	var jc *sfu.JanusClient
	for {
		// Читаем сообщение от клиента
		msgType, raw, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error (%s): %v", roomID, err)
			break
		}

		// Парсим JSON для получения поля type
		var msg WSMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Printf("invalid JSON (%s): %v", roomID, err)
			sendError(conn, "invalid JSON")
			continue
		}

		switch msg.Type {
		case "join":
			// парсим опциональные флаги аудио/видео
			var p struct {
				Audio bool `json:"audio"`
				Video bool `json:"video"`
			}
			if err := json.Unmarshal(msg.Data, &p); err != nil {
				sendError(conn, "invalid join payload")
				break
			}

			// инициализируем SFU-клиент
			janusURL := os.Getenv("JANUS_URL") // например "http://localhost:8088"
			jc = sfu.NewJanusClient(janusURL)
			if err := jc.CreateSession(); err != nil {
				sendError(conn, "SFU create session failed")
				break
			}
			if err := jc.Attach("janus.plugin.echotest"); err != nil {
				sendError(conn, "SFU attach failed")
				break
			}
			offer, err := jc.SetupMedia(p.Audio, p.Video)
			if err != nil {
				sendError(conn, "SFU setup failed")
				break
			}

			// шлём клиенту SDP-offer
			resp := map[string]interface{}{
				"type": "offer",
				"data": offer,
			}
			b, _ := json.Marshal(resp)
			conn.WriteMessage(websocket.TextMessage, b)

		case "offer":
			handleOffer(conn, roomID, db, msg.Data)
			conn.WriteMessage(msgType, raw)

		case "answer":
			if jc == nil {
				sendError(conn, "session not initialized")
				continue
			}
			var ans sfu.JSEP
			if err := json.Unmarshal(msg.Data, &ans); err != nil {
				sendError(conn, "invalid answer payload")
				continue
			}
			if err := jc.SendAnswer(&ans); err != nil {
				sendError(conn, "SFU answer failed")
				continue
			}
			// если нужно эхо — раскомментируй
			conn.WriteMessage(msgType, raw)

		case "candidate":
			if jc == nil {
				sendError(conn, "session not initialized")
				continue
			}
			var cand sfu.ICECandidate
			if err := json.Unmarshal(msg.Data, &cand); err != nil {
				sendError(conn, "invalid candidate payload")
				continue
			}
			if err := jc.Tricklet(&cand); err != nil {
				sendError(conn, "SFU trickle failed")
				continue
			}
			// если нужно эхо — раскомментируй
			conn.WriteMessage(msgType, raw)

		case "leave":
			handleLeave(conn, roomID, db, msg.Data)
			conn.WriteMessage(msgType, raw)

		default:
			log.Printf("unknown message type (%s): %s", roomID, msg.Type)
			sendError(conn, "unknown message type")
		}
	}
	log.Printf("WebSocket closed: %s", roomID)
}

// sendError шлёт клиенту сообщение типа error с текстом ошибки
func sendError(conn *websocket.Conn, errorMsg string) {
	resp := map[string]string{"type": "error", "message": errorMsg}
	if b, err := json.Marshal(resp); err == nil {
		conn.WriteMessage(websocket.TextMessage, b)
	}
}

// Ниже — заготовки обработчиков различных типов WS-сообщений
func handleJoin(conn *websocket.Conn, roomID string, db *gorm.DB, payload json.RawMessage) {
	log.Printf("JOIN (%s): %s", roomID, string(payload))
	// TODO: реализовать логику подключения нового участника
}

func handleOffer(conn *websocket.Conn, roomID string, db *gorm.DB, payload json.RawMessage) {
	log.Printf("OFFER (%s): %s", roomID, string(payload))
	// TODO: реализовать пересылку SDP-offer к SFU
}

func handleAnswer(conn *websocket.Conn, roomID string, db *gorm.DB, payload json.RawMessage) {
	log.Printf("ANSWER (%s): %s", roomID, string(payload))
	// TODO: реализовать обработку SDP-answer от клиента
}

func handleCandidate(conn *websocket.Conn, roomID string, db *gorm.DB, payload json.RawMessage) {
	log.Printf("CANDIDATE (%s): %s", roomID, string(payload))
	// TODO: реализовать обмен ICE-кандидатами
}

func handleLeave(conn *websocket.Conn, roomID string, db *gorm.DB, payload json.RawMessage) {
	log.Printf("LEAVE (%s): %s", roomID, string(payload))
	// TODO: реализовать логику выхода участника из комнаты
}
