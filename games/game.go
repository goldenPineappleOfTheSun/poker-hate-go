package games

import (
	"fmt"
	"math/rand"
	"poker/cards"
	"poker/chances"
	"poker/colors"
	"poker/combinations"
	"poker/players"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Game struct {
	name           string
	players        []*players.Player
	bb             int
	sb             int
	index          int
	progress       int
	playerprogress int
	stage          string
	deck           *cards.Deck
	shared         []cards.Card
	bets           []int
	folds          []bool
	smallBet       int
	consoleOutput  string
	winners        []*players.Player
	debug          bool
}

func CreateGame(people []*players.Player, smallBet int) *Game {
	rand.Seed(time.Now().UnixNano())
	bb := rand.Intn(len(people))
	sb := (bb + len(people) - 1) % len(people)
	index := sb
	deck := cards.CreateDeck()
	result := Game{
		"TheGame",
		people,
		bb, sb, index,
		0, 0, "init",
		deck, []cards.Card{},
		[]int{}, []bool{}, smallBet,
		"", []*players.Player{},
		true}
	result.NextGame()
	return &result
}

func (g *Game) ConsoleOutput() string {
	result := g.consoleOutput
	g.consoleOutput = ""
	return result
}

func (g *Game) Deck() *cards.Deck {
	return g.deck
}

func (g *Game) Done() bool {
	return g.stage == "done"
}

func (g *Game) IsStarted() bool {
	return g.stage != "init"
}

func (g *Game) NotStarted() bool {
	return g.stage == "init"
}

func (g *Game) PlayRound() {
	g.playRound()
}

func (g *Game) NextGame() {
	g.bb = (g.bb + 1) % len(g.players)
	g.sb = (g.sb + 1) % len(g.players)
	g.index = g.sb
	g.deck = cards.CreateDeck()
	g.shared = []cards.Card{g.deck.Draw(), g.deck.Draw(), g.deck.Draw(), g.deck.Draw(), g.deck.Draw()}
	g.bets = make([]int, len(g.players))
	g.folds = make([]bool, len(g.players))
	g.winners = []*players.Player{}
}

func (g *Game) StartGame() {
	g.stage = "preflop"

	bb := g.players[g.bb]
	sb := g.players[g.sb]

	if bb.Money() >= g.smallBet*2 {
		bb.RemoveMoney(g.smallBet * 2)
		g.makeBet(g.bb, g.smallBet*2)
	} else {
		bb.RemoveMoney(bb.Money())
		g.makeBet(g.bb, bb.Money())
	}
	if sb.Money() >= g.smallBet {
		sb.RemoveMoney(g.smallBet)
		g.makeBet(g.sb, g.smallBet)
	} else {
		sb.RemoveMoney(sb.Money())
		g.makeBet(g.sb, sb.Money())
	}

	g.next()
}

func (g *Game) IsPlayableCharacter() bool {
	return !g.players[g.index].IsNPC()
}

func (g *Game) RenderCommands() string {
	cBlack := colors.CBlack()
	cReset := colors.CReset()
	bWhite := colors.BWhite()
	bReset := colors.BReset()
	result := ""

	playerbet := g.bets[g.index]
	maxbet := 0
	for _, bet := range g.bets {
		if bet > maxbet {
			maxbet = bet
		}
	}

	if playerbet == maxbet {
		result += cBlack + bWhite
		result += "[C]"
		result += cReset + bReset
		result += " Check "
	}

	if playerbet < maxbet {
		result += cBlack + bWhite
		result += "[C]"
		result += cReset + bReset
		result += " Call "
	}

	result += cBlack + bWhite
	result += "[R]"
	result += cReset + bReset
	result += " Raise "

	result += cBlack + bWhite
	result += "[F]"
	result += cReset + bReset
	result += " Fold "

	result += cBlack + bWhite
	result += "[S]"
	result += cReset + bReset
	result += " Sense "

	result += cBlack + bWhite
	result += "[L]"
	result += cReset + bReset
	result += " Lookup "

	return result
}

func (g *Game) IsValidCommand(input string) bool {
	char := input[0]
	/* _____________ c   C   r    R   f    F   s    S   l    L   \n */
	valids := []byte{99, 67, 114, 82, 102, 70, 115, 83, 108, 76, 13}
	for _, key := range valids {
		if key == char {
			return true
		}
	}
	return false
}

/* returns true if another input needed*/
func (g *Game) Command(input string) bool {

	char := input[0]
	player := g.players[g.index]
	playerbet := g.bets[g.index]
	maxbet := 0
	for _, bet := range g.bets {
		if bet > maxbet {
			maxbet = bet
		}
	}

	if g.playerprogress == 1 {
		raise, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			fmt.Println("cant parse number")
			return true
		}
		if raise < maxbet {
			fmt.Println("Input a new max bet. The value must be larger than current max bet")
			return true
		}
		g.playerprogress = 0
		diff := raise - g.bets[g.index]
		if diff > player.Money() {
			diff = player.Money()
		}
		player.RemoveMoney(diff)
		g.makeBet(g.index, diff)

		g.next()
		return false
	}

	/* check or call */
	if char == 99 || char == 67 {
		if playerbet >= maxbet {
			g.next()
			return false
		} else {
			diff := maxbet - playerbet
			player.RemoveMoney(diff)
			g.makeBet(g.index, diff)
			g.next()
			return false
		}
	}

	/* raise */
	if char == 114 || char == 82 {
		g.playerprogress = 1
		fmt.Println("Input a new max bet:")
		return true
	}

	/* fold */
	if char == 102 || char == 70 {
		g.foldPlayer(g.index)
		g.next()
		return false
	}

	/* sense */
	if char == 115 || char == 83 {
		competantsCards := [][]cards.Card{}
		for _, player := range g.players {
			if player.IsNPC() {
				competantsCards = append(competantsCards, player.Cards())
			}
		}

		chance := chances.Calculate(1e4, g.players[g.index].Cards(), competantsCards, g.shared, g.deck)
		g.putOutputLn(fmt.Sprintf("%v", chance))
		return false
	}

	return true

}

