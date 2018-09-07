package server

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type SessionState int

const (
	StateWait SessionState = iota
	StateEstablished
	StatePutWhite
	StatePutBlack
	StatePassedWhite
	StatePassedBlack
	StateWonWhite
	StateWonBlack
	StateClose
)

const (
	EMPTY int = 0
	WHITE int = 1
	BLACK int = 2
)

const (
	MaxBoardSize int = 8
)

type Session struct {
	SessionID   string       `json:"sessionID"`
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

	sessionID := func(user User) string {
		for _, s := range sessionStore {
			// pairing
			if len(s.Players) == 1 && s.State == StateWait {
				s.Players = append(s.Players, user)
				s.State = StateEstablished
				return s.SessionID
			}
		}
		return ""
	}(user)

	if sessionID == "" {

		// create a new session
		sessionID = fmt.Sprintf("%x", sha256.Sum224([]byte((username + strconv.FormatInt(time.Now().UnixNano(), 10)))))

		sessionStore[sessionID] = &Session{
			SessionID: sessionID,
			Players:   []User{user},
			State:     StateWait,
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

	userStore[user.UserID] = user

	return user.UserID, sessionID
}

func RemoveSession(sessionID string) {
	delete(sessionStore, sessionID)
}

func GetSession() map[string]*Session {
	return sessionStore
}

func GetSessionInfo(sessionID string) *Session {
	return sessionStore[sessionID]
}

func GetBoard(sessionID string) [][]int {
	return sessionStore[sessionID].Board
}

func IsTurn(sessionID string, color int) bool {
	return sessionStore[sessionID].Turn == color
}

func PutDisc(sessionID string, color int, posX, posY int) int {

	if posY < 0 || posY >= MaxBoardSize || posX < 0 || posX >= MaxBoardSize {
		return 2
	}

	board := sessionStore[sessionID].Board
	if board[posY][posX] != EMPTY {
		return 3
	}

	// save previous state
	boardBackup := make([][]int, len(board))
	for i := range board {
		copy(boardBackup[i], board[i])
	}

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
		for y, x := posY-1, posX+1; y >= 0 && x < MaxBoardSize; y, x = y-1, x+1 {
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
		for x := posX + 1; x < MaxBoardSize; x++ {
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
		for y, x := posY+1, posX+1; y < MaxBoardSize && x < MaxBoardSize; y, x = y+1, x+1 {
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
		move := false
		side_y := -1
		for y := posY + 1; y < MaxBoardSize; y++ {
			if board[y][posX] == EMPTY {
				return false
			}
			if board[y][posX] == color {
				side_y = y
				break
			}
		}

		// flip discs
		if side_y != -1 {
			for y := posY + 1; y < side_y; y++ {
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
		for y, x := posY+1, posX-1; y < MaxBoardSize && x >= 0; y, x = y+1, x-1 {
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

	// put a disc
	board[posY][posX] = color

	return 0
}

//
// Find all candidate to put a disc on to the board
//
func FindCandidates(sessionID string, color int) [][]int {

	// save previous state
	board := sessionStore[sessionID].Board
	boardBackup := make([][]int, len(board))
	for i := range board {
		boardBackup[i] = make([]int, len(board[i]))
		copy(boardBackup[i], board[i])
	}

	candidates := make([][]int, 0)
	for y := 0; y < MaxBoardSize; y++ {
		for x := 0; x < MaxBoardSize; x++ {

			if PutDisc(sessionID, color, x, y) == 0 {
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

func RotateTurn(sessionID string) {
	// rotate a turn
	if sessionStore[sessionID].Turn == WHITE {
		sessionStore[sessionID].Turn = BLACK
	} else {
		sessionStore[sessionID].Turn = WHITE
	}
}

func UpdateSessionState(sessionID string, color int, posX int, posY int) {

	turn := sessionStore[sessionID].Turn

	// increment number of elapsed turn
	sessionStore[sessionID].ElapsedTurn++

	sessionStore[sessionID].LastMove = []int{color, posX, posY}
	sessionStore[sessionID].MoveLog = append(sessionStore[sessionID].MoveLog, []int{color, posX, posY})

	if turn == WHITE {
		fmt.Printf("black cand: %x, %d\n", FindCandidates(sessionID, BLACK), len(FindCandidates(sessionID, BLACK)))
		if len(FindCandidates(sessionID, BLACK)) == 0 {
			if len(FindCandidates(sessionID, WHITE)) == 0 {
				sessionStore[sessionID].State = StateWonWhite
			} else {
				// pass a turn
				sessionStore[sessionID].State = StatePassedBlack
			}
		} else {
			fmt.Println("Rotate a turn to BLACK")
			RotateTurn(sessionID)
			sessionStore[sessionID].State = StatePutWhite
		}

	} else if turn == BLACK {
		fmt.Printf("white cand: %d\n", len(FindCandidates(sessionID, WHITE)))
		if len(FindCandidates(sessionID, WHITE)) == 0 {
			if len(FindCandidates(sessionID, BLACK)) == 0 {
				sessionStore[sessionID].State = StateWonBlack
			} else {
				// pass a turn
				sessionStore[sessionID].State = StatePassedWhite
			}
		} else {
			fmt.Println("Rotate a turn to WHITE")
			RotateTurn(sessionID)
			sessionStore[sessionID].State = StatePutBlack
		}
	}
}
