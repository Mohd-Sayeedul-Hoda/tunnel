package natserver

import (
	"sync"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/shared/config"

	"github.com/hashicorp/yamux"
)

func HandleTcpStream(cfg config.Config, session *yamux.Session) {
	// handle authentication
	// after that update the pool for connection so we can connect have presistance connection and way to manage connect so need state
}

type Connection struct {
	session *yamux.Session
}

type ConnectionsPool struct {
	pool map[string]*Connection
	mu   sync.RWMutex
}

// \r\n
// functionality in connectionPool to check if user existis or not
// connectionString = (user_id + connection) base64 encoded
// connectionString = user_id\r

func (c *ConnectionsPool) AddConnection(conn Connection) {
	c.mu.Lock()

}
