package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/Tony-ArtZ/goxfer/cli"
	"github.com/Tony-ArtZ/goxfer/server"
	"github.com/Tony-ArtZ/goxfer/utils"
)

func main() {
	// Find available port
	var listener net.Listener
	var err error
	port := 8080

	for {
		listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			break
		}
		port++
	}

	portStr := strconv.Itoa(port)
	ip := utils.GetLocalIP()
	addr := fmt.Sprintf("http://%s:%s", ip, portStr)

	tm, err := server.NewTransferManager()
	if err != nil {
		fmt.Printf("Failed to initialize transfer manager: %v\n", err)
		os.Exit(1)
	}

	go func() {
		http.HandleFunc("/", server.HandleIndex)
		http.HandleFunc("/poll", tm.HandlePoll)
		http.HandleFunc("/download", tm.HandleDownload)
		http.HandleFunc("/upload", tm.HandleUpload)
		if err := http.Serve(listener, nil); err != nil {
			os.Exit(1)
		}
	}()

	cli.StartCLI(addr, portStr, tm)
}
