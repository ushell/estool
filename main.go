package main

import (
	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"

	"estool/internal/backend"
)

func main() {
	// pack file
	js := mewn.String("./ui/dist/app.js")
	css := mewn.String("./ui/dist/app.css")

	//wails.BuildMode = "debug"

	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 568,
		Title:  "ESTool",
		JS:     js,
		CSS:    css,
	})

	app.Bind(backend.NewApp())

	app.Run()
}
