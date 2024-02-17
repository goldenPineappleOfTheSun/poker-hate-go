package main

import (
	"fmt"
	"os"
    "os/exec"
	"runtime"

	//"time"
	"bufio"

	"github.com/mattn/go-colorable"

	//"poker/cards"
	"poker/colors"
	"poker/games"

	//"poker/chances"
	"poker/players"
	//"poker/combinations"
)

func main() {

	stdOut := bufio.NewWriter(colorable.NewColorableStdout())
	stdIn := bufio.NewReader(os.Stdin)

	//deck := cards.CreateDeck()
	player := players.CreatePlayer("Player", colors.CWhite(), colors.BBlack())

	//players.CreateAiPlayer("Jacob The Governour ♂", colors.CRed(), colors.BBlack())
	/* impulsive but with a very good intuition, which betrays him sometimes */
	//players.CreateAiPlayer("b1", colors.CYellow(), colors.BBlack())
	/* thougthful but nervous */
	//players.CreateAiPlayer("b2", colors.CGreen(), colors.BBlack())
	/* good player, calm, thougthful, good at counting, ok intuition */
	//players.CreateAiPlayer("b3", colors.CCyan(), colors.BBlack())

	/*black red
	green red
	yellow red
	blue red
	blue red
	cyan red
	white red
	black green
	red green
	yellow green
	blue green
	mag green
	white green
	black yellow
	red yellow
	green yellow
	blue yellow
	red yellow
	mag yellow
	white yellow
	black blue
	red blue
	yellow blue
	cyan blue
	white blue
	black mag
	red mag
	yellow mag
	cyan mag
	white mag
	black cyan
	red cyan
	blue cyan
	mag cyan
	white cyan
	black white
	red white
	green white
	blue white
	mag white
	cyan white*/

	/* a german. good at counting. risking */

	/* some gentleman. ok at counting and intuition. doesnt like to bluff and never risks */

	/* a cowboy. aggressive. likes to bluff and risk sometimes */

	/* an old cowboy, doesnt like to bluff, bot sometimes likes to risk big */

	/* a young girl. good intuition. everithing else is just random */

	/*
	impulsiveness  
	bluffing       
	counting       
	intuition      
	aggressiveness 
	*/

	/* a local psycho but not scary */
	//jason := players.CreateAiPlayer("Fast Anthony ♂", colors.CRed(), colors.BGreen())
	/* some wierdo */
	//mike := players.CreateAiPlayer("Wierd Mike ♂", colors.CRed(), colors.BYellow())
	/* a rich and dangerous man. a patriot */
	//nigel := players.CreateAiPlayer("Rich Nigel ♂", colors.CRed(), colors.BBlue())
	/* impulsive girl with a bad personality*/
	//betty := players.CreateAiPlayer("Wild Betty ♀", colors.CRed(), colors.BCyan())
	/* impulsive girl bad at a game*/
	emma := players.CreateAiPlayer("Деревенская Эмма ♀", colors.CYellow(), colors.BRed())
	emma.SetPersonality(map[string]float32{
		"impulsiveness": 0.5,
		"bluffing": 0.1,
		"counting": 0.05,
		"intuition": 0.1,
		"aggressiveness": 0.2,
	})
	emma.SetCommentaries(map[string][]string{
		"fold": []string{"О нет! Сдаюсь", "Не могу!  Сдаюсь", "Ой не знаю..  Сдаюсь", "Дальше без меня, я сдаюсь"},
		"keepgoing": []string{"Ха-ха!", "Думаю, я ещё могу победить"},
		"raise": []string{"Повышаю, держитесь, ребятки!", "Повышаю! Не ждали такого?"},
	})
	/* impulsive girl but with a good intuition */
	//zoe := players.CreateAiPlayer("Small Zoe ♀", colors.CWhite(), colors.BCyan())
	/* a very good player (blue + black + white. stronger then b2, but weaker than bb), a woman */
	jessy := players.CreateAiPlayer("Умная Джесси ♀", colors.BBlue(), colors.CWhite())
	jessy.SetPersonality(map[string]float32{
		"impulsiveness": 0.1,
		"bluffing": 0.3,
		"counting": 0.6,
		"intuition": 0.2,
		"aggressiveness": 0.0,
	})
	jessy.SetCommentaries(map[string][]string{
		"fold": []string{"Ну чтож, пас", "Хах, ну тогда пас", "Пас, отыграюсь в следующий раз"},
		"keepgoing":   []string{"Я в деле"},
		"raise": []string{"В данной ситуации... я бы повысила", "Повышаю. Но блефую ли я?"},
	})
	/* a local psycho and the scary one. do impulsive stuff but also have incredible intuition */
	//jason := players.CreateAiPlayer("Jason Abel ♂", colors.CBlack(), colors.BRed())
	/* antagonist */
	boss := players.CreateAiPlayer("Губернатор Якоб ♂", colors.CBlack(), colors.BWhite())
	boss.SetPersonality(map[string]float32{
		"impulsiveness": 0.1,
		"bluffing": 0.3,
		"counting": 1,
		"intuition": 0.9,
		"aggressiveness": 0.2,
	})
	boss.SetCommentaries(map[string][]string{
		"fold": []string{"Неа. Пас", "В таком случае я пас", "В следующий раз сьем тебя с потрохами. А пока пас"},
		"keepgoing":   []string{"*Кивает*", "Понятно..", "Да будет так"},
		"raise": []string{"Чтож, повышаю", "Ну давайте посмотрим.. Повышаю", "Да будет так, повышаю"},
	})
	/* boss. aggressive and impulsive */
	//boss2 := players.CreateAiPlayer("Billy The Bob ♂", colors.CRed(), colors.BWhite())
	/* boss. thoughtfull */
	//boss3 := players.CreateAiPlayer("Norman The Policeman ♂", colors.CBlue(), colors.BWhite())

	e1 := emma
	e2 := jessy
	e3 := boss

	game := games.CreateGame([]*players.Player{e1, e2, e3, player}, 4)
	deck := game.Deck()

	cBlack := "\033[30m"   //мафия
	cRed := "\033[31m"     //хорошая интуиция, эмоции скрывает чуть хуже среднего
	cGreen := "\033[32m"   //спокойный, не повышает сильно, редко колирует без хороших карт, эмоций почти никогда не выражает
	cYellow := "\033[33m"  //плохо играет, попеременно то не уверен в своих силах то слишком самоуверен, часто крадет префлопы, в половине случаев эмоции читаются очень легко а в половине очень хороший покерфейс
	cBlue := "\033[34m"    //неплохая интуиция, любит блефовать делает это чтобы повысить шансы когда ему кажется что они примерно средние, наслаждается победами
	cMagenta := "\033[35m" //агрессивный, спешит, сильно повышает, часто блефует, часто крадет префлопы, в половине случаев эмоции читаются очень легко а в половине очень трудно
	cCyan := "\033[36m"    //редко блефует, часто пугается, неплохо считает, эмоции скрывает чуть хуже среднего
	cWhite := "\033[37m"   //хаосит, многие решения принимает броском монеты, часто крадет префлопы, обычно веселый
	cReset := "\033[39m"   //хорошо считает, почти всегда хороший покерфейс, но может выдать себя когда шансы малы

	bBlack := "\033[40m"
	bRed := "\033[41m"
	bGreen := "\033[42m"
	bYellow := "\033[43m"
	bBlue := "\033[44m"
	bMagenta := "\033[45m"
	bCyan := "\033[46m"
	bWhite := "\033[47m"
	bReset := "\033[49m"

	//1. win sense + combos sense + lookup + straights lookup + nobluff
	//1. combos sense + lookup + straights lookup + nobluff
	//1 b3 combos sense + nobluff
	//2. combos sense + lookup + straights lookup + nobluff
	//2. win sense + combos sense + lookup + straights lookup
	//3. combos sense + lookup + straights lookup + nobluff
	//3. combos sense + lookup + straights lookup
	//3 b2. combos sense
	//4. lookup + straights lookup
	//4. straights lookup
	//4. none
	//4 b3. none. 1 to 1
	//boss. none

	//boss can visit any town if player is famous enough.
	//he will visit town two if he never visited town 1

	//b1 can visit any town
	//b2 can visit towns 2 - 3. is the boss of town 3 he ashamed after beaten in 3
	//b3 can visit towns 1 - 4. when he is not a boss, he leaves a table to never be 1 to 1. beaten as a boss in town 1 and asks for revanche in town 4 where he is also a boss

	//bosses always have some buff to their senses when they just visiting a game (when its not a boss battle)

	//the famousity treshold doubles with every visit for every boss

	for _, a := range []string{cBlack, cRed, cGreen, cYellow, cBlue, cMagenta, cCyan, cWhite, cReset} {
		for _, b := range []string{bBlack, bRed, bGreen, bYellow, bBlue, bMagenta, bCyan, bWhite, bReset} {
			fmt.Fprint(stdOut, a+b+"test")
		}
	}
		
	stdOut.Flush()
	clearTerminal()

	for true {

		fmt.Fprint(stdOut, "\n")
		fmt.Fprint(stdOut, game)
		fmt.Fprint(stdOut, "\n")
		stdOut.Flush()

		if game.NotStarted() || game.Done() {
			player.RemoveCards()
			e1.RemoveCards()
			e2.RemoveCards()
			e3.RemoveCards()
			
			game.NextGame()
			deck = game.Deck()

			player.Give(deck.Draw(), deck.Draw())
			e1.Give(deck.Draw(), deck.Draw())
			e2.Give(deck.Draw(), deck.Draw())
			e3.Give(deck.Draw(), deck.Draw())

			game.StartGame()
		}

		isPlayersMove := (game.IsPlayableCharacter() && !game.Done())

		if !isPlayersMove {
			game.PlayRound()
			input, _ := stdIn.ReadString('\n')
			fmt.Fprint(stdOut, input)
		}

		output := game.ConsoleOutput()
		if output != "" {
			fmt.Fprint(stdOut, output)
			stdOut.Flush()
			input, _ := stdIn.ReadString('\n')
			fmt.Fprint(stdOut, input)
		}

		if isPlayersMove {
			fmt.Fprint(stdOut, "Ваш ход:\n")
			fmt.Fprint(stdOut, game.RenderCommands())
			fmt.Fprint(stdOut, "\n")
			stdOut.Flush()
			input := ""
			needAnotherLine := true
			for needAnotherLine {
				input, _ = stdIn.ReadString('\n')
				needAnotherLine = game.Command(input)
			}
		}
		
		stdOut.Flush()
		clearTerminal()

	}

	stdOut.Flush()
}

func clearTerminal() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}