/*
 * Here description of my Dots game
 * Vit Ry <developer@bitthinker.com> (c) 2014
 *
 * ts: 4
 */
package main

import (
	"fmt"
	"time"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"bufio"
	"strconv"
	"strings"
	"math/rand"
)

// --------------------------------- CONSTS -------------------------------- //
const gameBoardSize			int = 9

const fieldEmptyCellChar	string = "."
const fieldUserCellChar		string = "*"
const fieldPCCellChar		string = "+"
const fieldCapturedCellChar	string = "#"

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

// global variable. A number of step (any player)
var stepNumber int = 0;

// ----------------------------- STRUCTS ----------------------------------- //
type GameBoardNode struct {
	value			int;		// 0, 1, 2 (Empty, User, PC)
	belongsToPlayer int;		// @id of player, who owns or captured this cell
	paintedId		int;		// @id of player, who paintOut this cell
	paintedOnStep	int;		// on the game's step @paintedOnStep
	captured		int;		// (0|1) - if this cell was captured by any player
}

type Player struct {
	stepId	int;		// 0, 1, 2 (Empty, User, PC)
	score 	int;		//  score of this player. 0 by default
}


// ----------------------------- FUNCTIONS --------------------------------- //
func d (debugMsg string) {
	fmt.Printf("D: ")
	fmt.Println(debugMsg)
	time.Sleep(200 * time.Millisecond)
}

/* ------------------------------------------------------------------------- */

func initGameBoard(size int) (gameBoard [][]GameBoardNode) {

	// Allocate the top-level slice, the same as before.
	gameBoard = make([][]GameBoardNode, size) // One row per unit of y.

	// Allocate one large slice to hold all the pixels.
	pixels := make([]GameBoardNode, size*size)

	// Loop over the rows,
	// slicing each row from the front of the remaining pixels slice.
	for i := range gameBoard {
		gameBoard[i], pixels = pixels[:size], pixels[size:]
	}

	/* a trap */
	//gameBoard[0][1].value = 1
	//gameBoard[1][0].value = 1
	//gameBoard[2][1].value = 1
	//gameBoard[1][2].value = 1
	//gameBoard[1][1].value = 2

	//gameBoard[0][1].belongsToPlayer = 1
	//gameBoard[1][0].belongsToPlayer = 1
	//gameBoard[2][1].belongsToPlayer = 1
	//gameBoard[1][2].belongsToPlayer = 1
	//gameBoard[1][1].belongsToPlayer = 2

	/* free sace */
	//gameBoard[0][0].value, gameBoard[0][0].belongsToPlayer  = 1,1
	//gameBoard[2][0].value, gameBoard[2][0].belongsToPlayer  = 1,1

	return
}

//func pause() {
	//for {

	//}
