package main

import (
	"github.com/sachatarba/rsoi_hotels/config"
	"github.com/sachatarba/rsoi_hotels/internal/payments/app"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Printf("Error loading config: %s. App exitting...", err)
		return
	}

	app.Run(cfg)
}
