/*
 * Here is description of my Dots game
 * Vit Ry <developer@bitthinker.com> (c) 2014-2015
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
	"runtime"
)

// --------------------------------- CONSTS -------------------------------- //
const gameBoardSize			int = 6

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

// ----------------------------- STRUCTS ----------------------------------- //
type GameBoardNode struct {
	value			int;		// 0, 1, 2 (Empty, User, PC)
	belongsToPlayer int;		// @id of player, who owns the cell (by drawing this dot, or by capturing it)
	paintedId		int;		// @id of player, who paintOut this cell
	captured		int;		// (0|1) - if this cell was captured by any player
}

type Player struct {
	stepId	int;		// 0, 1, 2 (Empty, User, PC)
	score 	int;		//  score of this player. 0 by default
}

type GameBoardCell struct {
					//                                  ( | )
	x	int;		// x coord from top left corner (i) ( v )
	y	int;		// y coord from top left corner (j) (-->)
}

// ----------------------------- FUNCTIONS --------------------------------- //
func d (debugMsg string) {
	fmt.Printf("D: ")
	fmt.Println(debugMsg)
	//time.Sleep(200 * time.Millisecond)
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

	// debug
	//gameBoard[0][0].value, gameBoard[0][0].belongsToPlayer  = 1,1
	//gameBoard[0][1].value, gameBoard[0][1].belongsToPlayer  = 1,1
	//gameBoard[0][2].value, gameBoard[0][2].belongsToPlayer  = 1,1
	//gameBoard[1][0].value, gameBoard[1][0].belongsToPlayer  = 1,1
	//gameBoard[1][2].value, gameBoard[1][2].belongsToPlayer  = 1,1
	//gameBoard[2][0].value, gameBoard[2][0].belongsToPlayer  = 1,1
	//gameBoard[2][3].value, gameBoard[2][3].belongsToPlayer  = 1,1
	//gameBoard[3][0].value, gameBoard[3][0].belongsToPlayer  = 1,1
	//gameBoard[3][1].value, gameBoard[3][1].belongsToPlayer  = 1,1
	//gameBoard[3][2].value, gameBoard[3][2].belongsToPlayer  = 1,1

	return
}

/* -------------------------------------------------------------------------- */
/*
 * Permanent pause
 */
func pause() {
	for {
		time.Sleep(1 * time.Second)
	}
}

/* ------------------------------------------------------------------------- */

/*
 * Clear terminal screen in Unix (tested on linux and mac)
 */
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

				// empty cell (gray .)
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

/* Determine, if given cell free for player step on given game board
 * 
 * name: isCellAvailableForStep
 * @param
 * 		gameBoard: game board on which look it
 * 		x, y : cell's coords
 * @return
 * 		true:  yes
 * 		false: no
 */
func isCellAvailableForStep(gameBoard [][]GameBoardNode, x int, y int) (cellIsAvailable bool) {

	cellIsAvailable = false

	// can't do step, if
	// * non empty cell
	// * cell is inside captured area with enemy dots
	/* TODO: handle the situation:
	 * allow to do step into just captured free space
	 * deny to do step into captured free space with enemy dots.
	 * NOW: allow to do step in all captured area,
	 * because it is usefull "score" and traps. */

	var emptyCell bool = false
	// var capturedCell bool = false

	if fieldEmptyCellId == gameBoard[x][y].value {
		emptyCell = true
	}
	//if 1 == gameBoard[x][y].captured {
		//capturedCell = true
	//}

	//fmt.Printf("Step to : %d %d. nonEmpty(%t), Captured(%t)\n", x, y, nonEmptyCell, capturedCell)

	//if emptyCell || ! capturedCell {
	if emptyCell {
		cellIsAvailable = true
	}

	return
}

/* -------------------------------------------------------------------------- */

/*
 * return 0 if fine
 * 	  1 if cell not empty already or error occured
 */
func doGameStep(gameBoard [][]GameBoardNode, x int, y int, symbol int) (result int) {

	result = 1

	canDoStep := isCellAvailableForStep(gameBoard, x, y);
	if true == canDoStep {
		gameBoard[x][y].value 			= symbol
		gameBoard[x][y].belongsToPlayer = symbol
		result = 0
	}

	return
}

/* -------------------------------------------------------------------------- */

/*
 * Undo any step on given game board in given cell
 */
