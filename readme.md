func main() {
const (
screenWidth  = 800
screenHeight = 450
)

	// Initialize Raylib
	rl.InitWindow(screenWidth, screenHeight, "Cropped Texture Example")
	defer rl.CloseWindow()

	// Load PNG sprite
	image, width, height := rl.LoadImageEx("sprite.png")
	if image.Data == nil {
		fmt.Println("Failed to load image.")
		os.Exit(1)
	}
	defer rl.UnloadImage(image)

	// Define the region to crop
	cropX := 50
	cropY := 50
	cropWidth := 100
	cropHeight := 100

	// Crop the sprite
	croppedImage := rl.ImageFromImage(image, rl.Rectangle{
		X: float32(cropX),
		Y: float32(cropY),
		Width: float32(cropWidth),
		Height: float32(cropHeight),
	})

	// Create texture from the cropped image
	croppedTexture := rl.LoadTextureFromImage(croppedImage)
	defer rl.UnloadTexture(croppedTexture)

	// Main loop
	for !rl.WindowShouldClose() {
		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawTexture(croppedTexture, 0, 0, rl.White)
		rl.EndDrawing()
	}