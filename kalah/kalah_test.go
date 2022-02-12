package kalah

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

// TestNewBoard makes sure the board gets setup right
func TestNewBoard(t *testing.T) {
	for _, test := range []uint64{1, 23, 4, 5, 6, 7, 7, 8} {
		t.Run(fmt.Sprintf("test%d", test), func(t *testing.T) {
			b := NewBoard(test, os.Stdout, false)
			if b.size != test {
				t.Errorf("expected %v got %v", test, b.size)
			}
			count := (test * 2) + 2
			if got := uint64(len(b.cups)); got != count {
				t.Errorf("expected %v got %v", got, count)
			}
			if count != b.count {
				t.Errorf("expected %v got %v", b.count, count)
			}
			if test >= 3 && b.cups[0] != test-2 {
				t.Errorf("expected %d got %d", test-2, b.cups[0])
			} else if test < 3 && b.cups[0] != 1 {
				t.Errorf("expected %d got %d", 1, b.cups[0])
			}
		})
	}
}

// TestGameOver makes sure the game ends when it's supposed to
func TestGameOver(t *testing.T) {
	cases := map[string]struct {
		board board
		over  bool
	}{
		"playerDone": {
			board: board{
				size: 6,
				cups: []uint64{
					4, 4, 4, 4, 4, 4, 4,
					0, 0, 0, 0, 0, 0, 13,
				},
				count: 14,
			},
			over: true,
		},
		"cpuDone": {
			board: board{
				size: 6,
				cups: []uint64{
					0, 0, 0, 0, 0, 0, 4,
					4, 4, 4, 4, 4, 4, 13,
				},
				count: 14,
			},
			over: true,
		},
		"gameStillGoing": {
			board: board{
				size: 6,
				cups: []uint64{
					4, 4, 4, 0, 0, 0, 4,
					4, 4, 4, 4, 4, 4, 13,
				},
				count: 14,
			},
			over: false,
		},
	}
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			if got := test.board.GameOver(); got != test.over {
				t.Errorf("expected %v got %v", test.over, got)
			}
		})
	}
}

