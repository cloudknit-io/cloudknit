package stream

import (
	"github.com/sirupsen/logrus"
)

const (
	SSEDataFormat = "data: %s\n\n"
)

type API interface {
	Notify(msg []byte)
	Clients() int
	NewMessageChannel() chan []byte
	NewConnection(messageChan chan []byte)
	CloseConnection(messageChan chan []byte)
}

type Broker struct {
	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan []byte
	// New client connections
	newClients chan chan []byte
	// Closed client connections
	closingClients chan chan []byte
	// Client connections registry
	clients map[chan []byte]bool
}

func NewService(l *logrus.Entry) (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	// Set it running - listening and broadcasting events
	go broker.listen(l)

	return
}

func (broker *Broker) listen(l *logrus.Entry) {
	l.Info("SSE Broker is listening for events")
	for {
		select {
		case s := <-broker.newClients:
			// A new client has connected.
			// Register their message channel
			broker.clients[s] = true
			l.Printf("Client added. %d registered clients", len(broker.clients))
		case s := <-broker.closingClients:
			// A client has detached, and we want to
			// stop sending them messages.
			delete(broker.clients, s)
			l.Printf("Removed client. %d registered clients", len(broker.clients))
		case event := <-broker.Notifier:
			// We got a new event from the outside!
			// Send event to all connected clients
			if len(broker.clients) > 0 {
				l.Printf("Sending notification to %d clients", len(broker.clients))
			}
			for clientMessageChan := range broker.clients {
				clientMessageChan <- event
			}
		}
	}
}

func (broker *Broker) Clients() int {
	return len(broker.clients)
}

func (broker *Broker) Notify(msg []byte) {
	select {
	case broker.Notifier <- msg:
		break
	default:
		break
	}
}

func (broker *Broker) NewMessageChannel() chan []byte {
	return make(chan []byte, 1)
}

func (broker *Broker) NewConnection(messageChan chan []byte) {
	broker.newClients <- messageChan
}

func (broker *Broker) CloseConnection(messageChan chan []byte) {
	broker.closingClients <- messageChan
}
