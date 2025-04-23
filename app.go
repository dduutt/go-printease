package main

import (
	"context"
	"encoding/json"
	"os/exec"
	"syscall"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) FileSelector(title, displayName, ext string) (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
		Filters: []runtime.FileFilter{
			{
				DisplayName: displayName,
				Pattern:     ext,
			},
		},
	})
}

func (a *App) GetPinters() ([]string, error) {
	// 执行powershell命令读取打印机名称列表，以json格式返回
	cmd := exec.Command("powershell", "-Command", "Get-Printer | Select-Object -Property Name | ConvertTo-Json")

	// 设置命令的创建标志以隐藏控制台窗口
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	type p struct {
		Name string `json:"Name"`
	}
	var printerNames []p

	if err := json.Unmarshal(out, &printerNames); err != nil {
		return nil, err
	}
	printers := make([]string, 0, len(printerNames))
	for _, p := range printerNames {
		printers = append(printers, p.Name)
	}
	return printers, nil

}
