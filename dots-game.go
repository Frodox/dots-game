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
const gameBoardSize			int = 4

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
	belongsToPlayer int;		// @id of player, who owns the cell (by drawing this dot, or by capturing it)
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

	gameBoard[0][0].value, gameBoard[0][0].belongsToPlayer  = 1,1
	gameBoard[0][1].value, gameBoard[0][1].belongsToPlayer  = 1,1
	gameBoard[0][2].value, gameBoard[0][2].belongsToPlayer  = 1,1
	gameBoard[1][0].value, gameBoard[1][0].belongsToPlayer  = 1,1
	gameBoard[1][3].value, gameBoard[1][3].belongsToPlayer  = 1,1
	gameBoard[2][0].value, gameBoard[2][0].belongsToPlayer  = 1,1
	gameBoard[2][3].value, gameBoard[2][3].belongsToPlayer  = 1,1
	gameBoard[3][0].value, gameBoard[3][0].belongsToPlayer  = 1,1
	gameBoard[3][1].value, gameBoard[3][1].belongsToPlayer  = 1,1
	gameBoard[3][2].value, gameBoard[3][2].belongsToPlayer  = 1,1

	return
}

func pause() {
	for {

	}
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

	result = 1

	canDoStep := isCellAvailableForStep(gameBoard, x, y);
	if true == canDoStep {
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
 * * didn't paint out yet
 * @param
 * 		@gameBoard -- game board
 * 		@i, @j -- index of a cell
 * 		@player -- player, for whome we count a score (who did a step)
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
 * Detect, if there any winner on given game board
 * @ return
 * 		0: winner does not exist yet
 * 		1: first User
 * 		2: second user
 * 		3: Toe
 */
func getWinner(gameBoard [][]GameBoardNode, player1 *Player, player2 *Player) (winner int) {
	d("get winner")

	// does gameBoard have any free space for step yet ?
	cellForStepExists := false

	FindFreeCell:
	for i, _ := range gameBoard {
		for j, _ := range gameBoard[i] {
			if true == isCellAvailableForStep(gameBoard, i, j) {
				cellForStepExists = true
				fmt.Printf("BREAAAAAAK");
				break FindFreeCell
			} else {
				fmt.Printf("[%d;%d] - busy\n", i, j);
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

		doAIStepRandom(mainGameBoard);
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

/* ========================================================================= */

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

/* --------------------------------------------------------------------------- */

func doAIStepRandom(gameBoard [][]GameBoardNode) {
	d("do random AI step")

	rand.Seed(time.Now().UTC().UnixNano())

	// loop untill do some step
	for {
		// suppose that field is square/rectangle
		var x int = rand.Intn(gameBoardSize)
		var y int = rand.Intn(gameBoardSize)
		fmt.Printf("values: x: %d, y: %d\n", x, y)

		res := doGameStep(gameBoard, x, y, fieldPCCellId)
		if 0 != res {
			continue
		}

		break
	}
}

/*
## проектирование алгоритма Просчёта хода компьютера
* уровень сложности - глубина просчёта
*/

/*
Функция, совершающая ход на предоставленном поле
за искусственный интеллект(комп)
*/
/*
func doAIStep(игровое_поле, level int, максимум/минимум)
{

	// TODO: в зависимости от игрока, ищем или максимум или минимум

	клетка_для_хода = nill (какая-то случайная клетка, или ничего)
	очки_клетки = -1000000
	первое_новое_очко = -1000000
	первый_ход = true
*/
	/* TODO: оптимизация скорости
	 * определяем игровую область,
	 * в которой будем проводить расчёты и прогнозирования
	 * HINT: для маленького поля может не потребоваться */
/*
	for текущая_клетка in всё_игровое_поле (клетки, куда можно ходить)
	{
		doGameStep(поле2, текущая_клетка, за_ПК);
		tmp_score = определить_ситуацию(поле2, глубина);

		если tmp_score > очки_клетки
		{
			клетка_для_хода = текущая_клетка
			очки_клетки = tmp_score
			если первый_ход
			{
				первое_новое_очко = очки_клетки
				первый_ход = false
			}
		}
	}

	если
		(клетка_для_хода == nill) || (очки_клетки == первое_новое_очко)<(ходы во все клетки - равнозначены)
	{
		// определяем ход используя методы эвристики (на сайте алгоритм)
		клетка_для_хода = определить_ход_используй_эвристику(игровое_поле);
	}

	doGameStep(игровое_поле, клетка_для_хода, за_ПК);
}
*/

/* --------------------------------------------------------------------------- */

/*
Функция определяет ситуацию (лучшую, или худшую. Зависит от того, для какого игрока вызываем)
на игровом поле на определённой глубине просчёта.

max/min - true/false
and xor it to change

func определить_ситуацию(поле2, глубина, max/min)
{
	лучшие_очки = -10000

	for текущая_клетка in всё_игровое_поле2 (клетки, куда можно ходить) // предусмотреть что конец поля (некуда ходить)
	{
		doGameStep(поле3, текущая_клетка, за_игрока(max/min));
		tmp_score = определить_ситуацию(поле3, глубина-1, xor(max/min));

		если tmp_score > лучшие_очки // тут добавить max/min
		{
			лучшие_очки = tmp_score
		}
	}

	если глубина = 0 или конец_поля
	{
		подсчитать_счёт_на_поле(поле3)
		лучшие_очки  = счёт_игрока1 - счёт_игрока2
	}

	return лучшие_очки
}
*/

/* --------------------------------------------------------------------------- */

/*
Функция вернёт 2 значения -
* количество точек, захваченных первым игроком
* количество точек, захваченных вторым игроком
func подсчитать_счёт_на_поле(поле3)
{
	цикл с полсчётом разницы - кладеет и захвачена :
	belongsToPlayer
	paintedId
	captured

	для первого  и второга игроков
}
*/
