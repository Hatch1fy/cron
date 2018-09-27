package main

import (
	"fmt"
	"time"

	"github.com/Hatch1fy/cron"
)

func main() {
	// Initialize new instance of cron
	c := cron.New()

	// After three seconds, display a success message
	c.Set("@ 3s", func() {
		fmt.Println("The server has been successfully started!")
	})

	// Every one minute, print the current time
	c.Set("every 1m", func() {
		fmt.Println("The current time is:", time.Now())
	})

	// Call empty select so we can keep the service open indefinitely
	select {}
}
