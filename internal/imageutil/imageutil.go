package imageutil

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/image/draw"
)

const (
	ThumbnailSize = 128
)

// SaveImage saves the image data to disk and creates a thumbnail
// Returns (fullImagePath, thumbnailPath, error)
func SaveImage(imageData []byte, format string) (string, string, error) {
	// Decode the image
	img, err := decodeImage(imageData, format)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Get or create the images directory
	exePath, err := os.Executable()
	if err != nil {
		return "", "", err
	}
	exeDir := filepath.Dir(exePath)
	imagesDir := filepath.Join(exeDir, "data", "images")

	// Create images directory if it doesn't exist
	if err := os.MkdirAll(imagesDir, os.ModePerm); err != nil {
		return "", "", fmt.Errorf("failed to create images directory: %w", err)
	}

	// Generate unique filename using hash and timestamp
	hash := sha256.Sum256(imageData)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%x_%d.%s", hash[:8], timestamp, format)

	fullPath := filepath.Join(imagesDir, filename)
	thumbnailFilename := fmt.Sprintf("thumb_%s", filename)
	thumbnailPath := filepath.Join(imagesDir, thumbnailFilename)

	// Save the full image
	if err := saveImageFile(fullPath, img, format); err != nil {
		return "", "", fmt.Errorf("failed to save full image: %w", err)
	}

	// Create and save thumbnail
	thumbnail := createThumbnail(img, ThumbnailSize)
	if err := saveImageFile(thumbnailPath, thumbnail, format); err != nil {
		// Clean up the full image if thumbnail creation fails
		os.Remove(fullPath)
		return "", "", fmt.Errorf("failed to save thumbnail: %w", err)
	}

	return fullPath, thumbnailPath, nil
}

// decodeImage decodes image data based on format
func decodeImage(data []byte, format string) (image.Image, error) {
	reader := &bytesReader{data: data}

	switch format {
	case "png":
		return png.Decode(reader)
	case "jpg", "jpeg":
		return jpeg.Decode(reader)
	case "gif":
		return gif.Decode(reader)
	default:
		// Try to auto-detect
		reader.Reset()
		img, _, err := image.Decode(reader)
		return img, err
	}
}

// saveImageFile saves an image to a file
func saveImageFile(path string, img image.Image, format string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch format {
	case "png":
		return png.Encode(file, img)
	case "jpg", "jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	case "gif":
		return gif.Encode(file, img, nil)
	default:
		return png.Encode(file, img)
	}
}

// createThumbnail creates a thumbnail of the specified size
func createThumbnail(img image.Image, size int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Calculate aspect ratio
	var newWidth, newHeight int
	if width > height {
		newWidth = size
		newHeight = (height * size) / width
	} else {
		newHeight = size
		newWidth = (width * size) / height
	}

	// Create thumbnail
	thumbnail := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(thumbnail, thumbnail.Bounds(), img, bounds, draw.Over, nil)

	return thumbnail
}

// bytesReader implements io.Reader for byte slices
type bytesReader struct {
	data []byte
	pos  int
}

func (r *bytesReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, nil
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func (r *bytesReader) Reset() {
	r.pos = 0
}

// DeleteImage deletes both the full image and thumbnail
func DeleteImage(imagePath, previewPath string) error {
	var errs []error

	if imagePath != "" {
		if err := os.Remove(imagePath); err != nil && !os.IsNotExist(err) {
			errs = append(errs, err)
		}
	}

	if previewPath != "" {
		if err := os.Remove(previewPath); err != nil && !os.IsNotExist(err) {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors deleting images: %v", errs)
	}

	return nil
}
