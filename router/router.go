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
		AllowOrigins:     []string{"http://localhost:3000"},
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
		library.POST("/upload_image", handler.UploadImage)

		// Wood Database (Collection) - RESTful APIs
		database := library.Group("/database")
		{
			database.GET("/list", handler.ListWoodDatabase)         // GET /database/list?limit=10&offset=0
			database.GET("/get", handler.GetWoodDatabase)           // GET /database/get?id=xxx
			database.POST("/create", handler.CreateWoodDatabase)    // POST /database/create
			database.PUT("/update/:id", handler.UpdateWoodDatabase) // PUT /database/update/:id
			database.DELETE("/delete", handler.DeleteWoodDatabase)  // DELETE /database/delete?id=xxx
		}

		// Wood Piece - RESTful APIs
		piece := library.Group("/piece")
		{
			piece.GET("/list", handler.ListWoodPiecesByDatabase) // GET /piece/list?database_id=xxx&limit=10&offset=0
			piece.GET("/get", handler.GetWoodPiece)              // GET /piece/get?id=xxx
			piece.POST("/create", handler.CreateWoodPiece)       // POST /piece/create
			piece.PUT("/update/:id", handler.UpdateWoodPiece)    // PUT /piece/update/:id
			piece.DELETE("/delete", handler.DeleteWoodPiece)     // DELETE /piece/delete?id=xxx
		}
	}

	return r
}
