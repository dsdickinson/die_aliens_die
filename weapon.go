package main

import (
	"fmt"
	"time"
	"math"
	_ "image/png"
//	"image"
//    "image/color"
 //   "image/draw"

//"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	goaseprite "github.com/solarlune/GoAseprite"
)

const (
	weaponImageH = 8
	weaponImageW = 8
	weaponImageHTrim = 1 
	weaponImageWTrim = 1
	weaponVelocity = 7
)

// Weapon ..
type Weapon struct {
	name			string
	weaponType		string
	imageH			float64 // image height in spritesheet
	imageW			float64 // image width in spriteshet
	imageHTrim		float64 // image height to trim (slim rect)
	imageWTrim		float64 // image width to trim (slim rect)
	X               float64 // X pos
	Y               float64 // Y pos
	pos             pixel.Vec
	velocity		int
	weaponFrames    []pixel.Rect
	weaponFrame     pixel.Rect
	weaponRect     pixel.Rect
	spritesheetPathMove string
	spritesheetDataMove pixel.Picture
	spritesheetPathIdle string
	spritesheetDataIdle pixel.Picture
	spritesheetPathDead string
	spritesheetDataDead pixel.Picture
	sprite          *pixel.Sprite
	direction       string
	newDirection    string
	directionCountN int
	directionCountS int
	directionCountW int
	directionCountE int
	Ase             *goaseprite.File
	Texture         *pixel.Sprite
	TextureRect     *pixel.Rect
	dead			bool
	active			bool
}

func newWeapon(player *Player) (*Weapon, error) {
	if debug { trace_func() }

	var err error
	weapon := Weapon{}
	weapon.name = "fireball"
	weapon.dead = true
	weapon.active = false
	weapon.weaponType = "projectile"
	weapon.imageH = weaponImageH
	weapon.imageW = weaponImageW
	weapon.imageHTrim = weaponImageHTrim
	weapon.imageWTrim = weaponImageWTrim
	weapon.X = player.pos.X
	weapon.Y = player.pos.Y
	weapon.pos = pixel.V(weapon.X, weapon.Y)
	weapon.velocity = weaponVelocity
	weapon.direction = player.direction
    weapon.newDirection = ""
    weapon.directionCountN = 0 
    weapon.directionCountS = 0 
    weapon.directionCountE = 0 
    weapon.directionCountW = 0 
	weapon.spritesheetPathMove = "assets/weapons/fireball_spritesheet.png"
	weapon.spritesheetDataMove, err = loadPicture(weapon.spritesheetPathMove)
	if err != nil {
		return nil, err
	}
	weapon.weaponFrames, err = weapon.getWeaponFrames(weapon.spritesheetDataMove)
	if err != nil {
		return nil, err
	}
	weapon.Ase = goaseprite.Open("assets/weapons/fireball_idle.json")
	weapon.spritesheetPathIdle = "assets/weapons/fireball_idle.png"
	weapon.spritesheetDataIdle, err = loadPicture(weapon.spritesheetPathIdle)
	if err != nil {
		return nil, err
	}
	weapon.sprite = pixel.NewSprite(weapon.spritesheetDataIdle, pixel.R(0, 0, float64(weapon.Ase.FrameWidth), float64(weapon.Ase.FrameHeight)))
	weapon.Ase.Play("Idle")

	weapon.spritesheetPathDead = "assets/weapons/fireball_dead.png"
	weapon.spritesheetDataDead, err = loadPicture(weapon.spritesheetPathDead)
	if err != nil {
		return nil, err
	}

	return &weapon, nil
}

func (weapon *Weapon) getWeaponFrames(spritesheetData pixel.Picture) ([]pixel.Rect, error) {
	if debug { trace_func() }

	var weaponFrames []pixel.Rect
	for x := spritesheetData.Bounds().Min.X; x < spritesheetData.Bounds().Max.X; x += weapon.imageW {
		if debug {
			fmt.Printf("x: %f", x)
		}
		for y := spritesheetData.Bounds().Min.Y; y < spritesheetData.Bounds().Max.Y; y += weapon.imageH {
			if debug {
				fmt.Printf("y: %f", y)
			}
			weaponFrames = append(weaponFrames, pixel.R(x, y, x + weapon.imageW, y + weapon.imageH))
		}
	}

	return weaponFrames, nil
}

