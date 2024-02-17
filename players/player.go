package players

import (
	"fmt"
	"math/rand"
	"poker/cards"
	"poker/chances"
	"poker/colors"
)

type Player struct {
	name      string
	color1    string
	color2    string
	painted   string
	card1     cards.Card
	card2     cards.Card
	ai        bool
	money     int
	p_normal  Personality // base
	p_preflop Personality // multiplier
	p_angry   Personality // multiplier
	/*p_winner Personality
	  p_looser Personality
	  state string // normal, angry, winner, looser, concentrated,*/
	//is_preflop bool
	is_angry     bool
	calculations chances.Chances
	possible_commentaries map[string][]string // fold, keepgoing, win, loose, raise
	commentary string
}

func CreatePlayer(name string, color1 string, color2 string) *Player {
	return &Player{name, color1, color2, "", cards.EmptyCard(), cards.EmptyCard(), false, 100,
		EmptyPersonality(), EmptyPersonality(), EmptyPersonality(), false, chances.EmptyChances(),
		make(map[string][]string), ""}
}

func CreateAiPlayer(name string, color1 string, color2 string) *Player {
	player := Player{
		name, color1, color2, "", cards.EmptyCard(), cards.EmptyCard(), true, 100,
		EmptyPersonality(), EmptyPersonality(), EmptyPersonality(), false, chances.EmptyChances(),
		make(map[string][]string), ""}
	return &player
}

/*
	func (p *Player) SetPersonality(value Personality) {
	    p.p_normal = value
	    p.p_preflop = CreatePersonality(value.impulsiveness * 0.5, value.bluffing * 0, value.counting * 0.1, value.intuition * 0.1, value.aggressiveness * 1.5)
	    p.p_badcards = CreatePersonality(value.impulsiveness * 0.7, value.bluffing * 1.5, value.counting * 1, value.intuition * 1, value.aggressiveness * 0.7)
	    p.p_angry = CreatePersonality(value.impulsiveness * 3, value.bluffing * 1, value.counting * 0.6, value.intuition * 0.6, value.aggressiveness * 3)
	    p.p_winner = CreatePersonality(value.impulsiveness * 3, value.bluffing * 1.5, value.counting * 0.7, value.intuition * 0.7, value.aggressiveness * 3)
	    p.p_looser = CreatePersonality(value.impulsiveness * 0.7, value.bluffing * 1.5, value.counting * 1.2, value.intuition * 1.2, value.aggressiveness * 0.7)
	}
*/

/*
	func (p *Player) MultiplyAngryPersonality(value Personality) {
	    p.p_angry.Multiply(value)
	}

	func (p *Player) MultiplyWinnerPersonality(value Personality) {
	    p.p_winner.Multiply(value)
	}

	func (p *Player) MultiplyLooserPersonality(value Personality) {
	    p.p_looser.Multiply(value)
	}
*/
func (p *Player) Card1() cards.Card {
	return p.card1
}

func (p *Player) Card2() cards.Card {
	return p.card2
}

func (p *Player) Cards() []cards.Card {
	return []cards.Card{p.Card1(), p.Card2()}
}

func (p *Player) Color1() string {
	return p.color1
}

func (p *Player) Color2() string {
	return p.color2
}

func (p *Player) Money() int {
	return p.money
}

func (p *Player) PaintIn(color string) {
	p.painted = color
}

func (p *Player) UnPaint() {
	p.painted = ""
}

func (p Player) String() string {
	cReset := colors.CReset()
	bReset := colors.BReset()
	if p.painted == "" {
		return fmt.Sprintf("(%s%s %s %s%s %+v %+v %d$)", p.color1, p.color2, p.name, cReset, bReset, p.card1, p.card2, p.money)
	} else {
		return fmt.Sprintf("(%s %s  %+v %+v %s%d$%s%s)", p.painted, p.name, p.card1, p.card2, p.painted, p.money, cReset, bReset)
	}
}

