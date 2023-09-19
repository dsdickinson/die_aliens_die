package main

import (
	"fmt"
	//	"math"
	_ "image/png"
	uuid "github.com/google/uuid"
//	"image"
//    "image/color"
 //   "image/draw"

//"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	goaseprite "github.com/solarlune/GoAseprite"
)

const (
/*
	imageH	= 32
	imageW	= 32
*/
	imageH	= 64
	imageW	= 64
	imageHTrim	= 0
	imageWTrim	= 16
	playerVelocity = 4
)

// Player ...
type Player struct {
	Id              uuid.UUID
	imageH			float64 // image height in spritesheet
	imageW			float64 // image width in spriteshet
	imageHTrim		float64 // image height to trim (slim rect)
	imageWTrim		float64 // image width to trim (slim rect)
	X               float64 // X pos
	Y               float64 // Y pos
	pos             pixel.Vec
	velocity		int
	health			int
	healthMin		int
	healthMax		int
	score			int
	scoreMin		int
	dead			bool
	playerFrames    []pixel.Rect
	playerFrame     pixel.Rect
	playerRect     pixel.Rect
	spritesheetPathMove string
	spritesheetDataMove pixel.Picture
	spritesheetPathIdle string
	spritesheetDataIdle pixel.Picture
	spritesheetPathDead string
	spritesheetDataDead pixel.Picture
	sprite          *pixel.Sprite
	direction       string
	directionCountN int
	directionCountS int
	directionCountW int
	directionCountE int
	Ase             *goaseprite.File
	Texture         *pixel.Sprite
	TextureRect     *pixel.Rect
}

func newPlayer(startX float64, startY float64) (*Player, error) {
	if debug { trace_func() }

	var err error
	player := Player{}
	player.Id = uuid.New()
	player.imageH = imageH
	player.imageW = imageW
	player.imageHTrim = imageHTrim
	player.imageWTrim = imageWTrim
	player.X = startX
	player.Y = startY
	player.pos = pixel.V(player.X, player.Y)
	player.velocity = playerVelocity
	player.health = 100
	player.healthMin = 1
	player.healthMax = 100
	player.score = 0
	player.scoreMin = 0
	player.dead = false
	player.direction = "S"
	player.directionCountN = 0
	player.directionCountS = 0
	player.directionCountE = 0
	player.directionCountW = 0
//	player.spritesheetPathMove = "assets/characters/antborg/move/antborg_move.png"
//	player.spritesheetDataMove, err = loadPicture(player.spritesheetPathMove)
//	player.spritesheetPathMove = "assets/characters/grit/grit_spritesheet02.png"
	player.spritesheetPathMove = "assets/characters/grit/64x64/grit_spritesheet.png"
	player.spritesheetDataMove, err = loadPicture(player.spritesheetPathMove)
	if err != nil {
		return nil, err
	}
	player.playerFrames, err = player.getPlayerFrames(player.spritesheetDataMove)
	if err != nil {
		return nil, err
	}
//	player.Ase = goaseprite.Open("assets/characters/antborg/idle/antborg_idle.json")
//	player.spritesheetPathIdle = "assets/characters/antborg/idle/antborg_idle.png"
//	player.Ase = goaseprite.Open("assets/characters/grit/grit_idle.json")
//	player.spritesheetPathIdle = "assets/characters/grit/grit_idle.png"
	player.Ase = goaseprite.Open("assets/characters/grit/64x64/grit_idle.json")
	player.spritesheetPathIdle = "assets/characters/grit/64x64/grit_idle.png"
	player.spritesheetDataIdle, err = loadPicture(player.spritesheetPathIdle)
	if err != nil {
		return nil, err
	}
	player.sprite = pixel.NewSprite(player.spritesheetDataIdle, pixel.R(0, 0, float64(player.Ase.FrameWidth), float64(player.Ase.FrameHeight)))
	player.Ase.Play("Idle")

	return &player, nil
}

func (player *Player) getPlayerFrames(spritesheetData pixel.Picture) ([]pixel.Rect, error) {
	if debug { trace_func() }

	// Store player graphics
	// 3, 7, 11 north
	// 2, 6, 10 south
	// 1, 5, 9 west
	// 0, 4, 8 east

	var playerFrames []pixel.Rect
	for x := spritesheetData.Bounds().Min.X; x < spritesheetData.Bounds().Max.X; x += player.imageW {
		for y := spritesheetData.Bounds().Min.Y; y < spritesheetData.Bounds().Max.Y; y += player.imageH {
			playerFrames = append(playerFrames, pixel.R(x, y, x + player.imageW, y + player.imageH))
		}
	}

	return playerFrames, nil
}

