package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWinner(t *testing.T) {
	horse := &Horse{
		Label: "A",
		Score: 10,
	}
	assert.Equal(t, "The horse winner is: A - Score 10", horse.Winner())
}
