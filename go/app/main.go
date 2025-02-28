package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const (
	imgDir = "images"
	port   = "9000"
)

type Response struct {
	Message string `json:"message"`
}

type Item struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	// Image    string `json:"image_name"`
}

type Items struct {
	Items []Item `json:"items"`
}

// hello is an endpoint to return a Hello, world! message.
func hello(c echo.Context) error {
	res := Response{
		Message: "Hello, world!",
	}
	return c.JSON(http.StatusOK, res)
}

var item Item
var items Items

func addItem(c echo.Context) error {
	// Get form data
	name := c.FormValue("name")
	c.Logger().Infof("Receive item: %s", name)

	category := c.FormValue("category")
	c.Logger().Infof("Receive category: %s", category)

	item = Item{Name: name, Category: category}
	fmt.Println(item)

	message := fmt.Sprintf("item received: %s, %s", item.Name, item.Category)
	res := Response{Message: message}

	items.Items = append(items.Items, item)
	log.Printf("%v", items)

	return c.JSON(http.StatusOK, res)
}

func getItems(c echo.Context) error {
	return c.JSON(http.StatusOK, items)
}

// getImage is an endpoint to return an image.
func getImage(c echo.Context) error {
	// Build image path
	imgPath := filepath.Join(imgDir, filepath.Clean(c.Param("imageFilename")))
	rel, err := filepath.Rel(imgDir, imgPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		res := Response{
			Message: "Invalid image path",
		}
		return c.JSON(http.StatusBadRequest, res)
	}

	c.Logger().Info("Image path: ", imgPath)

	if !strings.HasSuffix(imgPath, ".jpg") {
		res := Response{
			Message: "Image path does not end with .jpg",
		}
		return c.JSON(http.StatusBadRequest, res)
	}
	if _, err := os.Stat(imgPath); err != nil {
		c.Logger().Debugf("Image not found: %s", imgPath)
		imgPath = filepath.Join(imgDir, "default.jpg")
	}
	return c.File(imgPath)
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.INFO)

	frontURL := os.Getenv("FRONT_URL")
	if frontURL == "" {
		frontURL = "http://localhost:3000"
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{frontURL},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	// Routes
	e.GET("/", hello)
	e.POST("/items", addItem)
	e.GET("/items", getItems)
	e.GET("/image/:imageFilename", getImage)

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}
