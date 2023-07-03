package unusual_database_program

import (
	"log"
	"net"
	"strconv"
)

const PORT = 80
const MAX_BUFFER_SIZE = 65507
const VERSION = "Ken's Key-Value Store 1.0"

var db = map[string]string{}

func RunServer() {
	pc, err := net.ListenPacket("udp", ":"+strconv.Itoa(PORT))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer pc.Close()
	log.Printf("Listening on port %d", PORT)

	ReadFromClient(pc)
}

func ReadFromClient(pc net.PacketConn) {
	buffer := make([]byte, MAX_BUFFER_SIZE)
	for {
		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			log.Printf("Failed to read from client: %v", err)
			return
		}
		log.Printf("Received %d bytes from %s", n, addr.String())

		hasResponse, response := HandleRequest(buffer[:n])

		if hasResponse {
			log.Printf("Sending response: %s", response)
			_, err := pc.WriteTo([]byte(response), addr)
			if err != nil {
				log.Printf("Failed to send response: %v", err)
				return
			}
		}
	}
}

func FindFirst(s []byte, c byte) int {
	for i, b := range s {
		if b == c {
			return i
		}
	}
	return -1
}

func HandleRequest(request []byte) (bool, string) {
	log.Printf("Received request: %s", request)

	firstindexOfEqual := FindFirst(request, '=')
	if firstindexOfEqual == -1 {
		log.Printf("Request type: Retrieve")
		key := string(request)

		if key == "version" {
			return true, "version=Ken's Key-Value Store 1.0"
		}

		value, ok := db[key]
		if !ok {
			return true, key + "="
		}

		return true, key + "=" + value
	} else {
		log.Printf("Request type: Insert")
		key := string(request[:firstindexOfEqual])
		value := string(request[firstindexOfEqual+1:])

		if key == "version" {
			return false, ""
		}

		db[key] = value

		return false, ""
	}
}
