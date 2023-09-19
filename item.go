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
	itemImageH	= 32
	itemImageW	= 32
	itemImageHTrim	= 8
	itemImageWTrim	= 5
)

type Item struct {
	id				uuid.UUID
	itemType		string
	imageH			float64 // image height in spritesheet
	imageW			float64 // image width in spriteshet
	imageHTrim		float64 // image height to trim (slim rect)
	imageWTrim		float64 // image width to trim (slim rect)
	X               float64 // X pos
	Y               float64 // Y pos
	pos             pixel.Vec
	paused			bool
	moving			bool
	itemFrames    []pixel.Rect
	itemFrame     pixel.Rect
	itemRect     pixel.Rect
	spritesheetPathMove string
	spritesheetDataMove pixel.Picture
	spritesheetPathIdle string
	spritesheetDataIdle pixel.Picture
	sprite          *pixel.Sprite
	direction       string
	directionCountN int
	directionCountS int
	directionCountW int
	directionCountE int
	Ase             *goaseprite.File
	Texture         *pixel.Sprite
	TextureRect     *pixel.Rect
	slice			*imgSlice
	amount			int
}

func newItem(itemType string, startX float64, startY float64, direction string) (*Item, error) {
	if debug { trace_func() }

	var err error
	var imageMovePath, imageJSONPath, imageIdlePath string
	item := Item{}
	item.id = uuid.New()
	item.itemType = itemType
	item.imageH = imageH
	item.imageW = imageW
    item.imageHTrim = itemImageHTrim
    item.imageWTrim = itemImageWTrim
	item.X = startX
	item.Y = startY
	item.pos = pixel.V(item.X, item.Y)
	item.direction = direction
/*
	XXX trim up image before load
        itemMinX := item.X - (item.imageW/2) + item.imageWTrim
        itemMinY := item.Y + (item.imageH/2)
        itemMaxX := item.X - (item.imageW/2) + item.imageW - item.imageWTrim
        itemMaxY := itemMinY - item.imageH
        item.itemRect = pixel.R(itemMinX, itemMinY, itemMaxX, itemMaxY)

        item.direction = direction
        item.itemFrame = item.getItemFrame()
        item.pos = pixel.V(item.X, item.Y)
        item.sprite = pixel.NewSprite(item.spritesheetDataMove, item.itemFrame)
*/
	if item.itemType == "chest" {
		imageMovePath = "assets/items/chest_spritesheet.png"
		imageJSONPath = "assets/items/chest_idle.json"
		imageIdlePath = "assets/items/chest_idle.png"
		item.amount = 100
	} else if item.itemType == "potion" {
		imageMovePath = "assets/items/potion_spritesheet.png"
		imageJSONPath = "assets/items/potion_idle.json"
		imageIdlePath = "assets/items/potion_idle.png"
		item.amount = 5
	} else {
		return nil, err
	}
	item.spritesheetPathMove = imageMovePath
	item.spritesheetDataMove, err = loadPicture(item.spritesheetPathMove)
	if err != nil {
		return nil, err
	}
	item.itemFrames, err = item.getItemFrames(item.spritesheetDataMove)
	if err != nil {
		return nil, err
	}
	item.Ase = goaseprite.Open(imageJSONPath)
	item.spritesheetPathIdle = imageIdlePath
	item.spritesheetDataIdle, err = loadPicture(item.spritesheetPathIdle)
	if err != nil {
		return nil, err
	}
	item.sprite = pixel.NewSprite(item.spritesheetDataIdle, pixel.R(0, 0, float64(item.Ase.FrameWidth), float64(item.Ase.FrameHeight)))
	item.Ase.Play("Idle")

	return &item, nil

}

func (item *Item) getItemFrames(spritesheetData pixel.Picture) ([]pixel.Rect, error) {
	if debug { trace_func() }

	// Store item graphics
	// 3, 7, 11 north
	// 2, 6, 10 south
	// 1, 5, 9 west
	// 0, 4, 8 east

	var itemFrames []pixel.Rect
	for x := spritesheetData.Bounds().Min.X; x < spritesheetData.Bounds().Max.X; x += item.imageW {
		for y := spritesheetData.Bounds().Min.Y; y < spritesheetData.Bounds().Max.Y; y += item.imageH {
			itemFrames = append(itemFrames, pixel.R(x, y, x + item.imageW, y + item.imageH))
		}
	}

	return itemFrames, nil
}

