package main

import (
	"fmt"
	"github.com/JanBerktold/sse"
	"net/http"
	"time"
)

func HandleSSE(w http.ResponseWriter, r *http.Request) {
	conn, err := sse.Upgrade(w, r)

	if err != nil {
		// log error to console
		fmt.Printf("Error occured: %q", err.Error())
	}

	go func() {
		time.Sleep(5e9) // 5s
		conn.Close()
		fmt.Println("conn.Close()...")
		return
	}()

	go func() {
		for {
			if !conn.IsOpen() {
				fmt.Println("fast conn closed ...")
				return
			}
			// update time every second
			conn.WriteString("-----" + time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
			time.Sleep(500 * time.Millisecond)
		}
	}()

	go func() {
		for {
			if !conn.IsOpen() {
				fmt.Println("fastest conn closed ...")
				return
			}
			// update time every second
			conn.WriteString("----------" + time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
			time.Sleep(100 * time.Millisecond)
		}
	}()

	for {
		if !conn.IsOpen() {
			fmt.Println("slow conn closed ...")
			return
		}
		// update time every second
		conn.WriteString(time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
		time.Sleep(1 * time.Second)
	}
}

func main() {

	// handle server-sent events request
	http.HandleFunc("/event", HandleSSE)

	// serve HTML page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "main.html")
	})

	http.ListenAndServe(":12340", nil)
}
