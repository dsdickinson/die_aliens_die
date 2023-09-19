package main

import (
	"fmt"
	"math"
	_ "image/png"
	uuid "github.com/google/uuid"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	goaseprite "github.com/solarlune/GoAseprite"
//	"github.com/lafriks/go-tiled"
//	"github.com/lafriks/go-tiled/render"
)

const ( 
    mapPathJSON = "assets/maps/map00.json" // Path to save tiled map JSON file.
    mapPathPNG = "assets/maps/map00.png" // Path to save tiled map PNG file.
)

// Scene ...
type Scene struct {
	X float64 // Width
	Y float64 // Height
	pos pixel.Vec
	background *pixel.Sprite
	backgroundH float64
	backgroundW float64
	backgroundHHalf float64
	backgroundWHalf float64
	mapPathPNG		string
	mapPathJSON		string
    spritesheetPath string
    spritesheetData pixel.Picture
    sprite          *pixel.Sprite
	Ase             *goaseprite.File
    Texture         *pixel.Sprite
    TextureRect     *pixel.Rect
	itemCollisionObjects	[]*imgSlice
	foreignCollisionObjects	[]*imgSlice
	playerCollisionObjects	[]*imgSlice
	lastCollisionObject	*imgSlice
	IsActive		bool
	enemies			[]*Enemy
	items			[]*Item
}

func newScene(win *pixelgl.Window, world *World, mapPathPNG string, mapPathJSON string, IsActive bool, player *Player, enemies []*Enemy, items []*Item) (*Scene, error) {
	if debug { trace_func() }

	var err error
	scene := Scene{}
	scene.X = world.X
	scene.Y = world.Y
	scene.IsActive = IsActive
	scene.backgroundH = 0
	scene.backgroundW = 0
	scene.backgroundHHalf = 0
	scene.backgroundWHalf = 0
	scene.mapPathPNG = mapPathPNG
	scene.mapPathJSON = mapPathJSON
	scene.pos = pixel.V(scene.X, scene.Y)
	scene.background, err = scene.getBackground(mapPathJSON, mapPathPNG)
	if err != nil {
		return nil, err
	}
	scene.enemies = enemies
	scene.items = items

	scene.doCollisions(player, enemies, items, debug)

/*
	// Update scene
	err = scene.updateScene(win, player, enemy, weapon)
	if err != nil {
		return nil, err
	}
*/

	return &scene, nil
}

func (scene *Scene) getBackground(mapPathJSON string, mapPathPNG string) (*pixel.Sprite, error) {
	if debug { trace_func() }

	var err error
    gameMapPNG, err := loadPicture(mapPathPNG)
    if err != nil {
        return nil, err
    }

	scene.backgroundH = gameMapPNG.Bounds().H()
	scene.backgroundW = gameMapPNG.Bounds().W()
	scene.backgroundHHalf = scene.backgroundH / 2
	scene.backgroundWHalf = scene.backgroundW / 2

	// XXX need this?
    background := pixel.NewSprite(gameMapPNG, gameMapPNG.Bounds())

   // world.Ase = goaseprite.Open("assets/maps/map03.json")
   // world.spritesheetPath = "assets/maps/map03.png"
    scene.Ase = goaseprite.Open(mapPathJSON)
    scene.spritesheetPath = mapPathPNG
    scene.spritesheetData, err = loadPicture(scene.spritesheetPath)
    if err != nil {
        return nil, err
    }

    background = pixel.NewSprite(scene.spritesheetData, pixel.R(0, 0, float64(scene.Ase.FrameWidth), float64(scene.Ase.FrameHeight)))

	return background, nil
}

