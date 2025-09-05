package internal

import (
	"bytes"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetHorseLabel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"A", "A"},
		{"", HorseLabelDefault},
		{"AB", HorseLabelDefault},
	}

	for _, test := range tests {
		// Reset global state for this test
		horseLabel = HorseLabelDefault
		setHorseLabel(test.input)
		assert.Equal(t, test.expected, horseLabel, "setHorseLabel(%q)", test.input)
	}
}

func TestSetHorseQuantity(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{0, HorseQuantityDefault},
		{100, HorseQuantityDefault},
	}

	for _, test := range tests {
		horseQuantity = HorseQuantityDefault
		setHorseQuantity(test.input)
		assert.Equal(t, test.expected, horseQuantity, "setHorseQuantity(%q)", test.input)
	}
}

func TestSetScoreTarget(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{50, 50},
		{5, ScoreTargetDefault},
		{150, ScoreTargetDefault},
	}

	for _, test := range tests {
		scoreTarget = ScoreTargetDefault
		setScoreTarget(test.input)
		assert.Equal(t, test.expected, scoreTarget, "setScoreTarget(%q)", test.input)
	}
}

func TestSetGameTimeout(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"15s", "15s"},
		{"0s", GameTimeoutDefault},
		{"120s", GameTimeoutDefault},
		{"what", GameTimeoutDefault},
	}

	for _, test := range tests {
		gameTimeout = GameTimeoutDefault
		setGameTimeout(test.input)
		assert.Equal(t, test.expected, gameTimeout, "setGameTimeout(%q)", test.input)
	}
}

func TestSetGameTimeoutDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0s", GameTimeoutDefault},
		{"20s", "20s"},
		{GameTimeoutDefault, "10s"},
		{"120s", "10s"},
	}

	for _, test := range tests {
		gameTimeout = GameTimeoutDefault
		setGameTimeout(test.input)
		setGameTimeoutDuration()
		assert.Equal(t, test.expected, gameTimeoutDuration.String(), "setGameTimeoutDuration() with input %q", test.input)
	}
}

func TestLoadHorses(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{1, HorseQuantityDefault},
		{100, HorseQuantityDefault},
	}

	for _, test := range tests {
		clearHorses()
		loadHorses(test.input)
		assert.Equal(t, test.expected, len(horses), "loadHorses(%q)", test.input)
	}
}

func TestGenerateHorseTrack(t *testing.T) {
	tests := []struct {
		inputHorse Horse
		inputScore int
		expected   string
	}{
		{Horse{Label: "A01", Score: 5}, 75, "A01|.....A01                                                                     |"},
		{Horse{Label: "B01", Score: 30}, 30, "B01|..............................B01|"},
		{Horse{Label: "C", Score: 0}, 25, "C|C                        |"},
		{Horse{Label: "D101", Score: 3}, 25, "D101|...D101                     |"},
	}

	for _, test := range tests {
		track := generateHorseTrack(&test.inputHorse, test.inputScore)
		assert.Equal(t, test.expected, track, "generateHorseTrack(%+v, %v)", test.inputHorse, test.inputScore)
	}
}

func TestGenerateTrackMark(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{75, "   +---------|---------|---------|---------|---------|---------|---------|-------+"},
		{15, "   +---------|-------+"},
		{28, "   +---------|---------|---------|+"},
		{100, "   +---------|---------|---------|---------|---------|---------|---------|---------|---------|---------|--+"},
		{150, "   +---------|---------|---------|---------|---------|---------|---------|-------+"}, // Invalid input, falls back to default
	}

	for _, test := range tests {
		mark := generateTrackMark(test.input)
		assert.Equal(t, test.expected, mark, "generateTrackMark(%v)", test.input)
	}
}

func TestGetRaceStr(t *testing.T) {
	tests := []struct {
		inputHorses      []*Horse
		inputScoreTarget int
		expected         string
	}{
		{[]*Horse{
			{Label: "A01", Score: 68},
			{Label: "B02", Score: 75},
		},
			0,
			"   +---------|---------|---------|---------|---------|---------|---------|-------+\n" + "A01|....................................................................A01      |\n" + "B02|...........................................................................B02|\n" +
				"   +---------|---------|---------|---------|---------|---------|---------|-------+\n"},
		{[]*Horse{
			{Label: "A01", Score: 5},
			{Label: "B02", Score: 10},
			{Label: "C03", Score: 15},
		},
			20,
			"   +---------|---------|--+\n" +
				"A01|.....A01              |\n" +
				"B02|..........B02         |\n" +
				"C03|...............C03    |\n" +
				"   +---------|---------|--+\n"},
		{[]*Horse{
			{Label: "EEE", Score: 12},
			{Label: "AAA", Score: 20},
			{Label: "HHH", Score: 25},
		},
			22,
			"   +---------|---------|----+\n" +
				"EEE|............EEE         |\n" +
				"AAA|....................AAA |\n" +
				"HHH|.........................HHH|\n" +
				"   +---------|---------|----+\n"},
	}

	for _, test := range tests {
		// Setup test-specific state
		originalHorses := horses
		originalScoreTarget := scoreTarget
		horses = test.inputHorses
		scoreTarget = test.inputScoreTarget

		raceStr := getRaceStr()
		assert.Equal(t, test.expected, raceStr)

		// Teardown: Restore global state
		horses = originalHorses
		scoreTarget = originalScoreTarget
	}
}

func TestGoHorse(t *testing.T) {
	tests := []struct {
		inputHorse *Horse
		expected   int
		gameOver   bool
	}{
		{&Horse{Label: "A01", Score: -99}, 15, false},
		{&Horse{Label: "B02", Score: 20}, 25, false},
		{&Horse{Label: "C02", Score: 100}, 75, true}, // trigger timeout
	}

	chGameOver := make(chan bool, 1)
	for _, test := range tests {
		isGameOver := atomic.Bool{}
		isGameOver.Store(test.gameOver)
		scoreTarget = test.expected

		go goHorse(test.inputHorse, &isGameOver, chGameOver)

		select {
		case <-chGameOver:
			assert.GreaterOrEqual(t, test.inputHorse.Score, scoreTarget, "Horse score should be greater or equal to target score")
			assert.True(t, isGameOver.Load(), "isGameOver should be true")
		case <-time.After(2 * time.Second):
			assert.True(t, isGameOver.Load(), "goHorse() timed out after 2s with isGameOver false")
		}
	}
	close(chGameOver)
}

func TestClearTerminal(t *testing.T) {
	assert.Equal(t, "\033[H\033[2J", clearTerminal())
}

func TestDisplay(t *testing.T) {
	// Setup: Isolate global variables for this test
	originalHorseWinner := horseWinner
	originalScoreTarget := scoreTarget
	originalHorses := horses
	horseWinner = &Horse{Label: "A01", Score: 15}
	chGameOver := make(chan bool)
	scoreTarget = 15
	loadHorses(2)

	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
		r.Close()
		w.Close()
		// Teardown: Restore global state
		horseWinner = originalHorseWinner
		scoreTarget = originalScoreTarget
		horses = originalHorses
	}()

	go display()
	time.Sleep(DelayRefreshScreen * 2) // Allow display to run at least once
	close(chGameOver)
	w.Close()

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expected := "\x1b[H\x1b[2J\n   +---------|-------+\nH01|H01              |\nH02|H02              |\n   +---------|-------+\n\n"
	assert.Contains(t, output, expected)
}
