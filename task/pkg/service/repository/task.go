package repository

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"time"
)

type FSTaskInterface interface {
	Create(ctx context.Context, in Task) (Task, error)
	Get(ctx context.Context, userID, taskID string) (Task, error)
	Update(ctx context.Context, newTask Task, userID, taskID string) (Task, error)
	Delete(ctx context.Context, userID, taskID string) error
	GetLastN(ctx context.Context, userID string, n int32) (tasks []Task, err error)
	GetExpired(ctx context.Context, userID string) (expiredTasks []Task, err error)
	SearchForExpiringTasks(ctx context.Context) (map[string][]Task, error)
}

type FSTask struct {
	fs     *firestore.CollectionRef
	client *firestore.Client
}

func NewFSTask(fs *firestore.CollectionRef, client *firestore.Client) *FSTask {
	return &FSTask{
		fs:     fs,
		client: client,
	}
}

func (f *FSTask) Create(ctx context.Context, in Task) (Task, error) {
	// sub collection logic
	docRef := f.fs.Doc(in.UserID).Collection(CollectionTasks).NewDoc()
	in.TaskID = docRef.ID
	in.CreatedAt = time.Now().Unix()
	// todo validation for time
	// todo validation of input strings -> max length of name , desc
	_, err := docRef.Set(ctx, in)
	if err != nil {
		return Task{}, err
	}
	// redundant data for optimization
	_, err = f.client.Collection(TaskList).Doc(docRef.ID).Set(ctx, in)
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
	// todo add field updatedAt instead of updating createdAt
	newTask.CreatedAt = time.Now().Unix()
	_, err := f.fs.Doc(userID).Collection(CollectionTasks).Doc(taskID).Set(ctx, newTask)
	if err != nil {
		return Task{}, err
	}
	// redundant data for optimization
	_, err = f.client.Collection(TaskList).Doc(taskID).Set(ctx, newTask)
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
	// redundant operation for optimization
	_, err = f.client.Collection(TaskList).Doc(taskID).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (f *FSTask) GetLastN(ctx context.Context, userID string, n int32) (tasks []Task, err error) {
	// todo -> put cap on a number of tasks returned?
	taskQuery := f.fs.Doc(userID).Collection(CollectionTasks).OrderBy("createdAt", firestore.Desc).Limit(int(n)).Documents(ctx)
	for {
		doc, err := taskQuery.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		task := Task{}
		err = doc.DataTo(&task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (f *FSTask) GetExpired(ctx context.Context, userID string) (expiredTasks []Task, err error) {
	taskQuery := f.fs.Doc(userID).Collection(CollectionTasks).Where("time", "<=", time.Now().Unix()).Documents(ctx)
	for {
		taskDoc, err := taskQuery.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		expiredTask := Task{}
		err = taskDoc.DataTo(&expiredTask)
		if err != nil {
			return nil, err
		}
		expiredTasks = append(expiredTasks, expiredTask)
	}
	return expiredTasks, nil
}

func (f *FSTask) SearchForExpiringTasks(ctx context.Context) (map[string][]Task, error) {
	toRemind := make(map[string][]Task)
	taskDocs, err := f.client.Collection(TaskList).
		Where("reminderSent", "==", false).
		Where("time", ">", time.Now().Unix()).
		Where("time", "<", time.Now().Add(time.Minute*5).Unix()).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	for _, taskDoc := range taskDocs {
		task := Task{}
		err := taskDoc.DataTo(&task)
		if err != nil {
			return nil, err
		}
		toRemind[task.UserEmail] = append(toRemind[task.UserEmail], task)
	}
	return toRemind, nil
}