func undoGameStep(gameBoard [][]GameBoardNode, x int, y int) {
	gameBoard[x][y].value 			= fieldEmptyCellId
	gameBoard[x][y].belongsToPlayer = fieldEmptyCellId
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
 * desc: Paint a cell, if it:
 * * in gameBoard range
 * * doesn't have current user dot yet
 * * didn't paint out yet
 * @param
 * 		@gameBoard -- game board
 * 		@i, @j -- index of a cell
 * 		@playerStepID -- player step ID, for whome we count a score (who did a step)
 */
func paintOutACell(gameBoard [][]GameBoardNode, i int, j int, playerStepID int) {

	//fmt.Printf("[%d,%d] > %d <\n", i, j, gameBoard[i][j].value)

	var lastCellIndex int = len(gameBoard) -1

	// if indexes out of a gameBoard, return
	if i < 0 || j < 0 || i > lastCellIndex || j > lastCellIndex {
		return
	}

	// if already painted out, return
	if fieldEmptyCellId != gameBoard[i][j].paintedId {
		return
	}

	// if this cell containts a current player's dot and didn't captured, return
	if playerStepID == gameBoard[i][j].value && gameBoard[i][j].captured == 0 {
		return
	}

	// paint Out this cell
	gameBoard[i][j].paintedId 		= playerStepID

	// paint out neighboards
	paintOutACell(gameBoard, i-1,	j,	playerStepID)
	paintOutACell(gameBoard, i+1,	j,	playerStepID)
	paintOutACell(gameBoard, i,		j-1,playerStepID)
	paintOutACell(gameBoard, i,		j+1,playerStepID)
}

/* ------------------------------------------------------------------------- */

/*
 * Calculate score for one player
 */
func calculateScorePerPlayer(gameBoard [][]GameBoardNode, playerStepID int) {

	//fmt.Printf("D: Calculate score per player: %d\n", playerStepID);
	var lastCellIndex int = len(gameBoard) -1

	// Go over the gameBoard edge and paint over every cell

	// loop over all rows
	for index := range gameBoard {

		// take first and last row completely
		if index == 0 || index == lastCellIndex {
			for j := range gameBoard[index] {
				paintOutACell(gameBoard, index, j, playerStepID)
			}
		} else {
			// paint first and last element of row
			paintOutACell(gameBoard, index, 0,				playerStepID)
			paintOutACell(gameBoard, index, lastCellIndex,	playerStepID)
		}
	}

	//debug_print_gameBoard(gameBoard);
	//pause();
	//time.Sleep(2 * time.Second)

	/* Not painted cells -- captured cells.
	 * They may contain enemy's captured cells */
	for i := range gameBoard {
		for j, cell := range gameBoard[i] {

			// do nothing with 'current' cell, if
			// * it is painted
			// * it has value of current player
			// * is it empty
			alreadyPainted := false
			currentPlayersCell := false
			emptyCell := false

			if cell.paintedId != fieldEmptyCellId {
				alreadyPainted = true
			}
			if cell.value == playerStepID {
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
				// Capture this cell

				// if it is already this player's cell - do nothing
				if playerStepID == cell.belongsToPlayer {
					continue
				} else {
					//fmt.Printf("D: Capture enemy's cell [%d %c]\n\n", i, chars[j]);
					gameBoard[i][j].belongsToPlayer = playerStepID
					gameBoard[i][j].captured = 1
				}

			}

		}
	}

	// Reset painting because of calculating score for other player on same step
	for i := range gameBoard {
		for j := range gameBoard[i] {
			//fmt.Printf("clean %d %d, ", i, j);
			gameBoard[i][j].paintedId = fieldEmptyCellId
		}
	}

	//fmt.Printf("----- after calculating -----------\n");
	//debug_print_gameBoard(gameBoard);
}

/* -------------------------------------------------------------------------- */

/*
 * Функция вернёт количество точек, захваченных игроком
 */
func getScorePerPlayer(gameBoard [][]GameBoardNode, playerStepID int) (score int) {

	score = 0

	for _, row := range gameBoard {
		for _, cell := range row {
			if cell.captured == 1 && cell.belongsToPlayer == playerStepID && cell.value != playerStepID {
				score++
			}
		}
	}

	return
}

/* ------------------------------------------------------------------------- */

func debug_print_gameBoard(gameBoard [][]GameBoardNode) {
	// print all field in readable format

	fmt.Printf("val,BTP,PId,Capt\n");

	for _, row := range gameBoard {
		for _, cell := range row {
			fmt.Printf("%d,%d,%d,%-5d ", cell.value, cell.belongsToPlayer, cell.paintedId, cell.captured);
		}
		fmt.Println();
	}
}

/* ------------------------------------------------------------------------- */

/*
 * Calculate all game score for the game (for both players)
 */
func calculateScoreOnBoard(gameBoard [][]GameBoardNode, player1 *Player, player2 *Player) {

	calculateScorePerPlayer(gameBoard, player1.stepId)
	player1.score = getScorePerPlayer(gameBoard, player1.stepId);

	calculateScorePerPlayer(gameBoard, player2.stepId)
	player2.score = getScorePerPlayer(gameBoard, player2.stepId);
}

/* ------------------------------------------------------------------------- */

/*
 * Detect, if there any winner on given game board
 * @ return
 * 		0: winner does not exist yet
 * 		1: first User
 * 		2: second user
 * 		3: Toe
 */
func getWinner(gameBoard [][]GameBoardNode, player1 *Player, player2 *Player) (winner int) {
	d("f: get winner")

	// does gameBoard have any free space for step yet ?
	cellForStepExists := false

	FindFreeCell:
	for i, _ := range gameBoard {
		for j, _ := range gameBoard[i] {
			if true == isCellAvailableForStep(gameBoard, i, j) {
				cellForStepExists = true
				break FindFreeCell
			}
		}
	}

	if cellForStepExists {
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

	runtime.GOMAXPROCS(4)

	// handle ^C signal
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
			fmt.Printf("\nOkay, bye bye, Loser!\n")
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
		calculateScoreOnBoard(mainGameBoard, &userPlayer, &pcPlayer)

		clear_screen_linux()
		drawGameBoard(mainGameBoard, &userPlayer, &pcPlayer)
		isWin = getWinner(mainGameBoard, &userPlayer, &pcPlayer);
		if 0 != isWin {
			break
		}

		//doAIStepRandom(mainGameBoard);
		doAIStep(mainGameBoard, 4);
		calculateScoreOnBoard(mainGameBoard, &userPlayer, &pcPlayer)

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

/* ========================================================================= */

/* Do random AI step on given game board
 * 
 * name: doAIStepRandom
 * @param
 * 		gameBoard : game board on which do AI step
 */
func doAIStepRandom(gameBoard [][]GameBoardNode) {
	d("f: do random AI step")

	rand.Seed(time.Now().UTC().UnixNano())

	// loop untill do some step
	for {
		// suppose that field is square/rectangle
		var x int = rand.Intn(gameBoardSize)
		var y int = rand.Intn(gameBoardSize)
		fmt.Printf("D: AI: try to do step in [%d; %c]\n", x, chars[y])

		if 0 != doGameStep(gameBoard, x, y, fieldPCCellId) {
			continue
		}

		break
	}
}

/* -------------------------------------------------------------------------- */

/* Do AI step on given game board with pre-calculating on given game depth
 * 
 * name: doAIStep
 * @param
 * 		gameBoard : game board on which do AI step
 * 		depth     : depth of pre-calculating (aka level of difficulty)
 *                  bigger depth -> smarter AI
 */
func doAIStep(gameBoard [][]GameBoardNode, depth int) {

	cellToDoStepX := -1
	cellToDoStepY := -1
	cellToDoStepScore := -1000000		// default score (unreal game score)
	minScore := +1000000

	N := gameBoardSize*gameBoardSize
	sem := make(chan bool, N);  // semaphore pattern

	/* TODO: оптимизация скорости
	 * определяем игровую область,
	 * в которой будем проводить расчёты и прогнозирования
	 * HINT: для маленького поля может не потребоваться */

	// loop over free for step cells
	for i := range gameBoard {
		for j := range gameBoard[i] {

			go func (i, j int) {

				if true == isCellAvailableForStep(gameBoard, i, j) {

					// create a dublicate of game board, operate with it next
					gameBoardDuplicate := getGameBoardCopy(gameBoard);

					// do steps on fake game board and
					// look what'll happen on some depth
					doGameStep(gameBoardDuplicate, i, j, fieldPCCellId);
					tmp_score := determinePossibleGameSituation(gameBoardDuplicate, depth, true);

					fmt.Printf("=> If go to [%d, %c]: score may be %d  \n", i, chars[j], tmp_score);
					if tmp_score > cellToDoStepScore {
						cellToDoStepX = i
						cellToDoStepY = j
						cellToDoStepScore = tmp_score
					}
					if tmp_score < minScore {
						minScore = tmp_score
					}

					undoGameStep(gameBoardDuplicate, i, j);
				}

				sem <- true;

			} (i, j);


		}
	}

	// wait for all routins
	for i := 0; i < N; i++ {
		<-sem
	}

	/* TODO: If "best" cell does not exist
	 * or all cells are "same":
	 * determine step by heuristic (see web article) */
	/* for simpleness : use random now */
	if -1 == cellToDoStepX || cellToDoStepScore == minScore {
		////cellToDoStepX, cellToDoStepY = determinePCStepByHeuristic(gameBoard);
		fmt.Printf("All cells are same on given depth. No matter what to do\n");
		doAIStepRandom(gameBoard);
	} else {
		// do step on real game board
		doGameStep(gameBoard, cellToDoStepX, cellToDoStepY, fieldPCCellId);
	}
}

/* -------------------------------------------------------------------------- */

/* Return a dublicate of given game board
 * 
 * @param
 * 		gameBoard : game board to copy
 * @return
 * 		duplicate of game board (of nil, if passed so)
 */
func getGameBoardCopy (gameBoard [][]GameBoardNode) (duplicate [][]GameBoardNode) {

	if gameBoard == nil {
		return nil
	}

	duplicate = make([][]GameBoardNode, len(gameBoard))
	for i := range gameBoard {
		duplicate[i] = make([]GameBoardNode, len(gameBoard[i]))
		copy(duplicate[i], gameBoard[i])
	}

	//debug_print_gameBoard(duplicate);
	return
}

/* --------------------------------------------------------------------------- */

/*
Функция определяет ситуацию (лучшую, или худшую. в зав-сти от параметра)
на игровом поле на определённой глубине просчёта.

Ходим во все ячейки на заданную глубину и смотрим,
какой будет "счёт".
Определяем лучший счёт по всем ходам (что вообще может быть на заданную глубину)
и его возвращаем.

max/min - true/false - bestForPC/not
and xor it to change
*/
func determinePossibleGameSituation(gameBoard [][]GameBoardNode, depth int, findBestForPC bool) (gameSituation int) {

	// TODO: в зависимости от игрока (findBestForPC), ищем или максимум или минимум
	gameSituation = -10000 // maybe set it as current situation
	minScore := 10000

	if 0 != depth {

		// loop over free for step cells
		for i := range gameBoard {
			for j := range gameBoard[i] {
				if true == isCellAvailableForStep(gameBoard, i, j) {

					switch {
					case findBestForPC == true :
						// on this step we look for best situation for PC,
						// so on top level we've alredy done PC's step
						doGameStep(gameBoard, i, j, fieldUserCellId);
					case findBestForPC == false :
						doGameStep(gameBoard, i, j, fieldPCCellId);
					}

					tmp_score := determinePossibleGameSituation(gameBoard, depth-1, ! findBestForPC);

					if tmp_score > gameSituation {
						gameSituation = tmp_score
					}
					if tmp_score < minScore {
						minScore = tmp_score
					}
					//fmt.Println("=>> tmp. score: ", tmp_score, " ------------ ");

					undoGameStep(gameBoard, i, j);
				}
			}
		}

		if findBestForPC {
			gameSituation = minScore
			return gameSituation
		} else {
			return gameSituation
		}

	}

	// return current score
	userPlayer	:= Player{stepId: fieldUserCellId, score: 0}
	pcPlayer	:= Player{stepId: fieldPCCellId,   score: 0}
	calculateScoreOnBoard(gameBoard, &userPlayer, &pcPlayer)

	//if findBestForPC {
		//gameSituation = pcPlayer.score - userPlayer.score
	//} else {
		//gameSituation = userPlayer.score - pcPlayer.score
	//}
	//fmt.Println("D: dps: User score: ", userPlayer.score);
	//fmt.Println("D: dps: PC score: ", pcPlayer.score);
	gameSituation = pcPlayer.score - userPlayer.score

	return
}


/* --------------------------------------------------------------------------- */