func (weapon *Weapon) getWeaponFrame() pixel.Rect {
	if debug { trace_func() }

	var weaponFrame pixel.Rect
	weaponFrame = weapon.weaponFrames[0]
/*
	switch {
	case weapon.direction == "N":
		if weapon.directionCountN == 1 {
			weaponFrame = weapon.weaponFrames[7]
		} else if weapon.directionCountN == 2 {
			weaponFrame = weapon.weaponFrames[11]
		} else {
			weaponFrame = weapon.weaponFrames[3]
			weapon.directionCountN = 0
		}
	case weapon.direction == "S":
		if weapon.directionCountS == 1 {
			weaponFrame = weapon.weaponFrames[2]
		} else if weapon.directionCountS == 2 {
			weaponFrame = weapon.weaponFrames[6]
		} else {
			weaponFrame = weapon.weaponFrames[10]
			weapon.directionCountS = 0
		}
	case weapon.direction == "E":
		if weapon.directionCountE == 1 {
			weaponFrame = weapon.weaponFrames[0]
		} else if weapon.directionCountE == 2 {
			weaponFrame = weapon.weaponFrames[4]
		} else {
			weaponFrame = weapon.weaponFrames[8]
			weapon.directionCountE = 0
		}
	case weapon.direction == "W":
		if weapon.directionCountW == 1 {
			weaponFrame = weapon.weaponFrames[1]
		} else if weapon.directionCountW == 2 {
			weaponFrame = weapon.weaponFrames[5]
		} else {
			weaponFrame = weapon.weaponFrames[9]
			weapon.directionCountW = 0
		}
	}
*/

	return weaponFrame
}

func (weapon *Weapon) fire (win *pixelgl.Window, player *Player, debug bool) {
	if debug { trace_func() }

	// XXX Play blast sound HERE
	weapon.dead = false
	weapon.active = true
	weapon.X = player.X
	weapon.Y = player.Y
	weapon.pos = player.pos
	weapon.direction = player.direction
}

//func (weapon *Weapon) move (win *pixelgl.Window, direction string, debug bool) {
func (weapon *Weapon) move (win *pixelgl.Window, scene *Scene, debug bool) {
	if debug { trace_func() }

	switch {
	case weapon.direction == "N":
		weapon.Y = weapon.Y + float64(weapon.velocity)
		weapon.directionCountN++
	case weapon.direction == "S":
		weapon.Y = weapon.Y - float64(weapon.velocity)
		weapon.directionCountS++
	case weapon.direction == "E":
		weapon.X = weapon.X + float64(weapon.velocity)
		weapon.directionCountE++
	case weapon.direction == "W":
		weapon.X = weapon.X - float64(weapon.velocity)
		weapon.directionCountW++
	}
	
	weaponMinX := weapon.X - (weapon.imageW/2) + weapon.imageWTrim
	weaponMinY := weapon.Y + (weapon.imageH/2)
	weaponMaxX := weapon.X - (weapon.imageW/2) + weapon.imageW - weapon.imageWTrim
	weaponMaxY := weaponMinY - weapon.imageH
	weapon.weaponRect = pixel.R(weaponMinX, weaponMinY, weaponMaxX, weaponMaxY)
	
	//weapon.direction = direction
	weapon.weaponFrame = weapon.getWeaponFrame()
	weapon.pos = pixel.V(weapon.X, weapon.Y)
	weapon.sprite = pixel.NewSprite(weapon.spritesheetDataMove, weapon.weaponFrame)
	// XXX Add collision check which would also die
	if (scene.checkWeaponCollision(win, weapon, debug) ||
	    weapon.X >= win.Bounds().W() || weapon.Y >= win.Bounds().H() || 
		weapon.X <= 0 || weapon.Y <= 0) {
		weapon.stop(win, debug)
	}
}