func (player *Player) getPlayerFrame() pixel.Rect {
	if debug { trace_func() }

	var playerFrame pixel.Rect
	switch {
	case player.direction == "N":
		if player.directionCountN == 1 {
			playerFrame = player.playerFrames[7]
		} else if player.directionCountN == 2 {
			playerFrame = player.playerFrames[11]
		} else {
			playerFrame = player.playerFrames[3]
			player.directionCountN = 0
		}
	case player.direction == "S":
		if player.directionCountS == 1 {
			playerFrame = player.playerFrames[2]
		} else if player.directionCountS == 2 {
			playerFrame = player.playerFrames[6]
		} else {
			playerFrame = player.playerFrames[10]
			player.directionCountS = 0
		}
	case player.direction == "E":
		if player.directionCountE == 1 {
			playerFrame = player.playerFrames[0]
		} else if player.directionCountE == 2 {
			playerFrame = player.playerFrames[4]
		} else {
			playerFrame = player.playerFrames[8]
			player.directionCountE = 0
		}
	case player.direction == "W":
		if player.directionCountW == 1 {
			playerFrame = player.playerFrames[1]
		} else if player.directionCountW == 2 {
			playerFrame = player.playerFrames[5]
		} else {
			playerFrame = player.playerFrames[9]
			player.directionCountW = 0
		}
	}

	return playerFrame
}

func (player *Player) move(weapon *Weapon, direction string, debug bool) {
	if debug { trace_func() }

	if player.dead == false {
		switch {
		case direction == "N":
			player.Y = player.Y + float64(player.velocity)
			player.directionCountN = player.directionCountN + player.velocity
		case direction == "S":
			player.Y = player.Y - float64(player.velocity)
			player.directionCountS = player.directionCountS - player.velocity
		case direction == "E":
			player.X = player.X + float64(player.velocity)
			player.directionCountE = player.directionCountE + player.velocity
		case direction == "W":
			player.X = player.X - float64(player.velocity)
			player.directionCountW = player.directionCountW - player.velocity
		}
	
		playerMinX := player.X - (player.imageW/2) + player.imageWTrim
		playerMinY := player.Y + (player.imageH/2)
		playerMaxX := player.X - (player.imageW/2) + player.imageW - player.imageWTrim
		//playerMaxX := playerMinX + player.imageW - player.imageWTrim
		playerMaxY := playerMinY - player.imageH
		player.playerRect = pixel.R(playerMinX, playerMinY, playerMaxX, playerMaxY)
	
		//if ((math.Abs(player.X) < (world.backgroundW/4)) && (math.Abs(player.Y) < (world.backgroundH/4))) {
		//	if ((math.Abs(player.X) < 1024) && (math.Abs(player.Y) < 768)) {
		player.direction = direction
		player.playerFrame = player.getPlayerFrame()
		player.pos = pixel.V(player.X, player.Y)
		player.sprite = pixel.NewSprite(player.spritesheetDataMove, player.playerFrame)
	}
}

func (player *Player) jumpBack(weapon *Weapon, direction string, debug bool) {
	if debug { trace_func() }

	if player.dead == false {
		switch {
		case direction == "N":
			player.Y = player.Y - float64(player.velocity)
			player.directionCountN = player.directionCountN - player.velocity
		case direction == "S":
			player.Y = player.Y + float64(player.velocity)
			player.directionCountS = player.directionCountS + player.velocity
		case direction == "E":
			player.X = player.X - float64(player.velocity)
			player.directionCountE = player.directionCountE - player.velocity
		case direction == "W":
			player.X = player.X + float64(player.velocity)
			player.directionCountW = player.directionCountW + player.velocity
		}
	
		playerMinX := player.X - (player.imageW/2) + player.imageWTrim
		playerMinY := player.Y + (player.imageH/2)
		playerMaxX := player.X - (player.imageW/2) + player.imageW - player.imageWTrim
		//playerMaxX := playerMinX + player.imageW - player.imageWTrim
		playerMaxY := playerMinY - player.imageH
		player.playerRect = pixel.R(playerMinX, playerMinY, playerMaxX, playerMaxY)
	
		//if ((math.Abs(player.X) < (world.backgroundW/4)) && (math.Abs(player.Y) < (world.backgroundH/4))) {
		//	if ((math.Abs(player.X) < 1024) && (math.Abs(player.Y) < 768)) {
		player.direction = direction
		player.playerFrame = player.getPlayerFrame()
		player.pos = pixel.V(player.X, player.Y)
		player.sprite = pixel.NewSprite(player.spritesheetDataMove, player.playerFrame)
	}
}

