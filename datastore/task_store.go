package datastore

import (
	"context"

	"github.com/google/uuid"
	"github.com/localopsco/go-sample/ent"
	"github.com/localopsco/go-sample/ent/task"
	"github.com/localopsco/go-sample/models"
)

type TaskStore struct {
	client *ent.Client
}

func NewTaskStore(client *ent.Client) *TaskStore {
	return &TaskStore{
		client,
	}
}

func (store *TaskStore) CreateTask(task models.Task) (*models.Task, error) {
	entTask, err := store.client.Task.Create().
		SetTitle(task.Title).
		SetDescription(task.Description).
		SetIsCompleted(task.IsCompleted).
		Save(context.Background())
	if err != nil {
		return nil, err
	}
	return convertEntTask(entTask), nil
}

func (store *TaskStore) GetTask(taskID uuid.UUID) (*models.Task, error) {
	entTask, err := store.client.Task.Get(context.Background(), taskID)
	if err != nil {
		return nil, err
	}
	return convertEntTask(entTask), nil
}

func (store *TaskStore) DeleteTask(taskID uuid.UUID) error {
	return store.client.Task.DeleteOneID(taskID).Exec(context.Background())
}

func (store *TaskStore) ListTasks() ([]*models.Task, error) {
	entTasks, err := store.client.Task.Query().
		Order(ent.Desc(task.FieldCreatedAt)).
		All(context.Background())
	if err != nil {
		return nil, err
	}
	tasks := make([]*models.Task, 0, len(entTasks))
	for _, entTask := range entTasks {
		tasks = append(tasks, convertEntTask(entTask))
	}
	return tasks, nil
}

func (store *TaskStore) UpdateTask(task models.Task) (*models.Task, error) {
	entTask, err := store.client.Task.UpdateOneID(task.ID).
		SetTitle(task.Title).
		SetDescription(task.Description).
		SetIsCompleted(task.IsCompleted).
		Save(context.Background())
	if err != nil {
		return nil, err
	}
	return convertEntTask(entTask), nil
}
func (store *TaskStore) UpdateAttachmentURL(taskID uuid.UUID, url string) (*models.Task, error) {
	entTask, err := store.client.Task.UpdateOneID(taskID).
		SetAttachmentURL(url).
		Save(context.Background())
	if err != nil {
		return nil, err
	}
	return convertEntTask(entTask), nil
}

func convertEntTask(entTask *ent.Task) *models.Task {
	task := &models.Task{
		ID:          entTask.ID,
		Title:       entTask.Title,
		Description: entTask.Description,
		IsCompleted: entTask.IsCompleted,
		CreatedAt:   entTask.CreatedAt,
	}
	if entTask.AttachmentURL != "" {
		task.AttachmentURL = &entTask.AttachmentURL
	}
	return task
}
