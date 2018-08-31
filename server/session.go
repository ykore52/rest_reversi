package server

type SessionState int

const (
	SESSION_WAIT SessionState = iota
	SESSION_ESTABLISHED
	SESSION_CLOSE
)

const (
	EMPTY int = 0
	WHITE int = 1
	BLACK int = 2
)

const (
	MAX_BOARD int = 8
)

type Session struct {
	SessionId string
	Players   []User
	State     SessionState
	Turn      int
	Board     [][]int
}

var sessionStore map[string]*Session

func InitSessionStore(force bool) {
	if sessionStore == nil || force {
		sessionStore = make(map[string]*Session)
	}
}

func CreateSession(name string) (string, string) {

	InitSessionStore(false)

	user := CreateUser(name)

	isPaired := func(user User) bool {
		for _, s := range sessionStore {
			// pairing
			if len(s.Players) == 1 && s.State == SESSION_WAIT {
				user.SessionId = s.SessionId
				s.Players = append(s.Players, user)
				s.State = SESSION_ESTABLISHED
				return true
			}
		}
		return false
	}(user)

	if !isPaired {
		// create a new session
		sessionStore[user.SessionId] = &Session{
			SessionId: user.SessionId,
			Players:   []User{user},
			State:     SESSION_WAIT,
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
		}
	}

	userStore[user.UserId] = user

	return user.UserId, user.SessionId
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

func PutDisc(sessionId string, color int, posY, posX int) int {

	if posY < 0 || posY >= MAX_BOARD || posX < 0 || posX >= MAX_BOARD {
		return 1
	}

	board := sessionStore[sessionId].Board
	if board[posY][posX] != EMPTY {
		return 1
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
		for y, x := posY-1, posX+1; y >= 0 && x < MAX_BOARD; y, x = y-1, x+1 {
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
		for x := posX + 1; x < MAX_BOARD; x++ {
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
		for y, x := posY+1, posX+1; y < MAX_BOARD && x < MAX_BOARD; y, x = y+1, x+1 {
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
		for y := posY + 1; y < MAX_BOARD; y++ {
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
		for y, x := posY+1, posX-1; y < MAX_BOARD && x >= 0; y, x = y+1, x-1 {
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
		return 1
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
	for y := 0; y < MAX_BOARD; y++ {
		for x := 0; x < MAX_BOARD; x++ {

			if PutDisc(sessionId, color, y, x) == 0 {
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
	sessionStore[sessionId].Turn = (^(sessionStore[sessionId].Turn - 1)) + 1
}
