package model

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

func QueryUser(db *pgxpool.Pool, args interface{}) (User, error) {
	sql := `SELECT name,uid FROM users.user_info WHERE name=$1`
	var user User
	err := db.QueryRow(context.Background(), sql, args).Scan(&user.Name, &user.Id)
	if err != nil {
		return user, err
	}
	return user, nil
}

func InsertUser(db *pgxpool.Pool, user User) (pgconn.CommandTag, error) {
	sql := `INSERT INTO users.user_info(
			uid,name,password_hash,create_at,update_at)
          VALUES($1,$2,$3,$4,$5)`
	tag, err := db.Exec(context.Background(), sql, user.Id, user.Name, user.Password, user.Created, user.Edited)
	if err != nil {
		return nil, err
	}
	return tag, nil
}
