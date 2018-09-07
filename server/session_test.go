package server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {

	userID, sessionID := CreateSession("test")
	userID2, sessionID2 := CreateSession("test2")

	assert.Equal(t, sessionID, sessionID2)

	s := GetSessionInfo(sessionID)
	assert.Equal(t, sessionID, s.SessionID)

	RemoveSession(sessionID)
	s = GetSessionInfo(sessionID)
	assert.Nil(t, s)

	u := GetUser(userID)
	assert.Equal(t, u.Name, "test")
	u = GetUser(userID2)
	assert.Equal(t, u.Name, "test2")
}

func TestGetBoard(t *testing.T) {
	_, sessionID := CreateSession("test")
	_, _ = CreateSession("test2")

	assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionID)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 1 2 0 0 0] [0 0 0 2 1 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
}

func TestPutDisc(t *testing.T) {

	if true {
		_, sessionID := CreateSession("test")
		_, _ = CreateSession("test2")

		// put a disc to out of the board
		assert.NotEqual(t, 0, PutDisc(sessionID, WHITE, -2, 4))
		assert.NotEqual(t, 0, PutDisc(sessionID, WHITE, -1, 4))
		assert.NotEqual(t, 0, PutDisc(sessionID, WHITE, MaxBoardSize, 4))
		assert.NotEqual(t, 0, PutDisc(sessionID, WHITE, MaxBoardSize+1, 4))
		assert.NotEqual(t, 0, PutDisc(sessionID, WHITE, 4, -2))
		assert.NotEqual(t, 0, PutDisc(sessionID, WHITE, 4, -1))
		assert.NotEqual(t, 0, PutDisc(sessionID, WHITE, 4, MaxBoardSize))
		assert.NotEqual(t, 0, PutDisc(sessionID, WHITE, 4, MaxBoardSize+1))

		// put a disc to grid existing any disc
		assert.NotEqual(t, 0, PutDisc(sessionID, WHITE, 3, 3))
		assert.NotEqual(t, 0, PutDisc(sessionID, WHITE, 3, 4))
		assert.NotEqual(t, 0, PutDisc(sessionID, WHITE, 4, 3))
		assert.NotEqual(t, 0, PutDisc(sessionID, WHITE, 4, 4))
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionID := CreateSession("test")
		_, _ = CreateSession("test2")

		// put a disc to available grid
		assert.Equal(t, 0, PutDisc(sessionID, WHITE, 4, 2))
		assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionID)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 1 0 0 0] [0 0 0 1 1 0 0 0] [0 0 0 2 1 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionID := CreateSession("test")
		_, _ = CreateSession("test2")

		// put a disc to available grid
		assert.Equal(t, 0, PutDisc(sessionID, WHITE, 5, 3))
		assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionID)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 1 1 1 0 0] [0 0 0 2 1 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionID := CreateSession("test")
		_, _ = CreateSession("test2")

		// put a disc to available grid
		assert.Equal(t, 0, PutDisc(sessionID, WHITE, 2, 4))
		assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionID)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 1 2 0 0 0] [0 0 1 1 1 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionID := CreateSession("test")
		_, _ = CreateSession("test2")

		// put a disc to available grid
		assert.Equal(t, 0, PutDisc(sessionID, WHITE, 3, 5))
		assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionID)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 1 2 0 0 0] [0 0 0 1 1 0 0 0] [0 0 0 1 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionID := CreateSession("test")
		_, _ = CreateSession("test2")

		assert.Equal(t, 0, PutDisc(sessionID, WHITE, 3, 5))
		RotateTurn(sessionID)
		assert.Equal(t, 0, PutDisc(sessionID, BLACK, 2, 5))
		assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionID)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 1 2 0 0 0] [0 0 0 2 1 0 0 0] [0 0 2 1 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionID := CreateSession("test")
		_, _ = CreateSession("test2")

		assert.Equal(t, 0, PutDisc(sessionID, WHITE, 5, 3))
		RotateTurn(sessionID)
		assert.Equal(t, 0, PutDisc(sessionID, BLACK, 5, 2))
		assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionID)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 2 0 0] [0 0 0 1 2 1 0 0] [0 0 0 2 1 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionID := CreateSession("test")
		_, _ = CreateSession("test2")

		RotateTurn(sessionID)
		assert.Equal(t, 0, PutDisc(sessionID, BLACK, 2, 3))
		RotateTurn(sessionID)
		assert.Equal(t, 0, PutDisc(sessionID, WHITE, 2, 2))
		assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionID)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 1 0 0 0 0 0] [0 0 2 1 2 0 0 0] [0 0 0 2 1 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionID := CreateSession("test")
		_, _ = CreateSession("test2")

		RotateTurn(sessionID)
		assert.Equal(t, 0, PutDisc(sessionID, BLACK, 5, 4))
		RotateTurn(sessionID)
		assert.Equal(t, 0, PutDisc(sessionID, WHITE, 5, 5))
		assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionID)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 1 2 0 0 0] [0 0 0 2 1 2 0 0] [0 0 0 0 0 1 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionID := CreateSession("test")
		_, _ = CreateSession("test2")

		// put a disc to unavailable grid
		PutDisc(sessionID, WHITE, 2, 4)
		RotateTurn(sessionID)
		assert.NotEqual(t, 0, PutDisc(sessionID, BLACK, 3, 5))
	}
}

func TestSearchCandidates(t *testing.T) {
	/*
	   01234567
	  0
	  1
	  2    x
	  3   12x
	  4  x21
	  5   x
	  6
	  7
	*/
	InitSessionStore(true)
	InitUserStore(true)
	_, sessionID := CreateSession("test")
	_, _ = CreateSession("test2")

	cand := FindCandidates(sessionID, WHITE)
	assert.Equal(t, fmt.Sprintf("%x", cand), "[[2 4] [3 5] [4 2] [5 3]]")
}

func TestSearchCandidates2(t *testing.T) {
	InitSessionStore(true)
	InitUserStore(true)
	_, sessionID := CreateSession("test")
	_, _ = CreateSession("test2")

	fmt.Println("TestSearchCandidate2")
	PutDisc(sessionID, WHITE, 4, 2)
	RotateTurn(sessionID)
	cand := FindCandidates(sessionID, BLACK)
	fmt.Println(cand)
	assert.Equal(t, fmt.Sprintf("%x", cand), "[[2 3] [2 5] [4 5]]")
}
