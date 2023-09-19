package main

import (
	"fmt"
	_ "image/png"
//	"strconv"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
//	"github.com/lafriks/go-tiled"
//	"github.com/lafriks/go-tiled/render"
)

// World ...
type World struct {
	X float64 // Width
	Y float64 // Height
	background *pixel.Sprite
	pos pixel.Vec
//	maps			map[string]bool
	scenes			[]*Scene
	currentScene	*Scene
}

func newWorld(win *pixelgl.Window, data Config, player *Player, weapon *Weapon, startX float64, startY float64) (*World, error) {
	if debug { trace_func() }

	var err error
	world := World{}
	world.X = startX
	world.Y = startY
	world.pos = pixel.V(world.X, world.Y)
	world.currentScene = nil
//	world.maps = worldMaps

	world.scenes, err = world.getScenes(win, player, data)
	if err != nil {
		return nil, err
	}

	// Update world
	err = world.Update(win, player)
	//err = world.Update(win, player, weapon)
	if err != nil {
		return nil, err
	}

	return &world, nil
}

func (world *World) getScenes(win *pixelgl.Window, player *Player, data Config) ([]*Scene, error) {
	if debug { trace_func() }

	var err error
	var scenes []*Scene
	var enemies []*Enemy
	var items []*Item

/*
	enemy, err := newEnemy(500, 350)
    if err != nil {
        panic(err)
    }   
	enemies = append(enemies, enemy)

	enemy2, err := newEnemy(200, 350)
    if err != nil {
        panic(err)
    }   
	enemies = append(enemies, enemy2)
*/

    for i := 0; i < len(data.Config.Maps); i++ {
        //fmt.Println("Map Names: ", data.Config.Maps[i].Name)
        //fmt.Println("Map Names: ", data.Config.Maps[i].IsActive)
		enemies = nil
		items = nil

		mapPathPNG := data.Config.MapDir + data.Config.Maps[i].Name + ".png"
		mapPathJSON := data.Config.MapDir + data.Config.Maps[i].Name + ".json"

    	for j := 0; j < len(data.Config.Maps[i].Enemies); j++ {
/*
			fmt.Println("X: ", data.Config.Maps[i].Enemies[j].X)
			fmt.Println("Y: ", data.Config.Maps[i].Enemies[j].Y)
*/
			enemy, err := newEnemy(float64(data.Config.Maps[i].Enemies[j].X), float64(data.Config.Maps[i].Enemies[j].Y), data.Config.Maps[i].Enemies[j].Direction)
			if data.Config.Maps[i].Enemies[j].IsMoving {
				if debug {
					fmt.Println("ENEMY MOVING!")
				}
				enemy.moving = true
				enemy.move(enemy.direction, debug)
			}
			if err != nil {
				panic(err)
			}   
			enemies = append(enemies, enemy)
		}

    	for j := 0; j < len(data.Config.Maps[i].Items); j++ {
/*
			fmt.Println("X: ", data.Config.Maps[i].Enemies[j].X)
			fmt.Println("Y: ", data.Config.Maps[i].Enemies[j].Y)
*/
			item, err := newItem(data.Config.Maps[i].Items[j].Type, float64(data.Config.Maps[i].Items[j].X), float64(data.Config.Maps[i].Items[j].Y), data.Config.Maps[i].Items[j].Direction)
			if data.Config.Maps[i].Items[j].IsMoving {
				if debug {
					fmt.Println("ITEM MOVING!")
				}
				item.moving = true
				item.move(item.direction, debug)
			}
			if err != nil {
				panic(err)
			}   
			items = append(items, item)
		}

		scene, err := newScene(win, world, mapPathPNG, mapPathJSON, data.Config.Maps[i].IsActive, player, enemies, items)
    	if err != nil {
        	return nil, err
		}
		//fmt.Printf("\nThe Scene: %+v\n\n", scene)
		if scene.IsActive {
			if debug {
				fmt.Printf("Scene IsActive!!: %+v", scene)
			}
			world.setCurrentScene(scene)
			world.setBackground(scene.background)
			//world.background = scene.background
		}
		scene.doCollisions(player, enemies, items, debug)
		scenes = append(scenes, scene)
	}

//	fmt.Printf("THESE SCENES: %+v", scenes)

/*
	for mapName, IsActive := range world.maps {
		mapPathPNG := mapDir + mapName + ".png"
		mapPathJSON := mapDir + mapName + ".json"
		scene, err := newScene(win, world, mapPathPNG, mapPathJSON, IsActive, enemies)
    	if err != nil {
        	return nil, err
		}
		if scene.IsActive {
			world.currentScene = scene
			world.background = scene.background
		}
		scenes = append(scenes, scene)
	}
*/
	//fmt.Printf("Current Scene: %+v", world.currentScene)

	return scenes, err
}

