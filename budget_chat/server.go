package budget_chat

import (
	"log"
	"net"
	"strconv"

	"arunangshu.dev/protohackers/budget_chat/chat"
)

const PORT = 80

var room = chat.NewRoom()

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

	participant, err := room.RegisterParticipant(conn)
	if err != nil {
		log.Printf("Failed to register participant: %v", err)
		return
	}

	room.ListenForMessages(participant)

	log.Printf("Participant %s disconnected", participant.Name)
}