func (scene *Scene) doCollisions(player *Player, enemies []*Enemy, items []*Item, debug bool) {
	if debug { trace_func() }

	var thisEnemySlices []*imgSlice
	var thisItemSlices []*imgSlice
	var enemyCollisionObjects []*imgSlice
	var itemCollisionObjects []*imgSlice
    mapCollisionObjects := scene.getMapCollisionObjects(debug)
	for i := 0; i < len(enemies); i++ {
		thisEnemySlices = scene.getEnemyCollisionObjects(enemies[i], debug)
		for j := 0; j < len(thisEnemySlices); j++ {
			enemyCollisionObjects = append(enemyCollisionObjects, thisEnemySlices[j])
		}
	}
	for i := 0; i < len(items); i++ {
		thisItemSlices = scene.getItemCollisionObjects(items[i], debug)
		for j := 0; j < len(thisItemSlices); j++ {
			itemCollisionObjects = append(itemCollisionObjects, thisItemSlices[j])
		}
	}
	scene.foreignCollisionObjects = scene.setForeignCollisionObjects(mapCollisionObjects, enemyCollisionObjects, itemCollisionObjects, debug)
	scene.playerCollisionObjects = scene.getPlayerCollisionObjects(player, debug)
}

type imgSlice struct {
	Id		uuid.UUID
	Name	string
	Type	string
	Data 	string
	Color	int64
	X		int32
	Y		int32
	H		int32
	W		int32
	Rect	pixel.Rect
}

func (scene *Scene) getMapCollisionObjects(debug bool) []*imgSlice {
	if debug { trace_func() }

	var collisionObjects []*imgSlice

	for i := 0; i < len(scene.Ase.Slices); i++ {
		for j := 0; j < len(scene.Ase.Slices[i].Keys); j++ {
			//X := float64(scene.Ase.Slices[i].Keys[j].X)
			X := scene.Ase.Slices[i].Keys[j].X
			Y := float64(scene.Ase.Slices[i].Keys[j].Y)
			H := float64(scene.Ase.Slices[i].Keys[j].H)
			W := float64(scene.Ase.Slices[i].Keys[j].W)
			thisSlice := imgSlice{
				Id: uuid.New(),
				Name: scene.Ase.Slices[i].Name,
				Type: "map",
				Data: scene.Ase.Slices[i].Data,
				Color: scene.Ase.Slices[i].Color,
				X: X,
				Y: scene.Ase.Slices[i].Keys[j].Y,
				H: scene.Ase.Slices[i].Keys[j].H,
				W: scene.Ase.Slices[i].Keys[j].W,
				Rect: pixel.R(float64(X), scene.backgroundH-Y, float64(X)+W, scene.backgroundH-(Y+H)),
			}
            collisionObjects = append(collisionObjects, &thisSlice)
		}   
	}   

	if debug {
		for i := range(collisionObjects) {
			fmt.Printf("mapCollisionObjects: %+v\n", collisionObjects[i])
		}
	}

	return collisionObjects
}

func (scene *Scene) getEnemyCollisionObjects(enemy *Enemy, debug bool) []*imgSlice {
	if debug { trace_func() }

	var enemyCollisionObjects []*imgSlice

	for i := 0; i < len(enemy.Ase.Slices); i++ {
		for j := 0; j < len(enemy.Ase.Slices[i].Keys); j++ {
			X := enemy.Ase.Slices[i].Keys[j].X
			enemyMinX := enemy.X - (enemy.imageW/2) + enemy.imageWTrim
			enemyMinY := enemy.Y + (enemy.imageH/2)
			//enemyMaxX := enemyMinX + enemy.imageW
			enemyMaxX := enemy.X - (enemy.imageW/2) + enemy.imageW - enemy.imageWTrim
			enemyMaxY := enemyMinY - enemy.imageH

			thisSlice := imgSlice{
				Id: enemy.Id,
				Name: enemy.Ase.Slices[i].Name,
				Type: "enemy",
				Data: enemy.Ase.Slices[i].Data,
				Color: enemy.Ase.Slices[i].Color,
				X: X,
				Y: enemy.Ase.Slices[i].Keys[j].Y,
				H: enemy.Ase.Slices[i].Keys[j].H,
				W: enemy.Ase.Slices[i].Keys[j].W,
				//Rect: pixel.R(float64(X), scene.backgroundH-Y, float64(X)+W, scene.backgroundH-(Y+H)),
				Rect: pixel.R(float64(enemyMinX), enemyMinY, float64(enemyMaxX), enemyMaxY),
			}
			enemyCollisionObjects = append(enemyCollisionObjects, &thisSlice)
		}
	}

	if debug {
		for i := range(enemyCollisionObjects) {
			fmt.Printf("enemyCollisionObjects: %+v\n", enemyCollisionObjects[i])
		}
	}

	return enemyCollisionObjects
}

