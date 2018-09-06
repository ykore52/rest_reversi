package server

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type SessionState int

const (
	STATE_WAIT SessionState = iota
	STATE_ESTABLISHED
	STATE_PUT_WHITE
	STATE_PUT_BLACK
	STATE_PASSED_WHITE
	STATE_PASSED_BLACK
	STATE_WON_WHITE
	STATE_WON_BLACK
	STATE_CLOSE
)

const (
	EMPTY int = 0
	WHITE int = 1
	BLACK int = 2
)

const (
	MAX_BOARD_SIZE int = 8
)

type Session struct {
	SessionId   string       `json:"sessionId"`
	Players     []User       `json:"players"`
	State       SessionState `json:"state"`
	Turn        int          `json:"turn"`
	Board       [][]int      `json:"board"`
	ElapsedTurn int          `json:"elapsed_turn"`
	LastMove    []int        `json:"last_move"`
	MoveLog     [][]int      `json:"move_log"`
}

var sessionStore map[string]*Session

func InitSessionStore(force bool) {
	if sessionStore == nil || force {
		sessionStore = make(map[string]*Session)
	}
}

func CreateSession(username string) (string, string) {

	InitSessionStore(false)

	user := CreateUser(username)

	sessionId := func(user User) string {
		for _, s := range sessionStore {
			// pairing
			if len(s.Players) == 1 && s.State == STATE_WAIT {
				s.Players = append(s.Players, user)
				s.State = STATE_ESTABLISHED
				return s.SessionId
			}
		}
		return ""
	}(user)

	if sessionId == "" {

		// create a new session
		sessionId = fmt.Sprintf("%x", sha256.Sum224([]byte((username + strconv.FormatInt(time.Now().UnixNano(), 10)))))

		sessionStore[sessionId] = &Session{
			SessionId: sessionId,
			Players:   []User{user},
			State:     STATE_WAIT,
			Turn:      1,
			Board: [][]int{
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, WHITE, BLACK, 0, 0, 0},
				{0, 0, 0, BLACK, WHITE, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
			},
			ElapsedTurn: 1,
		}
	}

	userStore[user.UserId] = user

	return user.UserId, sessionId
}

func RemoveSession(sessionId string) {
	delete(sessionStore, sessionId)
}

func GetSession() map[string]*Session {
	return sessionStore
}

func GetSessionInfo(sessionId string) *Session {
	return sessionStore[sessionId]
}

func GetBoard(sessionId string) [][]int {
	return sessionStore[sessionId].Board
}

func IsTurn(sessionId string, color int) bool {
	return sessionStore[sessionId].Turn == color
}

