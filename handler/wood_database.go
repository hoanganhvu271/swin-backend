package handler

import (
	"net/http"

	"backend/firestore"
	"backend/models"

	"github.com/gin-gonic/gin"
)

// SaveWoodDatabase Create or Update WoodDatabase
func SaveWoodDatabase(c *gin.Context) {
	var db models.WoodDatabase
	if err := c.BindJSON(&db); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := firestore.SetDocument("wood_database", db.ID, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Saved"})
}

// GetWoodDatabase Get WoodDatabase by ID
func GetWoodDatabase(c *gin.Context) {
	id := c.Query("id")
	data, err := firestore.GetDocument("wood_database", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// ListWoodDatabase List all wood_database
func ListWoodDatabase(c *gin.Context) {
	docs, err := firestore.GetCollection("wood_database")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, docs)
}
