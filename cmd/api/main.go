package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/localopsco/go-sample/datastore"
	"github.com/localopsco/go-sample/handler"
	"github.com/localopsco/go-sample/service"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	entClient, err := datastore.NewEntClient(
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_USERNAME"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_NAME"),
	)
	if err != nil {
		log.Fatalf("error connecting to database: %w", err)
	}
	defer entClient.Close()

	taskStore := datastore.NewTaskStore(entClient)

	sdkConfig, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(viper.GetString("S3_BUCKET_REGION")))
	if err != nil {
		log.Fatalf("Couldn't load aws session. Please set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables")
		return
	}
	s3Client := s3.NewFromConfig(sdkConfig)

	taskSvc := service.NewTaskService(taskStore, s3Client)
	handler := handler.NewHandler(taskSvc)
	router := gin.New()

	router.GET("/health", handler.Health)

	apiV1RouterGroup := router.Group("/api/v1/")
	apiV1RouterGroup.POST("/tasks", handler.CreateTask)
	apiV1RouterGroup.GET("/tasks", handler.ListTasks)
	apiV1RouterGroup.GET("/tasks/:task_id", handler.GetTask)
	apiV1RouterGroup.PATCH("/tasks/:task_id", handler.UpdateTask)
	apiV1RouterGroup.DELETE("/tasks/:task_id", handler.DeleteTask)
	apiV1RouterGroup.POST("/tasks/:task_id/attach", handler.AddAttachment)

	router.Run(":" + viper.GetString("APP_PORT"))
}
