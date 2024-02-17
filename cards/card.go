package cards

import (
	"fmt"
	"poker/colors"
)

type Card struct {
	rank  int
	suit  string
	shown bool
}

func CreateCard(rank int, suit string) Card {
	return Card{rank, suit, false}
}

func EmptyCard() Card {
	return Card{0, "X", false}
}

func (c Card) IsRevealed() bool {
	return c.shown
}

func (c Card) Rank() int {
	return c.rank
}

func (c Card) Suit() string {
	return c.suit
}

func (c Card) IsEmpty() bool {
	return c.suit == "X"
}

func (c Card) String() string {
	ranks := [...]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	cRed := colors.CRed()
	cMagenta := colors.CMagenta()
	cYellow := colors.CYellow()
	cBlue := colors.CBlue()
	cBlack := colors.CBlack()
	cReset := colors.CReset()
	bWhite := colors.BWhite()
	bYellow := colors.BYellow()
	bReset := colors.BReset()

	suitColor := cBlack
	if c.suit == "♥" {
		suitColor = cRed
	}
	if c.suit == "♦" {
		suitColor = cMagenta
	}
	if c.suit == "♣" {
		suitColor = cBlue
	}
	if c.shown {
		return fmt.Sprintf(bWhite+suitColor+"%s"+"%s"+cBlack+cReset+bReset, ranks[c.rank], c.suit)
	} else {
		return fmt.Sprintf(bYellow+cYellow+"%s"+"%s"+cBlack+cReset+bReset, "◊", "◊")
	}
}

func (c *Card) Reveal() {
	c.shown = true
}

func (c *Card) Copy() Card {
	return CreateCard(c.rank, c.suit)
}

func (c *Card) Lookup() Card {
	card := c.Copy()
	card.Reveal()
	return card
}
