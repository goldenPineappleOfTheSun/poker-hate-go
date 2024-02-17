package combinations

import (
	"fmt"
	"sort"
	"strings"
	"poker/cards"
)

type Combination struct {
	Name string
	High int
	Cards []cards.Card
}

type ofKindCombination struct {
	highRank int
	pairs int
	pairRank int
	threes int
	threeRank int
	fours int
	fourRank int
}

type straight struct {
	rank int
	size int
	flush bool
}

type flush struct {
	rank int
	size int
	suit string
}

func checkName(name string) bool {
	names := []string {"royal flush", "straight flush", "four", "full house", "flush", "straight", "three", "two pairs", "pair", "high"}
    for _, a := range names {
        if a == name {
            return true
        }
    }
    return false
}

func CreateCombination(name string, high string, inputCards []cards.Card) Combination {
	if (!checkName(name)) {
		// golang has no exceptions?
		fmt.Println("ERROR: Invalid combination name \"" + name + "\"")
	}
	return Combination{name, rankStrToInt(high), inputCards}
}

func FindCombinations(list []cards.Card) []Combination {
	result := []Combination{}
	
	sort.Slice(list, func(i, j int) bool {
		return list[i].Rank() > list[j].Rank()
	})

	result = append(result, FindStraight(list)...)
	result = append(result, FindFlush(list)...)
	result = append(result, FindOfKindCombinations(list)...)
	result = append(result, FindHighCards(list)...)

	sort.Slice(result, func(i, j int) bool {
		return result[i].HigherThan(result[j]) > 0
	})

	return result
}

func FindOfKindCombinations(list []cards.Card) []Combination {
	result := []Combination{}
	counts := [15]int {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	
	for _, card := range list {
		counts[card.Rank()] += 1
	}

	for i, _ := range counts {
		if (counts[i] >= 2) {
			result = append(result, CreateCombination("pair", rankIntToStr(i), list))
			for j, _ := range counts {
				if (i != j && counts[j] >= 2) {
					if (i > j) {
						result = append(result, CreateCombination("two pairs", rankIntToStr(i), list))
					} else {
						result = append(result, CreateCombination("two pairs", rankIntToStr(j), list))
					}
				}
				if (i != j && counts[j] >= 3) {
					if (i > j) {
						result = append(result, CreateCombination("full house", rankIntToStr(i), list))
					} else {
						result = append(result, CreateCombination("full house", rankIntToStr(j), list))
					}
				}
			}
		}
		if (counts[i] >= 3) {
			result = append(result, CreateCombination("three", rankIntToStr(i), list))
		}
		if (counts[i] >= 4) {
			result = append(result, CreateCombination("four", rankIntToStr(i), list))
		}
	}

	return result
}

func FindStraight(list []cards.Card) []Combination {
	result := []Combination{}
	inrow := 0

	ace := find("?", 14, list)
	aced := !ace.IsEmpty()

	for i := 2; i <= 13; i++ {
		found := find("?", i, list)
		if (!found.IsEmpty()) {
			inrow += 1
			if ((inrow >= 5) || (inrow >= 4 && aced)) {
				result = append(result, CreateCombination("straight", rankIntToStr(i), list))
			}
		} else {
			inrow = 0
		}
	}

	return result
}

func FindFlush(list []cards.Card) []Combination {
	result := []Combination{}
	suits := map[string]int {"♥":0, "♦":0, "♣":0, "♠":0}
	
	for _, card := range list {
		suits[card.Suit()] += 1
		if (suits[card.Suit()] >= 5) {
			result = append(result, CreateCombination("flush", rankIntToStr(card.Rank()), list)) 
		}
	}

	return result
}

func FindHighCards(list []cards.Card) []Combination {
	result := []Combination{}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Rank() > list[j].Rank()
	})
	for _, card := range list {
		result = append(result, CreateCombination("high", rankIntToStr(card.Rank()), list))
	}
	return result
}

func (a *Combination) HigherThan(b Combination) int {
	names := map[string]int {"royal flush":9, "straight flush":8, "four":7, "full house":6, "flush":5, "straight":4, "three":3, "two pairs":2, "pair":1, "high":0}
	if (names[a.Name] > names[b.Name]) {
		return 1
	}	
	if (names[a.Name] < names[b.Name]) {
		return -1
	}
	if (a.High > b.High) {
		return 1
	}	
	if (a.High < b.High) {
		return -1
	}
	return 0
}

func (c Combination) String() string {
	return fmt.Sprintf("<%s:%s>", c.Name, rankIntToStr(c.High))
}

func (c Combination) Title() string {
	if (c.Name == "high") {
		return fmt.Sprintf("%s (%s)", "High card", rankIntToStr(c.High))
	} else {
		return fmt.Sprintf("%s (%s)", strings.Title(c.Name), rankIntToStr(c.High))
	}
}

func rankIntToStr(rank int) string {
	ranks := map[int]string {2:"2", 3:"3", 4:"4", 5:"5", 6:"6", 7:"7", 8:"8", 9:"9", 10:"10", 11:"J", 12:"Q", 13:"K", 14:"A"}
	return ranks[rank]
}

func rankStrToInt(rank string) int {
	ranks := map[string]int {"2":2, "3":3, "4":4, "5":5, "6":6, "7":7, "8":8, "9":9, "10":10, "J":11, "Q":12, "K":13, "A":14}
	return ranks[rank]
}

func find(suit string, rank int, list []cards.Card) cards.Card {
	for _, card := range list {
		if (card.Rank() == rank && (card.Suit() == suit || suit == "?")) {
			return card
		}
	}
	return cards.EmptyCard()
}