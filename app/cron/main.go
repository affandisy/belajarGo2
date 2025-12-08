package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()
	_, _ = c.AddFunc(" * * * *", func() {
		fmt.Println("1 menitan" + time.Now().Format("2006-01-02 15:04:05"))
	})

	_, _ = c.AddFunc("*/5 * * *", func() {
		fmt.Println("5 menitan" + time.Now().Format("2006-01-02 15:04:05"))
	})

	c.Start()
	select {}
}
