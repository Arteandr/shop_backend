package service

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
	fn "shop_backend/pkg/filename"
)

type ImagesService struct {
	repo repository.Images
}

func NewImagesService(repo repository.Images) *ImagesService {
	return &ImagesService{repo: repo}
}

func (s *ImagesService) Upload(image *multipart.FileHeader) (int, error) {
	ext := filepath.Ext(image.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return 0, errors.New("wrong file extension")
	}

	var filename string
LOOP:
	for {
		filename = fn.Generate() + ext
		if _, err := os.Stat("./files/" + filename); errors.Is(err, os.ErrNotExist) {
			break LOOP
		} else {
			continue LOOP
		}
	}

	if err := s.saveFile(image, "./files/"+filename); err != nil {
		return 0, err
	}

	return s.repo.Upload(filename)
}

func (s *ImagesService) saveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func (s *ImagesService) GetAll() ([]models.Image, error) {
	return s.repo.GetAll()
}