// TestMoves makes sure the move validation is working right at least on a 5
// cup board
func TestMoves(t *testing.T) {
	for name, test := range map[string]struct {
		cup      uint64
		turnOver bool
		who      string
		testCups []uint64
		resCups  []uint64
		board    *board
	}{
		"validPlayerMove": {
			cup:      5,
			turnOver: true,
			who:      "player",
			testCups: []uint64{
				0, 3, 3, 3, 3, 0,
				3, 0, 3, 3, 3, 0,
			},
			resCups: []uint64{
				0, 3, 3, 3, 3, 0,
				0, 1, 4, 4, 3, 0,
			},
		},
		"validPlayerCapture": {
			cup:      2,
			turnOver: false,
			who:      "player",
			testCups: []uint64{
				3, 3, 3, 3, 3, 0,
				0, 1, 4, 1, 0, 0,
			},
			resCups: []uint64{
				0, 3, 3, 3, 3, 0,
				0, 1, 4, 0, 0, 4,
			},
		},
		"invalidPlayerMove": {
			cup:      4,
			turnOver: false,
			who:      "player",
			testCups: []uint64{
				0, 3, 3, 3, 3, 0,
				3, 0, 3, 3, 3, 0,
			},
			resCups: []uint64{
				0, 3, 3, 3, 3, 0,
				3, 0, 3, 3, 3, 0,
			},
		},
		"invalidCpuMove": {
			cup:      6,
			turnOver: false,
			who:      "cpu",
			testCups: []uint64{
				0, 3, 3, 3, 3, 0,
				3, 0, 3, 3, 3, 0,
			},
			resCups: []uint64{
				0, 3, 3, 3, 3, 0,
				3, 0, 3, 3, 3, 0,
			},
		},
		"invalidCpuMove2": {
			cup:      5,
			turnOver: false,
			who:      "cpu",
			testCups: []uint64{
				0, 3, 3, 3, 3, 0,
				3, 0, 3, 3, 3, 0,
			},
			resCups: []uint64{
				0, 3, 3, 3, 3, 0,
				3, 0, 3, 3, 3, 0,
			},
		},
		"validCpuMove1": {
			cup:      4,
			turnOver: true,
			who:      "cpu",
			testCups: []uint64{
				0, 3, 3, 3, 3, 0,
				3, 0, 3, 3, 3, 0,
			},
			resCups: []uint64{
				0, 0, 4, 4, 4, 0,
				3, 0, 3, 3, 3, 0,
			},
		},
		"validCpuMove2": {
			cup:      5,
			turnOver: true,
			who:      "cpu",
			testCups: []uint64{
				3, 3, 3, 3, 3, 0,
				3, 0, 3, 3, 3, 0,
			},
			resCups: []uint64{
				0, 4, 4, 4, 3, 0,
				3, 0, 3, 3, 3, 0,
			},
		},
		"validCpuCapture": {
			cup:      3,
			turnOver: false,
			who:      "cpu",
			testCups: []uint64{
				1, 2, 1, 0, 5, 0,
				3, 3, 3, 3, 3, 0,
			},
			resCups: []uint64{
				1, 2, 0, 0, 5, 4,
				3, 0, 3, 3, 3, 0,
			},
		},
		"badInput": {
			cup:      0,
			turnOver: false,
			who:      "your mom",
			testCups: []uint64{
				1, 2, 1, 0, 5, 0,
				3, 3, 3, 3, 3, 0,
			},
			resCups: []uint64{
				1, 2, 1, 0, 5, 0,
				3, 3, 3, 3, 3, 0,
			},
		},
		"playerCapturesWith1": {
			cup:      2,
			turnOver: false,
			who:      "player",
			testCups: []uint64{
				8, 0, 1, 1, 0, 13, 8,
				2, 1, 0, 0, 1, 0, 13,
			},
			resCups: []uint64{
				0, 0, 1, 1, 0, 13, 8,
				2, 1, 0, 0, 0, 0, 22,
			},
			board: NewBoard(6, os.Stdout, false),
		},
		"playerMoves1Piece4": {
			cup:      4,
			turnOver: true,
			who:      "player",
			testCups: []uint64{
				7, 0, 7, 7, 0, 8, 6,
				0, 0, 1, 1, 0, 0, 20,
			},
			resCups: []uint64{
				7, 0, 7, 7, 0, 8, 6,
				0, 0, 0, 2, 0, 0, 20,
			},
			board: NewBoard(6, os.Stdout, false),
		},
		"playerMoves1Piece1": {
			cup:      1,
			turnOver: false,
			who:      "player",
			testCups: []uint64{
				0, 0, 0, 0, 1, 0, 6,
				3, 0, 1, 1, 1, 1, 34,
			},
			resCups: []uint64{
				0, 0, 0, 0, 1, 0, 6,
				3, 0, 1, 1, 1, 0, 35,
			},
			board: NewBoard(6, os.Stdout, false),
		},
		"cpuCapturesEndsGame": {
			cup:      2,
			turnOver: true,
			who:      "cpu",
			testCups: []uint64{
				0, 0, 0, 0, 1, 0, 6,
				3, 0, 1, 1, 1, 0, 35,
			},
			resCups: []uint64{
				0, 0, 0, 0, 0, 0, 10,
				0, 0, 1, 1, 1, 0, 35,
			},
			board: NewBoard(6, os.Stdout, false),
		},
		"playerMovesAPieceCaptures": {
			cup:      3,
			turnOver: true,
			who:      "player",
			testCups: []uint64{
				0, 4, 0, 1, 2, 1, 28,
				0, 0, 0, 1, 0, 0, 11,
			},
			resCups: []uint64{
				0, 0, 0, 1, 2, 1, 28,
				0, 0, 0, 0, 0, 0, 16,
			},
			board: NewBoard(6, os.Stdout, false),
		},
		"playerAlmostCreatesACaptureButDoesnt": {
			cup:      2,
			turnOver: true,
			who:      "player",
			testCups: []uint64{
				1, 2, 0, 9, 0, 1, 9,
				0, 0, 2, 2, 11, 2, 9,
			},
			resCups: []uint64{
				2, 3, 1, 10, 1, 2, 9,
				1, 1, 3, 2, 0, 3, 10,
			},
			board: NewBoard(6, os.Stdout, false),
		},
	} {
		t.Run(name, func(t *testing.T) {
			var b *board
			if test.board == nil {
				b = NewBoard(5, os.Stdout, false)
			} else {
				b = test.board
			}
			if len(test.testCups) > 0 {
				b.cups = test.testCups
			}
			if got := b.move(test.cup, test.who); got != test.turnOver {
				t.Errorf("expected %v got %v", test.turnOver, got)
			}
			for i, cup := range b.cups {
				if test.resCups[i] != cup {
					t.Errorf("on idx %v: expected %v got %v", i, test.resCups[i], cup)
				}
			}
		})
	}
}

// TestHandle makes sure I'm handling things okay
func TestHandle(t *testing.T) {
	for name, expected := range map[string]bool{
		"exit": true,
		"quit": true,
		"q":    true,
		"x":    true,
		"e":    false,
		"💩":    false,
		" ":    false,
		"-":    false,
		"5":    false,
		"6":    false,
		"3":    false,
	} {
		t.Run(name, func(t *testing.T) {
			b := NewBoard(5, os.Stdout, false)
			b.cups[0] = 0 // get rid of cup number 5 for cpu
			b.cups[7] = 0 // get rid of cup number 4 for player
			if got := b.Handle(name); got != expected {
				t.Errorf("expected %v got %v", expected, got)
			}
		})
	}
}

// TestAbs makes sure i'm getting my crunches in and also that i can do basic math
func TestAbs(t *testing.T) {
	for _, test := range []struct {
		in  int
		out uint64
	}{
		{12345, 12345},
		{-12345, 12345},
		{0, 0},
	} {
		t.Run(fmt.Sprintf("%d", test.in), func(t *testing.T) {
			if got := abs(test.in); got != test.out {
				t.Errorf("expected %v got %v", test.out, got)
			}
		})
	}
}

// TestWut makes sure the error message is nice and clear
func TestWut(t *testing.T) {
	var buf bytes.Buffer
	b := NewBoard(2, &buf, false)
	b.wut()

	if buf.String() != "?\n" {
		t.Errorf("expected '?\\n', got '%v'", buf.String())
	}
}
