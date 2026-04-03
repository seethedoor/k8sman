package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/k8s-dashboard/internal/kubernetes"
	"github.com/k8s-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/remotecommand"
)

type ExecHandler struct {
	clientManager *kubernetes.ClientManager
}

func NewExecHandler(clientManager *kubernetes.ClientManager) *ExecHandler {
	return &ExecHandler{clientManager: clientManager}
}

type terminalSizeQueue struct {
	resizeChan chan remotecommand.TerminalSize
}

func (t *terminalSizeQueue) Next() *remotecommand.TerminalSize {
	size, ok := <-t.resizeChan
	if !ok {
		return nil
	}
	return &size
}

func (h *ExecHandler) Exec(c *gin.Context) {
	client, _ := c.Get("client")
	clientSet := client.(*kubernetes.ClientSet)

	namespace := c.Param("namespace")
	podName := c.Param("name")

	container := c.Query("container")
	if container == "" {
		pod, err := clientSet.Clientset.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
		if err != nil {
			response.NotFound(c, "Pod not found")
			return
		}
		if len(pod.Spec.Containers) > 0 {
			container = pod.Spec.Containers[0].Name
		}
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	req := clientSet.Clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   []string{"/bin/sh", "-c", "TERM=xterm-256color; export TERM; [ -e /bin/bash ] && /bin/bash || /bin/sh"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, metav1.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(clientSet.RestConfig, "POST", req.URL())
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Error creating executor: "+err.Error()))
		return
	}

	sizeQueue := &terminalSizeQueue{
		resizeChan: make(chan remotecommand.TerminalSize, 10),
	}

	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()
	stderrReader, stderrWriter := io.Pipe()

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		buf := make([]byte, 1024)
		for {
			n, err := stdoutReader.Read(buf)
			if err != nil {
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		buf := make([]byte, 1024)
		for {
			n, err := stderrReader.Read(buf)
			if err != nil {
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		defer stdinWriter.Close()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var data map[string]interface{}
			if err := json.Unmarshal(msg, &data); err == nil {
				if resize, ok := data["resize"].(map[string]interface{}); ok {
					cols, _ := resize["cols"].(float64)
					rows, _ := resize["rows"].(float64)
					sizeQueue.resizeChan <- remotecommand.TerminalSize{
						Width:  uint16(cols),
						Height: uint16(rows),
					}
					continue
				}
			}

			stdinWriter.Write(msg)
		}
	}()

	streamOptions := remotecommand.StreamOptions{
		Stdin:             stdinReader,
		Stdout:            stdoutWriter,
		Stderr:            stderrWriter,
		Tty:               true,
		TerminalSizeQueue: sizeQueue,
	}

	err = executor.StreamWithContext(context.Background(), streamOptions)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Session ended: "+err.Error()))
	}

	stdoutWriter.Close()
	stderrWriter.Close()
	close(sizeQueue.resizeChan)
	wg.Wait()
}
