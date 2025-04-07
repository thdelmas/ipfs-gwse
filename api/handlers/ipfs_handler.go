package handlers

import (
	"bytes"
	"encoding/json"
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

func HandleMetadata(c *gin.Context) {
	cid := c.Param("cid")

	meta, fileData, err := processCID(cid)
	if err != nil {
		log.Printf("Error processing CID %s: %v", cid, err)
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

	// Check if the CID is a directory
	links, err := sh.List(cid)
	if err == nil && len(links) > 0 {
		// Build metadata for a directory
		dirMeta := &Metadata{
			CID:         cid,
			Name:        cid, // or some default name
			Size:        0,
			ContentType: "application/json",
			Path:        "", // not needed
		}

		// Create JSON buffer of the directory entries
		var dirContents []map[string]interface{}
		for _, link := range links {
			dirContents = append(dirContents, map[string]interface{}{
				"name": link.Name,
				"cid":  link.Hash,
				"size": link.Size,
				"type": link.Type,
			})
		}

		jsonBuf := new(bytes.Buffer)
		if err := json.NewEncoder(jsonBuf).Encode(dirContents); err != nil {
			return nil, nil, fmt.Errorf("failed to encode directory listing: %w", err)
		}

		dirMeta.Size = jsonBuf.Len()
		return dirMeta, jsonBuf, nil
	}

	// If not a directory, attempt to Cat the content
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
