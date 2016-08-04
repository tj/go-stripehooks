package stripehooks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stripe/stripe-go"
)

func handle(e *stripe.Event) error {
	return nil
}

func TestManager_Registered(t *testing.T) {
	m := New()
	assert.False(t, m.Registered(ChargeCaptured))
	m.Handle(ChargeCaptured, HandlerFunc(handle))
	assert.True(t, m.Registered(ChargeCaptured))
}

func TestManager_HandleFunc(t *testing.T) {
	m := New()
	assert.False(t, m.Registered(ChargeCaptured))
	m.HandleFunc(ChargeCaptured, handle)
	assert.True(t, m.Registered(ChargeCaptured))
}

func TestManager_HandleEvent(t *testing.T) {
	m := New()
	m.Verify = false

	called := false

	m.HandleFunc(ChargeCaptured, func(e *stripe.Event) error {
		assert.Equal(t, "evt_18eWOZH5vuhJUCbaZ7He0cFF", e.ID)
		called = true
		return nil
	})

	err := m.HandleEvent(&stripe.Event{
		ID:   "evt_18eWOZH5vuhJUCbaZ7He0cFF",
		Type: ChargeCaptured,
	})

	assert.NoError(t, err)
	assert.True(t, called, "should invoke the handler")
}
