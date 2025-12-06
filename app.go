package main

import (
	"click-boy/click"
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx   context.Context
	click *click.Click
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.click = click.NewClick(ctx)
}

func (a *App) LogPrint(prints ...any) {
	for _, p := range prints {
		fmt.Println("======", p)
	}
}

func (a *App) ConnectDevice() {
	a.click.ConnectDevice()
}

func (a *App) ScreenShot() string {
	return a.click.ScreenShot()
}

func (a *App) StartClick(points []*click.Point, params *click.Params) {
	a.click.StartClick(points, params)
}

func (a *App) Pause() click.ClickStatus {
	return a.click.Pause()
}

func (a *App) Resume() click.ClickStatus {
	return a.click.Resume()
}

func (a *App) Stop() click.ClickStatus {
	return a.click.Stop()
}

func (a *App) Status() click.ClickStatus {
	return a.click.GetStatus()
}
