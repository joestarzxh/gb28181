// Code generated by godddx, DO AVOID EDIT.
package sms

// Storer data persistence
type Storer interface {
	MediaServer() MediaServerStorer
}

// Core business domain
type Core struct {
	storer Storer
	*NodeManager
}

// NewCore create business domain
func NewCore(store Storer) Core {
	return Core{
		storer: store,

		NodeManager: NewNodeManager(store),
	}
}
