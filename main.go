package main

import (
	"context"
	"embed"
	"fmt"

	. "go-printease/internal"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	// Create an instance of the app structure
	app := NewApp()
	db := InitDB()
	defer func() {
		err := db.Disconnect(context.TODO())
		if err != nil {
			fmt.Println(err)
		}
	}()
	template := &Template{}

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
			template,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
