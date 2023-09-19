package main

import (
	"runtime"
	"encoding/json"
    "io/ioutil"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"time"
//	"github.com/faiface/pixel/imdraw"
//	"strconv"	
)

type Config struct {
	Config ConfigDef	`json:"Config"`
}

type ConfigDef struct {
	MapDir	string		`json:"MapDir"`
	Maps	[]MapDef	`json:"Maps"`
	Window	WindowDef	`json: "Window"`
}

type MapDef struct {
	Name 		string		`json: "name"`
	IsActive	bool		`json: "isactive"`
	Weather		string		`json: "weather"`
	Enemies		[]EnemyDef	`json: "enemies"`
	Items		[]ItemDef	`json: "items"`
}

type EnemyDef struct {
	Name		string	`json: "name"`
	IsMoving	bool	`json: "ismoving"`
	Direction	string	`json: "direction"`
	X			int		`json: "x"`
	Y			int		`json: "y"`
}

type ItemDef struct {
	Name		string	`json: "name"`
	Type		string	`json: "name"`
	IsMoving	bool	`json: "ismoving"`
	Direction	string	`json: "direction"`
	X			int		`json: "x"`
	Y			int		`json: "y"`
}

type WindowDef struct {
	Title		string	`json: "title"`
	Height		int		`json: "height"`
	Width		int		`json: "width"`
}

var (
	debug bool = false
	direction = "N"
)

type ImgSlices struct {
	ImgSlices	[]ImgSlice	`json:"Slices"`
}

type ImgSlice struct {
	Name	string		`json:"Name"`
	Data	string		`json:"Data"`
	Color	int			`json:"Color"`
	Keys 	[]ImgKey	`json:"Keys"`
}

type ImgKey struct {
	Frame	int	`json:"Frame"`
	X		int	`json:"X"`
	Y		int	`json:"Y"`
	W		int	`json:"W"`
	H		int	`json:"H"`
}

func trace_func() {
    pc := make([]uintptr, 15)
    n := runtime.Callers(2, pc)
    frames := runtime.CallersFrames(pc[:n])
    frame, _ := frames.Next()
    fmt.Printf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
}

func run() {

Start:
	// Read config
	file, err := ioutil.ReadFile("./game.conf")
	if err != nil {
		panic(err)
	}

	data := Config{}
	if debug {
		fmt.Printf("\n\nBASE DATA: %+v\n\n", data)
	}
 
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		panic(err)
	}

	if debug {
		fmt.Printf("JSON UNMARSHED DATA: %+v", data)
	}
