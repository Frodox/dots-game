/*
 * Here description of my Dots game
 *
 */
package main

import (
	"fmt"
	"time"
	"os"
	"os/exec"
	"bufio"
	"strconv"
	"strings"
	"math/rand"
)

// ---------------------------------- CONSTS -------------------------------- //
const gameBoardSize			int = 10

const fieldEmptyCellChar	string = "."
const fieldUserCellChar		string = "*"
const fieldPCCellChar		string = "+"

const fieldEmptyCellId		int = 0
const fieldUserCellId 		int = 1
const fieldPCCellId			int = 2

const CLR_0 = "\x1b[30;1m"
const CLR_R = "\x1b[31;1m"		// red
const CLR_G = "\x1b[32;1m"
const CLR_Y = "\x1b[33;1m"
const CLR_B = "\x1b[34;1m"		// blue
const CLR_M = "\x1b[35;1m"
const CLR_C = "\x1b[36;1m"
const CLR_W = "\x1b[37;1m"
const CLR_N = "\x1b[0m"			// reset color

const chars  string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+"


// ------------------------------ STRUCTS ----------------------------------- //
type GameBoardNode struct {
	value			int;		// 0, 1, 2 (Empty, User, PC)
	belongsToPlayer int;		// @id of player, who owns or captured this cell
}

type Player struct {
	stepId	int;		// 0, 1, 2 (Empty, User, PC)
	score 	int;		//  score of this player. 0 by default
}


// ------------------------------ FUNCTIONS --------------------------------- //
func d (debugMsg string) {
	fmt.Printf("D: ")
	fmt.Println(debugMsg)
	time.Sleep(400 * time.Millisecond)
}

/* ------------------------------------------------------------------------- */

func initGameBoard(size int) (gameBoard [][]GameBoardNode) {

	// Allocate the top-level slice, the same as before.
	gameBoard = make([][]GameBoardNode, size) // One row per unit of y.

	// Allocate one large slice to hold all the pixels.
	pixels := make([]GameBoardNode, size*size)

	// Loop over the rows, slicing each row from the front of the remaining pixels slice.
	for i := range gameBoard {
		gameBoard[i], pixels = pixels[:size], pixels[size:]
	}

	//gameBoard[0][1].value = 1
	//gameBoard[1][0].value = 1
	//gameBoard[2][1].value = 1
	//gameBoard[1][2].value = 1

	return
}

/* ------------------------------------------------------------------------- */

func clear_screen_linux() {
        cmd := exec.Command("clear") //Linux example, its tested
        cmd.Stdout = os.Stdout
        cmd.Run()
    }

/* ------------------------------------------------------------------------- */

/*
 * name: drawGameBoard
 * @param
 * @return
 *
 * 0 mean empty cell -- draw fieldEmptyCellChar
 * 1 mean user  cell -- draw fieldUserCellChar (blue)
 * 2 mean PC    cell -- draw fieldPCCellChar (red)
 */
func drawGameBoard(gameBoard [][]GameBoardNode, userPlayer *Player, pcPlayer *Player) {

	var length int = len(gameBoard)

	fmt.Printf("    ")
	for i := 0; i < length; i++ {
		fmt.Printf("%c ", chars[i])
    }
	fmt.Println()
	for i := range gameBoard {
		fmt.Printf("%2d  ", i)

		for j:= range gameBoard[i] {

			if fieldEmptyCellId == gameBoard[i][j].value {
				fmt.Printf("%s ", fieldEmptyCellChar)
			} else if fieldUserCellId == gameBoard[i][j].value {
				fmt.Printf("%s%s%s ", CLR_B, fieldUserCellChar, CLR_N)
			} else if fieldPCCellId == gameBoard[i][j].value {
				fmt.Printf("%s%s%s ", CLR_R, fieldPCCellChar, CLR_N)
			}
		}
		fmt.Println("")
	}

	fmt.Println("Game score")
	fmt.Printf("User: %s%d%s\t\tPC: %s%d%s\n",
			CLR_B, userPlayer.score, CLR_N,
			CLR_R, pcPlayer.score,   CLR_N)

	// TODO: если ячейка захвачена другим игроком -- сменить значок на '#'
}

/* ------------------------------------------------------------------------- */

func doUserStep(gameBoard [][]GameBoardNode) {

	var firstIndex  int = 0
	var secondIndex int = 0

	var fineInput int = 0;
	for fineInput != 1 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Your turn: [int char] > ")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n ")

		var tmp []string = strings.Split(text, " ");
		if len(tmp) != 2 {
			fmt.Println("Please, input 2 coords!");
			continue
		}

		// TODO: check for error
		tmpInt64, _ := strconv.ParseInt(tmp[0], 10, 0);
		firstIndex = int(tmpInt64)

		secondIndex = strings.Index(chars, tmp[1])

		if firstIndex  >= gameBoardSize || firstIndex < 0 ||
		   secondIndex >= gameBoardSize || secondIndex < 0 {
			fmt.Printf("Both index should be in game field range!\n");
			continue
		}

		res := doGameStep(gameBoard, firstIndex, secondIndex, fieldUserCellId);
		if 0 != res {
			fmt.Printf("Look's like cell is not empty or error occured\n");
			continue
		}
		fineInput = 1;
	}

	//fmt.Printf("First id %d , second id: %d\n", firstIndex, secondIndex);
}

