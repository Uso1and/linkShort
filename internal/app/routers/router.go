package routers

import (
	"log"

	"github.com/gin-gonic/gin"

	"lnkshrt/internal/app/handlers"
	"lnkshrt/internal/app/middleware"
	"lnkshrt/internal/domain/config"
	"lnkshrt/internal/domain/infrastructure/database"
	"lnkshrt/internal/domain/repo"
)

func SetupRoute() *gin.Engine {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	userRepo := repo.NewUserRepo(database.DB)

	userHandler := handlers.NewUserHandler(userRepo, cfg)

	r := gin.Default()

	r.LoadHTMLGlob("template/*")
	r.Static("static", "./static")

	r.GET("/", handlers.IndexPageHandler)
	r.GET("/register", handlers.RegisterPageHandler)
	r.GET("/login", handlers.LoginPageHandler)

	r.POST("/register", userHandler.CreateUserHandler)
	r.POST("/login", userHandler.LoginHandler)

	protected := r.Group("/")
	protected.Use(func(c *gin.Context) {
		if c.Request.Method == "GET" {
			if token := c.Query("token"); token != "" {
				c.Request.Header.Set("Authorization", "Bearer "+token)
			}
		}
		c.Next()
	})
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.GET("/main", handlers.MainPageHandler)
	}

	return r
}
