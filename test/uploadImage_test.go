package test

import (
	"fmt"
	"testing"

	"github.com/nimasrn/uploadsvc-go/domain"

	"github.com/nimasrn/uploadsvc-go/image/repository/redis"
	"github.com/nimasrn/uploadsvc-go/image/usecase"

	goredis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/spf13/viper"
)

func TestUploadImage(t *testing.T) {
	viper.SetConfigFile("./config.test.json")
	err := viper.ReadInConfig()
	assert.NoError(t, err)
	storage := viper.GetString(`disk`)
	redisHost := viper.GetString(`redis.host`)
	redisPort := viper.GetString(`redis.port`)
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	redisOpt := &goredis.Options{
		Addr: redisAddr,
		DB:   0,
	}
	redisConn := goredis.NewClient(redisOpt)
	assert.NotNil(t, redisConn)
	imageRepo := redis.NewRedisImageRepository(redisConn)
	iu := usecase.NewImageUsecase(imageRepo, storage)

	sha256 := uuid.New().String()
	chunk := &domain.ImageChunk{
		Id:   0,
		Size: 256,
		Data: "hello\ngo\n",
	}
	imageReg := &domain.ImageReg{
		Sha256:    sha256,
		Size:      int64(len([]byte(chunk.Data))),
		ChunkSize: 1,
	}
	err = iu.RegisterImage(imageReg)
	assert.NoError(t, err)
	err = iu.RegisterImage(imageReg)
	assert.Error(t, err)
	err = iu.StoreImageChunk(sha256, chunk)
	assert.NoError(t, err)
	err = iu.StoreImageChunk(sha256, chunk)
	assert.Error(t, err)
	image, err := iu.GetImage(sha256)
	assert.NoError(t, err)
	assert.Equal(t, image, []byte(chunk.Data))
	image, err = iu.GetImage(sha256 + "nothing")
	assert.Error(t, err)
	assert.Nil(t, image)
}
