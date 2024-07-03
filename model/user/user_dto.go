package user

type SignInDTO struct {
	UserAccount string `json:"userAccount" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Category    string `json:"category" validate:"required,oneof= username email phone"`
}

type RegisterDTO struct {
	UserAccount string `json:"userAccount" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Category    string `json:"category" validate:"required,oneof= username email phone"`
	Role        string `json:"role" validate:"required,oneof=admin ordinary_user guest"`
}

type UserDTO struct {
	UserId      string   `json:"userId"`
	UserAccount string   `json:"userAccount"`
	Category    string   `json:"category"`
	Role        RoleCode `json:"role"`
	Name        string   `json:"name"`
	Avatar      []byte   `json:"avatar"`
	CellPhone   string   `json:"cellPhone"`
}
