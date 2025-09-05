package internal

import "fmt"

// Input represents the input data for the game
type Input struct {
	HorseLabel     string
	HorsesQuantity int
	ScoreTarget    int
	GameTimeout    string
}

// Horse represent the horse entity
type Horse struct {
	Label string
	Score int
}

var horses = []*Horse{}

// Winner greeting to the champion horse
func (h *Horse) Winner() string {
	return fmt.Sprintf("The horse winner is: %s - Score %d", h.Label, h.Score)
}
