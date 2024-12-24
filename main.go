package main

import (
	"fmt"
	"net/http"
)

func main() {
	videoController, err := InitVideoController()

	if err != nil {
		fmt.Println("Error initializing:", err)
		return
	}

	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		videoController.ShowListResolution(w, r)
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
