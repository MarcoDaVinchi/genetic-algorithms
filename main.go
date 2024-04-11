package main

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	squareSize = 64
)

type CellType int64
type CellImages map[CellType]rl.Image
type AssetTextures map[CellType]rl.Texture2D

const (
	ActiveCell CellType = iota
	RootCell
	LeafCell
)

type CellDirection int8

const (
	Up CellDirection = iota
	UR
	Right
	RD
	Down
	DL
	Left
	LU
)

// Cell type
type Cell struct {
	Position rl.Vector2
	Size     rl.Vector2
	Alive    bool
	Next     bool
	Visited  bool
}

// Game type
type Game struct {
	ScreenWidth   int32
	ScreenHeight  int32
	Cols          int32
	Rows          int32
	FramesCounter int32
	Playing       bool
	Cells         [][]*Cell
}

func main() {
	game := Game{}
	game.Init(false)

	rl.InitWindow(game.ScreenWidth, game.ScreenHeight, "Cell automata sim")
	defer rl.CloseWindow()

	var camera rl.Camera2D
	camera.Zoom = 1

	//botCell := rl.LoadTexture("assets.Basic_cell.png")
	//position := rl.NewVector2(350.0, 280.9)
	//frameRec := rl.NewRectangle(0, 0, float32(botCell.Width), float32(botCell.Height))
	//currentFrame := float32(0)
	//
	rl.SetTargetFPS(20)

	basicCellImage := rl.LoadImage("assets/basic_cell.png")
	if basicCellImage == nil {
		fmt.Println("Failed to load basicCellImage.")
		os.Exit(1)
	}
	defer rl.UnloadImage(basicCellImage)

	cropX := 0
	cropY := 30
	cropWidth := 32
	cropHeight := 32

	croppedImage := rl.ImageFromImage(*basicCellImage, rl.Rectangle{
		X:      float32(cropX),
		Y:      float32(cropY),
		Width:  float32(cropWidth),
		Height: float32(cropHeight),
	})

	croppedTexture := rl.LoadTextureFromImage(&croppedImage)
	defer rl.UnloadTexture(croppedTexture)

	for !rl.WindowShouldClose() {
		if game.Playing {
			game.Update()
		}

		game.Input(&camera)
		game.Draw(&camera, &croppedTexture)
	}

}

// Init - Initialize game
func (g *Game) Init(clear bool) {
	g.ScreenWidth = 800
	g.ScreenHeight = 450
	g.FramesCounter = 0

	g.Cols = g.ScreenWidth / squareSize
	g.Rows = g.ScreenHeight / squareSize

	g.Cells = make([][]*Cell, g.Cols+1)
	for i := int32(0); i <= g.Cols; i++ {
		g.Cells[i] = make([]*Cell, g.Rows+1)
	}

	/*for x := int32(0); x <= g.Cols; x++ {
		for y := int32(0); y <= g.Rows; y++ {
			g.Cells[x][y] = &Cell{}
			g.Cells[x][y].Position = rl.NewVector2(float32(x)*squareSize, (float32(y)*squareSize)+1)
			g.Cells[x][y].Size = rl.NewVector2(squareSize-1, squareSize-1)
			if rand.Float64() < 0.1 && clear == false {
				g.Cells[x][y].Alive = true
			}
		}
	}*/
}

// Input - Game input
func (g *Game) Input(camera *rl.Camera2D) {
	// control
	if rl.IsKeyPressed(rl.KeyR) {
		g.Init(false)
	}
	if rl.IsKeyPressed(rl.KeyC) {
		g.Init(true)
	}
	if rl.IsKeyDown(rl.KeyU) && !g.Playing {
		g.Update()
	}

	if camera.Zoom != 1.0 {
		if rl.IsKeyDown(rl.KeyUp) {
			camera.Target = rl.Vector2Add(camera.Target, rl.NewVector2(0, -10))
		}
		if rl.IsKeyDown(rl.KeyDown) {
			camera.Target = rl.Vector2Add(camera.Target, rl.NewVector2(0, 10))
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			camera.Target = rl.Vector2Add(camera.Target, rl.NewVector2(-10, 0))
		}
		if rl.IsKeyDown(rl.KeyRight) {
			camera.Target = rl.Vector2Add(camera.Target, rl.NewVector2(10, 0))
		}
	}

	if rl.IsKeyPressed(rl.KeySpace) {
		g.Playing = !g.Playing
	}

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		g.Click(rl.GetMouseX(), rl.GetMouseY())
		if g.Playing {
			delta := rl.GetMouseDelta()
			delta = rl.Vector2Scale(delta, -1.0/camera.Zoom)

			camera.Target = rl.Vector2Add(camera.Target, delta)
		}
	}

	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		// Get the world point that is under the mouse
		mouseWorldPos := rl.GetScreenToWorld2D(rl.GetMousePosition(), *camera)

		// Set the offset to where the mouse is
		camera.Offset = rl.GetMousePosition()

		// Set the target to match, so that the camera maps the world space point
		// under the cursor to the screen space point under the cursor at any zoom
		camera.Target = mouseWorldPos

		// Zoom increment
		const zoomIncrement float32 = 0.125

		camera.Zoom += wheel * zoomIncrement
		if camera.Zoom < 1 {
			camera.Zoom = 1
			//camera.Target = rl.Vector2{}
			//camera.Offset = rl.Vector2{}
			camera.Target = rl.Vector2{X: float32(g.ScreenWidth / 2), Y: float32(g.ScreenHeight / 2)}
			camera.Offset = rl.Vector2{X: float32(g.ScreenWidth / 2), Y: float32(g.ScreenHeight / 2)}
		}
	}

	g.FramesCounter++
}

