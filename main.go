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
calculationTime(6.337055566s)
Results:
175,199,223
197,205,213
209,217,225
*/

const turns = 11
const nbWorkers = 4

type RandomSequence []int

func NewRandomSequence(nbTurns, seed int) RandomSequence {
	rs := make([]int, nbTurns)
	for i := 0; i < nbTurns; i++ {
		rs[i] = seed % 3
		seed = seed / 3
	}
	return rs
}

func game(p1, p2, p3 int, sequence RandomSequence) bool {
	for t, rest := range sequence {
		if p1 == p2 || p2 == p3 || p3 == p1 {
			return t == turns
		}

		switch rest {
		case 0:
			if p2 > p3 {
				p2 = p2 - p3
				p3 = p3 * 2
			} else {
				p3 = p3 - p2
				p2 = p2 * 2
			}
		case 1:
			if p1 > p3 {
				p1 = p1 - p3
				p3 = p3 * 2
			} else {
				p3 = p3 - p1
				p1 = p1 * 2
			}
		case 2:
			if p1 > p2 {
				p1 = p1 - p2
				p2 = p2 * 2
			} else {
				p2 = p2 - p1
				p1 = p1 * 2
			}
		default:
			fmt.Println("error rest:", rest)
		}
	}
	return true
}

// MAIN
func main() {
	start := time.Now()
	input := make(chan [3]int, 1000)
	output := make(chan [3]int, 1000)

	simulator := NewSimulator()
	for i := 0; i < nbWorkers; i++ {
		go simulator.simulate(input, output)
	}

	go func() {
		counter := 0
		for p1 := 1; p1 < 256; p1++ {
			for p2 := p1; p2 < 256; p2++ {
				for p3 := p2; p3 < 256; p3++ {
					input <- [3]int{p1, p2, p3}
				}
			}
			counter++
			percent := (counter * 100) / 256
			fmt.Printf("Percent Complete: %v\n", percent)
		}
		close(input)
		close(output)
	}()

	// process output
	var results [][3]int
	for result := range output {
		fmt.Printf("Valid: %v,%v,%v\n", result[0], result[1], result[2])
		results = append(results, result)
	}
	fmt.Printf("calculationTime(%v)\n", time.Since(start))
	fmt.Println("Results:")
	for _, result := range results {
		fmt.Printf("%v,%v,%v\n", result[0], result[1], result[2])
	}
}

//////////////
// Simulator
type Simulator struct {
	sequences []RandomSequence
}

func NewSimulator() Simulator {
	length := int(math.Pow(3, float64(turns)))
	fmt.Println(length)
	sequences := make([]RandomSequence, length)
	for p := 0; p < length; p++ {
		sequences[p] = NewRandomSequence(turns, p)
	}
	return Simulator{
		sequences: sequences,
	}
}

func (this *Simulator) simulate(input, output chan [3]int) {
	for in := range input {
		if this.simulatePermutations(in[0], in[1], in[2]) {
			output <- in
		}
	}
}

func (this *Simulator) simulatePermutations(p1, p2, p3 int) bool {
	for _, sequence := range this.sequences {
		if !game(p1, p2, p3, sequence) {
			return false
		}
	}
	return true
}
