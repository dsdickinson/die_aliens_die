package main

import (
    _ "image/png"
	"os"

    "github.com/faiface/pixel"
    "github.com/faiface/pixel/pixelgl"
)

const (
	quitPathPNG = "assets/menus/quit01.png" // Path to save tiled map PNG file.
)

func showQuit(win *pixelgl.Window) (string, error) {
    var err error
    quitMapPNG, err := loadPicture(quitPathPNG)
    if err != nil {
        return "", err 
    }   

    quit := pixel.NewSprite(quitMapPNG, quitMapPNG.Bounds())

	quit.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	var option string
	for !win.Closed() {
		if win.Pressed(pixelgl.KeyY) || win.Pressed(pixelgl.KeyN) {
			if win.Pressed(pixelgl.KeyN) {
				option = "No"
			}
			if win.Pressed(pixelgl.KeyY) {
				option = "Yes"
				os.Exit(0)
			}
			
			break
		}
		win.Update()
	}

    return option, nil
}

