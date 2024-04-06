package model

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"tinamic/model/user"
)

func QueryUser(db *pgxpool.Pool, args interface{}) (user.User, error) {
	sql := `SELECT name,uid FROM users.user_info WHERE name=$1`
	var user user.User
	err := db.QueryRow(context.Background(), sql, args).Scan(&user.Name, &user.UID)
	if err != nil {
		return user, err
	}
	return user, nil
}

func InsertUser(db *pgxpool.Pool, user user.User) (pgconn.CommandTag, error) {
	sql := `INSERT INTO users.user_info(
			uid,name,password_hash,create_at,update_at)
          VALUES($1,$2,$3,$4,$5)`
	tag, err := db.Exec(context.Background(), sql, user.UID, user.Name, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return tag, nil
}