func (g *Game) playRound() {
	maxbet := 0
	for _, bet := range g.bets {
		if bet > maxbet {
			maxbet = bet
		}
	}

	/*equilibrium := true
	for _, bet := range g.bets {
		if bet < maxbet {
			equilibrium = false
		}
	}*/

	ask(g, g.players[g.index])

	activePlayers := 0
	for i, _ := range g.players {
		if (!g.folds[i]) {
			activePlayers += 1
		}
	}

	if (activePlayers <= 1) {
		g.endRound()
	} else {
		if g.progress >= len(g.players) {
			g.nextStage()
		}
	}

	g.next()
}

func (g *Game) removePlayer(i int) {
	g.players = append(g.players[:i], g.players[i+1:]...)
	g.bets = append(g.bets[:i], g.bets[i+1:]...)
	g.index = (g.index + 1) % len(g.players)
	if g.bb >= i {
		g.bb -= 1
	}
	if g.sb >= i {
		g.sb -= 1
	}
	if g.sb == g.bb {
		g.sb = (g.sb + len(g.players) - 1) % len(g.players)
	}
}

func (g *Game) foldPlayer(i int) {
	g.folds[i] = true
}

func (g *Game) next() {
	stop := false
	for !stop {
		g.index = (g.index + 1) % len(g.players)
		g.progress += 1
		stop = !g.folds[g.index]
	}
}

func (g *Game) nextStage() {
	if g.stage == "preflop" {
		g.stage = "flop"
		g.progress = 0
		g.shared[0].Reveal()
		g.shared[1].Reveal()
		g.shared[2].Reveal()
		return
	}
	if g.stage == "flop" {
		g.stage = "turn"
		g.progress = 0
		g.shared[3].Reveal()
		return
	}
	if g.stage == "turn" {
		g.stage = "river"
		g.progress = 0
		g.shared[4].Reveal()
		g.endRound()
		return
	}
	/*if g.stage == "win" {
		g.stage = "done"
		g.progress = 0
		return
	}*/
}

func (g *Game) endRound() {
	/*g.stage = "win"*/
	g.progress = 0
	winners := g.findWinners()
	g.winners = winners
	for _, player := range g.players {
		player.Reveal()
	}
	g.shared[0].Reveal()
	g.shared[1].Reveal()
	g.shared[2].Reveal()
	g.shared[3].Reveal()
	g.shared[4].Reveal()
	g.win()
	return
}

func (g *Game) win() {
	bank := 0
	for _, bet := range g.bets {
		bank += bet
	}

	if (len(g.winners) > 0) {
		bank = bank / len(g.winners)
	}

	for _, winner := range g.winners {
		winner.AddMoney(bank)
	}

	bankrupt := -1
	for i, player := range g.players {
		if player.Money() <= 0 {
			bankrupt = i
		}
	}
	for bankrupt != -1 {
		g.removePlayer(bankrupt)
		bankrupt = -1
		for i, player := range g.players {
			if player.Money() <= 0 {
				bankrupt = i
			}
		}
	}

}

