package middleware

import (
	"github.com/k8s-dashboard/internal/kubernetes"
	"github.com/k8s-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
)

type SessionMiddleware struct {
	clientManager *kubernetes.ClientManager
}

func NewSessionMiddleware(clientManager *kubernetes.ClientManager) *SessionMiddleware {
	return &SessionMiddleware{
		clientManager: clientManager,
	}
}

func (m *SessionMiddleware) Required() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.GetHeader("X-Session-ID")
		if sessionID == "" {
			response.Unauthorized(c, "Session ID is required")
			c.Abort()
			return
		}

		client, ok := m.clientManager.GetClient(sessionID)
		if !ok {
			response.Unauthorized(c, "Session expired or invalid")
			c.Abort()
			return
		}

		c.Set("client", client)
		c.Set("sessionID", sessionID)
		c.Next()
	}
}
