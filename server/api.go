package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	// APIUser is an API endpoint of User
	APIUser string = "/user"

	// APISession is an API endpoint of Session
	APISession string = "/session"

	// APISessionBoard is an API endpoint of the board information
	APISessionBoard string = "/board"

	// APISessionCand is an API endpoint that you get candidates for putting discs on the board
	APISessionCand string = "/cand"
)

// GeneralMessageResponse ...
type GeneralMessageResponse struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

// GetUserRequest ...
type GetUserRequest struct {
	Status string `json:"status"`
	Name   string `json:"name"`
}

// PostBoardRequest ...
type PostBoardRequest struct {
	UserID string `json:"userID"`
	PosX   int    `json:"posX"`
	PosY   int    `json:"posY"`
}

// GetSessionInfoResponse ...
type GetSessionInfoResponse struct {
	Status    string `json:"status"`
	Username  string `json:"username"`
	UserID    string `json:"userID"`
	SessionID string `json:"sessionID"`
}

// GetBoardResponse ...
type GetBoardResponse struct {
	Status string  `json:"status"`
	Board  [][]int `json:"board"`
}

// GetCandidatesResponse ...
type GetCandidatesResponse struct {
	Status     string  `json:"status"`
	Candidates [][]int `json:"candidates"`
}

// APIRoute ...
func APIRoute(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" && r.RequestURI == APIUser {
		APIPostUser(w, r)
	} else if match, _ := regexp.Match(APIUser+"/[a-zA-Z0-9]+", []byte(r.RequestURI)); r.Method == "GET" && match {
		APIGetUser(w, r)
	} else if match, _ := regexp.Match(APISession+"/[a-zA-Z0-9]+"+APISessionBoard, []byte(r.RequestURI)); r.Method == "GET" && match {
		APIGetBoard(w, r)
	} else if match, _ := regexp.Match(APISession+"/[a-zA-Z0-9]+"+APISessionCand, []byte(r.RequestURI)); r.Method == "GET" && match {
		APIGetCandidates(w, r)
	} else if match, _ := regexp.Match(APISession+"/[a-zA-Z0-9]+", []byte(r.RequestURI)); r.Method == "POST" && match {
		APIPostBoard(w, r)
	} else if match, _ := regexp.Match(APISession+"/[a-zA-Z0-9]+", []byte(r.RequestURI)); r.Method == "GET" && match {
		APIGetSession(w, r)
	} else {
		returnJSONMessage(w, http.StatusInternalServerError, &GeneralMessageResponse{
			Status:      "fail",
			Description: "Invalid call",
		})
	}

}

// APIPostUser ...
func APIPostUser(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		returnJSONMessage(w, http.StatusInternalServerError, &GeneralMessageResponse{
			Status:      "fail",
			Description: "Cannot read body",
		})
		return
	}

	var reqBody GetUserRequest
	fmt.Println(string(body))
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		returnJSONMessage(w, http.StatusInternalServerError, &GeneralMessageResponse{
			Status:      "fail",
			Description: "Cannot parse to json",
		})
		return
	}

	username := reqBody.Name
	userID, sessionID := CreateSession(username)

	returnJSONMessage(w, http.StatusOK, &GetSessionInfoResponse{
		Status:    "success",
		Username:  username,
		UserID:    userID,
		SessionID: sessionID,
	})

}

// APIGetUser ...
func APIGetUser(w http.ResponseWriter, r *http.Request) {

	re := regexp.MustCompile(APIUser + "/([a-zA-Z0-9]+)")
	userIDMatch := re.FindStringSubmatch(r.URL.Path)
	if len(userIDMatch) < 2 {
		returnJSONMessage(w, http.StatusInternalServerError, &GeneralMessageResponse{
			Status:      "fail",
			Description: "Invalid username",
		})
		return
	}

	userID := userIDMatch[1]
	user := GetUser(userID)
	returnJSONMessage(w, http.StatusOK, user)

}

// APIGetSession ...
func APIGetSession(w http.ResponseWriter, r *http.Request) {

	re := regexp.MustCompile(APISession + "/([a-zA-Z0-9]+)")
	sessionIDMatch := re.FindStringSubmatch(r.URL.Path)
	if len(sessionIDMatch) < 2 {
		returnJSONMessage(w, http.StatusInternalServerError, &GeneralMessageResponse{
			Status:      "fail",
			Description: "Invalid session",
		})
		return
	}

	sessionID := sessionIDMatch[1]
	session := GetSessionInfo(sessionID)
	returnJSONMessage(w, http.StatusOK, session)

}

