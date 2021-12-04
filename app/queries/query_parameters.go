package queries

import (
	"github.com/gofiber/fiber/v2"
)

type queryParameters struct {
	Limit      int
	Properties []string
	Resolution int
	Buffer     int
}


func (qp *queryParameters)getQueryParameter(ctx *fiber.Ctx) error{
	err := ctx.QueryParser(qp)
	if err != nil {
		return err
	}
	return nil
}
