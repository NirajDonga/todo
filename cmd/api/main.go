package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/NirajDonga/todo/internal/config"
	"github.com/NirajDonga/todo/internal/database"

	"github.com/NirajDonga/todo/internal/auth"
	todopkg "github.com/NirajDonga/todo/internal/modules/todo"
	userpkg "github.com/NirajDonga/todo/internal/modules/user"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	ctx := context.Background()
	dbPool, err := database.Connect(ctx, cfg.DB_URI)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	defer dbPool.Close()
	// initialize services, repos and handlers
	authSvc := auth.NewService(cfg.JWT_SECRET, 24*time.Hour)

	userRepo := userpkg.NewUserRepository(dbPool)
	userSvc := userpkg.NewUserService(userRepo, authSvc)
	userHandler := userpkg.NewHandler(userSvc)

	todoRepo := todopkg.NewTodoRepo(dbPool)
	todoSvc := todopkg.NewTodoService(todoRepo)
	todoHandler := todopkg.NewHandler(todoSvc)

	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.LogIn)
	}

	// todo routes (will be protected by wrapping server with auth middleware)
	api.GET("/todos", todoHandler.GetTodos)
	api.POST("/todos", todoHandler.Create)

	// attach Gin auth middleware for protected routes
	authGin := func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}
		token := parts[1]
		claims, err := authSvc.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		if claims.UserID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}
		c.Set("userID", claims.UserID)
		c.Next()
	}

	protected := api.Group("")
	protected.Use(authGin)
	protected.GET("/todos", todoHandler.GetTodos)
	protected.POST("/todos", todoHandler.Create)

	// start Gin server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
