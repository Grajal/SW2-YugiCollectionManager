package main

import (
	"fmt"
	"os"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/internal/router"
)

var port = os.Getenv("PORT")

func main() {
	if port == "" {
		port = "8080"
	}

	e := router.New()
	fmt.Println("Starting server on port 8080...")

	e.Logger.Fatal(e.Start(":" + port))
}