func ask(g *Game, p *players.Player) {
	maxbet := 0
	playerbet := g.bets[g.index]
	shared := []cards.Card{}

	for _, bet := range g.bets {
		if bet > maxbet {
			maxbet = bet
		}
	}

	for i := 0; i < 5; i++ {
		if len(g.shared) > i {
			shared = append(shared, g.shared[i])
		} else {
			shared = append(shared, cards.EmptyCard())
		}
	}

	if p.IsNPC() {

		isPreflop := g.stage == "preflop"

		others := []*players.Player{}
		for _, player := range g.players {
			if player != p {
				others = append(others, player)
			}
		}

		added := 0
		if playerbet >= maxbet {
			added = p.RequestCheck(maxbet, others, shared, g.deck, isPreflop)
		} else {
			if p.RequestFold(maxbet, g.bets[g.index], others, shared, g.deck, isPreflop) {
				g.foldPlayer(g.index)
				return
			}
			added = p.RequestCall(maxbet, g.bets[g.index], others, shared, g.deck, isPreflop)
		}
		g.players[g.index].RemoveMoney(added)
		g.makeBet(g.index, added)

	}
}

func (g *Game) makeBet(index int, money int) {
	g.bets[index] += money
}

func (g *Game) findWinners() []*players.Player {
	result := []*players.Player{}
	winnerCombo := []combinations.Combination{}

	for i, player := range g.players {

		if (g.folds[i]) {
			continue;
		} 

		usedcards := []cards.Card{}
		usedcards = append(usedcards, g.shared[0], g.shared[1], g.shared[2], g.shared[3], g.shared[4])
		usedcards = append(usedcards, player.Card1(), player.Card2())
		combo := combinations.FindCombinations(usedcards)

		if len(result) == 0 {
			result = append(result, player)
			winnerCombo = combo
			continue
		}

		i := 0
		calculated := 0
		for i < len(winnerCombo) && i < len(combo) && calculated == 0 {
			calculated = winnerCombo[i].HigherThan(combo[i])
			i += 1
		}

		if calculated == 0 {

			var resultHighCard cards.Card
			var resultLowCard cards.Card
			if result[0].Card1().Rank() > result[0].Card2().Rank() {
				resultHighCard = result[0].Card1()
				resultLowCard = result[0].Card2()
			} else {
				resultHighCard = result[0].Card2()
				resultLowCard = result[0].Card1()
			}

			var playerHighCard cards.Card
			var playerLowCard cards.Card
			if player.Card1().Rank() > player.Card2().Rank() {
				playerHighCard = player.Card1()
				playerLowCard = player.Card2()
			} else {
				playerHighCard = player.Card2()
				playerLowCard = player.Card1()
			}

			if resultHighCard.Rank() > playerHighCard.Rank() {
				calculated = 1
			} else if resultHighCard.Rank() < playerHighCard.Rank() {
				calculated = -1
			} else {
				if resultLowCard.Rank() > playerLowCard.Rank() {
					calculated = 1
				} else if resultLowCard.Rank() < playerLowCard.Rank() {
					calculated = -1
				}
			}

		}

		if calculated == 0 {
			result = append(result, player)
		} else if calculated < 0 {
			result = []*players.Player{player}
			winnerCombo = combo
		} else if calculated > 0 {
			/* nothing changed */
		}

	}

	return result
}

func (g *Game) putOutputLn(txt string) {
	g.consoleOutput += txt
	g.consoleOutput += "\n"
}