func (p Player) LookupString() string {
	if p.painted == "" {
		return fmt.Sprintf("(%s %+v %+v %d$)", p.name, p.LookupCard1(), p.LookupCard2(), p.money)
	} else {
		cReset := colors.CReset()
		return fmt.Sprintf("(%s%s %+v %+v %s%d$%s)", p.painted, p.name, p.LookupCard1(), p.LookupCard2(), p.painted, p.money, cReset)
	}
}

func (p Player) IsNPC() bool {
	return p.ai
}

func (p *Player) Reveal() {
	p.card1.Reveal()
	p.card2.Reveal()
}

func (p *Player) LookupCard1() cards.Card {
	return p.card1.Lookup()
}

func (p *Player) LookupCard2() cards.Card {
	return p.card2.Lookup()
}

func (p *Player) GetCurrentPersonality(is_preflop bool, is_angry bool) Personality {
	result := p.p_normal
	if is_preflop {
		result.Multiply(p.p_preflop)
	}
	if is_angry {
		result.Multiply(p.p_angry)
	}
	return result
}

func (p *Player) MakeCalculations(personality Personality, others []*Player, shared []cards.Card, deck *cards.Deck) {
	tries_counting := 5000 * personality.counting
	tries_intuition := 2000 * personality.intuition
	tries_impulsiveness := 1 / personality.impulsiveness
	tries_aggressiveness := 1 / personality.aggressiveness
	tries := int(tries_counting + tries_intuition - tries_impulsiveness - tries_aggressiveness)
	if tries < 1000 {
		tries = 1000
	}

	if personality.intuition > 0 {
		for _, other := range others {
			dice := rand.Float32()
			if dice > personality.intuition {
				card := other.Card1()
				(&card).Reveal()
			}
		}
	}

	competantsCards := [][]cards.Card{}
	for _, player := range others {
		competantsCards = append(competantsCards, player.Cards())
	}

	p.calculations = chances.Calculate(tries, p.Cards(), competantsCards, shared, deck)
}

/* returns true if fold */
func (p *Player) RequestFold(maxbet int, currentbet int, others []*Player, shared []cards.Card, deck *cards.Deck, is_preflop bool) bool {
	personality := p.GetCurrentPersonality(is_preflop, p.is_angry)
	if p.calculations.IsEmpty() {
		p.MakeCalculations(personality, others, shared, deck)
	}
	pot := float32(maxbet) / float32(currentbet)
	if pot > 1.5 {
		pot = 1.5
	}
	if (is_preflop) {
		pot = 1.1
	}
	
	risky := pot == 1.5 // high pot chance
	crazy := false
	if (!is_preflop && maxbet > currentbet) {
		crazy = p.money / (maxbet - currentbet) < 2 // more than half of money for a bet
	}
	wealth := p.money / maxbet
	poor := wealth < 5
	rich := wealth > 100
	make_comment := rand.Float32() < personality.impulsiveness + personality.aggressiveness * 0.5
	
	if (crazy && rand.Float32() > personality.impulsiveness * 3) {
		if (make_comment) {
			p.MakeCommentary("fold")
		}
		return true
	}

	if (crazy && rand.Float32() < personality.aggressiveness) {
		if (make_comment) {
			p.MakeCommentary("fold")
		}
		return true
	}
	
	if (poor && rand.Float32() > personality.impulsiveness * 2) {
		if (make_comment) {
			p.MakeCommentary("fold")
		}
		return true
	}

	if (risky && rand.Float32() > personality.impulsiveness) {
		if (make_comment) {
			p.MakeCommentary("fold")
		}
		return true
	}

	if (rand.Float32() < personality.impulsiveness * 0.1) {
		if (make_comment) {
			p.MakeCommentary("fold")
		}
		return true
	}

	if (rich && p.calculations.Win() > 0.5 * float64(personality.aggressiveness)) {
		return false
	}

	if (rich && p.calculations.Win() > (1 - float64(personality.bluffing))) {
		return false
	}

	if (pot > 1 && p.calculations.Win() < float64(pot - 1)) {
		if (make_comment) {
			p.MakeCommentary("fold")
		}
		return true
	}

	return false
}

