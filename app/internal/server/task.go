package server

import (
	"context"

	"github.com/go-co-op/gocron"
	"go.uber.org/zap"

	"app/pkg/log"
)

type Task struct {
	log       *log.Logger
	scheduler *gocron.Scheduler
	//activityHandler *handler.ActivityHandler
}

func NewTask(log *log.Logger) *Task {
	return &Task{
		log: log,
		//activityHandler: activityHandler,
	}
}
func (t *Task) Start(ctx context.Context) error {
	gocron.SetPanicHandler(
		func(jobName string, recoverData interface{}) {
			t.log.Error("Task Panic", zap.String("job", jobName), zap.Any("recover", recoverData))
		},
	)

	// eg: crontab task
	//t.scheduler = gocron.NewScheduler(time.UTC)
	//_, err := t.scheduler.Every("5s").Do(
	//	func() {
	//		t.log.Info("Task1 start...")
	//		t.activityHandler.CheckActivity(ctx)
	//		t.log.Info("Task1 done.")
	//	},
	//)
	//if err != nil {
	//	t.log.Error("Task1 error", zap.Error(err))
	//}
	//
	//_, err = t.scheduler.Every("45s").Do(
	//	func() {
	//		t.log.Info("Task2 start...")
	//		_ = t.activityHandler.UpdateActivityCache(ctx)
	//		t.log.Info("Task2 done.")
	//	},
	//)
	//if err != nil {
	//	t.log.Error("Task2 error", zap.Error(err))
	//}
	//
	//// 定时每天凌晨0点执行
	//_, err = t.scheduler.Every(1).Day().At("00:00").Do(
	//	func() {
	//		t.log.Info("Task3 start...")
	//		_ = t.activityHandler.CollectEarnings(ctx)
	//		t.log.Info("Task3 done.")
	//	},
	//)
	t.scheduler.StartBlocking()
	return nil
}

func (t *Task) Stop(ctx context.Context) error {
	t.scheduler.Stop()
	t.log.Info("Task stop...")
	return nil
}
