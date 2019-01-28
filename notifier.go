package main

import (
	"fmt"
	"strings"
)

// AllowedCommands is an array with all possible commands
// TODO: notifier must show either which category changed, or whether the list of categories changed
var AllowedCommands = [1]string{"UPDATE"}

// Notifier is used to send messages between the streaming handler and the regular handlers
type Notifier struct {
	incoming  chan string
	listeners map[int]chan string
	Commands  map[string]bool
	maxSlots  int
	pool      []int
}

// NewNotifier creates and returns a notifier
func NewNotifier() *Notifier {
	notifier := &Notifier{
		incoming:  make(chan string),
		Commands:  make(map[string]bool),
		listeners: make(map[int]chan string),
		maxSlots:  100,
	}
	// create map with allowed commands
	for _, item := range AllowedCommands {
		notifier.Commands[item] = true
	}
	// initialize pool of available subscriber slots:
	for i := 0; i < notifier.maxSlots; i++ {
		notifier.pool = append(notifier.pool, i+1)
	}
	go notifier.Dispatcher()
	return notifier
}

// Dispatcher waits on new messages on inocmimg channel and dispatches them to all listening clients
func (n *Notifier) Dispatcher() {
	for {
		msg := <-n.incoming
		for chanID := range n.listeners {
			go n.sendToChannel(chanID, msg)
		}
	}
}

// go func to send notification to a given receiver
// this will silently fail if the channel is closed
func (n *Notifier) sendToChannel(chanID int, msg string) {
	if channel, ok := n.listeners[chanID]; ok {
		channel <- msg
	}
}

// Send a message to all listening receivers
func (n *Notifier) Send(msg string, categoryID int) error {
	upper := strings.ToUpper(msg)
	if _, ok := n.Commands[upper]; !ok {
		return fmt.Errorf("Not a Valid Command: %s", msg)
	}
	n.incoming <- upper

	return nil
}

// Listen for the next message from the notifier to a given receiver
func (n *Notifier) Listen(chanID int) string {
	cmd := ""
	if channel, ok := n.listeners[chanID]; ok {
		cmd = <-channel
	}
	return cmd
}

// NewReceiver creates a new Listening client and returns its id
func (n *Notifier) NewReceiver() (int, error) {
	if len(n.pool) == 0 {
		return 0, fmt.Errorf("Can not accept more than %d connections at the same time", n.maxSlots)
	}

	channel := make(chan string)

	// get id from pool:
	var chanID int
	chanID, n.pool = n.pool[0], n.pool[1:]

	n.listeners[chanID] = channel
	return chanID, nil
}

// RemoveReceiver deletes a given client and returns the id back to the pool of open slots
func (n *Notifier) RemoveReceiver(chanID int) {
	delete(n.listeners, chanID)
	n.pool = append(n.pool, chanID)
}
