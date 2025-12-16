# 🏴‍☠️ Sky Pirates!

A small arcade-style game written in **Go** using **raylib-go**, inspired by classic side-scrolling obstacle games like in Flappy Bird.\

---

## 🎮 Gameplay Overview

* Control a hot air balloon pirate ship flying across the screen
* Avoid collisions with:

  * Birds
  * Houses
  * Trees
* Each successfully passed obstacle increases your score
* The game ends immediately upon collision
* Top scores are saved locally using **SQLite**

---

## 🕹️ Controls

| Key | Action                |
| --- | --------------------- |
| ↑   | Boost balloon upward  |
| ↓   | Push balloon downward |

Gravity is always applied, so timing your movement is key.

---

## 🧱 Tech Info

* Written in **Go**
* Graphics and input handled via raylib-go
* Persistent high scores using SQLite
* Simple animation system for bird flapping
* Game state managed through a dedicated `GameState` struct

---

## 🗃️ High Score System

* Scores are stored in a local SQLite database (`pirate.db`)
* Each score entry is tagged with the date it was achieved
* Displays the top 3 scores on game over
* Automatically records the score once per run

---

## 🚀 How to Run

### Prerequisites

* Go 1.21+
* Raylib dependencies installed for your OS

###💻 Dependencies

```bash
go get github.com/gen2brain/raylib-go/raylib
go get github.com/glebarez/go-sqlite
---

## 📸 Screenshots

*(Yet to be added)*

Arrgh! Fair winds and high scores! 🏴‍☠️
