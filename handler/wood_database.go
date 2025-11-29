package handler

import (
	"net/http"
	"strconv"

	"backend/firestore"
	"backend/models"

	"github.com/gin-gonic/gin"
)

// CreateWoodDatabase tạo mới WoodDatabase
func CreateWoodDatabase(c *gin.Context) {
	var db models.WoodDatabase
	if err := c.BindJSON(&db); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate ID
	if db.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	// Kiểm tra ID đã tồn tại chưa
	exists, err := firestore.DocumentExists("wood_database", db.ID)
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Collection with this ID already exists"})
		return
	}

	err = firestore.CreateDocument("wood_database", db.ID, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Created successfully",
		"data":    db,
	})
}

// UpdateWoodDatabase cập nhật WoodDatabase
func UpdateWoodDatabase(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var db models.WoodDatabase
	if err := c.BindJSON(&db); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Đảm bảo ID trong body khớp với URL
	db.ID = id

	// Kiểm tra document có tồn tại không
	exists, err := firestore.DocumentExists("wood_database", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	err = firestore.UpdateDocument("wood_database", id, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Updated successfully",
		"data":    db,
	})
}

// GetWoodDatabase Get WoodDatabase by ID
func GetWoodDatabase(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	data, err := firestore.GetDocument("wood_database", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// ListWoodDatabase List all wood_database với phân trang
func ListWoodDatabase(c *gin.Context) {
	// Parse query params
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	orderBy := c.DefaultQuery("order_by", "title")
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

	result, err := firestore.GetCollectionPaginated("wood_database", params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteWoodDatabase xóa WoodDatabase
func DeleteWoodDatabase(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	// Kiểm tra document có tồn tại không
	exists, err := firestore.DocumentExists("wood_database", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	err = firestore.DeleteDocument("wood_database", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
