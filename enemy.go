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
//	"github.com/albertorestifo/dijkstra"
)

const (
	enemyImageH	= 64
	enemyImageW	= 64
	enemyImageHTrim	= 0
	enemyImageWTrim	= 16
	enemyVelocity = 3
)

type Enemy struct {
	Id				uuid.UUID
	imageH			float64 // image height in spritesheet
	imageW			float64 // image width in spriteshet
	imageHTrim		float64 // image height to trim (slim rect)
	imageWTrim		float64 // image width to trim (slim rect)
	X               float64 // X pos
	Y               float64 // Y pos
	pos             pixel.Vec
	velocity		int
	health			int
	healthMin       int
	healthMax       int
	dead			bool
	damage			int
	paused			bool
	moving			bool
	enemyFrames    []pixel.Rect
	enemyFrame     pixel.Rect
	enemyRect     pixel.Rect
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

func newEnemy(startX float64, startY float64, direction string) (*Enemy, error) {
	if debug { trace_func() }

	var err error
	enemy := Enemy{}
	enemy.Id = uuid.New()
	enemy.imageH = imageH
	enemy.imageW = imageW
    enemy.imageHTrim = imageHTrim
    enemy.imageWTrim = imageWTrim
	enemy.X = startX
	enemy.Y = startY
	enemy.pos = pixel.V(enemy.X, enemy.Y)
	enemy.velocity = enemyVelocity
	enemy.health = 100
	enemy.healthMin = 1
	enemy.healthMax = 100
	enemy.dead = false
	enemy.damage = 10
	enemy.direction = direction
	enemy.moving = false
	enemy.paused = false
	enemy.directionCountN = 0
	enemy.directionCountS = 0
	enemy.directionCountE = 0
	enemy.directionCountW = 0
	enemy.spritesheetPathMove = "assets/characters/ninja/64x64/ninja_spritesheet.png"
	enemy.spritesheetDataMove, err = loadPicture(enemy.spritesheetPathMove)
	if err != nil {
		return nil, err
	}
	enemy.enemyFrames, err = enemy.getEnemyFrames(enemy.spritesheetDataMove)
	if err != nil {
		return nil, err
	}
	enemy.Ase = goaseprite.Open("assets/characters/ninja/64x64/ninja_idle.json")
	enemy.spritesheetPathIdle = "assets/characters/ninja/64x64/ninja_idle.png"
	enemy.spritesheetDataIdle, err = loadPicture(enemy.spritesheetPathIdle)
	if err != nil {
		return nil, err
	}
	enemy.sprite = pixel.NewSprite(enemy.spritesheetDataIdle, pixel.R(0, 0, float64(enemy.Ase.FrameWidth), float64(enemy.Ase.FrameHeight)))
	enemy.Ase.Play("Idle")

	return &enemy, nil
}

func (enemy *Enemy) getEnemyFrames(spritesheetData pixel.Picture) ([]pixel.Rect, error) {
	if debug { trace_func() }

	// Store enemy graphics
	// 3, 7, 11 north
	// 2, 6, 10 south
	// 1, 5, 9 west
	// 0, 4, 8 east

	var enemyFrames []pixel.Rect
	for x := spritesheetData.Bounds().Min.X; x < spritesheetData.Bounds().Max.X; x += enemy.imageW {
		for y := spritesheetData.Bounds().Min.Y; y < spritesheetData.Bounds().Max.Y; y += enemy.imageH {
			enemyFrames = append(enemyFrames, pixel.R(x, y, x + enemy.imageW, y + enemy.imageH))
		}
	}

	return enemyFrames, nil
}

func (enemy *Enemy) getEnemyFrame() pixel.Rect {
	if debug { trace_func() }

	var enemyFrame pixel.Rect
	switch {
	case enemy.direction == "N":
		if enemy.directionCountN == 1 {
			enemyFrame = enemy.enemyFrames[7]
		} else if enemy.directionCountN == 2 {
			enemyFrame = enemy.enemyFrames[11]
		} else {
			enemyFrame = enemy.enemyFrames[3]
			enemy.directionCountN = 0
		}
	case enemy.direction == "S":
		if enemy.directionCountS == 1 {
			enemyFrame = enemy.enemyFrames[2]
		} else if enemy.directionCountS == 2 {
			enemyFrame = enemy.enemyFrames[6]
		} else {
			enemyFrame = enemy.enemyFrames[10]
			enemy.directionCountS = 0
		}
	case enemy.direction == "E":
		if enemy.directionCountE == 1 {
			enemyFrame = enemy.enemyFrames[0]
		} else if enemy.directionCountE == 2 {
			enemyFrame = enemy.enemyFrames[4]
		} else {
			enemyFrame = enemy.enemyFrames[8]
			enemy.directionCountE = 0
		}
	case enemy.direction == "W":
		if enemy.directionCountW == 1 {
			enemyFrame = enemy.enemyFrames[1]
		} else if enemy.directionCountW == 2 {
			enemyFrame = enemy.enemyFrames[5]
		} else {
			enemyFrame = enemy.enemyFrames[9]
			enemy.directionCountW = 0
		}
	}

	return enemyFrame
}

func (enemy *Enemy) move(direction string, debug bool) {
	if debug { trace_func() }

	if enemy.moving == true && enemy.dead == false {
		switch {

		case direction == "N":
			enemy.Y = enemy.Y + float64(enemy.velocity)
			enemy.directionCountN = enemy.directionCountN + enemy.velocity
		case direction == "S":
			enemy.Y = enemy.Y - float64(enemy.velocity)
			enemy.directionCountS = enemy.directionCountS - enemy.velocity
		case direction == "E":
			enemy.X = enemy.X + float64(enemy.velocity)
			enemy.directionCountE = enemy.directionCountE + enemy.velocity
		case direction == "W":
			enemy.X = enemy.X - float64(enemy.velocity)
			enemy.directionCountW = enemy.directionCountW - enemy.velocity
		}  
	
		enemyMinX := enemy.X - (enemy.imageW/2) + enemy.imageWTrim
		enemyMinY := enemy.Y + (enemy.imageH/2)
		enemyMaxX := enemy.X - (enemy.imageW/2) + enemy.imageW - enemy.imageWTrim
		//enemyMaxX := enemyMinX + enemy.imageW - enemy.imageWTrim
		enemyMaxY := enemyMinY - enemy.imageH
		enemy.enemyRect = pixel.R(enemyMinX, enemyMinY, enemyMaxX, enemyMaxY)
	
		enemy.direction = direction
		enemy.enemyFrame = enemy.getEnemyFrame()
		enemy.pos = pixel.V(enemy.X, enemy.Y)
		enemy.sprite = pixel.NewSprite(enemy.spritesheetDataMove, enemy.enemyFrame)
	}
}

func (enemy *Enemy) jumpBack(direction string, debug bool) {
	if debug { trace_func() }

	if enemy.moving == true && enemy.dead == false {
		switch {

		case direction == "N":
			enemy.Y = enemy.Y - float64(enemy.velocity)
			enemy.directionCountN = enemy.directionCountN - enemy.velocity
		case direction == "S":
			enemy.Y = enemy.Y + float64(enemy.velocity)
			enemy.directionCountS = enemy.directionCountS + enemy.velocity
		case direction == "E":
			enemy.X = enemy.X - float64(enemy.velocity)
			enemy.directionCountE = enemy.directionCountE - enemy.velocity
		case direction == "W":
			enemy.X = enemy.X + float64(enemy.velocity)
			enemy.directionCountW = enemy.directionCountW + enemy.velocity
		}  
	
		enemyMinX := enemy.X - (enemy.imageW/2) + enemy.imageWTrim
		enemyMinY := enemy.Y + (enemy.imageH/2)
		enemyMaxX := enemy.X - (enemy.imageW/2) + enemy.imageW - enemy.imageWTrim
		//enemyMaxX := enemyMinX + enemy.imageW - enemy.imageWTrim
		enemyMaxY := enemyMinY - enemy.imageH
		enemy.enemyRect = pixel.R(enemyMinX, enemyMinY, enemyMaxX, enemyMaxY)
	
		enemy.direction = direction
		enemy.enemyFrame = enemy.getEnemyFrame()
		enemy.pos = pixel.V(enemy.X, enemy.Y)
		enemy.sprite = pixel.NewSprite(enemy.spritesheetDataMove, enemy.enemyFrame)
	}
}

func (enemy *Enemy) turnAround(direction string, debug bool) {
	if debug { trace_func() }

    switch {
        case direction == "N":
            enemy.direction = "S" 
        case direction == "S":
            enemy.direction = "N" 
        case direction == "E":
            enemy.direction = "W" 
        case direction == "W":
            enemy.direction = "E" 
    }   
    enemy.move(enemy.direction, debug)
}

func (enemy *Enemy) pause(debug bool) {
	if debug { trace_func() }

	enemy.paused = true
	enemy.moving = false
}

func (enemy *Enemy) stop(debug bool) {
	if debug { trace_func() }

	enemy.paused = false
	enemy.moving = false
}

func (enemy *Enemy) trackPlayer(player *Player, debug bool) {
	if debug { trace_func() }

/*
	playerX := player.pos.X
	playerY := player.pos.Y
	enemyX := enemy.pos.X
	enemyY := enemy.pos.Y
*/

/*
	p = 150, 150
	e = 500, 350
	Graph{
		"e": {"b": 10, "c": 20},
		"p": {"a": 50},
	}
*/
}

func (enemy *Enemy) increaseHealth (amount int, debug bool) int {
	if debug { trace_func() }

	health := enemy.health + amount
	if health > enemy.healthMax {
		health = enemy.healthMax
	}

	return health
}

func (enemy *Enemy) decreaseHealth (amount int, debug bool) int {
	if debug { trace_func() }

	health := enemy.health - amount
	if health < enemy.healthMin {
		health = 0
	}

	return health
}

func (enemy *Enemy) die (debug bool) (bool, error)  {
	if debug { trace_func() }

	if debug {
		fmt.Printf("Enemy DIE DIE DIE!!!")
	}

	var err error
	enemy.dead = true

    enemy.Ase = goaseprite.Open("assets/characters/ninja/64x64/ninja_dead.json")
    enemy.spritesheetPathDead = "assets/characters/ninja/64x64/ninja_dead.png"
    enemy.spritesheetDataDead, err = loadPicture(enemy.spritesheetPathDead)
    if err != nil {
        return false, err 
    }   
    enemy.sprite = pixel.NewSprite(enemy.spritesheetDataDead, pixel.R(0, 0, float64(enemy.Ase.FrameWidth), float64(enemy.Ase.FrameHeight)))
    enemy.Ase.Play("Dead")

	return true, err
}

func (enemy *Enemy) Update(dt float64, debug bool) {
	if debug { trace_func() }

	enemy.Ase.Update(float32(dt))

	// Set up the source rectangle for drawing the sprite (on the sprite sheet). File.GetFrameXY() will return the X and Y position
	// of the current frame of animation for the File.
	x, y := enemy.Ase.GetFrameXY()

	enemyFrame := pixel.R(float64(x), float64(y), float64(x) + enemy.imageW, float64(y) + enemy.imageH)
	if debug {
		fmt.Printf("dt: %f\n", dt)
		fmt.Printf("Enemy Spritesheet X: %d\n", x)
		fmt.Printf("Enemy Spritesheet Y: %d\n", y)
	}
	enemy.sprite.Set(enemy.spritesheetDataIdle, enemyFrame)
}

func (enemy *Enemy) Draw(win *pixelgl.Window) {
	if debug { trace_func() }

	enemy.sprite.Draw(win, pixel.IM.Moved(enemy.pos))
}
