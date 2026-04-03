package handler

import (
	"context"

	"github.com/k8s-dashboard/internal/kubernetes"
	"github.com/k8s-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterHandler struct {
	clientManager *kubernetes.ClientManager
}

func NewClusterHandler(clientManager *kubernetes.ClientManager) *ClusterHandler {
	return &ClusterHandler{clientManager: clientManager}
}

type ClusterInfo struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	NodeCount       int    `json:"nodeCount"`
	NamespaceCount  int    `json:"namespaceCount"`
	PodCount        int    `json:"podCount"`
	DeploymentCount int    `json:"deploymentCount"`
	ServiceCount    int    `json:"serviceCount"`
}

type ClusterHealth struct {
	NodesHealthy    int `json:"nodesHealthy"`
	NodesNotReady   int `json:"nodesNotReady"`
	PodsRunning     int `json:"podsRunning"`
	PodsPending     int `json:"podsPending"`
	PodsFailed      int `json:"podsFailed"`
	PodsSucceeded   int `json:"podsSucceeded"`
	PodsUnknown     int `json:"podsUnknown"`
}

func (h *ClusterHandler) GetInfo(c *gin.Context) {
	client, _ := c.Get("client")
	clientSet := client.(*kubernetes.ClientSet)

	version, err := clientSet.Clientset.ServerVersion()
	if err != nil {
		response.InternalError(c, "Failed to get server version")
		return
	}

	nodes, err := clientSet.Clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		response.InternalError(c, "Failed to list nodes")
		return
	}

	namespaces, err := clientSet.Clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		response.InternalError(c, "Failed to list namespaces")
		return
	}

	pods, err := clientSet.Clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		response.InternalError(c, "Failed to list pods")
		return
	}

	deployments, err := clientSet.Clientset.AppsV1().Deployments("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		response.InternalError(c, "Failed to list deployments")
		return
	}

	services, err := clientSet.Clientset.CoreV1().Services("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		response.InternalError(c, "Failed to list services")
		return
	}

	info := ClusterInfo{
		Name:            version.GitVersion,
		Version:         version.String(),
		NodeCount:       len(nodes.Items),
		NamespaceCount:  len(namespaces.Items),
		PodCount:        len(pods.Items),
		DeploymentCount: len(deployments.Items),
		ServiceCount:    len(services.Items),
	}

	response.Success(c, info)
}

func (h *ClusterHandler) GetHealth(c *gin.Context) {
	client, _ := c.Get("client")
	clientSet := client.(*kubernetes.ClientSet)

	nodes, err := clientSet.Clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		response.InternalError(c, "Failed to list nodes")
		return
	}

	health := ClusterHealth{}
	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady {
				if condition.Status == corev1.ConditionTrue {
					health.NodesHealthy++
				} else {
					health.NodesNotReady++
				}
			}
		}
	}

	pods, err := clientSet.Clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		response.InternalError(c, "Failed to list pods")
		return
	}

	for _, pod := range pods.Items {
		switch pod.Status.Phase {
		case corev1.PodRunning:
			health.PodsRunning++
		case corev1.PodPending:
			health.PodsPending++
		case corev1.PodFailed:
			health.PodsFailed++
		case corev1.PodSucceeded:
			health.PodsSucceeded++
		default:
			health.PodsUnknown++
		}
	}

	response.Success(c, health)
}