func (world *World) getCurrentScene() *Scene {
	if debug { trace_func() }

	return world.currentScene
}

func (world *World) setCurrentScene(scene *Scene) error {
	if debug { trace_func() }

	world.currentScene = scene

	return nil
}

func (world *World) getBackground() *pixel.Sprite {
	if debug { trace_func() }

	return world.background
}

func (world *World) setBackground(background *pixel.Sprite) error {
	if debug { trace_func() }

	world.background = background

	return nil
}

func (world *World) nextScene(win *pixelgl.Window, player *Player, weapon *Weapon, enemy *Enemy, nextMapPathJSON string, nextMapPathPNG string, debug bool) (bool, error) {
	if debug { trace_func() }

	for i, scene := range world.scenes {
		if debug {
			fmt.Printf("\nSCENE %d: %+v\n", i, scene)
		}
		if scene.mapPathPNG == nextMapPathPNG {
			//fmt.Printf("\nnextMapPathPNG: %s\n", nextMapPathPNG)
			world.scenes[i].IsActive = true
			world.setCurrentScene(scene)
			world.setBackground(scene.background)
			for j, _ := range world.scenes[i].enemies {
			//	fmt.Printf("\n\nENEMY\n: %+v", world.scenes[i].enemies[j])
				if world.scenes[i].enemies[j].paused {
					world.scenes[i].enemies[j].moving = true
					world.scenes[i].enemies[j].move(world.scenes[i].enemies[j].direction, debug)
				}
			}
			// Make sure we recalculate collisions since moving
			scene.doCollisions(player, world.scenes[i].enemies, world.scenes[i].items, debug)
			//fmt.Printf("\nTHIS SCENE %d: %+v\n", i, scene)
		} else {
			world.scenes[i].IsActive = false
			//fmt.Printf("\n\nENEMY SCENE\n: %+v", world.scenes[i])
			for j, _ := range world.scenes[i].enemies {
			//	fmt.Printf("\n\nENEMY\n: %+v", world.scenes[i].enemies[j])
				world.scenes[i].enemies[j].pause(debug)
			}
		}
	}

/*
	lastScene := world.currentScene
	lastScene.IsActive = false

	thisScene.IsActive = true
*/
	
	var err error
	/*
	world.background, err = scene.getBackground(scene.mapPathJSON, scene.mapPathPNG)
    if err != nil {
        return false, err
    }
	*/
//	world.background = thisScene.background


	player.sceneReposition(win, world.getCurrentScene(), weapon, player.direction, debug)

	// Update world
	err = world.Update(win, player)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (world *World) updateScore(win *pixelgl.Window, player *Player) {
	if debug { trace_func() }

/*
	distance := player.Y
	score := distance * -1 / 3
*/

/*
	if distance > 0 {
		score = 0
	}
*/

//	score := points

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	//basicTxt := text.New(scene.CameraPosition.Add(pixel.V(750, 725)), basicAtlas)
	basicTxt := text.New(pixel.V(30, 30), basicAtlas)
	basicTxt.Color = colornames.Black
	//fmt.Fprintf(basicTxt, "Score: %s\n", strconv.FormatFloat(player.score, 'f', 0, 64))
	fmt.Fprintf(basicTxt, "Score: %d\n", player.score)
	fmt.Fprintf(basicTxt, "Health: %d\n", player.health)
	fmt.Fprintf(basicTxt, "Level: %v\n", "1")
	basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))
}

func (world *World) Update(win *pixelgl.Window, player *Player) (error) {
	if debug { trace_func() }

	world.background.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
/*
    err := loadBackgroundMusic("./assets/music/test2.mp3")
    if err != nil {
        panic(err)
    }   
*/
	world.updateScore(win, player)
	//player.sprite.Draw(win, pixel.IM.Moved(player.pos))
		
	return nil
}

