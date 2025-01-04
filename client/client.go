package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type ChatState struct {
	messages []string
	nickname string
}

func (cs *ChatState) addMessage(msg string) {
	cs.messages = append(cs.messages, msg)
}

func (cs *ChatState) displayChat() {
	fmt.Print("\033[H\033[2J")
	fmt.Println("Connected to the chat! Start typing messages.")
	for _, msg := range cs.messages {
		fmt.Println(msg)
	}
	fmt.Print("\nEnter message: ")
}

func getNickname(reader *bufio.Reader) string {
	for {
		fmt.Print("Enter your nickname: ")
		nick, _ := reader.ReadString('\n')
		nick = strings.TrimSpace(nick)
		
		if nick != "" {
			return nick
		}
		fmt.Println("Nickname cannot be empty. Please try again.")
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	
	chatState := &ChatState{
		messages: make([]string, 0),
	}

	nick := getNickname(reader)
	chatState.nickname = nick
	conn.Write([]byte(nick + "\n"))

	serverReader := bufio.NewReader(conn)
	response, _ := serverReader.ReadString('\n')
	if strings.TrimSpace(response) == "NICKNAME_TAKEN" {
		fmt.Println("This nickname is already taken. Please try again.")
		return
	}

	chatState.displayChat()

	go func() {
		serverReader := bufio.NewReader(conn)
		for {
			message, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("\nDisconnected from server.")
				os.Exit(0)
			}

			chatState.addMessage(strings.TrimSpace(message))
			chatState.displayChat()
		}
	}()

	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if message != "" {
			conn.Write([]byte(message + "\n"))
			chatState.addMessage(fmt.Sprintf("%s: %s", nick, message))
			chatState.displayChat()
		}
	}
}