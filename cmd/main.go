package main

import (
	"fmt"
	// "log"
	// "nazartaraniuk/alertsProject/internal/app/client"
	"nazartaraniuk/alertsProject/internal/config"
	// "time"
)

const(
	CONFIG_PATH = "configs/config.dev.yaml"
)

func main() {
	cfg := config.LoadConfig(CONFIG_PATH)
	fmt.Println(cfg.Client.APIBaseURL)

	// client, err := client.NewClient(
	// 	cfg.Client.APIBaseURL,
	// 	time.Second,
	// 	cfg.Client.APIKey,
	// )


	// if err != nil {
	// 	log.Fatalf("Some error with creating Client, %v", err)
	// }
}
