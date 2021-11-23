package domain

// Article ...
type Image struct {
	Chunks     int  `json:"chunks" validate:"required"`
	Completed  bool `json:"completed" validate:"required"`
	Proccessed bool `json:"proccessed" validate:"required"`
}

type ImageReg struct {
	Sha256    string `json:"sha256"  validate:"required"`
	Size      int64  `json:"size" validate:"required,gt=0"`
	ChunkSize int64  `json:"chunk_size" validate:"required"`
}

type ImageChunk struct {
	Id   int64  `json:"id"  `
	Size int64  `json:"size" validate:"required"`
	Data string `json:"data" validate:"required"`
}

type ImageUsecase interface {
	RegisterImage(i *ImageReg) error
	StoreImageChunk(sha256 string, chunk *ImageChunk) error
	GetImage(sha256 string) ([]byte, error)
}

type ImageRepository interface {
	StoreImageInfo(sha256 string, chunks int64, chunkSize int64) error
	StoreImageChunk(sha256 string, chunkId string) error
	DeleteImageChunk(sha256 string, chunkId string) error
	GetImageChunksNumber(sha256 string) (string, error)
	GetImageChunkSize(sha256 string) (string, error)
	CheckImageExist(sha256 string) (int64, error)
	CheckChunkExist(sha256 string, chunkId string) (int64, error)
}
