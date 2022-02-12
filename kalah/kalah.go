package kalah

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// board keeps track of the state of the board
type board struct {
	size   uint64    // how many non-store cups are there on each side of the board
	cups   []uint64  // the actual cups
	count  uint64    // the number of actual cups
	w      io.Writer // where to print output
	sleepy bool      // whether to sleep on b.wut()
}

// NewBoard constructs a new board of a given size
func NewBoard(size uint64, w io.Writer, sleepy bool) *board {
	c := (size * 2) + 2
	b := board{
		size:   size,
		cups:   make([]uint64, c),
		count:  c,
		w:      w,
		sleepy: sleepy,
	}

	var val uint64
	if size >= 3 {
		val = size - 2
	} else {
		val = 1
	}

	for i := uint64(0); i < c; i++ {
		if !(i == b.size || i == b.count-1) {
			b.cups[i] = val
		}
	}

	return &b
}

// Print prints the board in BEAUTIFUL MODERN TECHNICOLOR (on ANSI compatible
// machines) otherwise it just looks like garbage
func (b *board) Print() {
	//	b.log("\033[2J\033[H")
	b.log("\n")
	for i := int(b.size) - 1; i >= 0; i-- {
		b.log("\t\033[1;32m%d\033[0m", b.cups[i])
	}

	b.log("\n\033[1;36mcpu: %d |\033[0m", b.cups[b.size])

	p := "-\t"
	for i := uint64(0); i < b.size; i++ {
		if i == b.size-1 {
			p = "-"
		}
		b.log("\033[1;36m%s\033[0m", p)
	}

	b.log("\033[1;36m| player: %d\033[0m\n\t", b.cups[b.count-1])
	for i := b.size + 1; i < b.count-1; i++ {
		b.log("\033[1;32m%d\033[0m\t", b.cups[i])
	}

	b.log("\n")
}

// GameOver tells whether one side of the board is empty
func (b *board) GameOver() bool {
	cpuSum := -1
	playerSum := -1
	for i := uint64(0); i < b.size; i++ {
		cpuSum += int(b.cups[i])
		playerSum += int(b.cups[i+1+b.size])
	}

	if cpuSum < 0 || playerSum < 0 {
		return true
	}

	return false
}

// PrintWinner prints who the winner is
func (b *board) PrintWinner() {
	msg := "cpu wins!"
	if player, cpu := b.cups[b.count-1], b.cups[b.size]; player > cpu {
		msg = "you win!"
	} else if player == cpu {
		msg = "it's a tie!"
	}

	b.log("%s\n", msg)
}

// Handle is used to handle raw input
func (b *board) Handle(input string) bool {
	i := strings.TrimSpace(input)
	switch i {
	case "exit", "quit", "q", "x":
		return true
	}

	cup, err := strconv.ParseUint(i, 10, 64)
	if err != nil {
		b.wut()
		return false
	}

	if b.move(cup, "player") {
		if !b.GameOver() {
			b.cpuPlays()
		}
	}
	return b.GameOver()
}

// move moves the cup indicated by `c` around the board
// and returns whether the turn is over
func (b *board) move(c uint64, who string) bool {
	var store, zeroth uint64
	switch who {
	case "player":
		store = b.count - 1
		zeroth = b.size + 1
	case "cpu":
		store = b.size
		zeroth = 0
	default:
		return false
	}
	c = store - c

	// make sure the cup is within the bounds of play
	if c >= store || c < zeroth || b.cups[c] == 0 {
		if who == "player" {
			b.wut()
		}
		return false
	}

	// move the pieces
	i := c
	for b.cups[c] > 0 {
		i++
		if i > b.count-1 {
			i = 0
		}
		// don't place pieces in the opponents store cup
		if (who == "player" && i == b.size) || (who == "cpu" && i == b.count-1) {
			continue
		}
		b.cups[i] += 1
		b.cups[c] -= 1
	}

	// your turn is only over if you did not land in your store cup
	turnOver := (i != store)

	if turnOver && b.cups[i] == 1 {
		// this means you did not land in your store cup and you only have one
		// piece in the last cup you landed in, then you have a capture
		var opp uint64
		var canCapture bool
		switch who {
		case "player":
			canCapture = i > b.size
			opp = abs(int(i) - int(b.size*2))
		case "cpu":
			canCapture = i < b.size
			opp = uint64(int(b.size) + (int(b.size) - int(i)))
		}

		if canCapture && b.cups[opp] != 0 {
			// take the opponent's pieces
			val := b.cups[opp]
			b.cups[opp] -= val
			b.cups[store] += val

			// take your pieces
			val = b.cups[i]
			b.cups[i] -= val
			b.cups[store] += val

			// if you capture a piece, your turn is not over
			turnOver = false
		}
	}

	// if the game is over, so is your turn
	return turnOver || b.GameOver()
}

// cpuPlays handles the cpu's side of things
func (b *board) cpuPlays() {
	var pick uint64
	for {
		b.Print()
		pick = b.cpuPick()
		b.log("cpu: %v\n", pick)
		if b.move(pick, "cpu") {
			break
		}
	}
}

// cpuPick is the brain of the computer's play.
// it picks moves in this order:
// 1. the capture move with the most points potential (if any)
// 2. the first piece to land in the store
// 3. the first non-empty piece
func (b *board) cpuPick() uint64 {
	var landInStore, captures, capturePts, nonEmpty []uint64
	for cup := uint64(1); cup <= b.size; cup++ {
		idx := b.size - cup
		if b.cups[idx] == cup {
			landInStore = append(landInStore, cup)
		}
		if b.cups[idx] != 0 {
			nonEmpty = append(nonEmpty, cup)
			endCup := (idx + b.cups[idx]) % b.count
			opposite := uint64(int(b.size) + (int(b.size) - int(idx)))
			if b.cups[endCup] == 0 && endCup < b.size && b.cups[opposite] != 0 {
				captures = append(captures, cup)
				capturePts = append(capturePts, b.cups[opposite]+1)
			}
		}
	}

	if len(captures) > 0 {
		maxIdx := 0
		var maxPts uint64
		for i, points := range capturePts {
			if points > maxPts {
				maxPts = points
				maxIdx = i
			}
		}
		return captures[maxIdx]
	}

	if len(landInStore) > 0 {
		return landInStore[0]
	}

	if len(nonEmpty) > 0 {
		return nonEmpty[0]
	}

	return uint64(rand.Intn(int(b.size-1))) + 1
}

// log is a convenient way to print stuff to our writer
func (b *board) log(format string, args ...interface{}) {
	fmt.Fprintf(b.w, format, args...)
}

// wut prints our Very Helpful™ help message
func (b *board) wut() {
	b.log("?\n")
	if b.sleepy {
		time.Sleep(500 * time.Millisecond)
	}
}

// absolute value convenience function
func abs(i int) uint64 {
	if i < 0 {
		return uint64(i * -1)
	}
	return uint64(i)
}
