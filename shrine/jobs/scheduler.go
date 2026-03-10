package jobs

import (
	"shrine/services"
	"shrine/utils/logger"

	"github.com/robfig/cron/v3"
)

var scheduler *cron.Cron

func Start() {
	scheduler = cron.New()

	scheduler.AddFunc("0 4 * * *", services.GenerateThumbnails)

	scheduler.Start()
	logger.Successf("Jobs", "Scheduler started")
}

func Stop() {
	if scheduler != nil {
		scheduler.Stop()
		logger.Infof("Jobs", "Scheduler stopped")
	}
}