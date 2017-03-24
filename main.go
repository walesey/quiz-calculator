package main

import (
	"fmt"
	"math"
	"time"
)

/*
Three people are playing the following betting game.
Every five minutes, a turn takes place in which a random player rests and the other two bet
against one another with all of their money.
The player with the smaller amount of money always wins,
doubling his money by taking it from the loser.
For example, if the initial amounts of money are 1, 4, and 6,
then the result of the first turn can be either
2,3,6 (1 wins against 4);
1,8,2 (4 wins against 6); or
2,4,5 (1 wins against 6).
If two players with the same amount of money play against one another,
the game immediately ends for all three players.
Find initial amounts of money for the three players, where none of the three has more than 255,
and in such a way that the game cannot end in less than one hour. (So at least 12 turns)
In the example above (1,4,6), there is no way to end the game in less than 15 minutes.
All numbers must be positive integers.
*/

type GameState struct {
	turnCount      int
	players        [3]int
	randomSequence *RandomSequence
}

func NewGame(p1, p2, p3 int, randomSequence *RandomSequence) *GameState {
	return &GameState{
		players:        [3]int{p1, p2, p3},
		randomSequence: randomSequence,
	}
}

type RandomSequence struct {
	index    int
	sequence []int
}

func NewRandomSequence(turns, seed int) *RandomSequence {
	rs := &RandomSequence{
		sequence: make([]int, turns),
	}
	for i := 0; i < turns; i++ {
		rs.sequence[i] = seed % 3
		seed = seed / 3
	}
	return rs
}

func (this *RandomSequence) next() int {
	result := this.sequence[this.index]
	this.index++
	return result
}

func (this *GameState) resolve(winner, loser int) {
	// fmt.Println("winner ", winner, this.players[winner])
	// fmt.Println("loser ", loser, this.players[loser])
	// fmt.Println("")
	this.players[loser] = this.players[loser] - this.players[winner]
	this.players[winner] = this.players[winner] * 2
}

func (this *GameState) turn() bool {
	this.turnCount = this.turnCount + 1
	rest := this.randomSequence.next() % 3
	p1, p2 := 0, 1
	if rest == 0 {
		p1, p2 = 1, 2
	} else if rest == 1 {
		p1, p2 = 0, 2
	}

	if this.players[p1] == this.players[p2] {
		return true // game ends
	}

	if this.players[p1] < this.players[p2] {
		this.resolve(p1, p2)
	} else {
		this.resolve(p2, p1)
	}
	return false
}

func main() {
	results := simulateAllStates(4)
	for turns := 5; turns <= 12; turns++ {
		results = simulateSubset(turns, results)
	}
}

func simulateAllStates(turns int) [][3]int {
	start := time.Now()
	states := [][3]int{}
	for p1 := 1; p1 < 256; p1++ {
		for p2 := 1; p2 < 256; p2++ {
			for p3 := 1; p3 < 256; p3++ {
				if simulatePermutations(turns, p1, p2, p3) {
					states = append(states, [3]int{p1, p2, p3})
				}
			}
		}
	}
	fmt.Printf("simulateAllStates: turns(%v) : results(%v) time(%v)\n", turns, len(states), time.Since(start))
	return states
}

func simulateSubset(turns int, subset [][3]int) [][3]int {
	start := time.Now()
	states := [][3]int{}
	for _, set := range subset {
		if simulatePermutations(turns, set[0], set[1], set[2]) {
			states = append(states, set)
		}
	}
	fmt.Printf("simulateAllStates: turns(%v) : results(%v) time(%v)\n", turns, len(states), time.Since(start))
	return states
}

func simulatePermutations(turns, p1, p2, p3 int) bool {
	length := int(math.Pow(3, float64(turns)))
Permutations:
	for p := 0; p < length; p++ {
		game := NewGame(p1, p2, p3, NewRandomSequence(turns, p))
		end := false
		for !end {
			end = game.turn()
			if game.turnCount == turns {
				continue Permutations
			}
		}
		return false
	}
	return true
}

// fmt.Printf("Turns: %v, Game Ended : %v : %v : %v\n", game.turnCount, p1, p2, p3)
// fmt.Printf("Final: : %v : %v : %v\n", game.players[0], game.players[1], game.players[2])
