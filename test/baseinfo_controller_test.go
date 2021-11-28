package test

import (
	"github.com/gofiber/fiber/v2"
	"testing"
	"tinamic/app/controllers"
)

func TestGetLayerBaseInfo(t *testing.T) {
	c := fiber.Ctx{}
	controllers.GetLayerBaseInfo(&c)
}