/*
 
	for i := 0; i < len(data.Config.Maps); i++ {
		fmt.Println("Map Names: ", data.Config.Maps[i].Name)
		fmt.Println("Map Names: ", data.Config.Maps[i].isActive)
	}

    prettyJSON, err := json.MarshalIndent(data, "", "    ")
    if err != nil {
        panic(err)
    }
    fmt.Printf("JSON DATA: %s\n", string(prettyJSON))
*/

	// Create main game window
	cfg := pixelgl.WindowConfig{
		Title:  data.Config.Window.Title,
		Bounds: pixel.R(0, 0, float64(data.Config.Window.Width), float64(data.Config.Window.Height)),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	option, err := showTitle(win)
	if option == "" {
	}

	//player, err := newPlayer(win.Bounds().Center().X, win.Bounds().Center().Y)
	player, err := newPlayer(150, 50)
	if err != nil {
		panic(err)
	}

	weapon, err := newWeapon(player)
	if err != nil {
		panic(err)
	}

	var weapons []*Weapon
/*
	for i := 0; i < 10; i++ {
		weapon, err := newWeapon(player)
		if err != nil {
			panic(err)
		}
		weapons = append(weapons, weapon)
	}
*/
/*
var (
	sprite  *pixel.Sprite
)
	batch := pixel.NewBatch(&pixel.TrianglesData{}, weapon.spritesheetDataIdle)
	sprite = pixel.NewSprite(weapon.spritesheetDataIdle, weapon.spritesheetDataIdle.Bounds())
*/

	enemy, err := newEnemy(500, 350, "W")
	if err != nil {
		panic(err)
	}

	world, err := newWorld(win, data, player, weapon, win.Bounds().Center().X, win.Bounds().Center().Y)
	if err != nil {
		panic(err)
	}


//	fmt.Printf("SCENES: %+v", world.scenes)

/*
	err = loadBackgroundMusic("./assets/music/test2.mp3")
	if err != nil {
		panic(err)
	}
*/

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	// XXX 
	enemy.moving = true

	weapon.dead = true
//	var thisWeapon *Weapon
	last := time.Now()
	for !win.Closed() {

		scene := world.getCurrentScene()

		dt := time.Since(last).Seconds()
		last = time.Now()
		if debug {
			fmt.Printf("Win H: %f\n", win.Bounds().Center().Y)
			fmt.Printf("Win W: %f\n", win.Bounds().Center().X)
			fmt.Printf("Win H: %f\n", win.Bounds().H())
			fmt.Printf("Win W: %f\n", win.Bounds().W())
			fmt.Printf("background H: %f\n", scene.backgroundH)
			fmt.Printf("background W: %f\n", scene.backgroundW)
			fmt.Printf("background H/2: %f\n", scene.backgroundHHalf)
			fmt.Printf("background W/2: %f\n", scene.backgroundWHalf)
			fmt.Printf("background H/4: %f\n", scene.backgroundH/4)
			fmt.Printf("background W/4: %f\n", scene.backgroundW/4)
			fmt.Printf("player pos: %f\n", player.pos)
			fmt.Printf("player X: %f\n", player.X)
			fmt.Printf("player Y: %f\n", player.Y)

			fmt.Printf("Player Ase: %+v\n", player.Ase)
			s, _ := json.MarshalIndent(player.Ase, "", "\t")
			fmt.Printf("Player JSON: %s\n", s)

			fmt.Printf("World Ase: %+v\n", scene.Ase.Slices)
			t, _ := json.MarshalIndent(scene.Ase, "", "\t")
			fmt.Printf("World JSONe: %s\n", t)

			for i := 0; i < len(scene.Ase.Slices); i++ {
				fmt.Printf("Ase: %+v\n", scene.Ase.Slices[i])
				fmt.Printf("Ase Name: %+v\n", scene.Ase.Slices[i].Name)
				fmt.Printf("Ase Data: %+v\n", scene.Ase.Slices[i].Data)
				fmt.Printf("Ase Color: %+v\n", scene.Ase.Slices[i].Color)
				for j := 0; j < len(scene.Ase.Slices[i].Keys); j++ {
					fmt.Printf("Ase X: %+v\n", scene.Ase.Slices[i].Keys[j].X)
					fmt.Printf("Ase Y: %+v\n", scene.Ase.Slices[i].Keys[j].Y)
					fmt.Printf("Ase W: %+v\n", scene.Ase.Slices[i].Keys[j].W)
					fmt.Printf("Ase H: %+v\n", scene.Ase.Slices[i].Keys[j].H)
				}	
			}
		}

		moved := false
		//if win.Pressed(pixelgl.KeySpace) {
		//if win.JustPressed(pixelgl.KeySpace) {
		if win.JustReleased(pixelgl.KeySpace) {
			// XXX change this pass in a player, we want thwe weapon to be part of the player object
			//player.fire(player.direction, debug)
/*
			weapon, err := newWeapon(player)
			if err != nil {
				panic(err)
			}
*/
			weapon.fire(win, player, debug)
			weapons = append(weapons, weapon)
			//weapon.move(win, player.direction, debug)
		}

		if win.Pressed(pixelgl.KeyQ) {
			option, err := showQuit(win)
			if err != nil {
				panic(err)
			}
			if option == "" {
			}
		}

		if weapons != nil {
			for i := 0; i < len(weapons); i++ {
				if debug {
					fmt.Printf("\nWEAPONS %d: %+v\n\n", i, weapons[i])
				}
				if weapons[i].dead == false {
					//var thisDirection string
					//thisDirection := new(player.direction)
					weapons[i].move(win, scene, debug)
					// XXX DId we hit anything???
					if scene.checkWeaponCollision(win, weapon, debug) {
						if scene.lastCollisionObject.Type == "enemy" {
							if debug {
								fmt.Printf("BLAM!!!")
							}
							nextSceneData := scene.lastCollisionObject.Data
							nextSceneJSON := nextSceneData + ".json"
							nextScenePNG := nextSceneData + ".png"
							for i := 0; i < len(scene.enemies); i++ {
								if scene.lastCollisionObject.Id == scene.enemies[i].Id {
									scene.enemies[i].health = scene.enemies[i].decreaseHealth(50, debug)
									if scene.enemies[i].health == 0 && scene.enemies[i].dead == false {
										scene.enemies[i].die(debug)
										scene.enemies[i].Draw(win)
									}
									if debug {
										fmt.Printf("enemy health: %d", scene.enemies[i].health)
									}
									break;
								}
							}
							if debug {
								fmt.Printf("JSON: %s\n", nextSceneJSON);
								fmt.Printf("PNG: %s\n", nextScenePNG);
							}
						}
					}
					if debug {
						fmt.Printf("Sprite: %+v", weapon.sprite)
					}
					//sprite.Draw(batch, pixel.IM.Moved(weapons[i].pos))
					//weapons[i].move(win, thisDirection, debug)
					//fmt.Printf("DIRECTION: %s\n", thisDirection)
				} else {
					// Make it invisible
					if debug {
						fmt.Printf("STOP!")
					}
					weapons[i].stop(win, debug)
					weapons = append(weapons[:i], weapons[i+1:]...)
					//weapon.Update(dt, debug)
				}
			}
		}

//		batch.Draw(win)
/*
		fmt.Printf("\nthisWeapon:\n")
		fmt.Printf("%+v\n", thisWeapon)
		fmt.Printf("All Weapons:\n")
		fmt.Printf("%+v\n", weapons)
	var batch *pixel.Batch
		if weapons != nil {
			fmt.Printf("WEAPONS MOVE:\n")
			for i := 0; i < len(weapons); i++ {
				fmt.Printf("array weapon %d %+v\n", i, weapons[i])
				if weapons[i].dead == false {
					weapons[i].move(win, player.direction, debug)
	batch := pixel.NewBatch(&pixel.TrianglesData{}, weapons[i].spritesheetDataMove)
					weapons[i].sprite.Draw(batch, pixel.IM.Moved(weapons[i].pos))
		fmt.Printf("BATCH0: %+v", batch)
		batch.Draw(win)
			//		weapons[i].Draw(win)
					fmt.Printf("array weapon %d moved\n", i)
				}
				if weapons[i].dead == true {
					weapons = append(weapons[:i], weapons[i+1:]...)
				}
			}
		}
	//	batch.Draw(win)
		fmt.Printf("BATCH: %+v", batch)
*/

		// This would be useful for camera scrolling,we wont be doing any of that for this game
		//if ((math.Abs(player.X) < world.backgroundW) && (math.Abs(player.Y) < world.backgroundH) &&
		if (math.Abs(player.X) < win.Bounds().W()) && (math.Abs(player.Y) < win.Bounds().H()) &&
			(math.Abs(player.X) > 0) && (math.Abs(player.Y) > 0) {
			if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyRight) ||
				win.Pressed(pixelgl.KeyW) || win.Pressed(pixelgl.KeyS) || win.Pressed(pixelgl.KeyA) || win.Pressed(pixelgl.KeyD) {
				if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
					direction = "N"
				}
				if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
					direction = "S"
				}
				if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD) {
					direction = "E"
				}
				if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA) {
					direction = "W"
				}

				if scene.checkPlayerCollision(win, player, debug) {
					if scene.lastCollisionObject.Name == "SceneExit" {
						nextSceneData := scene.lastCollisionObject.Data
						nextSceneJSON := nextSceneData + ".json"
						nextScenePNG := nextSceneData + ".png"
						if debug {
							fmt.Printf("\n\nlastCollisionObject: %+v\n", scene.lastCollisionObject);
							fmt.Printf("JSON: %s\n", nextSceneJSON);
							fmt.Printf("PNG: %s\n\n", nextScenePNG);
						}
						//player.move(direction, debug)
						moved = true
						// XXX do stuff
						var BLAH bool
						BLAH, err = world.nextScene(win, player, weapon, enemy, nextSceneJSON, nextScenePNG, debug)
						if BLAH {
						}
					} else {
						if scene.lastCollisionObject.Type == "item" {
							player.pickupItem(dt, world, scene, scene.lastCollisionObject, debug)
							player.move(weapon, direction, debug)
						} else {
							if scene.lastCollisionObject.Type == "enemy" {
								for i := 0; i < len(scene.enemies); i++ {
									if scene.lastCollisionObject.Id == scene.enemies[i].Id {
										if scene.enemies[i].dead == false {
											player.health = player.decreaseHealth(scene.enemies[i].damage, debug)
										}
										break
									}
								}
							}
							if player.health == 0 && player.dead == false {
								player.die(debug)
								player.Draw(win)
/*
								option, err := showContinue(win)
								if err != nil {
									panic(err)
								} 
								if option == "" {
								}
*/
							} else {
								fmt.Printf("Jump Back!")
								player.jumpBack(weapon, player.direction, debug)
								player.turnAround(weapon, player.direction, debug)
							}
						}
					}
				} else {
					player.move(weapon, direction, debug)
					moved = true
				}
			}
		} else {
			player.turnAround(weapon, direction, debug)
		}

		// Play idle animation
		// Only uncomment this to get player idle animation when there is no key input
		if !moved && player.dead == false {
		//	player.Update(dt, debug)
		}
		if debug {
			fmt.Printf("Ase: %v\n", player.Ase.Animations)
		}

		if player.dead == false { 
			if player.direction == "N" {
				player.Ase.Play("Idle")
			}
		} else {
			player.Ase.Play("Dead")
		}

		err = scene.Update(win, player, weapon, debug)
		if err != nil {
			panic(err)
		}   

		err := world.Update(win, player)
		if err != nil {
			panic(err)
		}   

		if player.dead == true {
			fmt.Printf("player health: %d", player.health)
			response, err := showContinue(win)
			if err != nil {
				panic(err)
			}   
			if response == "continue:no" {
				goto Start
			}
		}   

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
