package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)



type Song struct {
	Name string `json:"name"`
	File string `json:"file"`
}

func main() {
	// Create songs directory if it doesn't exist
	if err := os.MkdirAll("songs", 0755); err != nil {
		log.Fatal(err)
	}

	// Create static directory if it doesn't exist
	if err := os.MkdirAll("static", 0755); err != nil {
		log.Fatal(err)
	}

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("static/")))

	// API endpoint to list songs
	http.HandleFunc("/songs", listSongsHandler)

	// API endpoint to stream a song
	http.HandleFunc("/song/", streamSongHandler)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func listSongsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var songs []Song

	// Read songs directory
	err := filepath.WalkDir("songs", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Only include audio files
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".mp3" || ext == ".wav" || ext == ".ogg" || ext == ".flac" {
			relPath, _ := filepath.Rel("songs", path)
			song := Song{
				Name: strings.TrimSuffix(relPath, ext),
				File: relPath,
			}
			songs = append(songs, song)
		}

		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(songs)
}

func streamSongHandler(w http.ResponseWriter, r *http.Request) {
	// Extract filename from URL path
	filePath := strings.TrimPrefix(r.URL.Path, "/song/")

	// Security check to prevent directory traversal
	if strings.Contains(filePath, "..") {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Full path to the file
	fullPath := filepath.Join("songs", filePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set appropriate content type
	ext := strings.ToLower(filepath.Ext(fullPath))
	switch ext {
	case ".mp3":
		w.Header().Set("Content-Type", "audio/mpeg")
	case ".wav":
		w.Header().Set("Content-Type", "audio/wav")
	case ".ogg":
		w.Header().Set("Content-Type", "audio/ogg")
	case ".flac":
		w.Header().Set("Content-Type", "audio/flac")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}


	// Stream the file
	http.ServeFile(w, r, fullPath)
}
