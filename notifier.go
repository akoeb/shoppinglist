package main

import (
	"fmt"
	"strings"
)

// AllowedCommands is an array with all possible commands
var AllowedCommands = [1]string{"UPDATE"}

// Notifier is used to send messages between the streaming handler and the regular handlers
type Notifier struct {
	channel       chan string
	Commands      map[string]bool
	ReceiverCount int
}

// NewNotifier creates and returns a notifier
func NewNotifier() *Notifier {
	notifier := &Notifier{
		channel:       make(chan string),
		Commands:      make(map[string]bool),
		ReceiverCount: 0,
	}
	// create map with allowed commands
	for _, item := range AllowedCommands {
		notifier.Commands[item] = true
	}
	return notifier
}

// Send a message to the notifier
func (n *Notifier) Send(msg string) error {
	upper := strings.ToUpper(msg)
	if _, ok := n.Commands[upper]; !ok {
		return fmt.Errorf("Not a Valid Command: %s", msg)
	}
	if n.ReceiverCount > 0 {
		n.channel <- upper
	}

	return nil
}

// Wait for the next message from the notifier
func (n *Notifier) Wait() string {
	cmd := <-n.channel
	return cmd
}
func (n *Notifier) AddReceiver() {
	n.ReceiverCount++
}
func (n *Notifier) RemoveReceiver() {
	n.ReceiverCount--
}
