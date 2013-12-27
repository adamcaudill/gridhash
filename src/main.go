/**
 * Created by adamcaudill on 12/27/13.
 */
package main

import (
	"fmt"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

var grid [][][]byte

func main() {
	start := time.Now()
	hash := gridhash("test", 64, 1000, 1000, []byte("testtest"))
	elapsed := time.Since(start)

	fmt.Println("Final: ", hex.EncodeToString(hash))
	fmt.Printf("Hash took %s", elapsed)
}

func gridhash(password string, grid_size int, hmac_iter int, hash_iter int, salt []byte) []byte {
	//setup grid to hold hashes
	grid = make([][][]byte, grid_size)
	for i := range grid {
		grid[i] = make([][]byte, grid_size)
	}

	for i := 0; i < grid_size; i++ {
		round(i, password, hmac_iter, hash_iter, salt)
	}

	return grid[grid_size-1][grid_size-1]
}

func round(index int, password string, hmac_iter int, hash_iter int, salt []byte) {
	//fmt.Println("Index: ", index)
	if (index == 0) {
    //special case to bootstrap the process for later

		//set 0:0
		grid[index][index] = kdf([]byte(password), salt, password, hmac_iter)

		//set 0:1
		grid[index][index+1] = hash(append(grid[index][index], salt...), hash_iter)
	} else {
    //top down
		for y := 0; y <= index-1; y++ {
			process_cell(index, y, hmac_iter, hash_iter, salt, password)
		}

		//left to right
		for x := 0; x <= index; x++ {
			process_cell(x, index, hmac_iter, hash_iter, salt, password)
		}
	}
}

func process_cell(row int, column int, hmac_iter int, hash_iter int, salt []byte, password string) {
  //make sure this isn't a cell that we bootstrapped
	if (grid[row][column] == nil) {
		var value []byte
		var ret []byte

		//first, the pass going bottom up
		for idx := row - 1; idx >= 0; idx-- {
			value = append(value, grid[idx][column]...)
		}

		//next, go right to left
		for idx := column - 1; idx >= 0; idx-- {
			value = append(value, grid[row][idx]...)
		}

		//is this a key cell?
		if (row == column) {
			ret = kdf(value, salt, password, hmac_iter)
		} else {
			ret = hash(append(value, salt...), hash_iter)
		}

		grid[row][column] = ret
		//fmt.Println("Cell: ", row, column, hex.EncodeToString(grid[row][column]))
	}
}

func kdf(value []byte, salt []byte, password string, iterations int) []byte {
  for i := 0; i < iterations; i++ {
		mac := hmac.New(sha256.New, []byte(password))
		mac.Write(append(value, salt...))
		value = mac.Sum(nil)
	}

	return value
}

func hash(value []byte, iterations int) []byte {
	for i := 0; i < iterations; i++ {
		h := sha256.Sum256(value)
		s := string(h[:])
		value = []byte(s)
	}

	return value
}
