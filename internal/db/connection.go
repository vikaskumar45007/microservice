package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect() (*sql.DB,error){
	dbStr := "postgres://postgres:postgres@host.minikube.internal:5432/practice"
	db, err := sql.Open("pgx",dbStr)
	if err != nil{
		return nil,fmt.Errorf("open db : %w",err)
	}
	if err := db.Ping();err != nil{
		return nil,fmt.Errorf("ping db : %w",err)
	}
	return db,nil
}