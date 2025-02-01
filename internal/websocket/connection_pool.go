package websocket

import (
	"context"
	"sync"

	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/coder/websocket"
)

var connectionPool = NewConnectionPool()

// ConnectionPool struct to manage connections thread-safely
type ConnectionPool struct {
	mu          sync.RWMutex
	connections map[*websocket.Conn]*db_model.User
}

// NewConnectionPool creates a new ConnectionPool
func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		connections: make(map[*websocket.Conn]*db_model.User),
	}
}

// Add adds a connection to the ConnectionPool
func (cp *ConnectionPool) Add(conn *websocket.Conn, user *db_model.User) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.connections[conn] = user
}

// Remove removes a connection from the ConnectionPool
func (cp *ConnectionPool) Remove(conn *websocket.Conn) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	delete(cp.connections, conn)
}

// GetUser returns the user associated with a connection
func (cp *ConnectionPool) GetUser(conn *websocket.Conn) *db_model.User {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	return cp.connections[conn]
}

// Broadcast sends a message to all connections in the ConnectionPool
func (cp *ConnectionPool) Broadcast(ctx context.Context, message []byte) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	for conn := range cp.connections {
		conn.Write(ctx, websocket.MessageText, message)
	}
}

// CheckAlive checks if connections are still alive, removes dead connections
func (cp *ConnectionPool) CheckAlive(ctx context.Context) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	for conn := range cp.connections {
		err := conn.Ping(ctx)
		if err != nil {
			delete(cp.connections, conn)
		}
	}
}