/* ------------------------------------------------------------------------- */

/*
 * return 0 if fine
 * 	  1 if cell not empty already or error occured
 */
func doGameStep(gameBoard [][]GameBoardNode, x int, y int, symbol int) (result int) {

	if fieldEmptyCellId == gameBoard[x][y].value {
		gameBoard[x][y].value 			= symbol
		gameBoard[x][y].belongsToPlayer = symbol
		result = 0
	} else {
		result = 1
	}

	return
}


/* ------------------------------------------------------------------------- */

func doAIStep(gameBoard [][]GameBoardNode) {
	d("do AI step")

	var stepIsDone = 0
	for 1 != stepIsDone {
		// suppose that field is square/rectangle
		var x int = rand.Intn(gameBoardSize)
		var y int = rand.Intn(gameBoardSize)
		fmt.Printf("values: x: %d, y: %d\n", x, y)

		res := doGameStep(gameBoard, x, y, fieldPCCellId)
		if 0 != res {
			continue
		}

		stepIsDone = 1
	}

}

/* ------------------------------------------------------------------------- */

/*
 * name: printWinner
 * @param
 * 		winnerNumber --- number of player, who win the game
 * 		0: game is playing
 * 		1: first player
 * 		2: second player
 * 		3: toe
 * @return
 *
 */
func printWinner(winnerNumber int) {
	d("we have winner")
}

/* ------------------------------------------------------------------------- */

/*
 * Calculate score for one player
 *
 */
func calculateScorePerPlayer(gameBoard [][]GameBoardNode, player *Player) {


	//цикл для всех ячеек по контуру игрового поля
	//{
		// закраситьЯчеку(i, j, step)


		//funct закраситьЯчеку(i, j, stepCount){
			// stepCount - номер хода
			// stepId - значение ячейки,
			//          соответствующее о принадлежности конкретному пользователю


			//если выходим за границы или
				 //закрашено на текущем шаге или
				 //тут точка текущегоИгрока
				//возврат

			//иначе
				//устанавливаем бит закраски ячейки на шаге @stepCount
				//вызываем функцию закраски для всех соседей
					//слева,
					//справа,
					//сверху
					//снизу
		//}
	//}

	// оставшиеся незакрашенные ячейки - это внутренности контуров,
	// которые требуется *захватить* текущему пользователю
	/*
	цикл по всем ячейкам поля {

		если текущая ячейка ЗАКРАШЕНА -
			возврат
		иначе
			захватить ячейку [i][j]


		// фнкция захвата_ячейки() //

		если ячейка принадлежит игроку @текущему
			возврат
		иначе
			сменить принадлежность ячейки на текущего игрока
			увеличить счёт текущего игрока (+1)

	}

	*/



}

/* ------------------------------------------------------------------------- */

/*
 * Calculate all game score for the game (for both players)
 */
func calculateGameScore(gameBoard [][]GameBoardNode, player1 *Player, player2 *Player) {

	calculateScorePerPlayer(gameBoard, player1)
	calculateScorePerPlayer(gameBoard, player2)

}

/* ------------------------------------------------------------------------- */

/*
 * Detect, if there any winner on game board
 */
func getWinner(gameBoard [][]GameBoardNode) (winner int) {
	d("get winner")

	/* if emptyCellExist()
	 * 		return (no winner)
	 * else
	 * 		return (scoreUser > scorePC) ? userWin : pcWin;
	 *
	 * */


	winner = 0
	return
}

/* ========================================================================= */

func main() {
	fmt.Println("\n\t = = = Greeting in 'Dots' game = = =\n");

	time.Sleep(2 * time.Second)

	var mainGameBoard [][]GameBoardNode = initGameBoard(gameBoardSize)
	var userPlayer 	= Player{stepId: fieldUserCellId, score: 0}
	var pcPlayer 	= Player{stepId: fieldPCCellId, score: 0}


	var isWin int = 0
	for /* empty */; isWin == 0; /* empty */ {

		clear_screen_linux()
		drawGameBoard(mainGameBoard, &userPlayer, &pcPlayer)

		doUserStep(mainGameBoard);
		calculateGameScore(mainGameBoard, &userPlayer, &pcPlayer)

		clear_screen_linux()
		drawGameBoard(mainGameBoard, &userPlayer, &pcPlayer)
		isWin = getWinner(mainGameBoard);
		if 0 != isWin {
			break
		}

		doAIStep(mainGameBoard);
		calculateGameScore(mainGameBoard, &userPlayer, &pcPlayer)

		clear_screen_linux()
		drawGameBoard(mainGameBoard, &userPlayer, &pcPlayer)
		isWin = getWinner(mainGameBoard);
		if 0 != isWin {
			break;
		}
		d("-----------------------------------------")
	}

	printWinner(isWin);
}
