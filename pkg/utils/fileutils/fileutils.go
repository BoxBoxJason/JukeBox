package fileutils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
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
