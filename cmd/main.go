package main

import (
	"log"

	"lnkshrt/internal/app/routers"
	"lnkshrt/internal/domain/infrastructure/database"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initilization database:%v", err)
	}
	defer func() {
		if err := database.CloseDB(); err != nil {
			log.Printf("error closing database connection:%v", err)
		}
	}()

	r := routers.SetupRoute()

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
