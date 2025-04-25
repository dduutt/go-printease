package main

import (
	"embed"
	"log"

	"go-printease/internal"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	defer func() {
		err := internal.Close()
		if err != nil {
			log.Println("Error closing DB:", err)
		}
	}()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "go-printease",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []any{
			app,
			&internal.PrintRecord{},
			&internal.Template{},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
