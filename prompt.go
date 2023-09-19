package main

import (
    _ "image/png"
	"os"

    "github.com/faiface/pixel"
    "github.com/faiface/pixel/pixelgl"
)

const (
	titlePathPNG = "assets/menus/splash01.png" // Path to save tiled map PNG file.
	continuePathPNG = "assets/menus/continue01.png" // Path to save tiled map PNG file.
	quitPathPNG = "assets/menus/quit01.png" // Path to save tiled map PNG file.
)

func showTitle(win *pixelgl.Window) (string, error) {
    var err error
    titleMapPNG, err := loadPicture(titlePathPNG)
    if err != nil {
        return "", err 
    }   

    title := pixel.NewSprite(titleMapPNG, titleMapPNG.Bounds())

	title.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	

	var option string
	for !win.Closed() {
		if win.Pressed(pixelgl.Key1) || win.Pressed(pixelgl.Key2) || win.Pressed(pixelgl.Key3) || win.Pressed(pixelgl.KeyQ) {
			if win.Pressed(pixelgl.Key1) {
				option = "New"
			} else if win.Pressed(pixelgl.Key2) {
		//		option, err := showSettings(win)
				if option == "" {
				}   
			} else if win.Pressed(pixelgl.Key3) {
				option = "About"
		//		option, err := showAbout(win)
				if option == "" {
				}   
			} else if win.Pressed(pixelgl.KeyQ) {
				option = "Quit"
				os.Exit(0)
			} else {
			}
			
			break
		}
		win.Update()
	}

    return option, nil
}

func showContinue(win *pixelgl.Window) (string, error) {
    var err error
    promptMapPNG, err := loadPicture(continuePathPNG)
    if err != nil {
        return "", err 
    }   

    prompt := pixel.NewSprite(promptMapPNG, promptMapPNG.Bounds())

	prompt.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	var response string
	for !win.Closed() {
		if win.Pressed(pixelgl.KeyY) || win.Pressed(pixelgl.KeyN) || win.Pressed(pixelgl.KeyQ) {
			if win.Pressed(pixelgl.KeyY) {
				response = "continue:yes"
				// start from where left off (respawn)
	 		} else if win.Pressed(pixelgl.KeyN) {
				// go back to title
				response = "continue:no"
			} else {
				os.Exit(0)
			}
			
			break
		}
		win.Update()
	}

    return response, nil
}

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

