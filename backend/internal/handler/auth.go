package handler

import (
	"github.com/k8s-dashboard/internal/kubernetes"
	"github.com/k8s-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	clientManager *kubernetes.ClientManager
}

func NewAuthHandler(clientManager *kubernetes.ClientManager) *AuthHandler {
	return &AuthHandler{clientManager: clientManager}
}

type ConnectRequest struct {
	Kubeconfig  string `json:"kubeconfig"`
	ContextName string `json:"contextName"`
}

type ConnectResponse struct {
	SessionID string `json:"sessionId"`
}

func (h *AuthHandler) Connect(c *gin.Context) {
	var req ConnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if req.Kubeconfig == "" {
		response.BadRequest(c, "Kubeconfig is required")
		return
	}

	sessionID, _, err := h.clientManager.CreateClient([]byte(req.Kubeconfig), req.ContextName)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, ConnectResponse{SessionID: sessionID})
}

func (h *AuthHandler) Disconnect(c *gin.Context) {
	sessionID, exists := c.Get("sessionID")
	if !exists {
		response.Unauthorized(c, "Session not found")
		return
	}

	h.clientManager.RemoveClient(sessionID.(string))
	response.Success(c, nil)
}
