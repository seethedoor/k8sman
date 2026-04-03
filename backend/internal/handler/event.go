package handler

import (
	"context"

	"github.com/k8s-dashboard/internal/kubernetes"
	"github.com/k8s-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EventHandler struct {
	clientManager *kubernetes.ClientManager
}

func NewEventHandler(clientManager *kubernetes.ClientManager) *EventHandler {
	return &EventHandler{clientManager: clientManager}
}

type EventInfo struct {
	Type           string `json:"type"`
	Reason         string `json:"reason"`
	Message        string `json:"message"`
	ObjectKind     string `json:"objectKind"`
	ObjectName     string `json:"objectName"`
	Namespace      string `json:"namespace"`
	Count          int32  `json:"count"`
	FirstTimestamp string `json:"firstTimestamp"`
	LastTimestamp  string `json:"lastTimestamp"`
}

func (h *EventHandler) ListEvents(c *gin.Context) {
	client, _ := c.Get("client")
	clientSet := client.(*kubernetes.ClientSet)

	namespace := c.Query("namespace")
	eventType := c.Query("type")

	var events *corev1.EventList
	var err error

	if namespace != "" {
		events, err = clientSet.Clientset.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{})
	} else {
		events, err = clientSet.Clientset.CoreV1().Events("").List(context.Background(), metav1.ListOptions{})
	}

	if err != nil {
		response.InternalError(c, "Failed to list events: "+err.Error())
		return
	}

	var eventInfos []EventInfo
	for _, event := range events.Items {
		if eventType != "" && event.Type != eventType {
			continue
		}

		eventInfos = append(eventInfos, EventInfo{
			Type:           event.Type,
			Reason:         event.Reason,
			Message:        event.Message,
			ObjectKind:     event.InvolvedObject.Kind,
			ObjectName:     event.InvolvedObject.Name,
			Namespace:      event.InvolvedObject.Namespace,
			Count:          event.Count,
			FirstTimestamp: event.FirstTimestamp.Format("2006-01-02 15:04:05"),
			LastTimestamp:  event.LastTimestamp.Format("2006-01-02 15:04:05"),
		})
	}

	response.Success(c, eventInfos)
}