func (player *Player) turnAround(weapon *Weapon, direction string, debug bool) {
	if debug { trace_func() }  

	switch {
		case direction == "N":
			player.direction = "S"
		case direction == "S":
			player.direction = "N"
		case direction == "E":
			player.direction = "W"
		case direction == "W":
			player.direction = "E"
	}
	player.move(weapon, player.direction, debug)
}

func (player *Player) pickupItem(dt float64, world *World, scene *Scene, slice *imgSlice, debug bool) {
	if debug { trace_func() }

/*
    playerMinX := player.X - (player.imageW/2) + player.imageWTrim
    playerMinY := player.Y + (player.imageH/2)
    playerMaxX := player.X - (player.imageW/2) + player.imageW - player.imageWTrim
    //playerMaxX := playerMinX + player.imageW - player.imageWTrim
    playerMaxY := playerMinY - player.imageH
    player.playerRect = pixel.R(playerMinX, playerMinY, playerMaxX, playerMaxY)

    //if ((math.Abs(player.X) < (world.backgroundW/4)) && (math.Abs(player.Y) < (world.backgroundH/4))) {
    //  if ((math.Abs(player.X) < 1024) && (math.Abs(player.Y) < 768)) {
    player.direction = direction
    player.playerFrame = player.getPlayerFrame()
    player.pos = pixel.V(player.X, player.Y)
    player.sprite = pixel.NewSprite(player.spritesheetDataMove, player.playerFrame)
*/
/*
	fmt.Printf("SLICE: %+v", slice)
	for i := 0; i < len(scene.items); i++ {
		if scene.items[i].id == slice.Id {
	slice.Rect = pixel.R(0,0,0,0)
	fmt.Printf("SLICE: %+v", slice)
			fmt.Printf("\n\n\nHIDE IT: %+v\n\n\n", scene.items[i])
			scene.items[i].hide(dt, debug)
//			scene.items[i].itemRect = pixel.R(0, 0, 0, 0)
			break
		}
	}
*/

	if debug {
		fmt.Printf("\nScene Items: %+v\n", scene.items)
	}

/*
			scene.items[1].itemRect = pixel.R(0, 0, 0, 0)
	for i := 0; i < len(scene.foreignCollisionObjects); i++ {
		fmt.Printf("i: %d\n", i)
		if scene.foreignCollisionObjects[i].Id == slice.Id {
			for j := 0; j < len(scene.items); j++ {
				if debug {
					fmt.Printf("item %d: %+v\n", j, scene.items[j])
				}
				if scene.items[j].id == slice.Id {
					//scene.items[j].itemRect = pixel.R(0, 0, 0, 0)
					scene.items[j].hide(dt, debug)
					if debug {
						fmt.Printf("item hidden: %+v\n", scene.items[j])
					}
					break;
				}
			}
			scene.foreignCollisionObjects[i].Rect = pixel.R(0,0,0,0)
			fmt.Printf("foreignCollisionObject %d: %+v\n", i, scene.foreignCollisionObjects[i])
			break;
		}
	}
*/

	var amount int = 0

	for j := 0; j < len(scene.items); j++ {
		if debug {
			fmt.Printf("item %d: %+v\n", j, scene.items[j])
		}
		if scene.items[j].id == slice.Id {
			scene.items[j].hide(dt, debug)
			amount = scene.items[j].amount
			if scene.items[j].itemType == "chest" {
				player.score = player.increaseScore(amount, debug)
			}
			if scene.items[j].itemType == "potion" {
				player.health = player.increaseHealth(amount, debug)
			}
			if debug {
				fmt.Printf("item done : %+v\n", scene.items[j])
			}
			break;
		}
	}

	//scene.doCollisions(player, scene.enemies, scene.items, debug)
	//panic(1)

	// XXX register in money total if chest
	//world.updateScore(win, 100)

/*
	slice.Rect = pixel.R(0,0,0,0)
	scene.items[i].Draw(win)
*/
//Rect: pixel.R(float64(itemMinX), itemMinY, float64(itemMaxX), itemMaxY),
}

