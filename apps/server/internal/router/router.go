package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/config"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/handler"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/response"
	"go.uber.org/zap"
)

func New(cfg *config.Config, log *zap.Logger) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(requestLogger(log))
	engine.NoRoute(NoRoute)

	engine.GET("/swagger/*any", swaggerPlaceholder)

	api := engine.Group(cfg.API.Prefix)
	{
		healthHandler := handler.NewHealthHandler(cfg)
		api.GET("/health", healthHandler.Check)
	}

	return engine
}

func requestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		log.Info(
			"http request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
		)
	}
}

func swaggerPlaceholder(c *gin.Context) {
	response.Success(c, gin.H{
		"message": "Swagger UI will be enabled when API annotations are generated.",
		"spec":    "/docs/swagger.yaml",
	})
}

func NoRoute(c *gin.Context) {
	response.Error(c, http.StatusNotFound, "not_found", "route not found")
}
