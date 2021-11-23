package redis

import (
	"context"

	"github.com/nimasrn/uploadsvc-go/domain"

	goredis "github.com/go-redis/redis/v8"
)

type redisImageRepository struct {
	Conn *goredis.Client
}

func NewRedisImageRepository(Conn *goredis.Client) domain.ImageRepository {
	return &redisImageRepository{Conn}
}

func (m *redisImageRepository) StoreImageInfo(sha256 string, chunks int64, chunkSize int64) error {
	st := m.Conn.HMSet(context.Background(), sha256, "chunks", chunks, "completed", 0, "proccessed", 0, "chunkSize", chunkSize)
	return st.Err()
}

func (m *redisImageRepository) GetImageChunksNumber(sha256 string) (string, error) {
	st := m.Conn.HGet(context.Background(), sha256, "chunks")
	res, err := st.Result()
	return res, err
}
func (m *redisImageRepository) GetImageChunkSize(sha256 string) (string, error) {
	st := m.Conn.HGet(context.Background(), sha256, "chunkSize")
	return st.Result()
}

func (m *redisImageRepository) StoreImageChunk(sha256 string, chunkId string) error {
	st := m.Conn.Set(context.Background(), sha256+":"+chunkId, 1, 0)
	return st.Err()
}
func (m *redisImageRepository) DeleteImageChunk(sha256 string, chunkId string) error {
	st := m.Conn.Del(context.Background(), sha256+":"+chunkId)
	return st.Err()
}

func (m *redisImageRepository) CheckImageExist(sha256 string) (int64, error) {
	st := m.Conn.Exists(context.Background(), sha256)
	return st.Result()
}

func (m *redisImageRepository) CheckChunkExist(sha256 string, chunkId string) (int64, error) {
	st := m.Conn.Exists(context.Background(), sha256+":"+chunkId)
	return st.Result()
}
