package server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {

	userId, sessionId := CreateSession("test")
	userId2, sessionId2 := CreateSession("test2")

	assert.Equal(t, sessionId, sessionId2)

	s := GetSessionInfo(sessionId)
	assert.Equal(t, sessionId, s.SessionId)

	RemoveSession(sessionId)
	s = GetSessionInfo(sessionId)
	assert.Nil(t, s)

	u := GetUser(userId)
	assert.Equal(t, u.Name, "test")
	u = GetUser(userId2)
	assert.Equal(t, u.Name, "test2")
}

func TestGetBoard(t *testing.T) {
	_, sessionId := CreateSession("test")
	_, _ = CreateSession("test2")

	assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionId)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 1 2 0 0 0] [0 0 0 2 1 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
}

func TestPutDisc(t *testing.T) {

	if true {
		_, sessionId := CreateSession("test")
		_, _ = CreateSession("test2")

		// put a disc to out of the board
		assert.NotEqual(t, 0, PutDisc(sessionId, WHITE, -2, 4))
		assert.NotEqual(t, 0, PutDisc(sessionId, WHITE, -1, 4))
		assert.NotEqual(t, 0, PutDisc(sessionId, WHITE, MAX_BOARD_SIZE, 4))
		assert.NotEqual(t, 0, PutDisc(sessionId, WHITE, MAX_BOARD_SIZE+1, 4))
		assert.NotEqual(t, 0, PutDisc(sessionId, WHITE, 4, -2))
		assert.NotEqual(t, 0, PutDisc(sessionId, WHITE, 4, -1))
		assert.NotEqual(t, 0, PutDisc(sessionId, WHITE, 4, MAX_BOARD_SIZE))
		assert.NotEqual(t, 0, PutDisc(sessionId, WHITE, 4, MAX_BOARD_SIZE+1))

		// put a disc to grid existing any disc
		assert.NotEqual(t, 0, PutDisc(sessionId, WHITE, 3, 3))
		assert.NotEqual(t, 0, PutDisc(sessionId, WHITE, 3, 4))
		assert.NotEqual(t, 0, PutDisc(sessionId, WHITE, 4, 3))
		assert.NotEqual(t, 0, PutDisc(sessionId, WHITE, 4, 4))
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionId := CreateSession("test")
		_, _ = CreateSession("test2")

		// put a disc to available grid
		assert.Equal(t, 0, PutDisc(sessionId, WHITE, 4, 2))
		assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionId)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 1 0 0 0] [0 0 0 1 1 0 0 0] [0 0 0 2 1 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionId := CreateSession("test")
		_, _ = CreateSession("test2")

		// put a disc to available grid
		assert.Equal(t, 0, PutDisc(sessionId, WHITE, 5, 3))
		assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionId)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 1 1 1 0 0] [0 0 0 2 1 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionId := CreateSession("test")
		_, _ = CreateSession("test2")

		// put a disc to available grid
		assert.Equal(t, 0, PutDisc(sessionId, WHITE, 2, 4))
		assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionId)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 1 2 0 0 0] [0 0 1 1 1 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionId := CreateSession("test")
		_, _ = CreateSession("test2")

		// put a disc to available grid
		assert.Equal(t, 0, PutDisc(sessionId, WHITE, 3, 5))
		assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionId)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 1 2 0 0 0] [0 0 0 1 1 0 0 0] [0 0 0 1 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionId := CreateSession("test")
		_, _ = CreateSession("test2")

		assert.Equal(t, 0, PutDisc(sessionId, WHITE, 3, 5))
		RotateTurn(sessionId)
		assert.Equal(t, 0, PutDisc(sessionId, BLACK, 2, 5))
		assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionId)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 1 2 0 0 0] [0 0 0 2 1 0 0 0] [0 0 2 1 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionId := CreateSession("test")
		_, _ = CreateSession("test2")

		assert.Equal(t, 0, PutDisc(sessionId, WHITE, 5, 3))
		RotateTurn(sessionId)
		assert.Equal(t, 0, PutDisc(sessionId, BLACK, 5, 2))
		assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionId)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 2 0 0] [0 0 0 1 2 1 0 0] [0 0 0 2 1 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionId := CreateSession("test")
		_, _ = CreateSession("test2")

		RotateTurn(sessionId)
		assert.Equal(t, 0, PutDisc(sessionId, BLACK, 2, 3))
		RotateTurn(sessionId)
		assert.Equal(t, 0, PutDisc(sessionId, WHITE, 2, 2))
		assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionId)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 1 0 0 0 0 0] [0 0 2 1 2 0 0 0] [0 0 0 2 1 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionId := CreateSession("test")
		_, _ = CreateSession("test2")

		RotateTurn(sessionId)
		assert.Equal(t, 0, PutDisc(sessionId, BLACK, 5, 4))
		RotateTurn(sessionId)
		assert.Equal(t, 0, PutDisc(sessionId, WHITE, 5, 5))
		assert.Equal(t, fmt.Sprintf("%x", GetBoard(sessionId)), "[[0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0] [0 0 0 1 2 0 0 0] [0 0 0 2 1 2 0 0] [0 0 0 0 0 1 0 0] [0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0]]")
	}

	if true {
		InitSessionStore(true)
		InitUserStore(true)
		_, sessionId := CreateSession("test")
		_, _ = CreateSession("test2")

		// put a disc to unavailable grid
		PutDisc(sessionId, WHITE, 2, 4)
		RotateTurn(sessionId)
		assert.NotEqual(t, 0, PutDisc(sessionId, BLACK, 3, 5))
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
	_, sessionId := CreateSession("test")
	_, _ = CreateSession("test2")

	cand := FindCandidates(sessionId, WHITE)
	assert.Equal(t, fmt.Sprintf("%x", cand), "[[2 4] [3 5] [4 2] [5 3]]")
}

func TestSearchCandidates2(t *testing.T) {
	InitSessionStore(true)
	InitUserStore(true)
	_, sessionId := CreateSession("test")
	_, _ = CreateSession("test2")

	fmt.Println("TestSearchCandidate2")
	PutDisc(sessionId, WHITE, 4, 2)
	RotateTurn(sessionId)
	cand := FindCandidates(sessionId, BLACK)
	fmt.Println(cand)
	assert.Equal(t, fmt.Sprintf("%x", cand), "[[2 3] [2 5] [4 5]]")
}
