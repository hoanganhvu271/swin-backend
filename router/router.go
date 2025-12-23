package router

import (
	"backend/handler"

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

	library := r.Group("/library-api")
	{
		// Upload image
		library.POST("/upload_image", handler.UploadImage)

		// Wood Database (Collection) - RESTful APIs
		library.GET("/database/list", handler.ListWoodDatabase)         // GET /library-api/database/list?limit=10&offset=0
		library.GET("/database/get", handler.GetWoodDatabase)           // GET /library-api/database/get?id=xxx
		library.POST("/database/create", handler.CreateWoodDatabase)    // POST /library-api/database/create
		library.PUT("/database/update/:id", handler.UpdateWoodDatabase) // PUT /library-api/database/update/:id
		library.DELETE("/database/delete", handler.DeleteWoodDatabase)  // DELETE /library-api/database/delete?id=xxx

		// Wood Piece - RESTful APIs
		library.GET("/piece/list", handler.ListWoodPiecesByDatabase) // GET /library-api/piece/list?database_id=xxx&limit=10&offset=0
		library.GET("/piece/get", handler.GetWoodPiece)              // GET /library-api/piece/get?id=xxx
		library.POST("/piece/create", handler.CreateWoodPiece)       // POST /library-api/piece/create
		library.PUT("/piece/update/:id", handler.UpdateWoodPiece)    // PUT /library-api/piece/update/:id
		library.DELETE("/piece/delete", handler.DeleteWoodPiece)     // DELETE /library-api/piece/delete?id=xxx
	}

	return r
}