// Click - Toggle if a cell is alive or dead on click
func (g *Game) Click(x, y int32) {
	for i := int32(0); i <= g.Cols; i++ {
		for j := int32(0); j <= g.Rows; j++ {
			cell := g.Cells[i][j].Position
			if int32(cell.X) < x && int32(cell.X)+squareSize > x && int32(cell.Y) < y && int32(cell.Y)+squareSize > y {
				g.Cells[i][j].Alive = !g.Cells[i][j].Alive
				g.Cells[i][j].Next = g.Cells[i][j].Alive
			}
		}
	}
}

// Update - Update game
func (g *Game) Update() {
	/*for i := int32(0); i <= g.Cols; i++ {
		for j := int32(0); j <= g.Rows; j++ {
			NeighborCount := g.CountNeighbors(i, j)
			if g.Cells[i][j].Alive {
				if NeighborCount < 2 {
					g.Cells[i][j].Next = false
				} else if NeighborCount > 3 {
					g.Cells[i][j].Next = false
				} else {
					g.Cells[i][j].Next = true
				}
			} else {
				if NeighborCount == 3 {
					g.Cells[i][j].Next = true
					g.Cells[i][j].Visited = true
				}
			}
		}

		for j := int32(0); j < g.Rows; j++ {
			g.Cells[i][j].Alive = g.Cells[i][j].Next
		}
	}*/
}

// Draw - Draw game
func (g *Game) Draw(camera *rl.Camera2D, texture *rl.Texture2D) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	rl.BeginMode2D(*camera)

	//Draw cells some of which are visited and some of which are not
	//for x := int32(0); x <= g.Cols; x++ {
	//	for y := int32(0); y <= g.Rows; y++ {
	//		if g.Cells[x][y].Alive {
	//			rl.DrawRectangleV(g.Cells[x][y].Position, g.Cells[x][y].Size, rl.Blue)
	//		} else if g.Cells[x][y].Visited {
	//			rl.DrawRectangleV(g.Cells[x][y].Position, g.Cells[x][y].Size, rl.Color{R: 128, G: 177, B: 136, A: 255})
	//		}
	//	}
	//}

	sourceRec := rl.NewRectangle(
		float32(30),
		float32(0),
		float32(35),
		float32(30),
	)
	destRec := rl.NewRectangle(squareSize*3, 0, squareSize, squareSize)
	rl.DrawTexturePro(*texture, sourceRec, destRec, rl.Vector2{X: 0, Y: 0}, 0, rl.White)

	// Draw grid lines
	for i := int32(0); i < g.Cols+1; i++ {
		rl.DrawLineV(
			rl.NewVector2(float32(squareSize*i), 0),
			rl.NewVector2(float32(squareSize*i), float32(g.ScreenHeight)),
			rl.LightGray,
		)
	}

	for i := int32(0); i < g.Rows+1; i++ {
		rl.DrawLineV(
			rl.NewVector2(0, float32(squareSize*i)),
			rl.NewVector2(float32(g.ScreenWidth), float32(squareSize*i)),
			rl.LightGray,
		)
	}
	rl.EndMode2D()
	rl.EndDrawing()
}

func initializeAssetTextures() (AssetTextures, error) {
	assetTextures := make(AssetTextures)

	assets := map[string]string{
		"Basic_cell": "assets/basic_cell.png",
	}

}

func unloadAssetTextures(assetTextures AssetTextures) {
	for _, texture := range assetTextures {
		rl.UnloadTexture(texture)
	}
}
