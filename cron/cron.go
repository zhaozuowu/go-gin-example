package cron

import (
	"github.com/gin2/pkg/logging"
	"github.com/robfig/cron"
	"log"
	"time"
)

func init() {
	log.Println("Starting...")

	c := cron.New()
	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllTag...")
		logging.Info("Run models.CleanAllTag")
	})
	c.Start()
	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}