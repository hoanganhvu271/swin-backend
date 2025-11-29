package router

import (
	"backend/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

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

		library.POST("/database/save", handler.SaveWoodDatabase)
		library.GET("/database/get", handler.GetWoodDatabase)
		library.GET("/database/list", handler.ListWoodDatabase)

		library.POST("/piece/save", handler.SaveWoodPiece)
		library.GET("/piece/get", handler.GetWoodPiece)
		library.GET("/piece/list", handler.ListWoodPiecesByDatabase)
		library.DELETE("/piece/delete", handler.DeleteWoodPiece)
	}

	return r
}
