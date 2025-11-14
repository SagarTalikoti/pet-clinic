package handlers

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"pet-clinic/utils"
	"strings"

	"github.com/gorilla/mux"
)

// UploadFile handles file upload
func UploadFile(w http.ResponseWriter, r *http.Request) {
	utils.Log.Debug("Received file upload request")

	// Parse up to 10 MB of incoming data
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.Log.WithError(err).Warn("Failed to parse multipart form data")
		http.Error(w, "Could not process file", http.StatusBadRequest)
		return
	}

	// Get uploaded file from form-data
	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.Log.WithError(err).Warn("Missing 'file' field in upload request")
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure uploads folder exists
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	// Save file to uploads/ folder
	filePath := filepath.Join("uploads", handler.Filename)
	dest, err := os.Create(filePath)
	if err != nil {
		utils.Log.WithError(err).Error("Failed to create file on disk")
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}
	defer dest.Close()

	// Copy file contents
	_, err = io.Copy(dest, file)
	if err != nil {
		utils.Log.WithError(err).Error("Error saving file data")
		http.Error(w, "File save failed", http.StatusInternalServerError)
		return
	}

	utils.Log.WithField("filename", handler.Filename).Info("File uploaded successfully")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully: %v", handler.Filename)
}

// DownloadFile handles file download
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["filename"]

	//  Decode URL-encoded parts (e.g., %20 → space)
	decodedFilename, err := url.QueryUnescape(filename)
	if err != nil {
		utils.Log.WithError(err).Error("Failed to decode filename")
		http.Error(w, "Invalid filename encoding", http.StatusBadRequest)
		return
	}

	// Remove any stray newlines or spaces
	decodedFilename = strings.TrimSpace(decodedFilename)

	utils.Log.WithField("filename", decodedFilename).Debug("Received download request")

	filePath := filepath.Join("uploads", decodedFilename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.Log.WithField("filename", decodedFilename).Warn("Requested file not found")
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Open and serve file
	file, err := os.Open(filePath)
	if err != nil {
		utils.Log.WithError(err).Error("Error opening file for download")
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	mimeType := "application/octet-stream"
	if detected := mime.TypeByExtension(filepath.Ext(decodedFilename)); detected != "" {
		mimeType = detected
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+decodedFilename)
	w.Header().Set("Content-Type", mimeType)

	http.ServeFile(w, r, filePath)

	utils.Log.WithField("filename", decodedFilename).Info("File downloaded successfully")
}