func PutDisc(sessionId string, color int, posX, posY int) int {

	if !IsTurn(sessionId, color) {
		return 1
	}

	if posY < 0 || posY >= MAX_BOARD_SIZE || posX < 0 || posX >= MAX_BOARD_SIZE {
		return 2
	}

	board := sessionStore[sessionId].Board
	if board[posY][posX] != EMPTY {
		return 3
	}

	// save previous state
	boardBackup := make([][]int, len(board))
	for i := range board {
		copy(boardBackup[i], board[i])
	}

	// put a disc
	board[posY][posX] = color

	// a flag whether move is available
	move := false

	// parse for up
	move = move || func() bool {
		move := false
		side_y := -1
		for y := posY - 1; y >= 0; y-- {
			if board[y][posX] == EMPTY {
				return false
			}

			if board[y][posX] == color {
				side_y = y
				break
			}
		}

		// flip discs between [posY] and [side_y]
		if side_y != -1 {
			for y := posY - 1; y > side_y; y-- {
				board[y][posX] = color
				move = true
			}
		}
		return move
	}()

	// parse for up right
	move = move || func() bool {
		move := false
		side_y, side_x := -1, -1
		for y, x := posY-1, posX+1; y >= 0 && x < MAX_BOARD_SIZE; y, x = y-1, x+1 {
			if board[y][x] == EMPTY {
				return false
			}
			if board[y][x] == color {
				side_y, side_x = y, x
				break
			}
		}

		// flip discs
		if side_y != -1 && side_x != -1 {
			for y, x := posY-1, posX+1; y > side_y && x < side_x; y, x = y-1, x+1 {
				board[y][x] = color
				move = true
			}
		}
		return move
	}()

	// parse for right
	move = move || func() bool {
		move := false
		side_x := -1
		for x := posX + 1; x < MAX_BOARD_SIZE; x++ {
			if board[posY][x] == EMPTY {
				return false
			}
			if board[posY][x] == color {
				side_x = x
				break
			}
		}

		// flip discs
		if side_x != -1 {
			for x := posX + 1; x < side_x; x++ {
				board[posY][x] = color
				move = true
			}
		}
		return move
	}()

	// parse for down right
	move = move || func() bool {
		move := false
		side_y, side_x := -1, -1
		for y, x := posY+1, posX+1; y < MAX_BOARD_SIZE && x < MAX_BOARD_SIZE; y, x = y+1, x+1 {
			if board[y][x] == EMPTY {
				return false
			}
			if board[y][x] == color {
				side_y, side_x = y, x
				break
			}
		}

		// flip discs
		if side_y != -1 && side_x != -1 {
			for y, x := posY+1, posX+1; y < side_y && x < side_x; y, x = y+1, x+1 {
				board[y][x] = color
				move = true
			}
		}
		return move
	}()

	// parse for down
	move = move || func() bool {
		fmt.Printf("begin %d %d %d\n", color, posX, posY)
		move := false
		side_y := -1
		for y := posY + 1; y < MAX_BOARD_SIZE; y++ {
			if board[y][posX] == EMPTY {
				fmt.Printf(" (%d %d) is empty\n", posX, y)
				return false
			}
			if board[y][posX] == color {
				fmt.Printf(" (%d %d) found\n", posX, y)
				side_y = y
				break
			}
		}

		// flip discs
		if side_y != -1 {
			for y := posY + 1; y < side_y; y++ {
				fmt.Printf(" flip (%d %d)\n", posX, y)
				board[y][posX] = color
				move = true
			}
		}
		return move
	}()

	// parse for down left
	move = move || func() bool {
		move := false
		side_y, side_x := -1, -1
		for y, x := posY+1, posX-1; y < MAX_BOARD_SIZE && x >= 0; y, x = y+1, x-1 {
			if board[y][x] == EMPTY {
				return false
			}
			if board[y][x] == color {
				side_y, side_x = y, x
				break
			}
		}

		// flip discs
		if side_y != -1 && side_x != -1 {
			for y, x := posY+1, posX-1; y < side_y && x > side_x; y, x = y+1, x-1 {
				board[y][x] = color
				move = true
			}
		}
		return move
	}()

	// parse for left
	move = move || func() bool {
		move := false
		side_x := -1
		for x := posX - 1; x >= 0; x-- {
			if board[posY][x] == EMPTY {
				return false
			}
			if board[posY][x] == color {
				side_x = x
				break
			}
		}

		// flip discs
		if side_x != -1 {
			for x := posX - 1; x > side_x; x-- {
				board[posY][x] = color
				move = true
			}
		}
		return move
	}()

	// parse for up left
	move = move || func() bool {
		move := false
		side_y, side_x := -1, -1
		for y, x := posY-1, posX-1; y >= 0 && x >= 0; y, x = y-1, x-1 {
			if board[y][x] == EMPTY {
				return false
			}
			if board[y][x] == color {
				side_y, side_x = y, x
				break
			}
		}

		// flip discs
		if side_y != -1 && side_x != -1 {
			for y, x := posY-1, posX-1; y > side_y && x > side_x; y, x = y-1, x-1 {
				board[y][x] = color
				move = true
			}
		}
		return move
	}()

	// rollback to previous state if no moves
	if !move {
		for i := range board {
			copy(board[i], boardBackup[i])
		}
		return 4
	}

	return 0
}

//
// Find all candidate to put a disc on to the board
//
func FindCandidates(sessionId string, color int) [][]int {

	// save previous state
	board := sessionStore[sessionId].Board
	boardBackup := make([][]int, len(board))
	for i := range board {
		boardBackup[i] = make([]int, len(board[i]))
		copy(boardBackup[i], board[i])
	}

	candidates := make([][]int, 0)
	for y := 0; y < MAX_BOARD_SIZE; y++ {
		for x := 0; x < MAX_BOARD_SIZE; x++ {

			if PutDisc(sessionId, color, x, y) == 0 {
				candidates = append(candidates, []int{y, x})
			}

			// rollback
			for i := range board {
				copy(board[i], boardBackup[i])
			}
		}
	}

	return candidates
}

func RotateTurn(sessionId string) {
	// rotate a turn
	if sessionStore[sessionId].Turn == WHITE {
		sessionStore[sessionId].Turn = BLACK
	} else {
		sessionStore[sessionId].Turn = WHITE
	}
}

func UpdateSessionState(sessionId string, color int, posX int, posY int) {

	turn := sessionStore[sessionId].Turn

	// increment number of elapsed turn
	sessionStore[sessionId].ElapsedTurn++

	sessionStore[sessionId].LastMove = []int{color, posX, posY}
	sessionStore[sessionId].MoveLog = append(sessionStore[sessionId].MoveLog, []int{color, posX, posY})

	if turn == WHITE {
		fmt.Printf("black cand: %d\n", len(FindCandidates(sessionId, BLACK)))
		if len(FindCandidates(sessionId, BLACK)) == 0 {
			if len(FindCandidates(sessionId, WHITE)) == 0 {
				sessionStore[sessionId].State = STATE_WON_WHITE
			} else {
				// pass a turn
				sessionStore[sessionId].State = STATE_PASSED_BLACK
			}
		} else {
			RotateTurn(sessionId)
			sessionStore[sessionId].State = STATE_PUT_WHITE
		}

	} else if turn == BLACK {
		fmt.Printf("white cand: %d\n", len(FindCandidates(sessionId, WHITE)))
		if len(FindCandidates(sessionId, WHITE)) == 0 {
			if len(FindCandidates(sessionId, BLACK)) == 0 {
				sessionStore[sessionId].State = STATE_WON_BLACK
			} else {
				// pass a turn
				sessionStore[sessionId].State = STATE_PASSED_WHITE
			}
		} else {
			RotateTurn(sessionId)
			sessionStore[sessionId].State = STATE_PUT_BLACK
		}
	}
}
