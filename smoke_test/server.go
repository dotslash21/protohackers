package smoke_test

import (
	"fmt"
	"io"
	"log"
	"net"
)

var PORT = 80

func RunServer() {
	for {
		// Listen for incoming connections
		log.Print("Listening on port ", PORT, " for incoming connections.")
		lst, err := net.Listen("tcp", ":"+fmt.Sprint(PORT))
		if err != nil {
			// Handle error
			log.Fatal(err)
		}

		// Echo all incoming data.
		for {
			// Accept a TCP connection
			conn, err := lst.Accept()
			if err != nil {
				// Handle error
				log.Fatal(err)
			}

			log.Print("Accepted connection from ", conn.RemoteAddr().String())
			go handleConnection(conn)
		}
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	io.Copy(c, c)

	log.Print("Closing connection to ", c.RemoteAddr().String())
}
