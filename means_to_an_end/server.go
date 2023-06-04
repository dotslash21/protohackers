package means_to_an_end

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"arunangshu.dev/protohackers/means_to_an_end/utils"
)

const PORT = 80
const TIMEOUT_SECONDS = 5
const MSG_LEN = 9

func RunServer() {
	lst, err := net.Listen("tcp", ":"+strconv.Itoa(PORT))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listening on port %d", PORT)

	for {
		conn, err := lst.Accept()
		if err != nil {
			log.Fatalf("Failed to accept connection: %v", err)
		}
		log.Printf("Accepted connection from %s", conn.RemoteAddr())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	timeouts := 0
	request_processor := utils.NewRequestProcessor()

	for {
		// Create a buffer to read into
		requestBytes := make([]byte, MSG_LEN)

		// Set a deadline for the ReadFull call
		conn.SetReadDeadline(time.Now().Add(TIMEOUT_SECONDS * time.Second))

		// Read exactly 9 bytes
		_, err := io.ReadFull(conn, requestBytes)
		if err != nil {
			if err == io.EOF {
				log.Printf("End of connection")
				break
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf("Timeout reading from connection")

				if timeouts > 0 {
					log.Printf("Too many timeouts, closing connection")
					break
				}

				timeouts++
				continue
			}
			log.Printf("Failed to read request: %v", err)
			continue
		}

		// Parse the request
		request, err := utils.ParseRequest(requestBytes)
		if err != nil {
			log.Printf("Failed to parse request: %v", err)
			break
		}

		// Process the request
		response := request_processor.ProcessRequest(request)

		if response.Method == 'Q' {
			// Send the response
			responseBytes := make([]byte, 4)
			binary.BigEndian.PutUint32(responseBytes, uint32(*response.Number))
			_, err := conn.Write(responseBytes)
			if err != nil {
				log.Printf("Failed to write response: %v", err)
				return
			}
		}
	}
}
