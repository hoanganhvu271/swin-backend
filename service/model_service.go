package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"backend/config"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ModelMetadata struct {
	Name        string `json:"name"`
	Version     int    `json:"version"`
	Size        int64  `json:"size"`
	Checksum    string `json:"checksum"`
	DownloadURL string `json:"download_url"`
}

type VersionEntry struct {
	Version  int    `json:"version"`
	File     string `json:"file"`
	Name     string `json:"name"`
	Checksum string `json:"checksum"`
	Size     int64  `json:"size"`
}

type VersionInfo struct {
	CurrentVersion int            `json:"current_version"`
	Versions       []VersionEntry `json:"versions"`
}

func ModelDir() string {
	return config.ModelDir
}

func VersionFilePath() string {
	return filepath.Join(ModelDir(), "version.json")
}

func CurrentModelPath() string {
	v, _ := ReadVersion()
	for _, ver := range v.Versions {
		if ver.Version == v.CurrentVersion {
			return filepath.Join(ModelDir(), ver.File)
		}
	}
	return ""
}

func ReadVersion() (*VersionInfo, error) {
	data, err := os.ReadFile(VersionFilePath())
	if os.IsNotExist(err) {
		return &VersionInfo{
			CurrentVersion: 0,
			Versions:       []VersionEntry{},
		}, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot read version.json: %v", err)
	}
	var v VersionInfo
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("invalid version.json: %v", err)
	}
	return &v, nil
}

func WriteVersion(v *VersionInfo) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(VersionFilePath(), data, 0644)
}

func UploadNewModel(filePath string, version int, name string) (*VersionEntry, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(data)
	checksum := hex.EncodeToString(sum[:])

	ctx := context.Background()
	publicID := fmt.Sprintf("models/model_v%d_%s", version, filepath.Base(filePath))
	uploadResp, err := config.CLD.Upload.Upload(ctx, filePath, uploader.UploadParams{
		PublicID:     publicID,
		ResourceType: "raw",
	})
	if err != nil {
		return nil, fmt.Errorf("cloudinary upload failed: %v", err)
	}

	entry := VersionEntry{
		Version:  version,
		File:     filepath.Base(filePath),
		Name:     name,
		Checksum: checksum,
		Size:     int64(len(data)),
	}

	// Lưu URL download
	entryDownloadURL := uploadResp.SecureURL
	fmt.Println("Model uploaded to Cloudinary:", entryDownloadURL)

	// Ghi metadata
	vInfo, err := ReadVersion()
	if err != nil {
		vInfo = &VersionInfo{
			CurrentVersion: 0,
			Versions:       []VersionEntry{},
		}
	}
	entry.File = entryDownloadURL
	vInfo.Versions = append(vInfo.Versions, entry)

	if err := WriteVersion(vInfo); err != nil {
		return nil, err
	}

	return &entry, nil
}

func ActivateVersion(version int) error {
	vInfo, err := ReadVersion()
	if err != nil {
		return err
	}

	var target *VersionEntry
	for _, ver := range vInfo.Versions {
		if ver.Version == version {
			target = &ver
			break
		}
	}
	if target == nil {
		return fmt.Errorf("version %d not found", version)
	}

	vInfo.CurrentVersion = version
	return WriteVersion(vInfo)
}

func GetCurrentModelMetadata() (*ModelMetadata, error) {
	vInfo, err := ReadVersion()
	if err != nil {
		return nil, err
	}

	var current *VersionEntry
	for _, ver := range vInfo.Versions {
		if ver.Version == vInfo.CurrentVersion {
			current = &ver
			break
		}
	}
	if current == nil {
		return nil, fmt.Errorf("current version not found")
	}

	meta := &ModelMetadata{
		Name:        config.ModelName,
		Version:     current.Version,
		Size:        current.Size,
		Checksum:    current.Checksum,
		DownloadURL: current.File,
	}
	return meta, nil
}

func ListVersions() ([]VersionEntry, error) {
	vInfo, err := ReadVersion()
	if err != nil {
		return nil, err
	}
	sort.Slice(vInfo.Versions, func(i, j int) bool { return vInfo.Versions[i].Version < vInfo.Versions[j].Version })
	return vInfo.Versions, nil
}

func GetModelFilePath(version int) (string, error) {
	vInfo, err := ReadVersion()
	if err != nil {
		return "", err
	}
	for _, ver := range vInfo.Versions {
		if ver.Version == version {
			return filepath.Join(ModelDir(), ver.File), nil
		}
	}
	return "", fmt.Errorf("version %d not found", version)
}
