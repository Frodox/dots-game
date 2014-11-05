/*
 * Here description of my game
 *
 */
package main

import (
	"fmt"
	"time"
)

const gameBoardSize int = 20

const fieldEmptyCellChar 	string = "."
const fieldUserCellChar 	string = "*"
const fieldPCCellChar 		string = "+"

const fieldEmptyCellId 	int = 0
const fieldUserCellId 	int = 1
const fieldPCCellId		int = 2

const CLR_0 = "\x1b[30;1m"
const CLR_R = "\x1b[31;1m"
const CLR_G = "\x1b[32;1m"
const CLR_Y = "\x1b[33;1m"
const CLR_B = "\x1b[34;1m"
const CLR_M = "\x1b[35;1m"
const CLR_C = "\x1b[36;1m"
const CLR_W = "\x1b[37;1m"
const CLR_N = "\x1b[0m"

func d (debugMsg string) {
	fmt.Printf("D: ")
	fmt.Println(debugMsg)
	time.Sleep(time.Second)
}

func initGameBoard(size int) (gameBoard [][]int) {
	d("Init some game field")

	// Allocate the top-level slice, the same as before.
	gameBoard = make([][]int, size) // One row per unit of y.

	// Allocate one large slice to hold all the pixels.
	pixels := make([]int, size*size) // Has type []uint8 even though picture is [][]uint8.

	// Loop over the rows, slicing each row from the front of the remaining pixels slice.
	for i := range gameBoard {
		gameBoard[i], pixels = pixels[:size], pixels[size:]
	}

	//test
	gameBoard[2][3] = 1 // user
	gameBoard[3][5] = 2 // PC



	return
}

func doUserStep(gameBoard [][]int) {
	d("do user step")
}

/*
 * name: drawGameBoard
 * @param
 * @return
 *
 * 0 mean empty cell -- draw fieldEmptyCellChar
 * 1 mean user  cell -- draw fieldUserCellChar (blue)
 * 2 mean PC    cell -- draw fieldPCCellChar (red)
 */
func drawGameBoard(gameBoard [][]int) {
	d("draw game board")

	for i := range gameBoard {
		for j:= range gameBoard[i] {

			if fieldEmptyCellId == gameBoard[i][j] {
				fmt.Printf("%s ", fieldEmptyCellChar)
			} else if fieldUserCellId == gameBoard[i][j] {
				fmt.Printf("%s%s%s ", CLR_B, fieldUserCellChar, CLR_N)
			} else if fieldPCCellId == gameBoard[i][j] {
				fmt.Printf("%s%s%s ", CLR_R, fieldPCCellChar, CLR_N)
			}
		}
		fmt.Println("")
	}
}

func getWinner(gameBoard [][]int) (winner int) {
	d("get winner")

	winner = 0
	return
}


func doAIStep(gameBoard [][]int) {
	d("do AI step")

}

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

/* ========================================================================= */

func main() {
	fmt.Println("\n\t = = = Greeting in 'Dots' game = = =\n");



	var mainGameBoard [][]int = initGameBoard(gameBoardSize)
	drawGameBoard(mainGameBoard)

	var isWin int = 0
	for /* empty */; isWin == 0; /* empty */ {
		d("test")

		doUserStep(mainGameBoard);
		drawGameBoard(mainGameBoard);
		isWin = getWinner(mainGameBoard);
		if 0 != isWin {
			break
		}

		doAIStep(mainGameBoard);
		drawGameBoard(mainGameBoard);
		isWin = getWinner(mainGameBoard);
		if 0 != isWin {
			break;
		}
		d("-----------------------------------------")
	}

	printWinner(isWin);
}
