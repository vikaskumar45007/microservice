package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct{
	DB *sql.DB
}

type NewUser struct{
	Name string `json:"name"`
	Uid int `json:"uid"`
}

func NewHandler(db *sql.DB) *Handler{
	return &Handler{DB: db}
}



func (h *Handler) Create(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w,"Method not allowed.",http.StatusMethodNotAllowed)
		return
	}
	var user NewUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil{
		http.Error(w,"Not valid body.",http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Uid == 0{
		http.Error(w,"Please pass correct values. Name/uid can't be empty or 0.",http.StatusBadRequest)
		return
	}

	_,err := h.DB.Exec(`INSERT INTO users (name, userid) VALUES($1,$2)`,user.Name,user.Uid)
	if err != nil{
		http.Error(w,fmt.Sprintf("Unable to insert data. %s",err),http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]any{
		"id" : user.Uid,
        "message": "User created successfully",
    })
	
}