// Reposition player in the new scene to the correct side of the screen
func (player *Player) sceneReposition(win *pixelgl.Window, scene *Scene, weapon *Weapon, direction string, debug bool) {
	if debug { trace_func() }

	player.directionCountN = 0
	player.directionCountS = 0
	player.directionCountE = 0
	player.directionCountW = 0

	switch {
		case direction == "N":
			//player.Y = player.Y - win.Bounds().H() + player.imageH + 40
			//player.Y = player.Y - win.Bounds().H() + adjustPlayer + 50
			//player.Y = player.Y - win.Bounds().H()
			player.Y = 0
			player.directionCountN++
		case direction == "S":
			//player.Y = player.Y + win.Bounds().H() - player.imageH - 40
			//player.Y = player.Y + win.Bounds().H() - player.imageH - adjustPlayer - 6
			//player.Y = win.Bounds().H() - player.imageH - adjustPlayer - 6
			//player.Y = win.Bounds().H() - adjustPlayer - 10
			player.Y = win.Bounds().H()
			if debug {
				fmt.Printf("imageH: %f\n", player.imageH)
				//fmt.Printf("adjustPlayer: %f\n", adjustPlayer)
				fmt.Printf("Y: %f\n\n", player.Y)
			}
			player.directionCountS++
		case direction == "E":
			player.X = player.X - win.Bounds().W() + player.imageW
			player.directionCountE++
		case direction == "W":
			player.X = player.X + win.Bounds().W() - player.imageW
			player.directionCountW++
	}

	player.move(weapon, direction, debug)

	// Adjust player in the new scene to get him out of the respective SceneExit zone
	var adjustPlayer float64 = 0
	for i := 0; i < len(scene.foreignCollisionObjects); i++ {
		if scene.foreignCollisionObjects[i].Name == "SceneExit" {
			if scene.checkPlayerSingleCollision(win, player, scene.foreignCollisionObjects[i], debug) {
				if debug {
					fmt.Printf("Slice: %+v\n", scene.foreignCollisionObjects[i])
				}
				if direction == "N" || direction == "S" {
					adjustPlayer = (player.imageH/2) + float64(scene.foreignCollisionObjects[i].H)
				} else {
					adjustPlayer = (player.imageH/2) + float64(scene.foreignCollisionObjects[i].W)
				}
				if debug {
					fmt.Printf("adjustPlayer: %f\n", adjustPlayer)
				}

				if direction == "N" {
					player.Y += adjustPlayer
				} else if direction == "S" {
					player.Y -= adjustPlayer
				} else if direction == "E" {
					player.X += adjustPlayer
				} else {
					player.X -= adjustPlayer
				}
				if debug {
					fmt.Printf("player adjusted: %f\n\n", player.Y)
				}

				break
			}
		}
	}

	player.move(weapon, direction, debug)

	// Make sure player doesn't get stuck in a collision when entering scene
	if scene.checkPlayerCollision(win, player, debug) {
		player.turnAround(weapon, player.direction, debug)

/*
		// Come back to this. Doesn't go back to previous scene because its not hitting exit properly
		if direction == "N" || direction == "S" {
			adjustPlayer = (player.imageH/2) + float64(scene.lastCollisionObject.H)
		} else {
			adjustPlayer = (player.imageH/2) + float64(scene.lastCollisionObject.W)
		}

		fmt.Printf("adjustPlayer: %f\n", adjustPlayer)

		if player.direction == "N" {
			player.Y += adjustPlayer
		} else if player.direction == "S" {
			player.Y -= adjustPlayer
		} else if player.direction == "E" {
			player.X += adjustPlayer
		} else {
			player.X -= adjustPlayer
		}
		fmt.Printf("adjusted: %f\n\n", player.Y)

		player.move(weapon, player.direction, debug)
*/
		player.move(weapon, player.direction, debug)
		player.move(weapon, player.direction, debug)
		player.move(weapon, player.direction, debug)
		player.move(weapon, player.direction, debug)
	}
}

func (player *Player) increaseHealth (amount int, debug bool) int {
	if debug { trace_func() }

	health := player.health + amount
	if health > player.healthMax {
		health = player.healthMax
	}

	return health
}

func (player *Player) decreaseHealth (amount int, debug bool) int {
	if debug { trace_func() }

	health := player.health - amount
	if health < player.healthMin {
		health = 0
	}

	return health
}

