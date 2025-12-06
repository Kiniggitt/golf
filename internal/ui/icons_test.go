package ui

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

func TestCreateGolfBallIcon(t *testing.T) {
	icon := CreateGolfBallIcon()

	if icon == nil {
		t.Fatal("Expected non-nil icon")
	}

	if icon.StaticName != "golfball.png" {
		t.Errorf("Expected name 'golfball.png', got %s", icon.StaticName)
	}

	if len(icon.StaticContent) == 0 {
		t.Error("Expected non-empty icon content")
	}
}

func TestCreateGolfBallIconValidPNG(t *testing.T) {
	icon := CreateGolfBallIcon()

	// Try to decode the PNG to verify it's valid
	reader := bytes.NewReader(icon.StaticContent)
	img, err := png.Decode(reader)
	if err != nil {
		t.Fatalf("Failed to decode icon as PNG: %v", err)
	}

	if img == nil {
		t.Fatal("Expected non-nil image")
	}
}

func TestCreateGolfBallIconSize(t *testing.T) {
	icon := CreateGolfBallIcon()

	// Decode and check dimensions
	reader := bytes.NewReader(icon.StaticContent)
	img, err := png.Decode(reader)
	if err != nil {
		t.Fatalf("Failed to decode icon: %v", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	expectedSize := 128

	if width != expectedSize {
		t.Errorf("Expected width %d, got %d", expectedSize, width)
	}

	if height != expectedSize {
		t.Errorf("Expected height %d, got %d", expectedSize, height)
	}

	// Verify it's square
	if width != height {
		t.Errorf("Expected square icon, got %dx%d", width, height)
	}
}

func TestCreateGolfBallIconHasContent(t *testing.T) {
	icon := CreateGolfBallIcon()

	reader := bytes.NewReader(icon.StaticContent)
	img, err := png.Decode(reader)
	if err != nil {
		t.Fatalf("Failed to decode icon: %v", err)
	}

	// Check that the icon has some variation in colors
	// (not just a solid color)
	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		t.Fatal("Expected RGBA image")
	}

	// Sample some pixels and verify there's color variation
	bounds := rgbaImg.Bounds()
	centerX := bounds.Dx() / 2
	centerY := bounds.Dy() / 2

	// Get center pixel (should be part of the golf ball)
	centerColor := rgbaImg.RGBAAt(centerX, centerY)

	// Get corner pixel (should be background)
	cornerColor := rgbaImg.RGBAAt(0, 0)

	// They should be different
	if centerColor.R == cornerColor.R &&
		centerColor.G == cornerColor.G &&
		centerColor.B == cornerColor.B {
		t.Error("Expected different colors in center vs corner (ball vs background)")
	}
}

func TestCreateGolfBallIconGreenBackground(t *testing.T) {
	icon := CreateGolfBallIcon()

	reader := bytes.NewReader(icon.StaticContent)
	img, err := png.Decode(reader)
	if err != nil {
		t.Fatalf("Failed to decode icon: %v", err)
	}

	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		t.Fatal("Expected RGBA image")
	}

	// Check corner pixel for green background
	// Background color is RGBA{34, 139, 34, 255} (green)
	cornerColor := rgbaImg.RGBAAt(0, 0)

	// Verify it's greenish (G channel should be dominant)
	if cornerColor.G <= cornerColor.R || cornerColor.G <= cornerColor.B {
		t.Error("Expected green background (G channel should be dominant)")
	}

	// Verify it's not completely black or white
	if (cornerColor.R == 0 && cornerColor.G == 0 && cornerColor.B == 0) ||
		(cornerColor.R == 255 && cornerColor.G == 255 && cornerColor.B == 255) {
		t.Error("Expected colored background, not pure black or white")
	}
}

func TestCreateGolfBallIconHasBall(t *testing.T) {
	icon := CreateGolfBallIcon()

	reader := bytes.NewReader(icon.StaticContent)
	img, err := png.Decode(reader)
	if err != nil {
		t.Fatalf("Failed to decode icon: %v", err)
	}

	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		t.Fatal("Expected RGBA image")
	}

	// Check center pixel - should be part of the golf ball
	// Ball colors are light (200-240 range)
	bounds := rgbaImg.Bounds()
	centerX := bounds.Dx() / 2
	centerY := bounds.Dy() / 2

	centerColor := rgbaImg.RGBAAt(centerX, centerY)

	// Golf ball should be light colored (high RGB values)
	avgColorValue := (int(centerColor.R) + int(centerColor.G) + int(centerColor.B)) / 3

	if avgColorValue < 100 {
		t.Errorf("Expected light colored ball in center, got avg color value %d", avgColorValue)
	}
}

func TestCreateGolfBallIconDimples(t *testing.T) {
	icon := CreateGolfBallIcon()

	reader := bytes.NewReader(icon.StaticContent)
	img, err := png.Decode(reader)
	if err != nil {
		t.Fatalf("Failed to decode icon: %v", err)
	}

	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		t.Fatal("Expected RGBA image")
	}

	// Count unique colors in the ball area
	// Should have at least 3: background, ball, dimples, shadows
	colorMap := make(map[uint32]bool)

	bounds := rgbaImg.Bounds()
	for y := 20; y < bounds.Dy()-20; y += 5 {
		for x := 20; x < bounds.Dx()-20; x += 5 {
			c := rgbaImg.RGBAAt(x, y)
			// Create a color key
			colorKey := uint32(c.R)<<16 | uint32(c.G)<<8 | uint32(c.B)
			colorMap[colorKey] = true
		}
	}

	// Should have multiple colors (background, ball, dimples, shadows)
	if len(colorMap) < 3 {
		t.Errorf("Expected at least 3 different colors, got %d", len(colorMap))
	}
}

func TestCreateGolfBallIconConsistency(t *testing.T) {
	// Test that creating the icon multiple times produces the same result
	icon1 := CreateGolfBallIcon()
	icon2 := CreateGolfBallIcon()

	if !bytes.Equal(icon1.StaticContent, icon2.StaticContent) {
		t.Error("Expected consistent icon generation")
	}

	if icon1.StaticName != icon2.StaticName {
		t.Error("Expected same icon name")
	}
}

func TestCreateGolfBallIconNotEmpty(t *testing.T) {
	icon := CreateGolfBallIcon()

	// Minimum size check - a valid 128x128 PNG should be at least a few KB
	minExpectedSize := 100 // bytes

	if len(icon.StaticContent) < minExpectedSize {
		t.Errorf("Icon seems too small: %d bytes (expected at least %d)",
			len(icon.StaticContent), minExpectedSize)
	}
}

func TestCreateGolfBallIconAlphaChannel(t *testing.T) {
	icon := CreateGolfBallIcon()

	reader := bytes.NewReader(icon.StaticContent)
	img, err := png.Decode(reader)
	if err != nil {
		t.Fatalf("Failed to decode icon: %v", err)
	}

	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		t.Fatal("Expected RGBA image")
	}

	// Check that alpha channel is fully opaque (255)
	bounds := rgbaImg.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 10 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 10 {
			c := rgbaImg.RGBAAt(x, y)
			if c.A != 255 {
				t.Errorf("Expected fully opaque pixels, got alpha=%d at (%d,%d)", c.A, x, y)
			}
		}
	}
}
