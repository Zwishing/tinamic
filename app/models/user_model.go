/*
 * @Author: your name
 * @Date: 2021-11-28 10:17:19
 * @LastEditTime: 2021-11-28 13:46:45
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \tinamic\app\models\user_model.go
 */
package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// User struct to describe User object.
type User struct {
	UID          uuid.UUID `db:"uid" json:"id" validate:"required,uuid"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	Email        string    `db:"email" json:"email" validate:"required,email,lte=255"`
	PasswordHash string    `db:"password_hash" json:"password_hash,omitempty" validate:"required,lte=255"`
	UserStatus   int       `db:"user_status" json:"user_status" validate:"required,len=1"`
	UserRole     string    `db:"user_role" json:"user_role" validate:"required,lte=25"`
}