func (scene *Scene) getItemCollisionObjects(item *Item, debug bool) []*imgSlice {
	if debug { trace_func() }

	var itemCollisionObjects []*imgSlice

	for i := 0; i < len(item.Ase.Slices); i++ {
		for j := 0; j < len(item.Ase.Slices[i].Keys); j++ {
			X := item.Ase.Slices[i].Keys[j].X
			itemMinX := item.X - (item.imageW/2) + item.imageWTrim
			itemMinY := item.Y + (item.imageH/2)
			//itemMaxX := itemMinX + item.imageW
			itemMaxX := item.X - (item.imageW/2) + item.imageW - item.imageWTrim
			itemMaxY := itemMinY - item.imageH

			if item.slice == nil {
				thisSlice := imgSlice{
					Id: item.id,
					Name: item.Ase.Slices[i].Name,
					Type: "item",
					Data: item.Ase.Slices[i].Data,
					Color: item.Ase.Slices[i].Color,
					X: X,
					Y: item.Ase.Slices[i].Keys[j].Y,
					H: item.Ase.Slices[i].Keys[j].H,
					W: item.Ase.Slices[i].Keys[j].W,
					//Rect: pixel.R(float64(X), scene.backgroundH-Y, float64(X)+W, scene.backgroundH-(Y+H)),
					Rect: pixel.R(float64(itemMinX), itemMinY, float64(itemMaxX), itemMaxY),
				}
				item.slice = &thisSlice
			}
			//itemCollisionObjects = append(itemCollisionObjects, &thisSlice)
			itemCollisionObjects = append(itemCollisionObjects, item.slice)
		}
	}

	if debug {
		for i := range(itemCollisionObjects) {
			fmt.Printf("itemCollisionObjects: %+v\n", itemCollisionObjects[i])
		}
	}

	return itemCollisionObjects
}

func (scene *Scene) getPlayerCollisionObjects(player *Player, debug bool) []*imgSlice {
	if debug { trace_func() }

	var playerCollisionObjects []*imgSlice

	for i := 0; i < len(player.Ase.Slices); i++ {
		for j := 0; j < len(player.Ase.Slices[i].Keys); j++ {
			X := player.Ase.Slices[i].Keys[j].X
			playerMinX := player.X - (player.imageW/2) + player.imageWTrim
			playerMinY := player.Y + (player.imageH/2)
			//playerMaxX := playerMinX + player.imageW
			playerMaxX := player.X - (player.imageW/2) + player.imageW - player.imageWTrim
			playerMaxY := playerMinY - player.imageH

			thisSlice := imgSlice{
				Id: player.Id,
				Name: player.Ase.Slices[i].Name,
				Type: "player",
				Data: player.Ase.Slices[i].Data,
				Color: player.Ase.Slices[i].Color,
				X: X,
				Y: player.Ase.Slices[i].Keys[j].Y,
				H: player.Ase.Slices[i].Keys[j].H,
				W: player.Ase.Slices[i].Keys[j].W,
				//Rect: pixel.R(float64(X), scene.backgroundH-Y, float64(X)+W, scene.backgroundH-(Y+H)),
				Rect: pixel.R(float64(playerMinX), playerMinY, float64(playerMaxX), playerMaxY),
			}
			playerCollisionObjects = append(playerCollisionObjects, &thisSlice)
		}
	}

	if debug {
		for i := range(playerCollisionObjects) {
			fmt.Printf("playerCollisionObjects: %+v\n", playerCollisionObjects[i])
		}
	}

	return playerCollisionObjects
}

