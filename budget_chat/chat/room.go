package chat

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type Participant struct {
	Name string
	Conn net.Conn
}

type Room struct {
	mutex        sync.Mutex
	Participants []Participant
}

func NewRoom() *Room {
	return &Room{
		Participants: make([]Participant, 0),
	}
}

func (r *Room) RegisterParticipant(conn net.Conn) (*Participant, error) {
	// Prompt the client for their name
	namePrompt := []byte("[SERVER] Welcome to budgetchat! What shall I call you?\n")
	_, err := conn.Write(namePrompt)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
		return nil, fmt.Errorf("failed to write response")
	}

	// Read the client's name
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		name := scanner.Text()
		log.Printf("Client's name is %s", name)

		// Validate the name
		// a. Must be between 1 and 16 characters
		if len(name) < 1 || len(name) > 16 {
			log.Printf("The name must be between 1 and 16 characters")
			conn.Write([]byte("[SERVER] The name must be between 1 and 16 characters\n"))
			return nil, fmt.Errorf("the name must be between 1 and 16 characters")
		}
		// b. Must contain only alphanumeric characters
		for _, c := range name {
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
				log.Printf("The name must contain only alphanumeric characters")
				conn.Write([]byte("[SERVER] The name must contain only alphanumeric characters\n"))
				return nil, fmt.Errorf("the name must contain only alphanumeric characters")
			}
		}
		// c. Must not already be taken
		for _, p := range r.Participants {
			if p.Name == name {
				log.Printf("The name %s is already taken", name)
				conn.Write([]byte("[SERVER] The name " + name + " is already taken\n"))
				return nil, fmt.Errorf("the name %s is already taken", name)
			}
		}

		// Create a new client
		participant := Participant{
			Name: name,
			Conn: conn,
		}

		r.AddParticipant(&participant)

		return &participant, nil
	} else {
		log.Printf("Failed to read name")
		return nil, fmt.Errorf("failed to read name")
	}
}

func (r *Room) AddParticipant(participant *Participant) {
	r.mutex.Lock()

	r.notifyJoin(participant)
	r.Participants = append(r.Participants, *participant)

	r.mutex.Unlock()
}

func (r *Room) RemoveParticipant(participant *Participant) {
	r.mutex.Lock()

	for i, p := range r.Participants {
		if p.Name == participant.Name {
			r.Participants = append(r.Participants[:i], r.Participants[i+1:]...)
			break
		}
	}
	r.notifyLeave(participant)

	r.mutex.Unlock()
}

func (r *Room) BroadcastMessage(participant *Participant, message string) {
	r.mutex.Lock()

	for _, p := range r.Participants {
		if p.Name != participant.Name {
			_, err := p.Conn.Write([]byte("[" + participant.Name + "] " + message + "\n"))
			if err != nil {
				log.Printf("Failed to notify client: %v", err)
			}
		}
	}

	r.mutex.Unlock()
}

func (r *Room) ListenForMessages(participant *Participant) {
	defer r.RemoveParticipant(participant)

	scanner := bufio.NewScanner(participant.Conn)
	for scanner.Scan() {
		message := scanner.Text()
		log.Printf("Received message from %s: %s", participant.Name, message)

		r.BroadcastMessage(participant, message)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error occurred while reading from %s: %s", participant.Name, scanner.Err())
	}
}

func (r *Room) notifyJoin(newParticipant *Participant) {
	// Notify all existing participants of the new participant's presence
	for _, existingParticipant := range r.Participants {
		_, err := existingParticipant.Conn.Write([]byte("* " + newParticipant.Name + " has has entered the room\n"))
		if err != nil {
			log.Printf("Failed to notify client: %v", err)
		}
	}

	// Notify the new client of all other clients' presence
	_, err := newParticipant.Conn.Write([]byte("* The room contains: " + strings.Join(r.getParticipantNames(), ", ") + "\n"))
	if err != nil {
		log.Printf("Failed to notify client: %v", err)
	}
}

func (r *Room) getParticipantNames() []string {
	names := make([]string, 0, len(r.Participants))
	for _, p := range r.Participants {
		names = append(names, p.Name)
	}
	return names
}

func (r *Room) notifyLeave(participant *Participant) {
	// Notify all remaining participants that the participant has left
	for _, p := range r.Participants {
		_, err := p.Conn.Write([]byte("* " + participant.Name + " has left the room\n"))
		if err != nil {
			log.Printf("Failed to notify client: %v", err)
		}
	}
}
