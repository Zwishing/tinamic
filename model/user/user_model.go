package user

type AccountCategory int8

const (
	UserName AccountCategory = iota + 1
	Email
	Phone
)

func (ac AccountCategory) String() string {
	switch ac {
	case UserName:
		return "username"
	case Email:
		return "email"
	case Phone:
		return "phone"
	default:
		return "invalid acoount category"
	}
}
func AccountCategoryId(s string) AccountCategory {
	switch s {
	case "username":
		return UserName
	case "email":
		return Email
	case "phone":
		return Phone
	default:
		return UserName
	}
}

type PermssionCategory int8

const (
	Write PermssionCategory = iota + 1
	Read
)

func (pc PermssionCategory) String() string {
	switch pc {
	case Write:
		return "Write"
	case Read:
		return "Read"
	default:
		return "invalid permssion category"
	}
}

type RoleCode string

const (
	SuperAdmin   = RoleCode("super_admin")
	Admin        = RoleCode("admin")
	OrdinaryUser = RoleCode("ordinary_user")
	Guest        = RoleCode("guest")
)

func RoleId(rc *RoleCode) int {
	switch *rc {
	case SuperAdmin:
		return 0
	case Admin:
		return 1
	case OrdinaryUser:
		return 2
	case Guest:
		return 3
	default:
		return 3
	}
}

func RoleId2Code(roleId int) RoleCode {
	switch roleId {
	case 0:
		return SuperAdmin
	case 1:
		return Admin
	case 2:
		return OrdinaryUser
	case 3:
		return Guest
	default:
		return Guest
	}
}

// Account 账号
type Account struct {
	UserId      string          `json:"user_id"`
	UserAccount string          `json:"user_account"`
	Category    AccountCategory `json:"category"`
	*BaseRecorder
}
type AccountOption func(*Account)

func WithAccountUserId(userId string) AccountOption {
	return func(account *Account) {
		account.UserId = userId
	}
}

func WithLoginAccount(userAccount string) AccountOption {
	return func(account *Account) {
		account.UserAccount = userAccount
	}
}

func WithCategory(category AccountCategory) AccountOption {
	return func(account *Account) {
		account.Category = category
	}
}

func NewAccount(opts ...AccountOption) *Account {
	account := &Account{
		BaseRecorder: NewBaseRecorder("", ""),
	}
	for _, opt := range opts {
		opt(account)
	}
	return account
}

// User 用户
type User struct {
	UserId    string `json:"user_id"`
	State     bool   `json:"state"`
	Name      string `json:"name"`
	Avatar    []byte `json:"avatar"`
	CellPhone string `json:"cell_phone"`
	Salt      string `json:"salt"`
	Password  string `json:"password"`
	*BaseRecord
	*BaseRecorder
}

type UserOption func(*User)

func WithUserId(userId string) UserOption {
	return func(u *User) {
		u.UserId = userId
	}
}

func WithName(name string) UserOption {
	return func(u *User) {
		u.Name = name
	}
}

func WithAvatar(avatar []byte) UserOption {
	return func(u *User) {
		u.Avatar = avatar
	}
}

func WithCellPhone(cellPhone string) UserOption {
	return func(u *User) {
		u.CellPhone = cellPhone
	}
}

func WithSalt(salt string) UserOption {
	return func(u *User) {
		u.Salt = salt
	}
}

func WithPassword(password string) UserOption {
	return func(u *User) {
		u.Password = password
	}
}

func NewUser(opts ...UserOption) *User {
	user := &User{
		BaseRecord:   NewBaseRecord(),
		BaseRecorder: NewBaseRecorder("", ""),
	}

	for _, opt := range opts {
		opt(user)
	}
	return user
}

// Permission 权限
type Permission struct {
	Id           int               `json:"id"`
	ParentId     int               `json:"parent_id"`
	Code         string            `json:"code"`
	Name         string            `json:"name"`
	Introduction string            `json:"introduction"`
	Category     PermssionCategory `json:"category"`
	Uri          int               `json:"uri"`
	*BaseRecord
	*BaseRecorder
}

// Role 角色
type Role struct {
	Id           int      `json:"id"`
	ParentId     int      `json:"parent_id"`
	Code         RoleCode `json:"code"`
	Name         string   `json:"name"`
	Introduction string   `json:"introduction"`
	*BaseRecord
	*BaseRecorder
}

// UserRole 用户角色
type UserRole struct {
	UserId string `json:"user_id"`
	RoleId int    `json:"role_id"`
}

func GetRoleString() string {
	return "role"
}

// RolePermission 角色权限
type RolePermission struct {
	Id           int `json:"id"`
	RoleId       int `json:"role_id"`
	PermissionId int `json:"permission_id"`
	*BaseRecord
	*BaseRecorder
}
