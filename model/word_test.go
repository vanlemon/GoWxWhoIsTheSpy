package model

import (
	"fmt"
	logs "lmf.mortal.com/GoLogs"
	"testing"
)

func TestWordDao_RandomQueryWord(t *testing.T) {
	for i := 0; i < 10; i++ {
		word, err := WordDaoInstance().RandomQueryWord(logs.TestCtx())
		fmt.Println(word, err)
	}
}

func TestWordDao_RandomQueryWordList(t *testing.T) {
	for i := 0; i < 10; i++ {
		wordList, err := WordDaoInstance().RandomQueryWordList(logs.TestCtx(), 3)
		for _, word := range wordList {
			fmt.Println(word)
		}
		fmt.Println(err)
	}
}
