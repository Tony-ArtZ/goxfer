package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

type TransferManager struct {
	PushChan chan string
	SaveDir  string
	fileMap  map[string]string
	mu       sync.RWMutex
}

func NewTransferManager() (*TransferManager, error) {
	var saveDir string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		saveDir = filepath.Join(homeDir, "Downloads", "Goxfer")
	case "darwin", "linux":
		saveDir = filepath.Join(homeDir, "Downloads", "Goxfer")
	default:
		saveDir = filepath.Join(homeDir, "Goxfer")
	}

	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create save directory: %w", err)
	}

	color.Green("📂 Files will be saved to: %s", saveDir)

	return &TransferManager{
		PushChan: make(chan string),
		SaveDir:  saveDir,
		fileMap:  make(map[string]string),
	}, nil
}

func (tm *TransferManager) QueueFile(path string) {
	id := uuid.New().String()
	tm.mu.Lock()
	tm.fileMap[id] = path
	tm.mu.Unlock()

	downloadURL := fmt.Sprintf("/download?id=%s", id)
	tm.PushChan <- downloadURL
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	contentBytes, err := os.ReadFile("./template/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(contentBytes)
}

func (tm *TransferManager) HandleDownload(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	tm.mu.RLock()
	path, ok := tm.fileMap[id]
	tm.mu.RUnlock()

	if !ok {
		http.Error(w, "File Not Found or Expired", http.StatusNotFound)
		return
	}

	file, err := os.Open(path)
	if err != nil {
		http.Error(w, "File Not Found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(path))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, path)
}

func (tm *TransferManager) HandlePoll(w http.ResponseWriter, r *http.Request) {
	url := <-tm.PushChan
	json.NewEncoder(w).Encode(map[string]string{"url": url})
}

func (tm *TransferManager) HandleUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<30)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	savePath := filepath.Join(tm.SaveDir, header.Filename)
	color.Cyan("\n📥 Receiving file: %s", header.Filename)

	dst, err := os.Create(savePath)
	if err != nil {
		color.Red("❌ Failed to create file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		color.Red("❌ Failed to save file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	color.Green("✅ File saved successfully to %s", savePath)
	w.Write([]byte("Upload successful"))
}
