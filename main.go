package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

var config string

func switchHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	confs, ok := r.URL.Query()["image"]

	if !ok || len(confs[0]) < 1 {
		config = "original"
		fmt.Fprintf(w, "original")
		return
	}

	conf := confs[0]
	if conf == "original" {
		config = "original"
	} else if conf == "overlay" {
		config = "overlay"
	} else {
		config = "original"
	}

	fmt.Fprintf(w, config)
}

func serveImageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	rawImageName := ps.ByName("image")
	arrs := strings.Split(rawImageName, ".")
	imageName := arrs[0]

	if config == "overlay" {
		imageName += "_overlay"
	}

	fileURL := fmt.Sprintf("https://dummyimage.com/400x400/000/fff.jpg&text=%s", imageName)
	if err := downloadFile(fmt.Sprintf("/tmp/%s", rawImageName), fileURL); err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Error occured when streaming image")
		return
	}

	img, err := os.Open(fmt.Sprintf("/tmp/%s", rawImageName))
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Error occured when opening image")
		return
	}
	defer img.Close()

	w.Header().Set("Content-Type", "image/jpeg")
	io.Copy(w, img)
}

func downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func router() http.Handler {
	router := httprouter.New()

	router.GET("/switch", switchHandler)
	router.GET("/image/:image", serveImageHandler)

	return router
}

func main() {
	config = "original"

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "PUT", "HEAD", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		MaxAge:         86400,
	})

	handler := c.Handler(router())
	httpServer := &http.Server{Addr: ":8474", Handler: handler}

	log.Println("[server] listening on :8474")
	httpServer.ListenAndServe()
}