//}

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

	fmt.Printf("    ")
	var length int = len(gameBoard)
	for i := 0; i < length; i++ {
		fmt.Printf("%c ", chars[i])
    }
	fmt.Println()

	for i, row := range gameBoard {
		fmt.Printf("%2d  ", i)

		for _, cell := range row {

			if fieldEmptyCellId == cell.value {

				// empty cell
				value := fieldEmptyCellChar
				if 1 == cell.captured {
					value = fieldCapturedCellChar
				}
				fmt.Printf("%s ", value)

			} else if fieldUserCellId == cell.value {

				// User cell (blue *)
				value := fieldUserCellChar
				if cell.value != cell.belongsToPlayer {
					value = fieldCapturedCellChar
				}
				fmt.Printf("%s%s%s ", CLR_B, value, CLR_N)

			} else if fieldPCCellId == cell.value {

				// PC cell (red +)

				value := fieldPCCellChar
				if cell.value != cell.belongsToPlayer {
					value = fieldCapturedCellChar
				}
				fmt.Printf("%s%s%s ", CLR_R, value, CLR_N)
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

	// can't do step, if
	// * non empty cell
	var nonEmptyCell bool = false

	/* TODO: handle the situation:
	 * allow to do step into just captured free space
	 * deny to do step into captured free space with enemy dots.
	 * NOW: allow to do step in all captured area,
	 * because it is usefull "score" and traps.
	 * Also shoud fic alghoritm of getWinner() -
	 * is gameField have any free space yet (toe or not) */
	// # * cell already captured
	// # var capturedCell bool = false

	if fieldEmptyCellId != gameBoard[x][y].value {
		nonEmptyCell = true
	}
	//if 1 == gameBoard[x][y].captured {
		//capturedCell = true
	//}

	//fmt.Printf("Step to : %d %d. nonEmpty(%t), Captured(%t)\n", x, y, nonEmptyCell, capturedCell)

	//if nonEmptyCell || capturedCell {
	if nonEmptyCell {
		result = 1
	} else {
		gameBoard[x][y].value 			= symbol
		gameBoard[x][y].belongsToPlayer = symbol
		result = 0
	}

	// increase global variable - step number
	// (needed for calculateGameScore() function)
	stepNumber++

	return
}


/* ------------------------------------------------------------------------- */

func doAIStep(gameBoard [][]GameBoardNode) {
	d("do AI step")

	rand.Seed(time.Now().UTC().UnixNano())

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
 * 		winnerNumber --- number of player, who wins the game
 * 		0: game is playing
 * 		1: first player (User)
 * 		2: second player(PC)
 * 		3: toe
 * @return
 *
 */
func printWinner(winnerNumber int) {
	playerName := "ERROR"

	switch {
	case 0 == winnerNumber:
		playerName = "None of (game continue)"
	case 1 == winnerNumber:
		playerName = "User"
	case 2 == winnerNumber:
		playerName = "PC"
	case 3 == winnerNumber:
		playerName = "None of (TOE)"
	}

	fmt.Printf("%s player wins!\n", playerName)
}

/* ------------------------------------------------------------------------- */

/*
 * name: paintOutACell
 * desc: Paint a cell on current stepNumber, if it:
 * * in gameBoard range
 * * doesn't have current user dot yet
 * * didnn't paint out yet
 * @param
 * 		@gameBoard -- game board
 * 		@i, @j -- index of a cell
 * 		@player -- player, fow whome we count a score (who did a step)
 * @return
 *
 */
func paintOutACell(gameBoard [][]GameBoardNode, i int, j int, player *Player) {

	//fmt.Printf("[%d,%d]:%d > %d <\n", i, j, stepNumber, gameBoard[i][j].value)

	var lastBoardIndex int = len(gameBoard) -1

	// if indexes out of a gameBoard, return
	if i < 0 || j < 0 || i > lastBoardIndex || j > lastBoardIndex {
		return
	}

	// if already painted out, return
	if 0 != gameBoard[i][j].paintedId && stepNumber == gameBoard[i][j].paintedOnStep {
		return
	}

	// if this cell containts a current player's dot and didn't captured, return
	if player.stepId == gameBoard[i][j].value && gameBoard[i][j].captured == 0 {
		return
	}

	// paint Out this cell
	gameBoard[i][j].paintedId 		= player.stepId
	gameBoard[i][j].paintedOnStep 	= stepNumber

	// paint out neighboards
	paintOutACell(gameBoard, i-1,	j,	player)
	paintOutACell(gameBoard, i+1,	j,	player)
	paintOutACell(gameBoard, i,		j-1,player)
	paintOutACell(gameBoard, i,		j+1,player)
}

/* ------------------------------------------------------------------------- */

/*
 * Calculate score for one player
 */
func calculateScorePerPlayer(gameBoard [][]GameBoardNode, player *Player) {

	//fmt.Printf("D: Calculate score per player: %d\n", player.stepId);
	var lastBoardIndex int = len(gameBoard) -1
	// go over the gameBoard edge and paint over every cell

	//fmt.Printf("D: index: %d\n", lastBoardIndex);

	// iterate over all rows
	for index, _ := range gameBoard {

		// take first and last row completely
		if index == 0 || index == lastBoardIndex {

			for j, _ := range gameBoard {
				paintOutACell(gameBoard, index, j, player)
			}

		} else {
			// paint first and last element of row
			paintOutACell(gameBoard, index, 0, 				player)
			paintOutACell(gameBoard, index, lastBoardIndex, player)
		}
	}

	//debug_print_gameBoard(gameBoard);

	//pause();

    //time.Sleep(2 * time.Second)

	// not painted cells may contain captured cells
	// reset painting because of calculating score for other player on same step
	for i, row := range gameBoard {
		for j, cell := range row {

			// return, if
			// * painted
			// * value - current player
			// * value - empty cell
			alreadyPainted := false
			currentPlayersCell := false
			emptyCell := false

			if 0 != cell.paintedId {
				alreadyPainted = true
			}
			if cell.value == player.stepId {
				currentPlayersCell = true
			}

			if cell.value == fieldEmptyCellId {
				emptyCell = true
			}


			if ! alreadyPainted && emptyCell {
				gameBoard[i][j].captured = 1
			}

			if alreadyPainted || currentPlayersCell || emptyCell {
				continue
			} else {

				// capture this cell

				// if it is already this player's cell - do nothing
				if player.stepId == cell.belongsToPlayer {
					continue
				} else {
					fmt.Printf("Capture cell %d %d\n\n", i, j);
					gameBoard[i][j].belongsToPlayer = player.stepId
					gameBoard[i][j].captured = 1
					player.score += 1
				}

			}

		}
	}

	// reset painting, because of another player
	for i, _ := range gameBoard {
		for j, _ := range gameBoard[i] {

			//fmt.Printf("clean %d %d, ", i, j);
			gameBoard[i][j].paintedId = fieldEmptyCellId
		}
	}

	//fmt.Printf("----- after calculating -----------\n");
	//debug_print_gameBoard(gameBoard);

}

/* ------------------------------------------------------------------------- */

func debug_print_gameBoard(gameBoard [][]GameBoardNode) {
	// print all field in readable format

	fmt.Printf("val,pla,painId,paintS,Capt\n");

	for _, row := range gameBoard {
		for _, cell := range row {
			fmt.Printf("%d,%d,%d,%d,%-5d ", cell.value, cell.belongsToPlayer, cell.paintedId, cell.paintedOnStep, cell.captured);
		}
		fmt.Println();
	}
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
 * @ return
 * 		0: winner does not exist yet
 * 		1: first User
 * 		2: second user
 * 		3: Toe
 */
func getWinner(gameBoard [][]GameBoardNode, player1 *Player, player2 *Player) (winner int) {
	d("get winner")

	// does gameBoard have any free space yet ?
	hasEmpty := false
	for _, row := range gameBoard {
		for _, cell := range row {
			if fieldEmptyCellId == cell.value {
				hasEmpty = true
				break
			}
		}

		//fmt.Printf("hasEmpty (%t)\n", hasEmpty)

		if hasEmpty {
			break
		}
	}

	if hasEmpty {
		winner = 0
	} else if player1.score > player2.score {
		winner = 1
	} else if player1.score < player2.score {
		winner = 2
	} else {
		winner = 3
	}

	return
}

/* ========================================================================= */

func main() {

	// catch ^C signal
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
        switch sig {
        case os.Interrupt:
            //handle SIGINT
			fmt.Printf("\nOkay, bye bye, Looser!\n")
			os.Exit(0)
        case syscall.SIGTERM:
            //handle SIGTERM
        }
	}()


	// Start game
	fmt.Println("\n\t = = = Greeting in 'Dots' game = = =\n");

	//time.Sleep(2 * time.Second)

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
		isWin = getWinner(mainGameBoard, &userPlayer, &pcPlayer);
		if 0 != isWin {
			break
		}

		doAIStep(mainGameBoard);
		calculateGameScore(mainGameBoard, &userPlayer, &pcPlayer)

		//clear_screen_linux()
		//drawGameBoard(mainGameBoard, &userPlayer, &pcPlayer)
		isWin = getWinner(mainGameBoard, &userPlayer, &pcPlayer);
		if 0 != isWin {
			break;
		}
		//d("-----------------------------------------")
	}

	clear_screen_linux()
	drawGameBoard(mainGameBoard, &userPlayer, &pcPlayer)
	printWinner(isWin);
}