func (scene *Scene) setForeignCollisionObjects(mapCollisionObjects []*imgSlice, enemyCollisionObjects []*imgSlice, itemCollisionObjects []*imgSlice, debug bool) []*imgSlice {
	if debug { trace_func() }

	var collisionObjects []*imgSlice

	for i := 0; i < len(mapCollisionObjects); i++ {
		collisionObjects = append(collisionObjects, mapCollisionObjects[i])
	}

	for i := 0; i < len(enemyCollisionObjects); i++ {
		collisionObjects = append(collisionObjects, enemyCollisionObjects[i])
	}

	for i := 0; i < len(itemCollisionObjects); i++ {
		collisionObjects = append(collisionObjects, itemCollisionObjects[i])
	}

	if debug {
		for i := range(collisionObjects) {
			fmt.Printf("setForeignCollisionObjects: %+v\n", collisionObjects[i])
		}
	}

	return collisionObjects
}

func (scene *Scene) checkPlayerCollision(win *pixelgl.Window, player *Player, debug bool) bool {
	if debug { trace_func() }

	if debug {
		fmt.Printf("foreignCollisionObjects: %+v\n", scene.foreignCollisionObjects)
		fmt.Printf("playerRect: %+v\n", player.playerRect)
		fmt.Printf("player: %+v\n", player)
	}

	for i := 0; i < len(scene.foreignCollisionObjects); i++ {
		if debug {
			fmt.Printf("foreignCollisionObject %d: %+v\n", i, scene.foreignCollisionObjects[i])
		}

		if player.playerRect.Min.X < float64(scene.foreignCollisionObjects[i].Rect.Max.X) && 
		   player.playerRect.Max.X > float64(scene.foreignCollisionObjects[i].Rect.Min.X) && 
		   player.playerRect.Min.Y > float64(scene.foreignCollisionObjects[i].Rect.Max.Y) &&
		   player.playerRect.Max.Y < float64(scene.foreignCollisionObjects[i].Rect.Min.Y) {
			if debug {
				fmt.Printf("playerRect: %+v\n", player.playerRect)
				fmt.Printf("player: %+v\n", player)
			}
			scene.lastCollisionObject = scene.foreignCollisionObjects[i]
			return true
		}
		
	}
	return false
}

func (scene *Scene) checkWeaponCollision(win *pixelgl.Window, weapon *Weapon, debug bool) bool {
	if debug { trace_func() }

	if debug {
		fmt.Printf("foreignCollisionObjects: %+v\n", scene.foreignCollisionObjects)
		fmt.Printf("weaponRect: %+v\n", weapon.weaponRect)
		fmt.Printf("weapon: %+v\n", weapon)
	}

	for i := 0; i < len(scene.foreignCollisionObjects); i++ {
		if debug {
			fmt.Printf("foreignCollisionObject %d: %+v\n", i, scene.foreignCollisionObjects[i])
		}

		if weapon.weaponRect.Min.X < float64(scene.foreignCollisionObjects[i].Rect.Max.X) && 
		   weapon.weaponRect.Max.X > float64(scene.foreignCollisionObjects[i].Rect.Min.X) && 
		   weapon.weaponRect.Min.Y > float64(scene.foreignCollisionObjects[i].Rect.Max.Y) &&
		   weapon.weaponRect.Max.Y < float64(scene.foreignCollisionObjects[i].Rect.Min.Y) {
			if debug {
				fmt.Printf("weaponRect: %+v\n", weapon.weaponRect)
				fmt.Printf("weapon: %+v\n", weapon)
			}
			scene.lastCollisionObject = scene.foreignCollisionObjects[i]
			return true
		}
		
	}
	return false
}

func (scene *Scene) checkPlayerSingleCollision(win *pixelgl.Window, player *Player, slice *imgSlice, debug bool) bool {
	if debug { trace_func() }

	if debug {
		fmt.Printf("foreignCollisionObjects: %+v\n", scene.foreignCollisionObjects)
		fmt.Printf("playerRect: %+v\n", player.playerRect)
		fmt.Printf("player: %+v\n", player)
	}

	if player.playerRect.Min.X < float64(slice.Rect.Max.X) && 
	   player.playerRect.Max.X > float64(slice.Rect.Min.X) && 
	   player.playerRect.Min.Y > float64(slice.Rect.Max.Y) &&
	   player.playerRect.Max.Y < float64(slice.Rect.Min.Y) {
		if debug {
			fmt.Printf("playerRect: %+v\n", player.playerRect)
			fmt.Printf("player: %+v\n", player)
		}
		//scene.lastCollisionObject = scene.foreignCollisionObjects[i]
		return true
		
	}
	return false
}

