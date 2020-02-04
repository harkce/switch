package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var config string

func switchHandler(w http.ResponseWriter, r *http.Request) {
	confs, ok := r.URL.Query()["image"]

	if !ok || len(confs[0]) < 1 {
		config = "1"
		fmt.Fprintf(w, "1")
		return
	}

	conf := confs[0]
	if conf == "1" {
		config = "1"
	} else if conf == "2" {
		config = "2"
	} else {
		config = "1"
	}

	fmt.Fprintf(w, config)
}

func serveImageHandler(w http.ResponseWriter, r *http.Request) {
	img, err := os.Open(fmt.Sprintf("%s.jpg", config))
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Error occured")
		return
	}
	defer img.Close()

	w.Header().Set("Content-Type", "image/jpeg")
	io.Copy(w, img)
}

func main() {
	config = "1"
	http.HandleFunc("/switch", switchHandler)
	http.HandleFunc("/gambar.jpg", serveImageHandler)

	log.Println("Served at :8474")
	log.Fatal(http.ListenAndServe(":8474", nil))
}
