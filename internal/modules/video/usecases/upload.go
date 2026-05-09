package usecases

import (
	"go-fake-flix/internal/apierrors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func UploadVideo(file multipart.File, fileName string) (filePath string, error *apierrors.AppError) {
	dst := filepath.Join("./files/", filepath.Base(fileName))

	// err := validateFile(file, fileName)
	// if err != nil {
	// 	return "", &apierrors.AppError{Code: "FILE_NOT_VALID", Message: "File not valid"}
	// }

	dstFile, err := os.Create(dst)
	if err != nil {
		return "", &apierrors.AppError{Code: "FILE_NOT_CREATED", Message: "File not created"}
	}
	defer dstFile.Close()
	io.Copy(dstFile, file)
	return dst, nil
}

// func validateFile(file multipart.File, fileName string) error {

// }
