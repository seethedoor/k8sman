package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/k8s-dashboard/internal/kubernetes"
	"github.com/k8s-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"
)

type ResourceHandler struct {
	clientManager *kubernetes.ClientManager
}

func NewResourceHandler(clientManager *kubernetes.ClientManager) *ResourceHandler {
	return &ResourceHandler{clientManager: clientManager}
}

var resourceGVRMap = map[string]schema.GroupVersionResource{
	"pods":                   {Group: "", Version: "v1", Resource: "pods"},
	"deployments":            {Group: "apps", Version: "v1", Resource: "deployments"},
	"statefulsets":           {Group: "apps", Version: "v1", Resource: "statefulsets"},
	"daemonsets":             {Group: "apps", Version: "v1", Resource: "daemonsets"},
	"replicasets":            {Group: "apps", Version: "v1", Resource: "replicasets"},
	"jobs":                   {Group: "batch", Version: "v1", Resource: "jobs"},
	"cronjobs":               {Group: "batch", Version: "v1beta1", Resource: "cronjobs"},
	"services":               {Group: "", Version: "v1", Resource: "services"},
	"ingresses":              {Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"},
	"endpoints":              {Group: "", Version: "v1", Resource: "endpoints"},
	"networkpolicies":        {Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"},
	"persistentvolumes":      {Group: "", Version: "v1", Resource: "persistentvolumes"},
	"persistentvolumeclaims": {Group: "", Version: "v1", Resource: "persistentvolumeclaims"},
	"storageclasses":         {Group: "storage.k8s.io", Version: "v1", Resource: "storageclasses"},
	"configmaps":             {Group: "", Version: "v1", Resource: "configmaps"},
	"secrets":                {Group: "", Version: "v1", Resource: "secrets"},
	"resourcequotas":         {Group: "", Version: "v1", Resource: "resourcequotas"},
	"limitranges":            {Group: "", Version: "v1", Resource: "limitranges"},
	"serviceaccounts":        {Group: "", Version: "v1", Resource: "serviceaccounts"},
	"roles":                  {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"},
	"rolebindings":           {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings"},
	"clusterroles":           {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"},
	"clusterrolebindings":    {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterrolebindings"},
	"nodes":                  {Group: "", Version: "v1", Resource: "nodes"},
	"namespaces":             {Group: "", Version: "v1", Resource: "namespaces"},
	"crds":                   {Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"},
}

func (h *ResourceHandler) getGVR(kind string) (schema.GroupVersionResource, bool) {
	kind = strings.ToLower(kind)
	gvr, ok := resourceGVRMap[kind]
	return gvr, ok
}

func (h *ResourceHandler) ListNamespaces(c *gin.Context) {
	client, _ := c.Get("client")
	clientSet := client.(*kubernetes.ClientSet)

	namespaces, err := clientSet.Clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		response.InternalError(c, "Failed to list namespaces")
		return
	}

	response.Success(c, namespaces.Items)
}

func (h *ResourceHandler) ListResources(c *gin.Context) {
	client, _ := c.Get("client")
	clientSet := client.(*kubernetes.ClientSet)

	kind := c.Param("kind")
	gvr, ok := h.getGVR(kind)
	if !ok {
		response.NotFound(c, "Resource kind not supported: "+kind)
		return
	}

	namespace := c.Query("namespace")
	labelSelector := c.Query("labelSelector")

	var list *unstructured.UnstructuredList
	var err error

	if namespace == "" && (kind == "nodes" || kind == "namespaces" || kind == "clusterroles" || kind == "clusterrolebindings" || kind == "persistentvolumes" || kind == "storageclasses" || kind == "crds") {
		list, err = clientSet.DynamicClient.Resource(gvr).List(context.Background(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
	} else if namespace == "" {
		list, err = clientSet.DynamicClient.Resource(gvr).List(context.Background(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
	} else {
		list, err = clientSet.DynamicClient.Resource(gvr).Namespace(namespace).List(context.Background(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
	}

	if err != nil {
		response.InternalError(c, "Failed to list resources: "+err.Error())
		return
	}

	response.Success(c, list.Items)
}

func (h *ResourceHandler) GetResource(c *gin.Context) {
	client, _ := c.Get("client")
	clientSet := client.(*kubernetes.ClientSet)

	kind := c.Param("kind")
	namespace := c.Param("namespace")
	name := c.Param("name")

	gvr, ok := h.getGVR(kind)
	if !ok {
		response.NotFound(c, "Resource kind not supported: "+kind)
		return
	}

	var resource *unstructured.Unstructured
	var err error

	if kind == "nodes" || kind == "namespaces" || kind == "clusterroles" || kind == "clusterrolebindings" || kind == "persistentvolumes" || kind == "storageclasses" || kind == "crds" {
		resource, err = clientSet.DynamicClient.Resource(gvr).Get(context.Background(), name, metav1.GetOptions{})
	} else {
		resource, err = clientSet.DynamicClient.Resource(gvr).Namespace(namespace).Get(context.Background(), name, metav1.GetOptions{})
	}

	if err != nil {
		response.NotFound(c, "Resource not found: "+err.Error())
		return
	}

	response.Success(c, resource)
}

func (h *ResourceHandler) GetResourceYAML(c *gin.Context) {
	client, _ := c.Get("client")
	clientSet := client.(*kubernetes.ClientSet)

	kind := c.Param("kind")
	namespace := c.Param("namespace")
	name := c.Param("name")

	gvr, ok := h.getGVR(kind)
	if !ok {
		response.NotFound(c, "Resource kind not supported: "+kind)
		return
	}

	var resource *unstructured.Unstructured
	var err error

	if kind == "nodes" || kind == "namespaces" || kind == "clusterroles" || kind == "clusterrolebindings" || kind == "persistentvolumes" || kind == "storageclasses" || kind == "crds" {
		resource, err = clientSet.DynamicClient.Resource(gvr).Get(context.Background(), name, metav1.GetOptions{})
	} else {
		resource, err = clientSet.DynamicClient.Resource(gvr).Namespace(namespace).Get(context.Background(), name, metav1.GetOptions{})
	}

	if err != nil {
		response.NotFound(c, "Resource not found: "+err.Error())
		return
	}

	jsonBytes, err := json.Marshal(resource)
	if err != nil {
		response.InternalError(c, "Failed to marshal resource")
		return
	}

	yamlBytes, err := yaml.JSONToYAML(jsonBytes)
	if err != nil {
		response.InternalError(c, "Failed to convert to YAML")
		return
	}

	response.Success(c, map[string]string{
		"yaml": string(yamlBytes),
	})
}

func (h *ResourceHandler) UpdateResource(c *gin.Context) {
	client, _ := c.Get("client")
	clientSet := client.(*kubernetes.ClientSet)

	kind := c.Param("kind")
	namespace := c.Param("namespace")
	name := c.Param("name")

	gvr, ok := h.getGVR(kind)
	if !ok {
		response.NotFound(c, "Resource kind not supported: "+kind)
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	yamlBytes, err := yaml.Marshal(req)
	if err != nil {
		response.BadRequest(c, "Failed to parse YAML: "+err.Error())
		return
	}

	decoder := serializer.NewCodecFactory(scheme.Scheme).UniversalDeserializer()
	obj, _, err := decoder.Decode(yamlBytes, nil, nil)
	if err != nil {
		response.BadRequest(c, "Failed to decode resource: "+err.Error())
		return
	}

	unstructuredObj, err := toUnstructured(obj)
	if err != nil {
		response.BadRequest(c, "Failed to convert resource: "+err.Error())
		return
	}

	metadata, ok := unstructuredObj["metadata"].(map[string]interface{})
	if !ok {
		response.BadRequest(c, "Invalid resource metadata")
		return
	}
	metadata["name"] = name
	metadata["namespace"] = namespace
	unstructuredObj["metadata"] = metadata

	var updated *unstructured.Unstructured
	if kind == "nodes" || kind == "namespaces" || kind == "clusterroles" || kind == "clusterrolebindings" || kind == "persistentvolumes" || kind == "storageclasses" || kind == "crds" {
		updated, err = clientSet.DynamicClient.Resource(gvr).Update(context.Background(), &unstructured.Unstructured{Object: unstructuredObj}, metav1.UpdateOptions{})
	} else {
		updated, err = clientSet.DynamicClient.Resource(gvr).Namespace(namespace).Update(context.Background(), &unstructured.Unstructured{Object: unstructuredObj}, metav1.UpdateOptions{})
	}

	if err != nil {
		response.InternalError(c, "Failed to update resource: "+err.Error())
		return
	}

	response.Success(c, updated)
}

func (h *ResourceHandler) DeleteResource(c *gin.Context) {
	client, _ := c.Get("client")
	clientSet := client.(*kubernetes.ClientSet)

	kind := c.Param("kind")
	namespace := c.Param("namespace")
	name := c.Param("name")

	gvr, ok := h.getGVR(kind)
	if !ok {
		response.NotFound(c, "Resource kind not supported: "+kind)
		return
	}

	var err error

	if kind == "nodes" || kind == "namespaces" || kind == "clusterroles" || kind == "clusterrolebindings" || kind == "persistentvolumes" || kind == "storageclasses" || kind == "crds" {
		err = clientSet.DynamicClient.Resource(gvr).Delete(context.Background(), name, metav1.DeleteOptions{})
	} else {
		err = clientSet.DynamicClient.Resource(gvr).Namespace(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	}

	if err != nil {
		response.InternalError(c, "Failed to delete resource: "+err.Error())
		return
	}

	response.Success(c, nil)
}

func toUnstructured(obj interface{}) (map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to map: %w", err)
	}

	return result, nil
}

func init() {
	_ = corev1.AddToScheme(scheme.Scheme)
}
