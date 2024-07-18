package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/localopsco/go-sample/datastore"
	"github.com/localopsco/go-sample/handler"
	"github.com/localopsco/go-sample/service"
)

func main() {
	entClient, err := datastore.NewEntClient(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		log.Fatalf("error connecting to database: %w", err)
	}
	defer entClient.Close()

	taskStore := datastore.NewTaskStore(entClient)

	sdkConfig, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("S3_BUCKET_REGION")))
	if err != nil {
		log.Fatalf("Couldn't load aws session. Please set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables")
		return
	}
	s3Client := s3.NewFromConfig(sdkConfig)

	taskSvc := service.NewTaskService(taskStore, s3Client)
	handler := handler.NewHandler(taskSvc)
	router := gin.New()

	apiV1RouterGroup := router.Group("/api/v1/")
	apiV1RouterGroup.GET("/health/", handler.Health)
	apiV1RouterGroup.POST("/tasks/", handler.CreateTask)
	apiV1RouterGroup.GET("/tasks/", handler.ListTasks)
	apiV1RouterGroup.GET("/meta/", handler.GetMetaInfo)
	apiV1RouterGroup.GET("/tasks/:task_id/", handler.GetTask)
	apiV1RouterGroup.PATCH("/tasks/:task_id/", handler.UpdateTask)
	apiV1RouterGroup.DELETE("/tasks/:task_id/", handler.DeleteTask)
	apiV1RouterGroup.POST("/tasks/:task_id/attach/", handler.AddAttachment)
	apiV1RouterGroup.DELETE("/tasks/:task_id/attach/", handler.DeleteAttachment)

	router.Run(":" + os.Getenv("APP_PORT"))
}