func (weapon *Weapon) stop (win *pixelgl.Window, debug bool) {
	if debug { trace_func() }

	// $$XXX Play thud or ricochet sound HERE if its a collision?
	last := time.Now()
	dt := time.Since(last).Seconds()
	weapon.dead = true
	weapon.active = false
	weapon.Update(dt, debug)
}

/*
func (weapon *Weapon) fire (win *pixelgl.Window, direction string, debug bool) {
	if debug { trace_func() }

	batch := pixel.NewBatch(&pixel.TrianglesData{}, weapon.spritesheetDataMove)
	batch.Clear()
	for {
		switch {
		case direction == "N":
			weapon.Y++
			weapon.directionCountN++
		case direction == "S":
			weapon.Y--
			weapon.directionCountS++
		case direction == "E":
			weapon.X++
			weapon.directionCountE++
		case direction == "W":
			weapon.X--
			weapon.directionCountW++
		}
	
		weaponMinX := weapon.X - (weapon.imageW/2) + weapon.imageWTrim
		weaponMinY := weapon.Y + (weapon.imageH/2)
		weaponMaxX := weapon.X - (weapon.imageW/2) + weapon.imageW - weapon.imageWTrim
		weaponMaxY := weaponMinY - weapon.imageH
		weapon.weaponRect = pixel.R(weaponMinX, weaponMinY, weaponMaxX, weaponMaxY)
	
		weapon.direction = direction
		weapon.weaponFrame = weapon.getWeaponFrame()
		weapon.pos = pixel.V(weapon.X, weapon.Y)
		weapon.sprite = pixel.NewSprite(weapon.spritesheetDataMove, weapon.weaponFrame)
		//weapon.sprite.Draw(batch, pixel.IM.Moved(weapon.pos))
		weapon.Draw(win)
		//weapon.sprite.Draw(win, pixel.IM.Moved(weapon.pos))
		//win.Update()
	//	batch.Draw(win)
		if (weapon.X >= win.Bounds().W() || weapon.Y >= win.Bounds().H()) {
			break
		}  
	}
	win.Update()
}
*/

// This function is for when the player is idle.
func (weapon *Weapon) Update(dt float64, debug bool) {
	if debug { trace_func() }

	weapon.Ase.Update(float32(dt))

	// Set up the source rectangle for drawing the sprite (on the sprite sheet). File.GetFrameXY() will return the X and Y position
	// of the current frame of animation for the File.
	x, y := weapon.Ase.GetFrameXY()

	weaponFrame := pixel.R(float64(x), float64(y), float64(x) + weapon.imageW, float64(y) + weapon.imageH)
	if debug {
		fmt.Printf("dt: %f\n", dt)
		fmt.Printf("Weapon Spritesheet X: %d\n", x)
		fmt.Printf("Weapon Spritesheet Y: %d\n", y)
	}
	//weapon.sprite.Set(weapon.spritesheetDataIdle, weaponFrame)
	weapon.sprite.Set(weapon.spritesheetDataDead, weaponFrame)
}

// Draw ...
func (weapon *Weapon) Draw(win *pixelgl.Window, direction string) {
	if debug { trace_func() }

	var radians float64
	switch {
	case direction == "N":
		radians = 0
	case direction == "S":
		radians = 90 * (math.Pi / 2)
	case direction == "E":
		radians = 135 * (math.Pi / 2)
	case direction == "W":
		radians = 45 * (math.Pi / 2)
	}
	
	mat := pixel.IM
	mat = pixel.IM.Moved(weapon.pos)
	mat = mat.Rotated(weapon.pos, radians)
	weapon.sprite.Draw(win, mat)

	//weapon.sprite.Draw(win, pixel.IM.Moved(weapon.pos))

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
