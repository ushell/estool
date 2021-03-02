package main

import (
  "github.com/leaanthony/mewn"
  "github.com/wailsapp/wails"

  "estool/internal/backend"
)

func main() {
  // pack statics file
  js := mewn.String("./frontend/dist/app.js")
  css := mewn.String("./frontend/dist/app.css")

  wails.BuildMode = "debug"

  app := wails.CreateApp(&wails.AppConfig{
    Width:  1024,
    Height: 568,
    Title:  "ESTool v1.0.0",
    JS:     js,
    CSS:    css,
  })

  app.Bind(backend.NewApp())

  app.Run()
}
