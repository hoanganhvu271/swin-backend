package router

import (
	"backend/handler"
	"backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	model := r.Group("/model-api")
	{
		model.GET("/version", handler.GetModelVersion)
		model.GET("/list_versions", handler.ListModelVersions)
		model.POST("/activate", handler.ActivateNewModel)
		model.POST("/upload", handler.UploadNewModel)
	}

	// Protected routes - yêu cầu đăng nhập
	library := r.Group("/library-api")
	library.Use(middleware.AuthMiddleware()) // Thêm middleware auth
	{
		// Upload image
		library.POST("/upload_image", handler.UploadImage)

		// Wood Database (Collection) - RESTful APIs
		library.GET("/database/list", handler.ListWoodDatabase)
		library.GET("/database/get", handler.GetWoodDatabase)
		library.POST("/database/create", handler.CreateWoodDatabase)
		library.PUT("/database/update/:id", handler.UpdateWoodDatabase)
		library.DELETE("/database/delete", handler.DeleteWoodDatabase)

		// Wood Piece - RESTful APIs
		library.GET("/piece/list", handler.ListWoodPiecesByDatabase)
		library.GET("/piece/get", handler.GetWoodPiece)
		library.POST("/piece/create", handler.CreateWoodPiece)
		library.PUT("/piece/update/:id", handler.UpdateWoodPiece)
		library.DELETE("/piece/delete", handler.DeleteWoodPiece)
	}

	return r
}
