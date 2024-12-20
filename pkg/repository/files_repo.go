package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
)

type FilesInServer struct {
	rootFolder	string
}

func NewFilesInServer() domain.FilesRepo{
	return &FilesInServer{
		rootFolder: "",
	}
}

func (r *FilesInServer) StoreImage(ctx context.Context, image *[]byte, group string, imageName string, imageType string, path string ) (string, error){
	targetPath, err := r.getPath(group, path)
	if err != nil{
		log.Println(err.Error())
		return "", domain.ErrInternal 
	}
	fileName := filepath.Join(targetPath, fmt.Sprintf("%s.%s", imageName, imageType))
	if err := os.WriteFile(fileName, *image, os.ModePerm); err != nil{
		log.Println(err.Error())
		return "", domain.ErrInternal
	}
	return fileName, nil
}

func (r *FilesInServer) DeleteFile(ctx context.Context, group string,url string) error{
	return nil
}

func (r *FilesInServer) getPath(group string, path string)(string, error){
	targetFolder := filepath.Join("..", "storage", group, path)
	err := os.MkdirAll(targetFolder, os.ModePerm)
	if err != nil{
		return "", err
	}
	return targetFolder, nil
}