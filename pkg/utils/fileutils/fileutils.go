package fileutils

import (
	"archive/zip"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
)

// CompressFiles compresses a list of files into a single archive file at the specified path.
// The archive file will be created if it does not exist. If the archive file already exists,
// The new files will be appended to the archive file.
func CompressFiles(archive_path string, files []string) error {
	archive, err := os.Create(archive_path)
	if err != nil {
		return err
	}
	defer archive.Close()

	zip_writer := zip.NewWriter(archive)
	defer zip_writer.Close()

	for _, file := range files {
		if err := AddFileToZip(zip_writer, file); err != nil {
			return err
		}
	}

	return nil
}

// AddFileToZip adds a file to a zip.Writer.
func AddFileToZip(zip_writer *zip.Writer, file_path string) error {
	file, err := os.Open(file_path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filepath.Base(file_path)

	writer, err := zip_writer.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

func DeleteFile(file_path string) error {
	return os.Remove(file_path)
}

// SaveImageFile saves an uploaded image file to the specified directory with the given filename.
// It returns the full path to the saved file or an error if the save fails.
func SaveImageFile(file multipart.File, destination_dir string, filename string) error {
	// Ensure the destination directory exists
	err := os.MkdirAll(destination_dir, os.ModePerm)
	if err != nil {
		return httputils.NewInternalServerError("Failed to create directory: " + destination_dir)
	}

	// Define the full path for the new file
	dest_path := filepath.Join(destination_dir, filename)

	// Create a new file in the destination directory
	dest_file, err := os.Create(dest_path)
	if err != nil {
		return httputils.NewInternalServerError("Failed to create file: " + dest_path)
	}
	defer dest_file.Close()

	// Copy the uploaded file's contents to the new file
	_, err = io.Copy(dest_file, file)
	if err != nil {
		return httputils.NewInternalServerError("Failed to save file to filesystem")
	}

	return nil
}
