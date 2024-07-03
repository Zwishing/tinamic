package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"tinamic/pkg/pgxadapter"
)

const (
	SCHEMA = "user_info"
	TABLE  = "user_permission"
)

func NewAdapter(pool *pgxpool.Pool) (*pgxadapter.Adapter, error) {
	a, err := pgxadapter.NewAdapter("",
		pgxadapter.WithSchema(SCHEMA),
		pgxadapter.WithTableName(TABLE),
		pgxadapter.WithTimeout(1*time.Minute),
		pgxadapter.WithConnectionPool(pool))
	return a, err
}

func NewEnforcer(pool *pgxpool.Pool) (*casbin.Enforcer, error) {
	a, err := NewAdapter(pool)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer("D:\\Code\\go-web\\tinamic\\conf\\model.conf", a)
	if err != nil {
		return nil, err
	}
	return e, err
}
