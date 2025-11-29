package handler

import (
	"backend/service"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /model/version
func GetModelVersion(c *gin.Context) {
	meta, err := service.GetCurrentModelMetadata()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, meta)
}

// GET /model/list_versions
func ListModelVersions(c *gin.Context) {
	versions, err := service.ListVersions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

// POST /model/activate
func ActivateNewModel(c *gin.Context) {
	versionStr := c.Query("version")
	if versionStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version required"})
		return
	}
	version, _ := strconv.Atoi(versionStr)
	err := service.ActivateVersion(version)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "version activated", "version": version})
}

// POST /model/upload
func UploadNewModel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not provided"})
		return
	}

	name := c.PostForm("name")
	if name == "" {
		name = file.Filename
	}

	log.Println(file.Filename)

	vInfo, err := service.ReadVersion()
	if err != nil {
		vInfo = &service.VersionInfo{
			CurrentVersion: 0,
			Versions:       []service.VersionEntry{},
		}
	}

	newVersion := 1
	if len(vInfo.Versions) > 0 {
		for _, ver := range vInfo.Versions {
			if ver.Version >= newVersion {
				newVersion = ver.Version + 1
			}
		}
	}

	dst := filepath.Join("models", fmt.Sprintf("model_v%d_%s", newVersion, file.Filename))
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save file"})
		return
	}

	// Upload lên Cloudinary và lưu URL
	entry, err := service.UploadNewModel(dst, newVersion, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "file uploaded",
		"version": entry.Version,
		"name":    name,
		"file":    entry.File,
	})
}
