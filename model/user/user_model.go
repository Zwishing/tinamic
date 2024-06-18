package user

import "tinamic/util"

type AccountCategory = int8

const (
	Name AccountCategory = iota
	Email
	Phone
)

type PermssionCategory = int8

const (
	Write PermssionCategory = iota + 1
	Read
)

// Account 账号
type Account struct {
	Id           int             `json:"id,omitempty"`
	UserId       int             `json:"user_id,omitempty"`
	LoginAccount string          `json:"login_account,omitempty"`
	Category     AccountCategory `json:"category,omitempty"`
	*BaseRecord
	*BaseRecorder
}

func NewAccount(userId int, loginAccount string, category int8) *Account {
	return &Account{UserId: userId, LoginAccount: loginAccount, Category: category}
}

// User 用户
type User struct {
	Id        int    `json:"id,omitempty"`
	State     bool   `json:"state,omitempty"`
	Name      string `json:"name,omitempty"`
	Avatar    []byte `json:"avatar,omitempty"`
	CellPhone string `json:"cell_phone,omitempty"`
	Salt      string `json:"salt,omitempty"`
	Password  string `json:"password,omitempty"`
	*BaseRecord
	*BaseRecorder
}

func NewUser(name string, password string, creator string, editor string) *User {
	baseRecord := NewBaseRecord()
	baseRecorder := NewBaseRecorder(creator, editor)
	return &User{
		State:        true,
		Name:         name,
		Salt:         util.RandomSalt(),
		Password:     password,
		BaseRecord:   baseRecord,
		BaseRecorder: baseRecorder,
	}
}

// Permission 权限
type Permission struct {
	Id           int               `json:"id,omitempty"`
	ParentId     int               `json:"parent_id,omitempty"`
	Code         string            `json:"code,omitempty"`
	Name         string            `json:"name,omitempty"`
	Introduction string            `json:"introduction,omitempty"`
	Category     PermssionCategory `json:"category,omitempty"`
	Uri          int               `json:"uri,omitempty"`
	*BaseRecord
	*BaseRecorder
}

// Role 角色
type Role struct {
	Id           int    `json:"id,omitempty"`
	ParentId     int    `json:"parent_id,omitempty"`
	Code         string `json:"code,omitempty"`
	Name         string `json:"name,omitempty"`
	Introduction string `json:"introduction,omitempty"`
	*BaseRecord
	*BaseRecorder
}

// UserRole 用户角色
type UserRole struct {
	Id     int `json:"id,omitempty"`
	UserId int `json:"user_id,omitempty"`
	RoleId int `json:"role_id,omitempty"`
	*BaseRecord
	*BaseRecorder
}

// RolePermission 角色权限
type RolePermission struct {
	Id           int `json:"id,omitempty"`
	RoleId       int `json:"role_id,omitempty"`
	PermissionId int `json:"permission_id,omitempty"`
	*BaseRecord
	*BaseRecorder
}