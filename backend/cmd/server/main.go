package main

import (
	"github.com/k8s-dashboard/internal/handler"
	"github.com/k8s-dashboard/internal/kubernetes"
	"github.com/k8s-dashboard/internal/middleware"
	"github.com/k8s-dashboard/pkg/config"
	"github.com/k8s-dashboard/pkg/response"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Session-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "Resource not found")
	})

	clientManager := kubernetes.NewClientManager()
	sessionMiddleware := middleware.NewSessionMiddleware(clientManager)

	authHandler := handler.NewAuthHandler(clientManager)
	resourceHandler := handler.NewResourceHandler(clientManager)
	clusterHandler := handler.NewClusterHandler(clientManager)
	logHandler := handler.NewLogHandler(clientManager)
	execHandler := handler.NewExecHandler(clientManager)
	eventHandler := handler.NewEventHandler(clientManager)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/connect", authHandler.Connect)
			auth.POST("/disconnect", sessionMiddleware.Required(), authHandler.Disconnect)
		}

		cluster := api.Group("/cluster")
		cluster.Use(sessionMiddleware.Required())
		{
			cluster.GET("/info", clusterHandler.GetInfo)
			cluster.GET("/health", clusterHandler.GetHealth)
		}

		namespaces := api.Group("/namespaces")
		namespaces.Use(sessionMiddleware.Required())
		{
			namespaces.GET("", resourceHandler.ListNamespaces)
		}

		resources := api.Group("/resources")
		resources.Use(sessionMiddleware.Required())
		{
			resources.GET("/:kind", resourceHandler.ListResources)
			resources.GET("/:kind/:namespace/:name", resourceHandler.GetResource)
			resources.PUT("/:kind/:namespace/:name", resourceHandler.UpdateResource)
			resources.DELETE("/:kind/:namespace/:name", resourceHandler.DeleteResource)
			resources.GET("/:kind/:namespace/:name/yaml", resourceHandler.GetResourceYAML)
		}

		pods := api.Group("/pods")
		pods.Use(sessionMiddleware.Required())
		{
			pods.GET("/:namespace/:name/logs", logHandler.StreamLogs)
			pods.GET("/:namespace/:name/exec", execHandler.Exec)
		}

		events := api.Group("/events")
		events.Use(sessionMiddleware.Required())
		{
			events.GET("", eventHandler.ListEvents)
		}
	}

	r.Run(":8080")
}
