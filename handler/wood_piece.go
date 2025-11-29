package handler

import (
	"backend/firestore"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SaveWoodPiece Create or Update WoodPiece
func SaveWoodPiece(c *gin.Context) {
	var piece models.WoodPiece
	if err := c.BindJSON(&piece); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := firestore.SetDocument("wood_piece", piece.ID, piece)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Saved"})
}

// GetWoodPiece Get WoodPiece by ID
func GetWoodPiece(c *gin.Context) {
	id := c.Query("id")
	data, err := firestore.GetDocument("wood_piece", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// ListWoodPiecesByDatabase List WoodPiece by DatabaseID
func ListWoodPiecesByDatabase(c *gin.Context) {
	dbID := c.Query("database_id")
	docs, err := firestore.GetCollection("wood_piece")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter by database_id
	var result []map[string]interface{}
	for _, d := range docs {
		if d["database_id"] == dbID {
			result = append(result, d)
		}
	}
	c.JSON(http.StatusOK, result)
}

// DeleteWoodPiece Delete WoodPiece
func DeleteWoodPiece(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id"})
		return
	}

	err := firestore.DeleteDocument("wood_piece", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
