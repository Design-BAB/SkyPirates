//Author: Design-BAB
//Date: 10-23-2025
//Description: Arrgh! 'Tis be me pirate game
//Goal: Keep improving the game until it reaches 268 lines of code
//Notes: Finished all the coding suggestions on pg 129

package main

import (
	"math/rand/v2"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// the book's example made these global variables, for best practice I will make them a struct
type GameState struct {
	GameOver        bool
	Score           int
	NumberOfUpdates int
	WindowXSize     int32
}

func newGameState() *GameState {
	return &GameState{GameOver: false, Score: 0, NumberOfUpdates: 0, WindowXSize: 800}
}

type Actor struct {
	Texture rl.Texture2D
	//this is the collision box``
	rl.Rectangle // This gives Actor all the fields of rl.Rectangle (X, Y, Width, Height)
	Speed        float32
}

func newActor(texture rl.Texture2D, x, y float32) *Actor {
	return &Actor{Texture: texture, Rectangle: rl.Rectangle{X: x, Y: y, Width: float32(texture.Width), Height: float32(texture.Height)}, Speed: 7.0}
}

// This function resets the object back to the beginning.
// The beginning being right side off screen
func placeIt(yourGame *GameState) float32 {
	return float32(int32(rand.IntN(int(yourGame.WindowXSize))) + yourGame.WindowXSize)
}

func draw(balloon, bird, house, tree *Actor, background rl.Texture2D, yourGame *GameState) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	rl.DrawTexture(background, 0, 0, rl.White)
	if yourGame.GameOver == false {
		rl.DrawTexture(balloon.Texture, int32(balloon.X), int32(balloon.Y), rl.White)
		rl.DrawTexture(bird.Texture, int32(bird.X), int32(bird.Y), rl.White)
		rl.DrawTexture(house.Texture, int32(house.X), int32(house.Y), rl.White)
		rl.DrawTexture(tree.Texture, int32(tree.X), int32(tree.Y), rl.White)
		rl.DrawText(strconv.Itoa(yourGame.Score), 10, 10, 24, rl.LightGray)
	} else {
		displayHighScore()
	}
	rl.EndDrawing()
}

func update(balloon, house, tree, bird *Actor, birdTextures *[2]rl.Texture2D, yourGame *GameState) {
	noHold := true
	if yourGame.GameOver == false {
		//here we are doing the game controls
		if rl.IsKeyPressed(rl.KeyUp) && noHold {
			balloon.Y -= 43 //remember, up is down because the game engine's position layout
			noHold = false
		} else {
			balloon.Y += .5
		}

		if rl.IsKeyPressed(rl.KeyDown) {
			balloon.Y += 32
		}

		//Coding Games in Py the book does bird first, but since it is the most complicated
		//I decided to put it towards the bottom of the function
		if house.X > -60 {
			house.X = house.X - 2
		} else {
			house.X = placeIt(yourGame)
			yourGame.Score += 1
		}
		if tree.X > -80 {
			tree.X = tree.X - 2
		} else {
			tree.X = placeIt(yourGame)
			yourGame.Score += 1
		}
		//bird logic
		if bird.X > -10 {
			bird.X -= 4
			if yourGame.NumberOfUpdates == 9 {
				bird = flap(bird, birdTextures)
				yourGame.NumberOfUpdates = 0
			} else {
				yourGame.NumberOfUpdates += 1
			}
		} else {
			//this is else if the bird already past the balloon
			bird.X = placeIt(yourGame)
			bird.Y = float32(rand.IntN(200)) + 10
			yourGame.Score += 1
			yourGame.NumberOfUpdates = 0
		}
		//collision with the window
		balloon.X = rl.Clamp(balloon.X, 0.0, float32(yourGame.WindowXSize)-balloon.Width)
		balloon.Y = rl.Clamp(balloon.Y, 0.0, float32(600)-balloon.Height)

		//handle collisions with obstacles
		if rl.CheckCollisionRecs(balloon.Rectangle, bird.Rectangle) ||
			rl.CheckCollisionRecs(balloon.Rectangle, house.Rectangle) ||
			rl.CheckCollisionRecs(balloon.Rectangle, tree.Rectangle) {
			yourGame.GameOver = true
		}

	}
}

// from step 9 & 20
func updateHighScore() {

}

func displayHighScore() {
	rl.DrawText("No high scores yet!", 190, 300, 20, rl.LightGray)
}

func flap(bird *Actor, textures *[2]rl.Texture2D) *Actor {
	if bird.Texture == textures[0] {
		bird.Texture = textures[1]
	} else {
		bird.Texture = textures[0]
	}
	return bird
}

func main() {
	//creating the game
	yourGame := newGameState()
	//creating our window
	rl.InitWindow(yourGame.WindowXSize, 600, "Sky Pirates!")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	balloonTexture := rl.LoadTexture("images/balloon.png")
	defer rl.UnloadTexture(balloonTexture)
	balloon := newActor(balloonTexture, 110, 300)

	//in the bird, we need 2 textures to make the bird flap
	var birdTextures [2]rl.Texture2D
	birdTextures[0] = rl.LoadTexture("images/bird-up.png")
	defer rl.UnloadTexture(birdTextures[0])
	birdTextures[1] = rl.LoadTexture("images/bird-down.png")
	defer rl.UnloadTexture(birdTextures[1])
	bird := newActor(birdTextures[0], float32(rand.IntN(800))+800, float32(rand.IntN(190))+10)

	houseTexture := rl.LoadTexture("images/house.png")
	defer rl.UnloadTexture(houseTexture)
	house := newActor(houseTexture, placeIt(yourGame), 450)

	treeTexture := rl.LoadTexture("images/tree.png")
	defer rl.UnloadTexture(treeTexture)
	tree := newActor(treeTexture, placeIt(yourGame), 400)

	background := rl.LoadTexture("images/background.png")
	defer rl.UnloadTexture(background)

	balloon.X = rl.Clamp(balloon.X, 0.0, float32(yourGame.WindowXSize)-balloon.Width)
	balloon.Y = rl.Clamp(balloon.Y, 0.0, float32(600)-balloon.Height)

	for !rl.WindowShouldClose() {
		draw(balloon, bird, house, tree, background, yourGame)
		update(balloon, house, tree, bird, &birdTextures, yourGame)
	}
}