func (item *Item) getItemFrame() pixel.Rect {
	if debug { trace_func() }

	var itemFrame pixel.Rect
/*
	switch {
	case item.direction == "N":
		if item.directionCountN == 1 {
			itemFrame = item.itemFrames[7]
		} else if item.directionCountN == 2 {
			itemFrame = item.itemFrames[11]
		} else {
			itemFrame = item.itemFrames[3]
			item.directionCountN = 0
		}
	case item.direction == "S":
		if item.directionCountS == 1 {
			itemFrame = item.itemFrames[2]
		} else if item.directionCountS == 2 {
			itemFrame = item.itemFrames[6]
		} else {
			itemFrame = item.itemFrames[10]
			item.directionCountS = 0
		}
	case item.direction == "E":
		if item.directionCountE == 1 {
			itemFrame = item.itemFrames[0]
		} else if item.directionCountE == 2 {
			itemFrame = item.itemFrames[4]
		} else {
			itemFrame = item.itemFrames[8]
			item.directionCountE = 0
		}
	case item.direction == "W":
		if item.directionCountW == 1 {
			itemFrame = item.itemFrames[1]
		} else if item.directionCountW == 2 {
			itemFrame = item.itemFrames[5]
		} else {
			itemFrame = item.itemFrames[9]
			item.directionCountW = 0
		}
	}
*/

	itemFrame = item.itemFrames[0]

	return itemFrame
}

func (item *Item) move(direction string, debug bool) {
	if debug { trace_func() }

	if item.moving == true {
		switch {
		case direction == "N":
			item.Y++
			item.directionCountN++
		case direction == "S":
			item.Y--
			item.directionCountS++
		case direction == "E":
			item.X++
			item.directionCountE++
		case direction == "W":
			item.X--
			item.directionCountW++
		}
	
		itemMinX := item.X - (item.imageW/2) + item.imageWTrim
		itemMinY := item.Y + (item.imageH/2)
		itemMaxX := item.X - (item.imageW/2) + item.imageW - item.imageWTrim
		itemMaxY := itemMinY - item.imageH
		item.itemRect = pixel.R(itemMinX, itemMinY, itemMaxX, itemMaxY)
	
		item.direction = direction
		item.itemFrame = item.getItemFrame()
		item.pos = pixel.V(item.X, item.Y)
		item.sprite = pixel.NewSprite(item.spritesheetDataMove, item.itemFrame)
	}
}

func (item *Item) turnAround(direction string, debug bool) {
	if debug { trace_func() }

    switch {
        case direction == "N":
            item.direction = "S" 
        case direction == "S":
            item.direction = "N" 
        case direction == "E":
            item.direction = "W" 
        case direction == "W":
            item.direction = "E" 
    }   
    item.move(item.direction, debug)
}

func (item *Item) hide(dt float64, debug bool) {
	if debug { trace_func() }

    //item.itemRect = pixel.R(0, 0, 0, 0)
	//fmt.Printf("item done : %+v\n", item.slice)
	item.slice.Rect = pixel.R(0, 0, 0, 0)
	//fmt.Printf("item done : %+v\n", item.slice)
    item.itemFrame = item.getItemFrame()
    item.pos = pixel.V(0, 0)
    item.sprite = pixel.NewSprite(item.spritesheetDataMove, item.itemFrame)
	item.Update(dt, debug)
	if debug {
		fmt.Printf("item done : %+v\n", item.slice)
	}
}

func (item *Item) pause(debug bool) {
	if debug { trace_func() }

	item.paused = true
	item.moving = false
}

func (item *Item) stop(debug bool) {
	if debug { trace_func() }

	item.paused = false
	item.moving = false
}

func (item *Item) Update(dt float64, debug bool) {
	if debug { trace_func() }

	item.Ase.Update(float32(dt))

	// Set up the source rectangle for drawing the sprite (on the sprite sheet). File.GetFrameXY() will return the X and Y position
	// of the current frame of animation for the File.
	x, y := item.Ase.GetFrameXY()

	itemFrame := pixel.R(float64(x), float64(y), float64(x) + item.imageW, float64(y) + item.imageH)
	if debug {
		fmt.Printf("dt: %f\n", dt)
		fmt.Printf("Item Spritesheet X: %d\n", x)
		fmt.Printf("Item Spritesheet Y: %d\n", y)
	}
	item.sprite.Set(item.spritesheetDataIdle, itemFrame)
}

func (item *Item) Draw(win *pixelgl.Window) {
	if debug { trace_func() }

	item.sprite.Draw(win, pixel.IM.Moved(item.pos))
}
