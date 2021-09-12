package main

import (
	"fmt"
	"os"
)

func main() {
	loadEnv()
	a := App{}
	a.Initialize(
		os.Getenv("HOST"),
		os.Getenv("APP_DB_PORT"),
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_NAME"),
	)
	a.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
