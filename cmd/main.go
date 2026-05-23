package main

import (
	"log/slog"

	app "github.com/Vlad-Ali/Virtual-chat/internal/app"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		slog.Error("Error creating application", "error", err)
	}

	err = application.Run()
	if err != nil {
		slog.Error("Error running application", "error", err)
	}
}
