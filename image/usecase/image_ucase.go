package usecase

import (
	"os"
	"strconv"

	"github.com/nimasrn/uploadsvc-go/domain"

	"github.com/sirupsen/logrus"
)

type imageUsecase struct {
	ImageRepo       domain.ImageRepository
	DiskStoragePath string
}

func NewImageUsecase(i domain.ImageRepository, disk string) domain.ImageUsecase {
	return &imageUsecase{
		ImageRepo:       i,
		DiskStoragePath: disk,
	}
}

func (a *imageUsecase) RegisterImage(image *domain.ImageReg) error {
	exist, err := a.ImageRepo.CheckImageExist(image.Sha256)
	if exist == 1 {
		return domain.ErrorImageConflict
	}
	if err != nil {
		return err
	}
	chunksNum := image.Size / image.ChunkSize
	err = a.ImageRepo.StoreImageInfo(image.Sha256, chunksNum, image.ChunkSize)
	if err != nil {
		return err
	}
	fd, err := os.Create(a.DiskStoragePath + image.Sha256)
	defer fd.Close()
	if err != nil {
		return err
	}
	// err = fallocate.Fallocate(fd, 0, image.Size)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (a *imageUsecase) StoreImageChunk(sha256 string, chunk *domain.ImageChunk) error {
	strId := strconv.FormatInt(chunk.Id, 10)
	exist, err := a.ImageRepo.CheckImageExist(sha256)
	if exist == 0 {
		return domain.ErrorImageNotFound
	}
	if err != nil {
		return err
	}
	exist, err = a.ImageRepo.CheckChunkExist(sha256, strId)
	if exist == 1 {
		return domain.ErrorChunkConflict
	}
	if err != nil {
		return err
	}
	chunkSize, _ := a.ImageRepo.GetImageChunkSize(sha256)
	chunkSizeInt, err := strconv.Atoi(chunkSize)
	if err != nil {
		return err
	}

	rune := []rune(chunk.Data)
	b := []byte(string(rune))
	offset := chunk.Id * int64(chunkSizeInt)
	f, err := os.OpenFile(a.DiskStoragePath+sha256, os.O_RDWR, 0644)
	if err != nil {
		logrus.Error("error opening file", "error:", err)
		return err
	}
	if _, err := f.WriteAt(b, offset); err != nil {
		return err
	}
	err = a.ImageRepo.StoreImageChunk(sha256, strId)
	if err != nil {
		err = a.ImageRepo.DeleteImageChunk(sha256, strId)
		if err != nil {
			logrus.Error(err)
		}
	}
	return err
}

func (a *imageUsecase) GetImage(sha256 string) ([]byte, error) {
	dat, err := os.ReadFile(a.DiskStoragePath + sha256)
	return dat, err
}
