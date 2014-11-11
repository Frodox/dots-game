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

const gameBoardSize 		int = 20

const fieldEmptyCellChar 	string = "."
const fieldUserCellChar 	string = "*"
const fieldPCCellChar 		string = "+"

const fieldEmptyCellId 		int = 0
const fieldUserCellId 		int = 1
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

const chars  string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+"

func d (debugMsg string) {
	fmt.Printf("D: ")
	fmt.Println(debugMsg)
	time.Sleep(500 * time.Millisecond)
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

	return
}

func clear_screen_linux() {
        cmd := exec.Command("clear") //Linux example, its tested
        cmd.Stdout = os.Stdout
        cmd.Run()
    }

func doUserStep(gameBoard [][]int) {

	var firstIndex int  = 0
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

	var length int = len(gameBoard)

	fmt.Printf("    ")
	for i := 0; i < length; i++ {
		fmt.Printf("%c ", chars[i])
    }
	fmt.Println()
	for i := range gameBoard {
		fmt.Printf("%2d  ", i)

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

/*
 * return 0 if fine
 * 	  1 if cell not empty already or error occured
 */
func doGameStep(gameBoard [][]int, x int, y int, symbol int) (result int) {

	if fieldEmptyCellId == gameBoard[x][y] {
		gameBoard[x][y] = symbol
		result = 0
	} else {
		result = 1
	}

	return
}

func getWinner(gameBoard [][]int) (winner int) {
	d("get winner")

	winner = 0
	return
}


func doAIStep(gameBoard [][]int) {
	d("do AI step")

	var stepIsDone = 0
	for 1 != stepIsDone {
		// suppose that field is square/rectangle
		var x int = rand.Intn(len(gameBoard))
		var y int = rand.Intn(len(gameBoard[0]))
		fmt.Printf("values: x: %d, y: %d\n", x, y)

		res := doGameStep(gameBoard, x, y, fieldPCCellId)
		if 0 != res {
			continue
		}

		stepIsDone = 1
	}

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

	var isWin int = 0
	for /* empty */; isWin == 0; /* empty */ {

		clear_screen_linux()
		drawGameBoard(mainGameBoard)

		doUserStep(mainGameBoard);

		clear_screen_linux()
		drawGameBoard(mainGameBoard);
		isWin = getWinner(mainGameBoard);
		if 0 != isWin {
			break
		}

		doAIStep(mainGameBoard);
		clear_screen_linux()
		drawGameBoard(mainGameBoard);
		isWin = getWinner(mainGameBoard);
		if 0 != isWin {
			break;
		}
		d("-----------------------------------------")
	}

	printWinner(isWin);
}
