package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_imageHttpDeliveryHandler "github.com/nimasrn/uploadsvc-go/image/delivery/http"
	"github.com/nimasrn/uploadsvc-go/image/delivery/http/middleware"
	"github.com/nimasrn/uploadsvc-go/image/repository/redis"
	"github.com/nimasrn/uploadsvc-go/image/usecase"

	goredis "github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	envFile := "config.json"
	for _, v := range os.Args {
		if strings.Contains(v, "--env=") {
			s := strings.Split(v, "=")
			if _, err := os.Open(s[1] + ".json"); err != nil {
				log.Println("failed to open the passed env file, got error" + err.Error())
				break
			}
			envFile = s[1] + ".json"
		}
	}
	viper.SetConfigFile(envFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	if viper.GetBool(`release`) {
		gin.SetMode(gin.ReleaseMode)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	redisHost := viper.GetString(`redis.host`)
	redisPort := viper.GetString(`redis.port`)
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	redisOpt := &goredis.Options{
		Addr: redisAddr,
		DB:   0,
	}
	redisConn := goredis.NewClient(redisOpt)
	if cmd := redisConn.Ping(context.Background()); cmd.Err() != nil {
		log.Fatal(cmd.Err())
	}

	defer func() {
		err := redisConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	g := gin.Default()

	baseURI := viper.GetString("server.baseURI")
	m := middleware.InitMiddleware()
	g.Use(m.TransferEncodingCheck)
	rg := g.Group(baseURI)
	imageRepo := redis.NewRedisImageRepository(redisConn)

	storage := viper.GetString(`disk`)
	iu := usecase.NewImageUsecase(imageRepo, storage)
	_imageHttpDeliveryHandler.NewImageHandler(rg, iu, storage)

	serverAddr := viper.GetString("server.address")
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: g,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting")
}
