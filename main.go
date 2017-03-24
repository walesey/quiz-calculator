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

/////////////////////////////////////////
/*
The following go program solves the above problem by brute force using a go worker pool
Output:
calculationTime(6m1.2539174s)
Results:
0: 197,205,213
1: 197,213,205
2: 199,175,223
3: 199,223,175
4: 205,197,213
5: 205,213,197
6: 209,217,225
7: 209,225,217
8: 213,197,205
9: 213,205,197
10: 217,209,225
11: 217,225,209
12: 223,175,199
13: 223,199,175
14: 225,209,217
15: 225,217,209
16: 175,199,223
17: 175,223,199
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

func (this *GameState) resolve(winner, loser int) {
	this.players[loser] = this.players[loser] - this.players[winner]
	this.players[winner] = this.players[winner] * 2
}

// MAIN
func main() {
	start := time.Now()
	output := make(chan [3]int, 1000)
	progress := make(chan int)

	go func() {
		counter := 0
		for p := range progress {
			counter++
			percent := (counter * 100) / 256
			fmt.Printf("Percent Complete: %v (%v)\n", percent, p)
			if counter == 255 {
				close(output)
				close(progress)
			}
		}
	}()

	nbWorkers := 8
	for i := 1; i < 256; i += (256 / nbWorkers) {
		j := i + (256 / nbWorkers)
		if j > 256 {
			j = 256
		}
		fmt.Printf("Simulating: %v to %v\n", i, j)
		go simulateAllStates(12, i, j, output, progress)
	}

	var results [][3]int
	for result := range output {
		fmt.Printf("Valid: %v,%v,%v\n", result[0], result[1], result[2])
		results = append(results, result)
	}
	fmt.Printf("calculationTime(%v)\n", time.Since(start))
	fmt.Println("Results:")
	for i, result := range results {
		fmt.Printf("%v: %v,%v,%v\n", i, result[0], result[1], result[2])
	}
}

func simulateAllStates(turns, from, to int, output chan [3]int, progress chan int) {
	for p1 := from; p1 < to; p1++ {
		for p2 := 1; p2 < 256; p2++ {
			for p3 := 1; p3 < 256; p3++ {
				if simulatePermutations(turns, p1, p2, p3) {
					output <- [3]int{p1, p2, p3}
				}
			}
		}
		progress <- p1
	}
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
