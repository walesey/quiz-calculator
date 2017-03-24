package main

import (
	"fmt"
	"math"
	"sort"
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
calculationTime(4m19.8700246s)
Results:
225,209,217
225,217,209
197,205,213
197,213,205
199,175,223
199,223,175
205,197,213
205,213,197
175,199,223
175,223,199
209,217,225
209,225,217
213,197,205
213,205,197
217,209,225
217,225,209
223,175,199
223,199,175
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
	results = removeDupes(results)
	fmt.Printf("calculationTime(%v)\n", time.Since(start))
	fmt.Println("Results:")
	for _, result := range results {

		fmt.Printf("%v,%v,%v\n", result[0], result[1], result[2])
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

// sort and Dedupe
type BySize [3]int

func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i] < a[j] }

func removeDupes(results [][3]int) [][3]int {
	for i, _ := range results {
		sort.Sort(BySize(results[i]))
	}
	deduped := [][3]int{}
OuterLoop:
	for _, r := range results {
		for _, d := range deduped {
			if r[0] == d[0] && r[1] == d[1] && r[2] == d[2] {
				continue OuterLoop
			}
		}
		deduped = append(deduped, r)
	}
	return deduped
}
