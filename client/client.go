package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/user"
	"strconv"
)

func main() {
	ip, port := os.Getenv("IP"), os.Getenv("PORT")
	port_int, err := strconv.Atoi(port)
	if err != nil {
		fmt.Printf("Failed to parse port ENV: %s\n", err)
		os.Exit(-1)
	}

	current_user, err := user.Current()
	if err != nil {
		fmt.Printf("Failed to get current user: %s\n", err)
		current_user = &user.User{Username: "unknown"}
	}

	remote_address, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port_int))
	if err != nil {
		fmt.Printf("Failed to resolve UDP address: %s\n", err)
		os.Exit(-1)
	}

	connection, err := net.DialUDP("udp", nil, remote_address)
	if err != nil {
		fmt.Printf("Failed to dial %v: %s\n", remote_address, err)
		os.Exit(-1)
	}
	defer connection.Close()

	_, err = connection.Write([]byte(fmt.Sprintf("Hello. A am %s\n", current_user.Username)))
	if err != nil {
		fmt.Printf("Failed to send hello message: %s\n", err)
	}

	go handle_connection(connection, current_user)

	for {
		buffer := make([]byte, 1024)

		_, err = connection.Read(buffer)
		if err != nil {
			fmt.Printf("Failed to read from UDP: %s\n", err)
			os.Exit(-1)
		}

		fmt.Printf("\rMessage: %s", buffer)
	}
}

func handle_connection(
	connection *net.UDPConn,
	current_user *user.User,
) {
	stdin := bufio.NewReader(os.Stdin)
	_ = current_user
	for {
		message, err := stdin.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Failed to read from user input: %s\n", err)
		}

		_, err = connection.Write(message)
		if err != nil {
			fmt.Printf("Error trying to write to a connection: %s\n", err)
		}
	}
}
