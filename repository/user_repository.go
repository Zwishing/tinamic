package repository

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/rs/zerolog/log"
	"time"
	"tinamic/model/user"
	icasbin "tinamic/pkg/casbin"
	"tinamic/pkg/pg"
	"tinamic/pkg/redis"
)

type UserRepository interface {
	QueryUserByUserId(user *user.User) error
	QueryUserRole(userRole *user.UserRole) error
	QueryAccount(account *user.Account) error
	QueryUser(account *user.Account, user *user.User, role *user.UserRole) error
	AddUser(account *user.Account, user *user.User, role *user.UserRole) error
	SaveToRedis(key string, val any, exp time.Duration) error
	GetFromRedis(key string) (string, error)
	AddPolicy(sub user.RoleCode, obj string, act string) bool
}

type UserRepositoryImpl struct {
	*pg.PGPool
	dialect  goqu.DialectWrapper
	redis    *redis.RedisStorage
	enforcer *casbin.Enforcer
}

func NewUserRepository() UserRepository {
	e, _ := icasbin.NewEnforcer(GetDbPoolInstance().Pool)
	return &UserRepositoryImpl{
		GetDbPoolInstance(),
		goqu.Dialect("postgres"),
		GetRedisInstance(),
		e,
	}
}

func (db *UserRepositoryImpl) AddPolicy(sub user.RoleCode, obj string, act string) bool {
	policy, _ := db.enforcer.AddPolicy(string(sub), obj, act)
	return policy
}

func (db *UserRepositoryImpl) SaveToRedis(key string, val any, exp time.Duration) error {
	err := db.redis.SetMap(key, val, exp)
	if err != nil {
		return err
	}
	return nil
}

func (db *UserRepositoryImpl) GetFromRedis(key string) (string, error) {
	get, err := db.redis.Get(key)
	if err != nil {
		return "", err
	}
	return get, nil
}

func (db *UserRepositoryImpl) QueryUser(account *user.Account, user *user.User, role *user.UserRole) error {
	// 构建查询
	query := db.dialect.From(goqu.I("user_info.user").As("u")).
		Select(
			goqu.I("a.user_id").As("account_user_id"),
			goqu.I("u.user_id"),
			goqu.I("u.state"),
			goqu.I("u.name"),
			goqu.I("u.avatar"),
			goqu.I("u.cell_phone"),
			goqu.I("u.password"),
			goqu.I("ur.user_id").As("role_user_id"),
			goqu.I("ur.role_id"),
		).
		Join(
			goqu.I("user_info.account").As("a"),
			goqu.On(goqu.Ex{
				"u.user_id": goqu.I("a.user_id"),
			}),
		).
		Join(
			goqu.I("user_info.user_role").As("ur"),
			goqu.On(goqu.Ex{
				"ur.user_id": goqu.I("a.user_id"),
			}),
		).
		Where(goqu.Ex{
			"a.deleted":      false,
			"u.state":        true,
			"a.user_account": account.UserAccount,
			"a.category":     account.Category,
		})

	// 生成 SQL 查询语句
	sql, _, err := query.ToSQL()
	
	if err != nil {
		return err
	}

	err = db.QueryRow(context.Background(), sql).Scan(
		&account.UserId,
		&user.UserId, &user.State, &user.Name, &user.Avatar, &user.CellPhone, &user.Password,
		&role.UserId, &role.RoleId)

	if err != nil {
		return err
	}
	return nil
}

func (db *UserRepositoryImpl) AddUser(account *user.Account, u *user.User, role *user.UserRole) error {
	// 需要增加三张表 用户表、账号表和用户权限表
	addAccountSQL, _, err := db.dialect.Insert("user_info.account").
		Cols("user_id", "user_account", "category", "creator", "editor").
		Vals(goqu.Vals{account.UserId, account.UserAccount, account.Category, account.Creator, account.Editor}).
		ToSQL()
	if err != nil {
		return err
	}

	addUserSQL, _, err := db.dialect.Insert("user_info.user").
		Cols("user_id", "name", "avatar", "cell_phone", "salt", "password", "creator", "editor").
		Vals(goqu.Vals{u.UserId, u.Name, u.Avatar, u.CellPhone, u.Salt, u.Password, u.Creator, u.Editor}).
		ToSQL()
	if err != nil {
		return err
	}

	addUserRoleSQL, _, err := db.dialect.Insert("user_info.user_role").
		Cols("user_id", "role_id").
		Vals(goqu.Vals{role.UserId, role.RoleId}).
		ToSQL()
	if err != nil {
		return err
	}

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction error: %v", err)
	}
	defer func() {
		if err != nil {
			// 回滚事务
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Error().Msgf("rollback transaction error: %v", rbErr)
			}
			return
		}
		// 提交事务
		if err = tx.Commit(ctx); err != nil {
			log.Error().Msgf("commit transaction error: %v", err)
		}
	}()

	// 执行插入操作
	_, err = tx.Exec(ctx, addAccountSQL)
	if err != nil {
		return fmt.Errorf("insert account error: %v", err)
	}

	_, err = tx.Exec(ctx, addUserSQL)
	if err != nil {
		return fmt.Errorf("insert user error: %v", err)
	}

	_, err = tx.Exec(ctx, addUserRoleSQL)
	if err != nil {
		return fmt.Errorf("insert user role error: %v", err)
	}
	return nil
}

func (db *UserRepositoryImpl) QueryAccount(account *user.Account) error {
	toSQL, _, err := db.dialect.From("user_info.account").
		Select("user_id", "category").
		Where(goqu.C("user_account").
			Eq(account.UserAccount)).
		ToSQL()
	if err != nil {
		return err
	}
	err = db.QueryRow(context.Background(), toSQL).Scan(&account.UserId, &account.Category)
	if err != nil {
		return err
	}
	return nil
}

func (db *UserRepositoryImpl) QueryUserRole(userRole *user.UserRole) error {
	// 通过用户id查询用户角色
	sql := `SELECT role_id FROM user_info.user_role WHERE user_id=$1`
	err := db.QueryRow(context.Background(), sql, userRole.UserId).Scan(&userRole.RoleId)
	if err != nil {
		return err
	}
	return nil
}

func (db *UserRepositoryImpl) QueryUserByUserId(user *user.User) error {
	sql := `SELECT state,name,avatar,cell_phone FROM user_info.user WHERE name=$1`
	err := db.QueryRow(context.Background(), sql, user.UserId).Scan(
		&user.State,
		&user.Name,
		&user.Avatar,
		&user.CellPhone,
	)
	if err != nil {
		return err
	}
	return nil
}
