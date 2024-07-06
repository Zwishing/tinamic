// wire.go

//go:build wireinject
// +build wireinject

// wire/inject.go

package wire

import (
	"github.com/google/wire"
	gonsq "github.com/nsqio/go-nsq"
	"tinamic/handler"
	"tinamic/pkg/nsq"
	"tinamic/repository"
	"tinamic/service"
)

type UserComponents struct {
	UserHandler *handler.UserHandler
}

func InitializeUserService() *UserComponents {
	wire.Build(
		repository.NewUserRepository,
		service.NewUserService,
		handler.NewUserHandler,
		wire.Struct(new(UserComponents), "UserHandler"),
	)
	return &UserComponents{}
}

type DataSourceComponents struct {
	Consumer          *nsq.NSQConsumer
	DataSourceHandler *handler.DataSourceHandler
}

func InitializeDataSourceService() (*DataSourceComponents, func(), error) {
	wire.Build(
		repository.NewDataSourceRepository,
		service.NewDataSourceService,
		handler.NewDataSourceHandler,
		nsq.NewNSQConfig,
		nsq.NewNSQConsumer,
		//wire.Struct(new(*handler.NewDataSourceHandler), "*"),
		wire.Bind(new(gonsq.Handler), new(*handler.DataSourceHandler)),
		wire.Struct(new(DataSourceComponents), "Consumer", "DataSourceHandler"),
	)
	return &DataSourceComponents{}, nil, nil
}
