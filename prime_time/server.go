package prime_time

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
)

const PORT = 80

type Request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

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

	// Create a scanner to read requests from the connection
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		// Read the request string
		requestString := scanner.Text()
		log.Printf("Received request: %s", requestString)

		// Parse the JSON
		request, err := parseRequest(requestString)
		if err != nil {
			log.Printf("Failed to parse request: %v", err)
			writeErrorResponse(conn)
			break
		}

		// Process the request
		response, err := processRequest(request)
		if err != nil {
			log.Printf("Failed to process request: %v", err)
			writeErrorResponse(conn)
			break
		}

		// Write the response to the connection
		if err := writeResponse(conn, response); err != nil {
			log.Printf("Failed to write response: %v", err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Failed to read from connection: %v", err)
	}

	log.Printf("Closing connection")
}

func parseRequest(requestString string) (Request, error) {
	var request Request
	if err := json.Unmarshal([]byte(requestString), &request); err != nil {
		return Request{}, fmt.Errorf("failed to parse JSON: %v", err)
	}
	return request, nil
}

func processRequest(request Request) (Response, error) {
	if request.Method != "isPrime" {
		return Response{}, fmt.Errorf("invalid method: %s", request.Method)
	}

	if request.Number == nil {
		return Response{}, fmt.Errorf("missing number")
	}

	isPrime := isPrime(int64(*request.Number))

	response := Response{
		Method: "isPrime",
		Prime:  isPrime,
	}

	return response, nil
}

func isPrime(n int64) bool {
	if n <= 1 {
		return false
	}

	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func writeResponse(conn net.Conn, response Response) error {
	responseString, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to convert response to JSON: %v", err)
	}
	responseString = append(responseString, '\n')

	log.Printf("Writing response: %s", responseString)

	nw, err := conn.Write(responseString)
	if err != nil {
		return fmt.Errorf("failed to write response to connection: %v", err)
	}
	log.Printf("Wrote %d bytes to connection", nw)

	return nil
}

func writeErrorResponse(conn net.Conn) {
	conn.Write([]byte("{\"message\": \"Invalid request\"}\n"))
}
