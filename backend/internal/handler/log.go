package handler

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/k8s-dashboard/internal/kubernetes"
	"github.com/k8s-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
)

type LogHandler struct {
	clientManager *kubernetes.ClientManager
}

func NewLogHandler(clientManager *kubernetes.ClientManager) *LogHandler {
	return &LogHandler{clientManager: clientManager}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *LogHandler) StreamLogs(c *gin.Context) {
	client, _ := c.Get("client")
	clientSet := client.(*kubernetes.ClientSet)

	namespace := c.Param("namespace")
	podName := c.Param("name")

	container := c.Query("container")
	follow := c.Query("follow") == "true"
	tailLines := int64(100)
	if tl := c.Query("tailLines"); tl != "" {
		if parsed, err := parseInt64(tl); err == nil {
			tailLines = parsed
		}
	}

	opts := &corev1.PodLogOptions{
		Container: container,
		Follow:    follow,
		TailLines: &tailLines,
	}

	req := clientSet.Clientset.CoreV1().Pods(namespace).GetLogs(podName, opts)
	stream, err := req.Stream(context.Background())
	if err != nil {
		response.InternalError(c, "Failed to get log stream: "+err.Error())
		return
	}
	defer stream.Close()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(stream)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}

		if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
			break
		}
	}
}

func (h *LogHandler) GetLogs(c *gin.Context) {
	client, _ := c.Get("client")
	clientSet := client.(*kubernetes.ClientSet)

	namespace := c.Param("namespace")
	podName := c.Param("name")

	container := c.Query("container")
	tailLines := int64(100)
	if tl := c.Query("tailLines"); tl != "" {
		if parsed, err := parseInt64(tl); err == nil {
			tailLines = parsed
		}
	}

	opts := &corev1.PodLogOptions{
		Container: container,
		Follow:    false,
		TailLines: &tailLines,
	}

	req := clientSet.Clientset.CoreV1().Pods(namespace).GetLogs(podName, opts)
	stream, err := req.Stream(context.Background())
	if err != nil {
		response.InternalError(c, "Failed to get logs: "+err.Error())
		return
	}
	defer stream.Close()

	logs, err := io.ReadAll(stream)
	if err != nil {
		response.InternalError(c, "Failed to read logs")
		return
	}

	response.Success(c, map[string]string{
		"logs": string(logs),
	})
}

func parseInt64(s string) (int64, error) {
	var result int64
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
