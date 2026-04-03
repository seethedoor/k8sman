package kubernetes

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type ClientSet struct {
	Clientset     *kubernetes.Clientset
	DynamicClient dynamic.Interface
	RestConfig    *rest.Config
	CreatedAt     time.Time
}

type ClientManager struct {
	clients sync.Map
}

func NewClientManager() *ClientManager {
	return &ClientManager{}
}

func (m *ClientManager) CreateClient(kubeconfig []byte, contextName string) (string, *ClientSet, error) {
	config, err := clientcmd.Load(kubeconfig)
	if err != nil {
		return "", nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	if contextName != "" {
		config.CurrentContext = contextName
	}

	restConfig, err := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{
		CurrentContext: config.CurrentContext,
	}).ClientConfig()
	if err != nil {
		return "", nil, fmt.Errorf("failed to create rest config: %w", err)
	}

	restConfig.TLSClientConfig.Insecure = true
	restConfig.TLSClientConfig.NextProtos = []string{"http/1.1"}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	restConfig.Transport = transport

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	_, err = clientset.ServerVersion()
	if err != nil {
		return "", nil, fmt.Errorf("failed to connect to cluster: %w", err)
	}

	sessionID := uuid.New().String()
	clientSet := &ClientSet{
		Clientset:     clientset,
		DynamicClient: dynamicClient,
		RestConfig:    restConfig,
		CreatedAt:     time.Now(),
	}

	m.clients.Store(sessionID, clientSet)

	return sessionID, clientSet, nil
}

func (m *ClientManager) GetClient(sessionID string) (*ClientSet, bool) {
	value, ok := m.clients.Load(sessionID)
	if !ok {
		return nil, false
	}
	return value.(*ClientSet), true
}

func (m *ClientManager) RemoveClient(sessionID string) {
	m.clients.Delete(sessionID)
}

func (m *ClientManager) CreateClientFromToken(serverURL, token, caCert []byte) (string, *ClientSet, error) {
	config := &rest.Config{
		Host:        string(serverURL),
		BearerToken: string(token),
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}

	if len(caCert) > 0 {
		config.TLSClientConfig.CAData = caCert
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	_, err = clientset.ServerVersion()
	if err != nil {
		return "", nil, fmt.Errorf("failed to connect to cluster: %w", err)
	}

	sessionID := uuid.New().String()
	clientSet := &ClientSet{
		Clientset:     clientset,
		DynamicClient: dynamicClient,
		RestConfig:    config,
		CreatedAt:     time.Now(),
	}

	m.clients.Store(sessionID, clientSet)

	return sessionID, clientSet, nil
}

func (m *ClientManager) CreateClientFromServiceAccount(serverURL string, sa *clientcmdapi.Config) (string, *ClientSet, error) {
	restConfig, err := clientcmd.NewDefaultClientConfig(*sa, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return "", nil, fmt.Errorf("failed to create rest config: %w", err)
	}

	if serverURL != "" {
		restConfig.Host = serverURL
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = clientset.ServerVersion()
	if err != nil {
		return "", nil, fmt.Errorf("failed to connect to cluster: %w", err)
	}
	_ = ctx

	sessionID := uuid.New().String()
	clientSet := &ClientSet{
		Clientset:     clientset,
		DynamicClient: dynamicClient,
		RestConfig:    restConfig,
		CreatedAt:     time.Now(),
	}

	m.clients.Store(sessionID, clientSet)

	return sessionID, clientSet, nil
}
