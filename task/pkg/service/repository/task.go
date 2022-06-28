package repository

import (
	"cloud.google.com/go/firestore"
	"context"
	"time"
)

type FSTaskInterface interface {
	Create(ctx context.Context, userID string, in Task) (Task, error)
	Get(ctx context.Context, userID, taskID string) (Task, error)
	Update(ctx context.Context, newTask Task, userID, taskID string) (Task, error)
	Delete(ctx context.Context, userID, taskID string) error
	//Update(ctx context.Context, newTask Task, docID string) (Task, error)
}

type FSTask struct {
	fs *firestore.CollectionRef
}

func NewFSTask(fs *firestore.CollectionRef) *FSTask {
	return &FSTask{
		fs: fs,
	}
}

func (f *FSTask) Create(ctx context.Context, userID string, in Task) (Task, error) {
	docRef := f.fs.Doc(userID).Collection(CollectionTasks).NewDoc()
	in.TaskID = docRef.ID
	in.CreatedAt = time.Now().Unix()
	//todo validation for time
	_, err := docRef.Set(ctx, in)
	if err != nil {
		return Task{}, err
	}
	return in, nil
}

func (f *FSTask) Get(ctx context.Context, userID, taskID string) (Task, error) {
	doc, err := f.fs.Doc(userID).Collection(CollectionTasks).Doc(taskID).Get(ctx)
	if err != nil {
		return Task{}, err
	}
	task := Task{}
	err = doc.DataTo(&task)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func (f *FSTask) Update(ctx context.Context, newTask Task, userID, taskID string) (Task, error) {
	newTask.CreatedAt = time.Now().Unix()
	_, err := f.fs.Doc(userID).Collection(CollectionTasks).Doc(taskID).Set(ctx, newTask)
	if err != nil {
		return Task{}, err
	}
	return newTask, nil
}

func (f *FSTask) Delete(ctx context.Context, userID, taskID string) error {
	_, err := f.fs.Doc(userID).Collection(CollectionTasks).Doc(taskID).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
