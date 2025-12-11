//Author: Design-BAB
//Date: 10-23-2025
//Description: Arrgh! 'Tis be me pirate game

package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	_ "github.com/glebarez/go-sqlite"
)

// the book's example made these global variables, for best practice I will make them a struct
type GameState struct {
	IsOver          bool
	Score           int
	NumberOfUpdates int
	XSize           int32
	HighScore       []*Scoreboard
	ScoreRecorded   bool
}

func newGameState() *GameState {
	initialScores := []*Scoreboard{}
	return &GameState{IsOver: false, Score: 0, NumberOfUpdates: 0, XSize: 800, HighScore: initialScores, ScoreRecorded: false}
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

type Scoreboard struct {
	Name  string
	Score int
}

func newScoreToBoard(name string, score int) *Scoreboard {
	return &Scoreboard{Name: name, Score: score}
}

// This function resets the object back to the beginning.
// The beginning being right side off screen
func placeIt(yourGame *GameState) float32 {
	return float32(int32(rand.IntN(int(yourGame.XSize))) + yourGame.XSize)
}

func CreateTable(db *sql.DB) (sql.Result, error) {
	sqlCommand := `CREATE TABLE IF NOT EXISTS pirates (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		score INTEGER NOT NULL);`
	return db.Exec(sqlCommand)
}

func draw(balloon, bird, house, tree *Actor, background rl.Texture2D, yourGame *GameState, db *sql.DB) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	rl.DrawTexture(background, 0, 0, rl.White)
	if yourGame.IsOver == false {
		rl.DrawTexture(balloon.Texture, int32(balloon.X), int32(balloon.Y), rl.White)
		rl.DrawTexture(bird.Texture, int32(bird.X), int32(bird.Y), rl.White)
		rl.DrawTexture(house.Texture, int32(house.X), int32(house.Y), rl.White)
		rl.DrawTexture(tree.Texture, int32(tree.X), int32(tree.Y), rl.White)
		rl.DrawText(strconv.Itoa(yourGame.Score), 10, 10, 24, rl.LightGray)
	} else {
		if yourGame.ScoreRecorded == false {
			var err error
			yourGame.HighScore, err = updateHighScore(db, *yourGame)
			if err != nil {
				log.Println("Arrgh, the database was unreachable")
			}
			yourGame.ScoreRecorded = true
		}
		displayHighScore(yourGame)
	}
	rl.EndDrawing()
}

func update(balloon, house, tree, bird *Actor, birdTextures *[2]rl.Texture2D, yourGame *GameState) {
	noHold := true
	if yourGame.IsOver == false {
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
		if rl.CheckCollisionRecs(balloon.Rectangle, bird.Rectangle) {
			fmt.Println("You hit the bird!")
		}
		//Coding Games in Py the book does bird first, but since it is the most complicated
		//I decided to put it towards the bottom of the function
		if house.X > -60 {
			house.X = house.X - 2
		} else {
			house.X = placeIt(yourGame)
			if rl.CheckCollisionRecs(tree.Rectangle, house.Rectangle) {
				tree.X += 30
			}
			yourGame.Score += 1
		}
		//Tree logic is here
		if tree.X > -80 {
			tree.X = tree.X - 2
		} else {
			tree.X = placeIt(yourGame)
			if rl.CheckCollisionRecs(tree.Rectangle, house.Rectangle) {
				tree.X += 30
			}
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
		balloon.X = rl.Clamp(balloon.X, 0.0, float32(yourGame.XSize)-balloon.Width)
		balloon.Y = rl.Clamp(balloon.Y, 0.0, float32(600)-balloon.Height)

		//handle collisions with obstacles
		if rl.CheckCollisionRecs(balloon.Rectangle, bird.Rectangle) ||
			rl.CheckCollisionRecs(balloon.Rectangle, house.Rectangle) ||
			rl.CheckCollisionRecs(balloon.Rectangle, tree.Rectangle) {
			yourGame.IsOver = true
		}
	}
}

// from step 9 & 20
func updateHighScore(db *sql.DB, yourGame GameState) ([]*Scoreboard, error) {
	var results []*Scoreboard
	//gonna add the current score into the data base
	now := time.Now()
	dateOfToday := now.Format("Monday, January 2, 2006")
	_, err := db.Exec(`INSERT INTO pirates (name, score) VALUES (?, ?)`, dateOfToday, yourGame.Score)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT name, score FROM pirates ORDER BY score DESC LIMIT 3;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var score int
		err := rows.Scan(&name, &score)
		if err != nil {
			log.Println("Error Scanning row: ", err)
			continue
		}
		results = append(results, newScoreToBoard(name, score))
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func displayHighScore(yourGame *GameState) {
	if len(yourGame.HighScore) == 0 {
		rl.DrawText("No high scores yet!", 190, 300, 20, rl.DarkGray)
		return
	}

	if yourGame.Score > yourGame.HighScore[0].Score {
		rl.DrawText("Wow! You made the high score!", 190, 250, 20, rl.DarkGray)
		rl.DrawText(strconv.Itoa(yourGame.Score), 190, 300, 20, rl.DarkGray)
	} else {
		var y int32 = 300
		for i, theScoreToDisplay := range yourGame.HighScore {
			mssg := strconv.Itoa(theScoreToDisplay.Score) + "    " + theScoreToDisplay.Name
			rl.DrawText(mssg, 190, y, 20, rl.DarkGray)
			y = int32(i+1)*20 + 300
		}
	}
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
	rl.InitWindow(yourGame.XSize, 600, "Sky Pirates!")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	//connect to SQLite database and getting scores sorted out
	db, err := sql.Open("sqlite", "./pirate.db?_pragma=foreign_keys(1)")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	//create the table just in case
	_, err = CreateTable(db)
	if err != nil {
		fmt.Println(err)
		return
	}
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

	balloon.X = rl.Clamp(balloon.X, 0.0, float32(yourGame.XSize)-balloon.Width)
	balloon.Y = rl.Clamp(balloon.Y, 0.0, float32(600)-balloon.Height)

	for !rl.WindowShouldClose() {
		draw(balloon, bird, house, tree, background, yourGame, db)
		update(balloon, house, tree, bird, &birdTextures, yourGame)
	}
}
