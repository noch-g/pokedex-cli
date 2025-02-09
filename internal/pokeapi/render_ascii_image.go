package pokeapi

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io"
	"net/http"
	"strings"

	"github.com/noch-g/pokedex-cli/internal/logger"
)

func (c *Client) RenderImage(pokemon *Pokemon) (string, error) {
	imageUrl := pokemon.Sprites.Versions.GenerationI.RedBlue.FrontGray

	img, err := c.downloadImage(imageUrl)
	if err != nil {
		return "", err
	}
	return generateASCII(img), nil
}

func (c *Client) downloadImage(imageUrl string) (image.Image, error) {
	url := imageUrl

	// Check cache before request
	if val, ok := c.cache.Get(url); ok {
		logger.Debug("FROM CACHE", "url", url)
		img, _, err := image.Decode(bytes.NewReader(val))
		if err != nil {
			return nil, err
		}
		return img, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	logger.Debug("GET", "url", url)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image, HTTP status: %s", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil, fmt.Errorf("response is not an image, Content-Type: %s", contentType)
	}

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(dat))
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %v", err)
	}

	c.cache.Add(url, dat)
	return img, nil
}

func pixelToASCII(c color.Color) string {
	r, g, b, _ := c.RGBA()

	r = r / 256
	g = g / 256
	b = b / 256

	gray := (float64(r) + float64(g) + float64(b)) / 3

	gray = gray / 255.0

	asciiChars := []string{"#", "8", "@", "&", "o", ":", "*", ".", " "}
	idx := int(gray * float64(len(asciiChars)-1))

	return asciiChars[idx]
}

func generateASCII(img image.Image) string {
	bounds := img.Bounds()
	height := bounds.Dy()
	width := bounds.Dx()

	asciiArt := strings.Builder{}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := img.At(x, y)
			asciiArt.WriteString(pixelToASCII(pixel) + " ")
		}
		asciiArt.WriteString("\n")
	}

	return asciiArt.String()
}
