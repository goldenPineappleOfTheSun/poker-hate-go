package chances

import (
	"fmt"
	"math"
	"poker/cards"
	"poker/colors"
	"poker/combinations"
)

type Chances struct {
	tries         int
	win           float64
	royalflush    float64
	straightflush float64
	four          float64
	fullhouse     float64
	flush         float64
	straight      float64
	three         float64
	twopairs      float64
	pair          float64
	high          float64
}

func EmptyChances() Chances {
	return Chances{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
}

func Calculate(tries int, selfCards []cards.Card, othersCards [][]cards.Card, common []cards.Card, deck *cards.Deck) Chances {
	result := Chances{tries, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	selfCard1 := selfCards[0]
	selfCard2 := selfCards[1]

	for i := 0; i < tries; i++ {
		list := []cards.Card{}
		commonRevealed := []cards.Card{}

		for _, card := range common {
			commonRevealed = append(commonRevealed, lookupCard(card, deck))
		}

		list = append(list, selfCard1)
		list = append(list, selfCard2)

		for _, card := range commonRevealed {
			list = append(list, card)
		}
		hand := combinations.FindCombinations(list)

		checked := map[string]bool{"royal flush": false, "straight flush": false, "four": false, "full house": false, "flush": false, "straight": false, "three": false, "two pairs": false, "pair": false}

		for _, h := range hand {
			switch h.Name {
			case "royal flush":
				if !checked[h.Name] {
					checked[h.Name] = true
					result.royalflush += float64(1) / float64(tries)
				}
			case "straight flush":
				if !checked[h.Name] {
					checked[h.Name] = true
					result.straightflush += float64(1) / float64(tries)
				}
			case "four":
				if !checked[h.Name] {
					checked[h.Name] = true
					result.four += float64(1) / float64(tries)
				}
			case "full house":
				if !checked[h.Name] {
					checked[h.Name] = true
					result.fullhouse += float64(1) / float64(tries)
				}
			case "flush":
				if !checked[h.Name] {
					checked[h.Name] = true
					result.flush += float64(1) / float64(tries)
				}
			case "straight":
				if !checked[h.Name] {
					checked[h.Name] = true
					result.straight += float64(1) / float64(tries)
				}
			case "three":
				if !checked[h.Name] {
					checked[h.Name] = true
					result.three += float64(1) / float64(tries)
				}
			case "two pairs":
				if !checked[h.Name] {
					checked[h.Name] = true
					result.twopairs += float64(1) / float64(tries)
				}
			case "pair":
				if !checked[h.Name] {
					checked[h.Name] = true
					result.pair += float64(1) / float64(tries)
				}
			case "high":
				result.high += 0
			}
		}

		win := true

		for _, other := range othersCards {

			cards1 := other[0]
			cards2 := other[0]

			complist := []cards.Card{}
			complist = append(complist, lookupCard(cards1, deck))
			complist = append(complist, lookupCard(cards2, deck))
			for _, card := range commonRevealed {
				complist = append(complist, card)
			}
			competition := combinations.FindCombinations(complist)

			i := 0
			calculated := 0
			for i < len(hand) && i < len(competition) && calculated == 0 {
				calculated = hand[i].HigherThan(competition[i])
				i += 1
			}

			if calculated == 0 {

				var otherHighCard cards.Card
				var otherLowCard cards.Card
				if cards1.Rank() > cards2.Rank() {
					otherHighCard = cards1
					otherLowCard = cards2
				} else {
					otherHighCard = cards2
					otherLowCard = cards1
				}

				var playerHighCard cards.Card
				var playerLowCard cards.Card
				if selfCard1.Rank() > selfCard2.Rank() {
					playerHighCard = selfCard1
					playerLowCard = selfCard2
				} else {
					playerHighCard = selfCard2
					playerLowCard = selfCard1
				}

				if playerHighCard.Rank() > otherHighCard.Rank() {
					calculated = 1
				} else if playerHighCard.Rank() < otherHighCard.Rank() {
					calculated = -1
				} else {
					if playerLowCard.Rank() > otherLowCard.Rank() {
						calculated = 1
					} else if playerLowCard.Rank() < otherLowCard.Rank() {
						calculated = -1
					}
				}

			}

			if calculated < 0 {
				win = false
			}
		}

		if win {
			result.win += float64(1) / float64(tries)
		}
	}

	return result
}

func (c Chances) IsEmpty() bool {
	return c.tries == 0
}

func (c Chances) Win() float64 {
	return c.win
}

func (c Chances) String() string {
	/*cBlack := colors.CBlack()
	cReset := colors.CReset()
	background := colors.BYellow()
	bReset := colors.BReset()*/

	result := ""
	//result += cBlack + background
	result += "~~~~~~~ Senses ~~~~~~~"
	result += "\n"
	result += fitString(fmt.Sprintf("win: %.2f%%", c.win*100))
	result += fitString(fmt.Sprintf("flush: %.2f%%", c.flush*100))
	result += fitString(fmt.Sprintf("four: %.2f%%", c.four*100))
	result += fitString(fmt.Sprintf("full house: %.2f%%", c.fullhouse*100))
	result += fitString(fmt.Sprintf("straight: %.2f%%", c.straight*100))
	result += fitString(fmt.Sprintf("three: %.2f%%", c.three*100))
	result += fitString(fmt.Sprintf("two pairs: %.2f%%", c.twopairs*100))
	result += fitString(fmt.Sprintf("pair: %.2f%%", c.pair*100))
	//result += cBlack + background
	result += "~~~~~~~~~~~~~~~~~~~~~~"
	//result += cReset + bReset
	result += "\n"
	return result
}

func fitString(text string) string {
	cBlack := colors.CBlack()
	cReset := colors.CReset()
	bWhite := colors.BWhite()
	bReset := colors.BReset()
	length := len(text)
	padding := 23 - length
	left := int(math.Floor(float64(padding) / float64(2)))
	right := int(math.Ceil(float64(padding) / float64(2)))

	result := ""
	result += cBlack + bWhite
	for i := 0; i < left; i++ {
		result += " "
	}
	result += text
	for i := 0; i < right; i++ {
		result += " "
	}
	result += cReset + bReset
	result += "\n"
	return result

}

func lookupCard(card cards.Card, deck *cards.Deck) cards.Card {
	if card.IsRevealed() {
		return card
	} else {
		return deck.LookupRandom()
	}
}
