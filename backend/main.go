package main

import (
	"log"
	"quotes-app/config"
	"quotes-app/handlers"
	"quotes-app/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Подключение к БД
	config.ConnectDatabase()
	config.InitJWT()

	// Автомиграция
	//err := config.DB.AutoMigrate(
	//	&models.User{},
	//	&models.Quote{},
	//	&models.Category{},
	//	&models.QuoteLike{},
	//	&models.Comment{},
	//	&models.CommentLike{},
	//)
	//if err != nil {
	//	log.Fatal("Failed to migrate database:", err)
	//}
	//
	//log.Println("Database migration completed successfully")

	log.Println("Database connected successfully. Using SQL migrations.")

	router := gin.Default()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Handlers
	authHandler := handlers.NewAuthHandler()
	quoteHandler := handlers.NewQuoteHandler()
	categoryHandler := handlers.NewCategoryHandler()
	commentHandler := handlers.NewCommentHandler()

	// Public routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.GET("/quotes", quoteHandler.GetQuotes)
	router.GET("/quotes/:id", quoteHandler.GetQuoteByID)
	router.GET("/categories", categoryHandler.GetCategories)
	router.GET("/quotes/:id/comments", commentHandler.GetComments)

	// Protected routes
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		// Цитаты
		auth.POST("/quotes", quoteHandler.CreateQuote)
		auth.PUT("/quotes/:id", quoteHandler.UpdateQuote)
		auth.DELETE("/quotes/:id", quoteHandler.DeleteQuote)
		auth.POST("/quotes/:id/like", quoteHandler.LikeQuote)
		auth.POST("/quotes/:id/dislike", quoteHandler.DislikeQuote)

		// Комментарии
		auth.POST("/quotes/:id/comments", commentHandler.AddComment)
		auth.POST("/comments/:id/like", commentHandler.LikeComment)
		auth.PUT("/comments/:id", commentHandler.UpdateComment)
		auth.DELETE("/comments/:id", commentHandler.DeleteComment)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "OK",
			"database": "connected",
		})
	})

	// Database check
	router.GET("/db-check", func(c *gin.Context) {
		var result struct {
			UsersCount      int64 `json:"users_count"`
			QuotesCount     int64 `json:"quotes_count"`
			CategoriesCount int64 `json:"categories_count"`
		}

		config.DB.Table("users").Count(&result.UsersCount)
		config.DB.Table("quotes").Count(&result.QuotesCount)
		config.DB.Table("categories").Count(&result.CategoriesCount)

		c.JSON(200, gin.H{
			"status": "Database accessible",
			"data":   result,
		})
	})

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