// APIGetBoard ...
func APIGetBoard(w http.ResponseWriter, r *http.Request) {

	re := regexp.MustCompile(APISession + "/([a-zA-Z0-9]+)" + APISessionBoard)
	sessionIDMatch := re.FindStringSubmatch(r.URL.Path)
	if len(sessionIDMatch) < 2 {
		returnJSONMessage(w, http.StatusInternalServerError, &GeneralMessageResponse{
			Status:      "fail",
			Description: "Invalid session",
		})
		return
	}

	sessionID := sessionIDMatch[1]
	board := GetBoard(sessionID)
	returnJSONMessage(w, http.StatusOK, &GetBoardResponse{
		Status: "success",
		Board:  board,
	})

}

// APIGetCandidates ...
func APIGetCandidates(w http.ResponseWriter, r *http.Request) {

	re := regexp.MustCompile(APISession + "/([a-zA-Z0-9]+)" + APISessionCand)
	sessionIDMatch := re.FindStringSubmatch(r.URL.Path)
	if len(sessionIDMatch) < 2 {
		returnJSONMessage(w, http.StatusInternalServerError, &GeneralMessageResponse{
			Status:      "fail",
			Description: "Invalid session",
		})
		return
	}

	sessionID := sessionIDMatch[1]
	session := GetSessionInfo(sessionID)
	cand := FindCandidates(sessionID, session.Turn)
	returnJSONMessage(w, http.StatusOK, &GetCandidatesResponse{
		Status:     "success",
		Candidates: cand,
	})

}

// APIPostBoard ...
func APIPostBoard(w http.ResponseWriter, r *http.Request) {

	re := regexp.MustCompile(APISession + "/([a-zA-Z0-9]+)")
	sessionIDMatch := re.FindStringSubmatch(r.URL.Path)
	if len(sessionIDMatch) < 2 {
		returnJSONMessage(w, http.StatusInternalServerError, &GeneralMessageResponse{
			Status:      "fail",
			Description: "Invalid session",
		})
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		returnJSONMessage(w, http.StatusInternalServerError, &GeneralMessageResponse{
			Status:      "fail",
			Description: "Cannot read body",
		})
		return
	}

	var reqBody PostBoardRequest
	fmt.Println(string(body))
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		returnJSONMessage(w, http.StatusInternalServerError, &GeneralMessageResponse{
			Status:      "fail",
			Description: "Cannot parse to json",
		})
		return
	}

	sessionID := sessionIDMatch[1]
	session := GetSessionInfo(sessionID)
	if session.State < StateEstablished {
		returnJSONMessage(w, http.StatusOK, &GeneralMessageResponse{
			Status:      "success",
			Description: "Session does not start yet",
		})
		return
	}

	user := GetUser(reqBody.UserID)
	if user.UserID != session.Players[session.Turn-1].UserID {
		returnJSONMessage(w, http.StatusOK, &GeneralMessageResponse{
			Status:      "success",
			Description: "Not your turn",
		})
		return
	}

	if ret := PutDisc(sessionID, session.Turn, reqBody.PosX, reqBody.PosY); ret != 0 {
		returnJSONMessage(w, http.StatusOK, &GeneralMessageResponse{
			Status:      "success",
			Description: "Cannot put a disc",
		})
		return
	}

	fmt.Println("API: call UpdateSessionState")
	// update board
	UpdateSessionState(sessionID, session.Turn, reqBody.PosX, reqBody.PosY)

	board := GetBoard(sessionID)
	returnJSONMessage(w, http.StatusOK, board)
}

func returnJSONMessage(w http.ResponseWriter, returnCode int, res interface{}) {

	w.WriteHeader(returnCode)
	resJSON, err := json.Marshal(res)
	if err != nil {
		w.Write([]byte("returnJSONMessage: Cannot parse to JSON"))
		fmt.Println("returnJSONMessage: Cannot parse to JSON")
		return
	}
	w.Write(resJSON)
	fmt.Println(string(resJSON))

}
