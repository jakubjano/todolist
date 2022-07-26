package service

import (
	"context"
	cloudscheduler "google.golang.org/api/cloudscheduler/v1beta1"
	"strings"
)

const (
	jobDir = "/jobs/"
)

type Scheduler interface {
	CreateScheduledJob(ctx context.Context, name, schedule, target, method, description string) (*cloudscheduler.Job, error)
	PatchScheduledJob(ctx context.Context, name, schedule, description string) (*cloudscheduler.Job, error)
	ListScheduledJobs(ctx context.Context) ([]*cloudscheduler.Job, error)
	PauseScheduledJob(ctx context.Context, name string) (*cloudscheduler.Job, error)
	ResumeScheduledJob(ctx context.Context, name string) (*cloudscheduler.Job, error)
	DeleteScheduledJob(ctx context.Context, name string) error
	RunScheduledJob(ctx context.Context, name string) (*cloudscheduler.Job, error)
}

type SchedulerService struct {
	schedulerService    *cloudscheduler.Service
	pathToScheduledJobs string
}

func NewSchedulerService(schedulerService *cloudscheduler.Service, pathToScheduledJobs string) *SchedulerService {
	return &SchedulerService{
		schedulerService:    schedulerService,
		pathToScheduledJobs: pathToScheduledJobs,
	}
}

func (s *SchedulerService) CreateScheduledJob(ctx context.Context, name, schedule, target, method, description string) (*cloudscheduler.Job, error) {
	rb := &cloudscheduler.Job{
		Name:     s.pathToScheduledJobs + jobDir + name,
		Schedule: schedule,
		HttpTarget: &cloudscheduler.HttpTarget{
			Body:       "",
			Uri:        target,
			HttpMethod: strings.ToUpper(method),
		},
		Description: description,
		TimeZone:    "CET",
	}
	resp, err := s.schedulerService.Projects.Locations.Jobs.Create(s.pathToScheduledJobs, rb).Context(ctx).Do()
	if err != nil {
		return &cloudscheduler.Job{}, err
	}

	return resp, nil
}

func (s *SchedulerService) PatchScheduledJob(ctx context.Context, name, schedule, description string) (*cloudscheduler.Job, error) {
	rb := &cloudscheduler.Job{
		Name:        s.pathToScheduledJobs + jobDir + name,
		Schedule:    schedule,
		Description: description,
	}
	fullPath := s.pathToScheduledJobs + name
	resp, err := s.schedulerService.Projects.Locations.Jobs.Patch(fullPath, rb).Context(ctx).Do()
	if err != nil {
		return &cloudscheduler.Job{}, err
	}
	return resp, nil
}

func (s *SchedulerService) ListScheduledJobs(ctx context.Context) ([]*cloudscheduler.Job, error) {
	req := s.schedulerService.Projects.Locations.Jobs.List(s.pathToScheduledJobs)
	var jobs []*cloudscheduler.Job
	if err := req.Pages(ctx, func(page *cloudscheduler.ListJobsResponse) error {
		for _, job := range page.Jobs {
			jobs = append(jobs, job)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (s *SchedulerService) PauseScheduledJob(ctx context.Context, name string) (*cloudscheduler.Job, error) {
	rb := &cloudscheduler.PauseJobRequest{}
	fullPath := s.pathToScheduledJobs + jobDir + name
	resp, err := s.schedulerService.Projects.Locations.Jobs.Pause(fullPath, rb).Context(ctx).Do()
	if err != nil {
		return &cloudscheduler.Job{}, err
	}
	return resp, nil
}

func (s *SchedulerService) ResumeScheduledJob(ctx context.Context, name string) (*cloudscheduler.Job, error) {
	rb := &cloudscheduler.ResumeJobRequest{}
	fullPath := s.pathToScheduledJobs + jobDir + name

	resp, err := s.schedulerService.Projects.Locations.Jobs.Resume(fullPath, rb).Context(ctx).Do()
	if err != nil {
		return &cloudscheduler.Job{}, err
	}
	return resp, nil
}

func (s *SchedulerService) DeleteScheduledJob(ctx context.Context, name string) error {
	fullPath := s.pathToScheduledJobs + jobDir + name

	_, err := s.schedulerService.Projects.Locations.Jobs.Delete(fullPath).Context(ctx).Do()
	if err != nil {
		return err
	}
	return nil
}

func (s *SchedulerService) RunScheduledJob(ctx context.Context, name string) (*cloudscheduler.Job, error) {
	rb := &cloudscheduler.RunJobRequest{}
	fullPath := s.pathToScheduledJobs + jobDir + name

	resp, err := s.schedulerService.Projects.Locations.Jobs.Run(fullPath, rb).Context(ctx).Do()
	if err != nil {
		return &cloudscheduler.Job{}, err
	}
	return resp, nil
}
