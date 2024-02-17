package cards

import (
	"time"
	"fmt"
    "math/rand"
)

type Deck struct {
    cards []Card
}

func shuffle(slice []Card) {
	n := len(slice)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(len(slice))
		slice[i], slice[j] = slice[j], slice[i]
	}

	
}

func CreateDeck() *Deck {
	result := Deck{ []Card{} }
	rand.Seed(time.Now().UnixNano())

	for _, s := range []string{"♥", "♣", "♦", "♠"} {
		for i := 2; i <= 14; i++ {
			result.cards = append(result.cards, CreateCard(i, s))
		}
	}

	shuffle(result.cards)

	return &result
}

func (d *Deck) Draw() Card {
	result := d.cards[0]
	d.cards = d.cards[1:]
	return result
}

func (d *Deck) DrawReveal() Card {
	result := d.cards[0]
	result.Reveal()
	d.cards = d.cards[1:]
	return result
}

func (d Deck) LookupRandom() Card {
	i := rand.Intn(len(d.cards))
	result := CreateCard(d.cards[i].Rank(), d.cards[i].Suit())
	result.Reveal()
	return result
}

func (d Deck) String() string {
	result := fmt.Sprintf("%d", len(d.cards))
	for _, card := range d.cards {
		card.Reveal()
		result += fmt.Sprintf("%+v", card)
	}
	return result
}

/*
func (d *Deck) So() []card {
	return d.cards
}*/