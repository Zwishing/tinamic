package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	. "tinamic/model/user"
)

func TestUserRepositoryImpl_QueryAccount(t *testing.T) {
	user := NewAccount(WithLoginAccount("admin"))
	repo := NewUserRepository()
	repo.QueryAccount(user)
	assert.Equal(t, user.UserId, "cc9de34f-731f-4c45-9ecd-ef84a07b7427")
}

func TestUserRepositoryImpl_AddUser(t *testing.T) {
	userId := uuid.New()
	//people := &BaseRecorder{
	//	Editor:  "zhang",
	//	Creator: "zhang",
	//}
	account := NewAccount(
		WithLoginAccount("zhangwe11144"),
		WithAccountUserId(userId.String()),
		WithCategory(UserName),
	)
	//account.BaseRecorder = people
	user := NewUser(
		WithUserId(userId.String()),
		WithSalt("3444"),
		WithPassword("1212121212"),
		WithName("zhangweixin345"),
	)
	guest := Guest
	role := &UserRole{
		UserId: userId.String(),
		RoleId: RoleId(&guest),
	}
	repo := NewUserRepository()
	err := repo.AddUser(account, user, role)
	if err != nil {
		fmt.Println(err)
	}

}

func TestUserRepositoryImpl_QueryUser(t *testing.T) {
	account := NewAccount(WithLoginAccount("admin"), WithCategory(UserName))
	user := NewUser()
	role := &UserRole{}
	repo := NewUserRepository()
	err := repo.QueryUser(account, user, role)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(account)
	fmt.Println(user)
	fmt.Println(role)
}

func TestUserRepositoryImpl_AddPolicy(t *testing.T) {
	repo := NewUserRepository()
	repo.AddPolicy(Admin, "/v1/user/profile", "GET")
}