/* returns raise (relative to currentbet) */
func (p *Player) RequestCheck(maxbet int, others []*Player, shared []cards.Card, deck *cards.Deck, is_preflop bool) int {
	personality := p.GetCurrentPersonality(is_preflop, p.is_angry)
	if p.calculations.IsEmpty() {
		p.MakeCalculations(personality, others, shared, deck)
	}

	win := float32(p.calculations.Win())
	make_comment := rand.Float32() < 0.3 + personality.impulsiveness + personality.aggressiveness
	dice := rand.Float32()

	if (dice < personality.impulsiveness * 0.05 + personality.aggressiveness * 0.05 + personality.bluffing * 0.2) {
		if (make_comment) {
			p.MakeCommentary("raise")
		}
		step := float32(int((maxbet / 4) / 4) * 4)
		raise := int((rand.Float32() * float32(maxbet)) / step * step)
		if (raise > p.money) {
			raise = p.money
		}
		return raise
	}

	if (dice > 0.8 && win > personality.impulsiveness * 0.05 + personality.aggressiveness * 0.05 + personality.bluffing * 0.2) {
		if (make_comment) {
			p.MakeCommentary("raise")
		}
		step := float32(int((maxbet / 4) / 4) * 4)
		raise := int((rand.Float32() * float32(maxbet)) / step * step)
		if (raise > p.money) {
			raise = p.money
		}
		return raise
	}

	if (make_comment) {
		p.MakeCommentary("keepgoing")
	}

	return 0
}

/* returns raise (relative to currentbet) */
func (p *Player) RequestCall(maxbet int, currentbet int, others []*Player, shared []cards.Card, deck *cards.Deck, is_preflop bool) int {
	personality := p.GetCurrentPersonality(is_preflop, p.is_angry)
	if p.calculations.IsEmpty() {
		p.MakeCalculations(personality, others, shared, deck)
	}
	raise := maxbet - currentbet + p.RequestCheck(maxbet, others, shared, deck, is_preflop)
	if (raise > p.money) {
		raise = p.money
	}
	return raise
}

func (p *Player) RemoveMoney(amount int) {
	p.money -= amount
}

func (p *Player) AddMoney(amount int) {
	p.money += amount
}

func (p *Player) Give(card1 cards.Card, card2 cards.Card) {
	p.card1 = card1
	p.card2 = card2
}

func (p *Player) RemoveCards() {
	p.card1 = cards.EmptyCard()
	p.card2 = cards.EmptyCard()
}

func (p *Player) SetCommentaries(possible_commentaries map[string][]string) {
	p.possible_commentaries = possible_commentaries
}

func (p *Player) SetPersonality(traits map[string]float32) {
	p.p_normal = Personality{traits["impulsiveness"], traits["bluffing"], traits["counting"], traits["intuition"], traits["aggressiveness"]}
	p.p_preflop = Personality{traits["impulsiveness"], traits["bluffing"], traits["counting"], traits["intuition"], traits["aggressiveness"]}
	p.p_angry = Personality{traits["impulsiveness"], traits["bluffing"], traits["counting"], traits["intuition"], traits["aggressiveness"]}
}

func (p *Player) MultiplyPreflopPersonality(value Personality) {

}

func (p *Player) MultiplyAngryPersonality(value Personality) {

}

func (p *Player) MakeCommentary(key string) {
	val, ok := p.possible_commentaries[key]
	if ok {
		p.commentary = val[rand.Intn(len(val))]
	} else {
		p.commentary = ""
	}
}

func (p *Player) AskCommentary() string {
	result := p.commentary
	p.commentary = ""
	return result
}

func (p *Player) AskToMakeCommentary(key string) {
	personality := p.GetCurrentPersonality(false, p.is_angry)

	if (key == "win" && rand.Float32() < personality.impulsiveness * 2 + personality.aggressiveness) {
		p.MakeCommentary(key)
	}
	if (key == "loose" && rand.Float32() < personality.impulsiveness + personality.aggressiveness * 0.5) {
		p.MakeCommentary(key)
	}
}
