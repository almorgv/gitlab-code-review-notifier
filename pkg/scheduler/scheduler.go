package scheduler

import (
	"time"

	"github.com/jasonlvhit/gocron"

	"gitlab-code-review-notifier/pkg/log"
)

type Scheduler struct {
	config   Config
	instance *gocron.Scheduler
	log.Loggable
}

func NewScheduler(config Config) *Scheduler {
	instance := gocron.NewScheduler()
	instance.ChangeLoc(config.TimeZone)
	return &Scheduler{config: config, instance: instance}
}

func (s *Scheduler) Submit(job func()) error {
	if s.config.RepeatInterval != 0 {
		return s.SubmitInterval(job)
	}
	for _, t := range s.config.FixedTimes {
		if err := s.SubmitFixed(job, t); err != nil {
			return err
		}
	}
	return nil
}

func (s *Scheduler) SubmitInterval(job func()) error {
	return s.instance.
		Every(s.config.RepeatInterval).
		Minute().
		From(gocron.NextTick()).
		Do(s.workdayJobWrapper(job))
}

func (s *Scheduler) SubmitFixed(job func(), timeAt string) error {
	return s.instance.
		Every(1).
		Day().
		At(timeAt).
		Do(s.workdayJobWrapper(job))
}

func (s *Scheduler) workdayJobWrapper(job func()) func() {
	return func() {
		now := time.Now().In(s.config.TimeZone)
		if s.isWorkday(now) && s.isWorkHour(now) {
			job()
		}
	}
}

func (s *Scheduler) isWorkday(t time.Time) bool {
	return t.Weekday() != time.Saturday && t.Weekday() != time.Sunday
}

func (s *Scheduler) isWorkHour(t time.Time) bool {
	return t.Hour() >= s.config.WorkdayStartAt && t.Hour() < s.config.WorkdayEndAt
}

func (s *Scheduler) Run() {
	<-s.instance.Start()
}