func (g *Game) String() string {
	cGreen := colors.CGreen()
	cRed := colors.CRed()
	cMagenta := colors.CMagenta()
	cReset := colors.CReset()

	result := ""
	buffer := ""

	if (g.debug) {
		result += fmt.Sprintf("~~~~~ %s ~~~~~\n", "DEBUG")
		result += fmt.Sprintf("stage:    %s\n", g.stage)
		result += fmt.Sprintf("progress: %d\n", g.progress)
		result += fmt.Sprintf("\n")
	}

	result += fmt.Sprintf("~~~~~ %s ~~~~~\n", g.name)

	maxplayertext := 0
	for _, player := range g.players {
		length := getPlayerStringSize(player)
		if length > maxplayertext {
			maxplayertext = length
		}
	}
	maxplayertext += 1

	maxbet := 0
	for i, _ := range g.players {
		if g.bets[i] > maxbet {
			maxbet = g.bets[i]
		}
	}

	for i, player := range g.players {
		if player.IsNPC() {
			
			if (g.folds[i]) {
				g.players[i].PaintIn(colors.CMagenta())
			}

			for _, winner := range g.winners {
				if (winner == player) {
					g.players[i].PaintIn(colors.CYellow())
				}
			}

			playertext := fmt.Sprintf(" %v", player)
			result += playertext
			for l := 0; l < maxplayertext-len(playertext); l++ {
				result += " "
			}

			result += cReset

			if g.folds[i] {
				result += cMagenta
			} else {
				if g.bets[i] >= maxbet {
					result += cGreen
				} else {
					result += cRed
				}
			}

			result += fmt.Sprintf(" %d$", g.bets[i])
			result += cReset

			if g.folds[i] {
				result += cMagenta
			}
			if g.bb == i {
				result += " (BB)"
			}
			if g.sb == i {
				result += " [SB]"
			}
			if g.index == i && g.IsStarted() {
				result += " ←"
			}
			result += cReset

			if g.stage == "river" {
				cards := []cards.Card{}
				cards = append(cards, g.shared[0], g.shared[1], g.shared[2], g.shared[3], g.shared[4])
				cards = append(cards, player.Card1(), player.Card2())
				result += fmt.Sprintf("   -   %s", combinations.FindCombinations(cards)[0].Title())
				for _, winner := range g.winners {
					if player == winner {
						player.AskToMakeCommentary("win")
						result += colors.CYellow() + "   WINNER!!!" + colors.CReset()
					} else {
						player.AskToMakeCommentary("loose")
					}
				}
			}

			commentary := player.AskCommentary()
			if (commentary != "") {
				result += fmt.Sprintf(" - says: \"%s\"", commentary)
			}

			result += "\n"

			player.UnPaint()
		}
	}

	buffer = ""
	for i := 0; i < 5; i++ {
		if len(g.shared) > i {
			buffer += fmt.Sprintf(" %v", g.shared[i])
		} else {
			buffer += fmt.Sprintf(" []")
		}
	}

	result += fmt.Sprintf("*")
	for i := 0; i < 16; i++ {
		result += fmt.Sprintf("-")
	}
	result += fmt.Sprintf("*\n")
	result += "|" + buffer + " |\n"
	result += fmt.Sprintf("*")
	for i := 0; i < 16; i++ {
		result += fmt.Sprintf("-")
	}
	result += fmt.Sprintf("*\n")

	for i, player := range g.players {
		if !player.IsNPC() {
			
			if (g.folds[i]) {
				g.players[i].PaintIn(colors.CMagenta())
			}

			for _, winner := range g.winners {
				if (winner == player) {
					g.players[i].PaintIn(colors.CYellow())
				}
			}

			result += fmt.Sprintf(" %s", player.LookupString())

			if g.folds[i] {
				result += cMagenta
			} else {
				if g.bets[i] >= maxbet {
					result += cGreen
				} else {
					result += cRed
				}
			}
			result += fmt.Sprintf(" %d$", g.bets[i])
			result += cReset

			if g.folds[i] {
				result += cMagenta
			}
			if g.bb == i {
				result += " (BB)"
			}
			if g.sb == i {
				result += " [SB]"
			}
			if g.index == i && g.IsStarted() {
				result += " ←"
			}
			result += cReset

			if g.stage == "river" {
				cards := []cards.Card{}
				cards = append(cards, g.shared[0], g.shared[1], g.shared[2], g.shared[3], g.shared[4])
				cards = append(cards, player.Card1(), player.Card2())
				result += fmt.Sprintf("   -   %s", combinations.FindCombinations(cards)[0].Title())
				for _, winner := range g.winners {
					if player == winner {
						result += colors.CYellow() + "   WINNER!!!" + colors.CReset()
					}
				}
			}

			result += "\n"

			player.UnPaint()
		}
	}

	result += fmt.Sprintf("~~~~~~")
	for i := 0; i < len(g.name); i++ {
		result += fmt.Sprintf("~")
	}
	result += fmt.Sprintf("~~~~~~\n")

	return result
}

func getPlayerStringSize(player *players.Player) int {
	colorizedString := fmt.Sprintf(" %v", player)
	colorCodeRegex := regexp.MustCompile("\x1b\\[[0-9;]*[mK]")
	plainString := colorCodeRegex.ReplaceAllString(colorizedString, "")
	return len(plainString)
}

func checkStage(name string) bool {
	names := []string{"preflop", "flop", "turn", "river"}
	for _, a := range names {
		if a == name {
			return true
		}
	}
	return false
}
