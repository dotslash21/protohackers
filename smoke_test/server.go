package smoke_test

import (
	"fmt"
	"log"
	"net"
)

var PORT = 8080

func RunServer() {
	for {
		// Listen for incoming connections
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

			go func(c net.Conn) {
				// Echo all incoming data.
				buf := make([]byte, 1024)
				for {
					n, err := c.Read(buf)
					if err != nil {
						// Handle error
						log.Fatal(err)
					}
					_, err = c.Write(buf[0:n])
					if err != nil {
						// Handle error
						log.Fatal(err)
					}
				}
			}(conn)
		}
	}
}
