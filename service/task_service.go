package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/localopsco/go-sample/datastore"
	"github.com/localopsco/go-sample/ent"
	"github.com/localopsco/go-sample/models"
	"github.com/spf13/viper"
)

const TaskNotFoundError = "Task not found"

type TaskService struct {
	store    *datastore.TaskStore
	s3Client *s3.Client
}

func NewTaskService(taskStore *datastore.TaskStore, s3Client *s3.Client) *TaskService {
	return &TaskService{
		taskStore,
		s3Client,
	}
}

func (svc *TaskService) CreateTask(title, desc string, isCompleted bool) (*models.Task, error) {
	task := models.Task{
		Title:       title,
		Description: desc,
		IsCompleted: isCompleted,
	}
	return svc.store.CreateTask(task)
}

func (svc *TaskService) ListTasks() ([]*models.Task, error) {
	return svc.store.ListTasks()
}

func (svc *TaskService) UpdateTask(taskID uuid.UUID, title, desc string, isCompleted bool) (*models.Task, error) {
	task := models.Task{
		ID:          taskID,
		Title:       title,
		Description: desc,
		IsCompleted: isCompleted,
	}
	updatedTask, err := svc.store.UpdateTask(task)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New(TaskNotFoundError)
		}
		return nil, err
	}
	return updatedTask, nil
}

func (svc *TaskService) GetTask(taskID uuid.UUID) (*models.Task, error) {
	task, err := svc.store.GetTask(taskID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New(TaskNotFoundError)
		}
		return nil, err
	}
	return task, nil
}

func (svc *TaskService) AddAttachment(taskID uuid.UUID, file *multipart.FileHeader) (*models.Task, error) {
	_, err := svc.store.GetTask(taskID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New(TaskNotFoundError)
		}
		return nil, err
	}
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("Error reading attachment file: %w", err)
	}
	defer src.Close()
	bucketName := viper.GetString("S3_BUCKET_NAME")
	key := filepath.Base(file.Filename)
	_, err = svc.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		ContentType: aws.String(file.Header.Get("Content-Type")),
		Body:        src,
		Key:         aws.String(key),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return nil, fmt.Errorf("Error uploading attachment to s3: %w", err)
	}
	attachmentURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, key)
	updatedTask, err := svc.store.UpdateAttachmentURL(taskID, attachmentURL)
	if err != nil {
		return nil, fmt.Errorf("Error while updating task's attachment URL: %w", err)
	}
	return updatedTask, nil
}

func (svc *TaskService) DeleteTask(taskID uuid.UUID) error {
	err := svc.store.DeleteTask(taskID)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New(TaskNotFoundError)
		}
		return err
	}
	return nil
}
