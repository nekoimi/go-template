package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	wslib "github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/nekoimi/go-project-template/internal/pkg/jwtutil"
)

type WSHandler struct {
	manager       *Manager
	jwtSecret     string
	logger        *zap.Logger
	allowedOrigins []string
}

func NewWSHandler(manager *Manager, jwtSecret string, logger *zap.Logger, allowedOrigins []string) *WSHandler {
	return &WSHandler{
		manager:       manager,
		jwtSecret:     jwtSecret,
		logger:        logger,
		allowedOrigins: allowedOrigins,
	}
}

// Upgrade handles WebSocket upgrade requests.
// GET /ws/v1/chat?token=<jwt>
func (h *WSHandler) Upgrade(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	userID, err := jwtutil.ValidateToken(tokenStr, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	upgrader := wslib.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return h.checkOrigin(r.Header.Get("Origin"))
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("websocket upgrade failed", zap.Error(err))
		return
	}

	client := newClient(h.manager, conn, userID, h.logger)
	h.manager.register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (h *WSHandler) checkOrigin(origin string) bool {
	// 如果未配置允许的来源，则允许所有（开发模式）
	if len(h.allowedOrigins) == 0 {
		return true
	}
	// 检查来源是否在允许列表中
	for _, allowed := range h.allowedOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}
