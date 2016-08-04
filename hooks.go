// Package stripehooks provides hook management with optional event fetching to verify that the origin is Stripe, this
// functionality is defaulted to true.
package stripehooks

import (
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/event"
)

// EventType is the type of event for the hook.
type EventType string

// Handler is a Stripe event handler.
type Handler interface {
	HandleStripeEvent(*stripe.Event) error
}

// HandlerFunc handles a stripe event.
type HandlerFunc func(*stripe.Event) error

// HandleStripeEvent implements Handler.
func (h HandlerFunc) HandleStripeEvent(e *stripe.Event) error {
	return h(e)
}

// Manager manages Stripe hooks.
type Manager struct {
	Verify bool // Verify the event by fetching from Stripe
	hooks  map[EventType]Handler
}

// New hook manager.
func New() *Manager {
	return &Manager{
		Verify: true,
		hooks:  make(map[EventType]Handler),
	}
}

// Handle registers an even handler.
func (m *Manager) Handle(kind EventType, h Handler) {
	m.hooks[kind] = h
}

// HandleFunc registers an even handler.
func (m *Manager) HandleFunc(kind EventType, h HandlerFunc) {
	m.hooks[kind] = h
}

// Registered returns true if a handler is registered.
func (m *Manager) Registered(kind EventType) bool {
	_, ok := m.hooks[kind]
	return ok
}

// HandleEvent handles a Stripe event, no-oping if a handler
// has not been defined. If Verify is true then the event is
// fetched from Stripe to validate its origin.
func (m *Manager) HandleEvent(e *stripe.Event) error {
	h, ok := m.hooks[EventType(e.Type)]
	if !ok {
		return nil
	}

	if m.Verify {
		verified, err := event.Get(e.ID, nil)
		if err != nil {
			return err
		}

		e = verified
	}

	return h.HandleStripeEvent(e)
}
