package repository

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"time"
)

type FSTaskInterface interface {
	Create(ctx context.Context, userID string, in Task) (Task, error)
	Get(ctx context.Context, userID, taskID string) (Task, error)
	Update(ctx context.Context, newTask Task, userID, taskID string) (Task, error)
	Delete(ctx context.Context, userID, taskID string) error
	GetLastN(ctx context.Context, userID string, n int32) (tasks []Task, err error)
	GetExpired(ctx context.Context, userID string) (expiredTasks []Task, err error)
	SearchForExpiringTasks(ctx context.Context) (map[string][]string, error)
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

func (f *FSTask) GetLastN(ctx context.Context, userID string, n int32) (tasks []Task, err error) {
	taskQuery := f.fs.Doc(userID).Collection(CollectionTasks).OrderBy("createdAt", firestore.Desc).Limit(int(n)).Documents(ctx)
	for {
		doc, err := taskQuery.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []Task{{}}, err
		}
		task := Task{}
		err = doc.DataTo(&task)
		if err != nil {
			return []Task{{}}, err
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
			return []Task{{}}, err
		}
		expiredTask := Task{}
		err = taskDoc.DataTo(&expiredTask)
		if err != nil {
			return []Task{{}}, nil
		}
		expiredTasks = append(expiredTasks, expiredTask)
	}
	return expiredTasks, nil
}

func (f *FSTask) SearchForExpiringTasks(ctx context.Context) (map[string][]string, error) {
	// todo more tasks that are expiring
	// iterate through all the users in user collection
	toRemind := make(map[string][]string)
	userDocs, err := f.fs.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	for _, userDoc := range userDocs {
		user := User{}
		err := userDoc.DataTo(&user)
		if err != nil {
			return nil, err
		}
		// lower and upper interval to avoid sending reminders to the same user more than once
		timeBuffLower := time.Now().Unix() + 300
		timeBuffUpper := time.Now().Unix() + 269 // inconsistent cron runs -> needs slightly longer interval
		taskDocs, err := userDoc.Ref.Collection(CollectionTasks).
			Where("time", ">", time.Now().Unix()).
			Where("time", "<=", timeBuffLower).
			Where("time", ">", timeBuffUpper).Documents(ctx).GetAll()
		if err != nil {
			return nil, err
		}
		for _, taskDoc := range taskDocs {
			task := Task{}
			err := taskDoc.DataTo(&task)
			if err != nil {
				return nil, err
			}
			toRemind[user.Email] = append(toRemind[user.Email], task.Name)
		}
	}
	return toRemind, nil
}
