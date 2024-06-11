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
		fmt.Printf("Failed to parse PORT env variable (must be int): %s\n", err)
	}

	current_user, err := user.Current()
	if err != nil {
		fmt.Printf("Failed to get current user from os: %s\n", err)
		current_user = &user.User{Username: "unknown"}
	}

	address, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port_int))
	if err != nil {
		fmt.Printf("Failed to resolv UDP address: %s\n", address)
	}

	connection, err := net.ListenUDP("udp", address)
	if err != nil {
		fmt.Printf("Failed to listen UDP %v: %s\n", address, err)
		os.Exit(-1)
	}
	defer connection.Close()

	buffer := make([]byte, 1024)
	_, remote_address, err := connection.ReadFromUDP(buffer)
	if err != nil {
		fmt.Printf("Failed to read from UDP: %s\n", err)
	}
	go handle_connection(connection, remote_address, current_user)

	fmt.Printf("\rMessage: %s", buffer)
	for {
		buffer = make([]byte, 1024)

		_, new_remote_address, err := connection.ReadFromUDP(buffer)
		*remote_address = *new_remote_address
		if err != nil {
			fmt.Printf("Failed to read from UDP: %s\n", err)
		}

		fmt.Printf("\rMessage: %s", buffer)
	}
}

func handle_connection(
	connection *net.UDPConn,
	remote_address *net.UDPAddr,
	current_user *user.User,
) {
	stdin := bufio.NewReader(os.Stdin)
	_ = current_user
	for {
		message, err := stdin.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Failed to read from user input: %s\n", err)
		}

		_, err = connection.WriteToUDP(message, remote_address)
		if err != nil {
			fmt.Printf("Failed to write to %v: %s\n", remote_address, err)
			return
		}
	}
}
