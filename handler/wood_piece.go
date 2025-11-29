package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"backend/firestore"
	"backend/models"

	"github.com/gin-gonic/gin"
)

// CreateWoodPiece
func CreateWoodPiece(c *gin.Context) {
	var piece models.WoodPiece
	if err := c.BindJSON(&piece); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if piece.DatabaseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database_id is required"})
		return
	}

	pieces, err := firestore.GetDocumentsByField("wood_piece", "database_id", piece.DatabaseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	maxIndex := 0
	for _, p := range pieces {
		var index int
		fmt.Sscanf(p.ID, piece.DatabaseID+"_%d", &index)
		if index > maxIndex {
			maxIndex = index
		}
	}

	newIndex := maxIndex + 1
	piece.ID = fmt.Sprintf("%s_%02d", piece.DatabaseID, newIndex)

	err = firestore.CreateDocument("wood_piece", piece.ID, piece)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Created successfully",
		"data":    piece,
	})
}

// UpdateWoodPiece cập nhật WoodPiece
func UpdateWoodPiece(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var piece models.WoodPiece
	if err := c.BindJSON(&piece); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Đảm bảo ID trong body khớp với URL
	piece.ID = id

	// Kiểm tra document có tồn tại không
	exists, err := firestore.DocumentExists("wood_piece", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Piece not found"})
		return
	}

	err = firestore.UpdateDocument("wood_piece", id, piece)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Updated successfully",
		"data":    piece,
	})
}

// GetWoodPiece Get WoodPiece by ID
func GetWoodPiece(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	data, err := firestore.GetDocument("wood_piece", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Piece not found"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// ListWoodPiecesByDatabase List WoodPiece by DatabaseID với phân trang
func ListWoodPiecesByDatabase(c *gin.Context) {
	dbID := c.Query("database_id")
	if dbID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database_id is required"})
		return
	}

	// Parse query params
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	orderBy := c.DefaultQuery("order_by", "name")
	descending := c.DefaultQuery("desc", "false") == "true"

	// Validate limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	params := firestore.PaginationParams{
		Limit:      limit,
		Offset:     offset,
		OrderBy:    orderBy,
		Descending: descending,
	}

	result, err := firestore.GetCollectionWithFilter("wood_piece", "database_id", dbID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteWoodPiece Delete WoodPiece
func DeleteWoodPiece(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	// Kiểm tra document có tồn tại không
	exists, err := firestore.DocumentExists("wood_piece", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Piece not found"})
		return
	}

	err = firestore.DeleteDocument("wood_piece", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