func (scene *Scene) checkEnemySingleCollision(win *pixelgl.Window, enemy *Enemy, slice *imgSlice, debug bool) bool {
	if debug { trace_func() }

	if debug {
		for i := range(scene.foreignCollisionObjects) {
			fmt.Printf("foreignCollisionObjects: %+v\n", scene.foreignCollisionObjects[i])
		}
		fmt.Printf("enemyRect: %+v\n", enemy.enemyRect)
		fmt.Printf("enemy: %+v\n", enemy)
	}

	// Make sure we aren't checking collisions against the same objects
	if enemy.Id != slice.Id {
		if enemy.enemyRect.Min.X < float64(slice.Rect.Max.X) && 
		   enemy.enemyRect.Max.X > float64(slice.Rect.Min.X) && 
		   enemy.enemyRect.Min.Y > float64(slice.Rect.Max.Y) &&
		   enemy.enemyRect.Max.Y < float64(slice.Rect.Min.Y) {
			if debug {
				fmt.Printf("enemyRect: %+v\n", enemy.enemyRect)
				fmt.Printf("enemy: %+v\n", enemy)
			}
			//scene.lastCollisionObject = scene.foreignCollisionObjects[i]
			return true
		}
	}
	return false
}

func (scene *Scene) Update(win *pixelgl.Window, player *Player, weapon *Weapon, debug bool) (error) {
	if debug { trace_func() }

	//var err error

//	fmt.Println("SCENE UPDATE!")
//	world.background.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
/*
    err := loadBackgroundMusic("./assets/music/test2.mp3")
    if err != nil {
        panic(err)
    }   
*/

	for i := 0; i < len(scene.enemies); i++ {
		var enemyTurnAround int = 0
		for j := 0; j < len(scene.foreignCollisionObjects); j++ {
			if scene.checkEnemySingleCollision(win, scene.enemies[i], scene.foreignCollisionObjects[j], debug) {
				enemyTurnAround = 1
				break
			}
		}
		for j := 0; j < len(scene.playerCollisionObjects); j++ {
			if scene.checkEnemySingleCollision(win, scene.enemies[i], scene.playerCollisionObjects[j], debug) {
				if scene.enemies[i].dead == false {
					player.health = player.decreaseHealth(scene.enemies[i].damage, debug)
				}
				if player.health == 0 && player.dead == false { 
					player.die(debug)
					//player.Update(dt)
					player.Draw(win)
				}
				enemyTurnAround = 1
				break
			}
		}
        if (math.Abs(scene.enemies[i].X) < win.Bounds().W()) && (math.Abs(scene.enemies[i].Y) < win.Bounds().H()) &&
           (math.Abs(scene.enemies[i].X) > 0) && (math.Abs(scene.enemies[i].Y) > 0) && enemyTurnAround == 0 {
			scene.enemies[i].move(scene.enemies[i].direction, debug)
		} else {
			if debug {
				fmt.Printf("ENEMY TURN AROUND! %s", scene.enemies[i].direction)
			}
			// Do the following to avoid a collision loop
			scene.enemies[i].jumpBack(scene.enemies[i].direction, debug)
			scene.enemies[i].turnAround(scene.enemies[i].direction, debug)
		}
		scene.enemies[i].Draw(win)
	}

	for i := 0; i < len(scene.items); i++ {
		scene.items[i].Draw(win)
	}

	scene.doCollisions(player, scene.enemies, scene.items, debug)
	player.Draw(win)
	weapon.Draw(win, player.direction)
	win.Update()
	//player.sprite.Draw(win, pixel.IM.Moved(player.pos))
		
	return nil

}

