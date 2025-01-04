package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	conn net.Conn
	nick string
}

func isNicknameTaken(nick string, clients []Client) bool {
	for _, client := range clients {
		if strings.EqualFold(client.nick, nick) {
			return true
		}
	}
	return false
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is running on port 8080...")
	var clients []Client
	var mu sync.Mutex

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleNewClient(conn, &clients, &mu)
	}
}

func handleNewClient(conn net.Conn, clients *[]Client, mu *sync.Mutex) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	nick, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	nick = strings.TrimSpace(nick)

	mu.Lock()
	if nick == "" || isNicknameTaken(nick, *clients) {
		mu.Unlock()
		conn.Write([]byte("NICKNAME_TAKEN\n"))
		return
	}

	client := Client{conn: conn, nick: nick}
	*clients = append(*clients, client)
	mu.Unlock()

	conn.Write([]byte("OK\n"))

	broadcast(fmt.Sprintf("%s has joined the chat\n", nick), conn, *clients)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			mu.Lock()
			*clients = removeClient(*clients, conn)
			mu.Unlock()
			broadcast(fmt.Sprintf("%s has left the chat\n", nick), conn, *clients)
			break
		}

		message = strings.TrimSpace(message)
		if message != "" {
			broadcast(fmt.Sprintf("%s: %s\n", nick, message), conn, *clients)
		}
	}
}

func broadcast(message string, sender net.Conn, clients []Client) {
	for _, client := range clients {
		if client.conn != sender {
			client.conn.Write([]byte(message))
		}
	}
}

func removeClient(clients []Client, conn net.Conn) []Client {
	for i, client := range clients {
		if client.conn == conn {
			return append(clients[:i], clients[i+1:]...)
		}
	}
	return clients
}