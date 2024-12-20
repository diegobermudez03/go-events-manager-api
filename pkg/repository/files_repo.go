package repository

import (
	"context"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
)

type FilesInServer struct {
}

func NewFilesInServer() domain.FilesRepo{
	return &FilesInServer{}
}

func (r *FilesInServer) StoreImage(ctx context.Context, image *[]byte, group string, imageName string, imageType string, path string ) (string, error){
	return "", nil
}

func (r *FilesInServer) DeleteFile(ctx context.Context, group string,url string) error{
	return nil
}