// wire.go

//go:build wireinject
// +build wireinject

//package wire

//import (
//	"github.com/google/wire"
//	"tinamic/handler"
//	"tinamic/repository"
//	"tinamic/service"
//)
//
//func InitializeUserRepository() repository.UserRepository {
//	wire.Build(
//		repository.NewUserRepository,
//	)
//	return nil
//}
//
//func InitializeUserService() service.UserService {
//	wire.Build(
//		service.NewUserService,
//		InitializeUserRepository,
//	)
//	return nil
//}
//
//func InitializeUserHandler() *handler.UserHandler {
//	wire.Build(
//		handler.NewUserHandler,
//		InitializeUserService,
//	)
//	return nil
//}

// wire/inject.go

package wire

import (
	"github.com/google/wire"
	"tinamic/handler"
	"tinamic/repository"
	"tinamic/service"
)

func InitializeUserService() *handler.UserHandler {
	wire.Build(
		repository.NewUserRepository,
		service.NewUserService,
		handler.NewUserHandler,
	)
	return &handler.UserHandler{}
}

func InitializeDataSourceService() *handler.DataSourceHandler {
	wire.Build(
		repository.NewDataSourceRepository,
		service.NewDataSourceService,
		handler.NewDataSourceHandler,
	)
	return &handler.DataSourceHandler{}
}
