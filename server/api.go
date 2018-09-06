package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	API_USER          string = "/user"
	API_SESSION       string = "/session"
	API_SESSION_BOARD string = "/board"
	API_SESSION_CAND  string = "/cand"
)

type GeneralMessageResponse struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

type GetUserRequest struct {
	Status string `json:"status"`
	Name   string `json:"name"`
}

type PostBoardRequest struct {
	UserId string `json:"userId"`
	PosX   int    `json:"posX"`
	PosY   int    `json:"posY"`
}

type GetSessionInfoResponse struct {
	Status    string `json:"status"`
	Username  string `json:"username"`
	UserId    string `json:"userId"`
	SessionId string `json:"sessionId"`
}

type GetBoardResponse struct {
	Status string  `json:"status"`
	Board  [][]int `json:"board"`
}

type GetCandidatesResponse struct {
	Status     string  `json:"status"`
	Candidates [][]int `json:"candidates"`
}

func ApiRoute(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" && r.RequestURI == API_USER {
		ApiPostUser(w, r)
	} else if match, _ := regexp.Match(API_USER+"/[a-zA-Z0-9]+", []byte(r.RequestURI)); r.Method == "GET" && match {
		ApiGetUser(w, r)
	} else if match, _ := regexp.Match(API_SESSION+"/[a-zA-Z0-9]+"+API_SESSION_BOARD, []byte(r.RequestURI)); r.Method == "GET" && match {
		ApiGetBoard(w, r)
	} else if match, _ := regexp.Match(API_SESSION+"/[a-zA-Z0-9]+"+API_SESSION_CAND, []byte(r.RequestURI)); r.Method == "GET" && match {
		ApiGetCandidates(w, r)
	} else if match, _ := regexp.Match(API_SESSION+"/[a-zA-Z0-9]+", []byte(r.RequestURI)); r.Method == "POST" && match {
		ApiPostBoard(w, r)
	} else if match, _ := regexp.Match(API_SESSION+"/[a-zA-Z0-9]+", []byte(r.RequestURI)); r.Method == "GET" && match {
		ApiGetSession(w, r)
	} else {
		returnJSONMessage(w, http.StatusInternalServerError, &GeneralMessageResponse{
			Status:      "fail",
			Description: "Invalid call",
		})
	}

}

func ApiPostUser(w http.ResponseWriter, r *http.Request) {

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
	userId, sessionId := CreateSession(username)

	returnJSONMessage(w, http.StatusOK, &GetSessionInfoResponse{
		Status:    "success",
		Username:  username,
		UserId:    userId,
		SessionId: sessionId,
	})

}

func ApiGetUser(w http.ResponseWriter, r *http.Request) {

	re := regexp.MustCompile(API_USER + "/([a-zA-Z0-9]+)")
	userIdMatch := re.FindStringSubmatch(r.URL.Path)
	if len(userIdMatch) < 2 {
		returnJSONMessage(w, http.StatusInternalServerError, &GeneralMessageResponse{
			Status:      "fail",
			Description: "Invalid username",
		})
		return
	}

	userId := userIdMatch[1]
	user := GetUser(userId)
	returnJSONMessage(w, http.StatusOK, user)

}

func ApiGetSession(w http.ResponseWriter, r *http.Request) {

	re := regexp.MustCompile(API_SESSION + "/([a-zA-Z0-9]+)")
	sessionIdMatch := re.FindStringSubmatch(r.URL.Path)
	if len(sessionIdMatch) < 2 {
		returnJSONMessage(w, http.StatusInternalServerError, &GeneralMessageResponse{
			Status:      "fail",
			Description: "Invalid session",
		})
		return
	}

	sessionId := sessionIdMatch[1]
	session := GetSessionInfo(sessionId)
	returnJSONMessage(w, http.StatusOK, session)

}

func ApiGetBoard(w http.ResponseWriter, r *http.Request) {

	re := regexp.MustCompile(API_SESSION + "/([a-zA-Z0-9]+)" + API_SESSION_BOARD)
	sessionIdMatch := re.FindStringSubmatch(r.URL.Path)
	if len(sessionIdMatch) < 2 {
		returnJSONMessage(w, http.StatusInternalServerError, &GeneralMessageResponse{
			Status:      "fail",
			Description: "Invalid session",
		})
		return
	}

	sessionId := sessionIdMatch[1]
	board := GetBoard(sessionId)
	returnJSONMessage(w, http.StatusOK, &GetBoardResponse{
		Status: "success",
		Board:  board,
	})

}

func ApiGetCandidates(w http.ResponseWriter, r *http.Request) {

	re := regexp.MustCompile(API_SESSION + "/([a-zA-Z0-9]+)" + API_SESSION_CAND)
	sessionIdMatch := re.FindStringSubmatch(r.URL.Path)
	if len(sessionIdMatch) < 2 {
		returnJSONMessage(w, http.StatusInternalServerError, &GeneralMessageResponse{
			Status:      "fail",
			Description: "Invalid session",
		})
		return
	}

	sessionId := sessionIdMatch[1]
	session := GetSessionInfo(sessionId)
	cand := FindCandidates(sessionId, session.Turn)
	returnJSONMessage(w, http.StatusOK, &GetCandidatesResponse{
		Status:     "success",
		Candidates: cand,
	})

}

func ApiPostBoard(w http.ResponseWriter, r *http.Request) {

	re := regexp.MustCompile(API_SESSION + "/([a-zA-Z0-9]+)")
	sessionIdMatch := re.FindStringSubmatch(r.URL.Path)
	if len(sessionIdMatch) < 2 {
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

	sessionId := sessionIdMatch[1]
	session := GetSessionInfo(sessionId)
	if session.State < STATE_ESTABLISHED {
		returnJSONMessage(w, http.StatusOK, &GeneralMessageResponse{
			Status:      "success",
			Description: "Session does not start yet",
		})
		return
	}

	user := GetUser(reqBody.UserId)
	if user.UserId != session.Players[session.Turn-1].UserId {
		returnJSONMessage(w, http.StatusOK, &GeneralMessageResponse{
			Status:      "success",
			Description: "Not your turn",
		})
		return
	}

	if ret := PutDisc(sessionId, session.Turn, reqBody.PosX, reqBody.PosY); ret != 0 {
		returnJSONMessage(w, http.StatusOK, &GeneralMessageResponse{
			Status:      "success",
			Description: "Cannot put a disc",
		})
		return
	}

	// update board
	UpdateSessionState(sessionId, session.Turn, reqBody.PosX, reqBody.PosY)

	board := GetBoard(sessionId)
	returnJSONMessage(w, http.StatusOK, board)
}

func returnJSONMessage(w http.ResponseWriter, returnCode int, res interface{}) {

	w.WriteHeader(returnCode)
	resJson, _ := json.Marshal(res)
	w.Write(resJson)
	fmt.Println(string(resJson))

}
