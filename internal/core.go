package internal

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	// DelayHorseStep is the length of the horse's "stride"
	DelayHorseStep time.Duration = time.Duration(300 * time.Millisecond)
	// DelayRefreshScreen is the time between screen refreshes
	DelayRefreshScreen time.Duration = time.Duration(100 * time.Millisecond)
	// HorseLabelDefault is the default label for the horses
	HorseLabelDefault string = "H"
	// HorseSpeedMax is the maximum speed of the horses
	HorseSpeedMax int = 5
	// HorseSpeedMin is the minimum speed of the horses
	HorseSpeedMin int = 1
	// HorseQuantityDefault is the default number of horses
	HorseQuantityDefault int = 2
	// HorseQuantityMax is the maximum number of horses
	HorseQuantityMax int = 99
	// HorseQuantityMin is the minimum number of horses
	HorseQuantityMin int = 2
	// ScoreTargetDefault is the default target score
	ScoreTargetDefault int = 75
	// ScoreTargetMin is the minimum target score
	ScoreTargetMin int = 15
	// ScoreTargetMax is the maximum target score
	ScoreTargetMax int = 100
	// GameTimeoutDefault is the default game timeout
	GameTimeoutDefault string = "10s"
	// GameTimeoutMin is the min game timeout
	GameTimeoutMin int = 10
	// GameTimeoutMax is the max game timeout
	GameTimeoutMax int = 90
	// GameTimeoutRegexPattern is the regex pattern for the game timeout
	GameTimeoutRegexPattern string = `^\d{1,2}s$`
)

var (
	horseLabel          = HorseLabelDefault
	horseQuantity       = HorseQuantityDefault
	horseWinner         = &Horse{}
	scoreTarget         = ScoreTargetDefault
	gameTimeout         = GameTimeoutDefault
	gameTimeoutDuration = time.Duration(GameTimeoutMax) * time.Second
)

// Run is the point of start the game engine
func Run(input Input) {
	setHorseLabel(input.HorseLabel)
	setHorseQuantity(input.HorsesQuantity)
	setScoreTarget(input.ScoreTarget)
	setGameTimeout(input.GameTimeout)
	setGameTimeoutDuration()
	loadHorses(horseQuantity)
	chGameOver := make(chan bool)
	isGameOver := atomic.Bool{}

	go display()
	for _, horse := range horses {
		go goHorse(horse, &isGameOver, chGameOver)
	}

	outStr := "\x01"
	select {
	case <-chGameOver:
		outStr += clearTerminal() + "\n" + getRaceStr()
	case <-time.After(gameTimeoutDuration):
		outStr += "\x01Today is a very hot day, the horses are tired!"
	}
	close(chGameOver)

	if horseWinner.Score > 0 {
		outStr += "\n" + horseWinner.Winner()
	}

	fmt.Println(outStr)
}

func setHorseLabel(inputHorseLabel string) {
	horseLabel = HorseLabelDefault
	if len(inputHorseLabel) == 1 {
		horseLabel = inputHorseLabel
	}
}

func setScoreTarget(inputScoreTarget int) {
	scoreTarget = ScoreTargetDefault
	if isValidScoreTarget(inputScoreTarget) {
		scoreTarget = inputScoreTarget
	}
}

func setHorseQuantity(inputHorsesQuantity int) {
	horseQuantity = HorseQuantityDefault
	if inputHorsesQuantity >= HorseQuantityMin && inputHorsesQuantity <= HorseQuantityMax {
		horseQuantity = inputHorsesQuantity
	}
}

func setGameTimeout(inputGameTimeout string) {
	gameTimeout = GameTimeoutDefault
	r, err := regexp.Compile(GameTimeoutRegexPattern)
	if err == nil && r.MatchString(inputGameTimeout) {
		value, err := strconv.ParseInt(strings.TrimSuffix(inputGameTimeout, "s"), 10, 32)
		if value32 := int(value); err == nil &&
			value32 >= GameTimeoutMin &&
			value32 <= GameTimeoutMax {
			gameTimeout = inputGameTimeout
		}
	}
}

func setGameTimeoutDuration() {
	gameTimeoutDuration = time.Duration(GameTimeoutMin) * time.Second
	tmp, err := time.ParseDuration(gameTimeout)
	if err == nil {
		gameTimeoutDuration = tmp
	}
}

func loadHorses(quantity int) {
	clearHorses()
	if !isValidHorsesQuantity(quantity) {
		quantity = HorseQuantityDefault
	}

	for i := 0; i < quantity; i++ {
		index := i + 1
		prefix := ""
		if index < 10 {
			prefix = "0"
		}

		horses = append(horses, &Horse{Label: horseLabel + prefix + strconv.Itoa(index)})
	}
}

func goHorse(target *Horse, isGameOver *atomic.Bool, chGameOver chan bool) {
	if target.Score < 0 {
		target.Score = 0
	}

	ticker := time.NewTicker(DelayHorseStep)
	defer ticker.Stop()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		select {
		case <-ticker.C:
			target.Score += r.Intn(HorseSpeedMax) + HorseSpeedMin
			if target.Score >= scoreTarget {
				if isGameOver.Load() {
					return
				}

				isGameOver.Store(true)
				chGameOver <- true
				horseWinner = target
				return
			}
		case <-chGameOver:
			return
		}
	}
}

func display() {
	for {
		fmt.Println(clearTerminal() + "\n" + getRaceStr())
		if horseWinner.Score > 0 {
			break
		}
		time.Sleep(DelayRefreshScreen)
	}
}

func clearTerminal() string {
	return "\033[H\033[2J"
}

func getRaceStr() string {
	msg := ""
	msg += generateTrackMark(scoreTarget) + "\n"
	for _, horse := range horses {
		msg += generateHorseTrack(horse, scoreTarget) + "\n"
	}
	msg += generateTrackMark(scoreTarget) + "\n"
	return msg
}

func generateHorseTrack(horse *Horse, scoreTarget int) string {
	more := ""
	if !isValidScoreTarget(scoreTarget) {
		scoreTarget = ScoreTargetDefault
	}

	if scoreTarget-horse.Score > 0 {
		more = strings.Repeat(" ", scoreTarget-horse.Score-1)
	}

	less := strings.Repeat(".", horse.Score)
	return fmt.Sprintf("%s|%v%v%v|", horse.Label, less, horse.Label, more)
}

func generateTrackMark(scoreTarget int) string {
	if !isValidScoreTarget(scoreTarget) {
		scoreTarget = ScoreTargetDefault
	}

	temp := ""
	for i := 1; i <= scoreTarget+2; i++ {
		if i%10 == 0 {
			temp += "|"
			continue
		}
		temp += "-"
	}

	return "   +" + temp + "+"
}

func clearHorses() {
	horses = []*Horse{}
}

func isValidScoreTarget(scoreTarget int) bool {
	return scoreTarget >= ScoreTargetMin && scoreTarget <= ScoreTargetMax
}

func isValidHorsesQuantity(quantity int) bool {
	return quantity >= HorseQuantityMin && quantity <= HorseQuantityMax
}
