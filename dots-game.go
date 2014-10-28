/*
 * Here description of my game
 *
 */
package main

import "fmt"

const fieldSize int = 20

func main() {
	fmt.Println("\t = = = Greeting in 'Dots' game = = =");

	array mainGameBoard[][]

	initGameBoard

	var isWin int = 0
	while (isWin == 0)
	{
		doUserStep(mainGameBoard);
		drawGameBoard(mainGameBoard);
		isWin = getWinner(mainGameBoard);
		if isWin : break;

		doAIStep(mainGameBoard);
		drawGameBoard(mainGameBoard);
		isWin = getWinner(mainGameBoard);
		if isWin : break;
	}

	printResults(isWin);
}