func (player *Player) increaseScore (amount int, debug bool) int {
	if debug { trace_func() }

	score := player.score + amount

	return score
}

func (player *Player) decreaseScore (amount int, debug bool) int {
	if debug { trace_func() }

	score := player.score - amount
	if score < player.scoreMin {
		score = 0
	}

	return score
}

func (player *Player) die (debug bool) (bool, error)  {
	if debug { trace_func() }

	if debug {
		fmt.Printf("Player DIE DIE DIE!!!")
	}

	var err error
	player.dead = true

    player.Ase = goaseprite.Open("assets/characters/grit/64x64/grit_dead.json")
    player.spritesheetPathDead = "assets/characters/grit/64x64/grit_dead.png"
    player.spritesheetDataDead, err = loadPicture(player.spritesheetPathDead)
    if err != nil {
        return false, err 
    }   
    player.sprite = pixel.NewSprite(player.spritesheetDataDead, pixel.R(0, 0, float64(player.Ase.FrameWidth), float64(player.Ase.FrameHeight)))
    player.Ase.Play("Dead")

	return true, err
}

/*
func (player *Player) fire (direction string, debug bool) {
	if debug { trace_func() }

	switch {
	case direction == "N":
		player.Y++
		player.directionCountN++
	case direction == "S":
		player.Y--
		player.directionCountS++
	case direction == "E":
		player.X++
		player.directionCountE++
	case direction == "W":
		player.X--
		player.directionCountW++
	}

	until (it hits a wall or leaves screen) {
	playerMinX := player.X - (player.imageW/2) + player.imageWTrim
	playerMinY := player.Y + (player.imageH/2)
	playerMaxX := player.X - (player.imageW/2) + player.imageW - player.imageWTrim
	//playerMaxX := playerMinX + player.imageW - player.imageWTrim
	playerMaxY := playerMinY - player.imageH
	player.playerRect = pixel.R(playerMinX, playerMinY, playerMaxX, playerMaxY)

	player.playerFrame = player.getPlayerFrame()
	player.pos = pixel.V(player.X, player.Y)
	player.sprite = pixel.NewSprite(player.spritesheetDataMove, player.playerFrame)
	player.sprite.Draw(win, pixel.IM.Moved(player.pos))
	}
}
*/


// This function is for when the player is idle.
func (player *Player) Update(dt float64, debug bool) {
	if debug { trace_func() }

	player.Ase.Update(float32(dt))

	// Set up the source rectangle for drawing the sprite (on the sprite sheet). File.GetFrameXY() will return the X and Y position
	// of the current frame of animation for the File.
	x, y := player.Ase.GetFrameXY()

	playerFrame := pixel.R(float64(x), float64(y), float64(x) + player.imageW, float64(y) + player.imageH)
	if debug {
		fmt.Printf("dt: %f\n", dt)
		fmt.Printf("Player Spritesheet X: %d\n", x)
		fmt.Printf("Player Spritesheet Y: %d\n", y)
	}
	player.sprite.Set(player.spritesheetDataIdle, playerFrame)
}

// Draw ...
func (player *Player) Draw(win *pixelgl.Window) {
	if debug { trace_func() }

	player.sprite.Draw(win, pixel.IM.Moved(player.pos))

	//upLeft := image.Point{int(player.playerRect.Min.X), int(player.playerRect.Min.Y)}
	//lowRight := image.Point{int(player.playerRect.Max.X), int(player.playerRect.Max.Y)}
/*
	imd := imdraw.New(nil)
	imd.Color = pixel.RGB(1, 0, 0)
	imd.Push(pixel.V(player.playerRect.Min.X, player.playerRect.Min.Y))
	imd.Color = pixel.RGB(0, 1, 0)
	imd.Push(pixel.V(player.playerRect.Max.X, player.playerRect.Max.Y))
	imd.Rectangle(2)
	imd.Draw(win)
*/

/*
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	cyan := color.RGBA{100, 200, 200, 0xff}
	draw.Draw(img, img.Bounds(), &image.Uniform{cyan}, image.ZP, draw.Src)
	player.sprite.Draw(win, pixel.IM.
                ScaledXY(pixel.ZV, pixel.V(
                        player.pos.X/player.sprite.Frame().W(),
                        player.pos.Y/player.sprite.Frame().H(),
                )).
                Moved(player.pos))
*/
}
