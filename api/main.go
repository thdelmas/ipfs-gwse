package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"

	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
	router := gin.Default()
	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	router.GET("/:cid", handleMetadata)
	router.Run(":8000")
}

func handleMetadata(c *gin.Context) {
	cid := c.Param("cid")

	meta, fileData, err := processCID(cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set metadata as headers
	c.Header("X-CID", meta.CID)
	c.Header("X-File-Name", meta.Name)
	c.Header("X-Size", fmt.Sprintf("%d", meta.Size))
	c.Header("X-Content-Type", meta.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", meta.Name))
	c.Data(http.StatusOK, meta.ContentType, fileData.Bytes())
}

type Metadata struct {
	CID         string `json:"cid"`
	Name        string `json:"name"`
	Size        int    `json:"size"`
	ContentType string `json:"content_type"`
	Path        string `json:"path"`
}

func processCID(cid string) (*Metadata, *bytes.Buffer, error) {
	sh := shell.NewShell("ipfs:5001")
	reader, err := sh.Cat(cid)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download CID %s: %w", cid, err)
	}
	defer reader.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		return nil, nil, fmt.Errorf("failed to read content: %w", err)
	}

	contentType := http.DetectContentType(buf.Bytes())
	size := buf.Len()

	tmpFile, err := saveToTempFile(buf)
	if err != nil {
		return nil, nil, err
	}
	defer removeFile(tmpFile.Name())

	info, err := tmpFile.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &Metadata{
		CID:         cid,
		Name:        info.Name(),
		Size:        size,
		ContentType: contentType,
		Path:        tmpFile.Name(),
	}, &buf, nil
}

func saveToTempFile(data bytes.Buffer) (*os.File, error) {
	tmpFile, err := os.CreateTemp("", "ipfsfile-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	if _, err = io.Copy(tmpFile, &data); err != nil {
		return nil, fmt.Errorf("failed to write to temp file: %w", err)
	}
	return tmpFile, nil
}

func openFile(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}

func removeFile(path string) {
	time.Sleep(5 * time.Second)
	if err := os.Remove(path); err != nil {
		log.Printf("Error removing file: %v", err)
	}
}
