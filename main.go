package main

import (
  "github.com/leaanthony/mewn"
  "github.com/wailsapp/wails"

  "es2log/internal/ui"
)

func main() {
  // pack statics file
  js := mewn.String("./frontend/dist/app.js")
  css := mewn.String("./frontend/dist/app.css")

  wails.BuildMode = "debug"

  app := wails.CreateApp(&wails.AppConfig{
    Width:  1024,
    Height: 568,
    Title:  "ESTool",
    JS:     js,
    CSS:    css,
  })

  app.Bind(ui.NewApp())

  app.Run()
}
