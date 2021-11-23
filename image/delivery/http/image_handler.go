package http

import (
	"net/http"

	"github.com/nimasrn/uploadsvc-go/domain"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

const errorMalformed = "Malformed"

type ResponseError struct {
	Message string `json:"message"`
}

type ImageHandler struct {
	IUsecase        domain.ImageUsecase
	DiskStoragePath string
}

func NewImageHandler(rg *gin.RouterGroup, us domain.ImageUsecase, disk string) {

	handler := &ImageHandler{
		IUsecase:        us,
		DiskStoragePath: disk,
	}

	rg.POST("/image", handler.Register)
	rg.POST("/image/:sha256/chunks", handler.UploadChunk)
	rg.GET("/image/:sha256", handler.GetImage)
}

func (i *ImageHandler) Register(c *gin.Context) {
	var imageReg domain.ImageReg
	err := c.BindJSON(&imageReg)
	if err != nil {
		c.String(http.StatusBadRequest, errorMalformed)
		c.Abort()
		return
	}
	_, err = isRequestValid(&imageReg)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		c.Abort()
		return
	}
	err = i.IUsecase.RegisterImage(&imageReg)
	if err != nil {
		c.String(getStatusCode(err), err.Error())
		c.Abort()
		return
	}
	c.String(201, "")
}

func (i *ImageHandler) UploadChunk(c *gin.Context) {
	sha256 := c.Param("sha256")
	var chunk domain.ImageChunk
	err := c.BindJSON(&chunk)
	if err != nil {
		c.JSON(200, err)
	}
	_, err = isRequestValidChunk(&chunk)
	if err != nil {
		c.JSON(400, err.Error())
		c.Abort()
		return
	}
	err = i.IUsecase.StoreImageChunk(sha256, &chunk)
	if err != nil {
		c.JSON(getStatusCode(err), err.Error())
		c.Abort()
		return
	}
	c.String(201, "")
}

func (i *ImageHandler) GetImage(c *gin.Context) {
	sha256 := c.Param("sha256")
	c.File(i.DiskStoragePath + sha256)
}

func isRequestValid(m *domain.ImageReg) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func isRequestValidChunk(m *domain.ImageChunk) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrorChunkConflict:
		return http.StatusConflict
	case domain.ErrorImageConflict:
		return http.StatusConflict
	case domain.ErrorImageNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
