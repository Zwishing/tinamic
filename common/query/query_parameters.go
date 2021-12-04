package query

import "github.com/gofiber/fiber/v2"

type QueryParameters struct {
	Limit      int
	Properties []string
	Resolution int
	Buffer     int
}


func (qp *QueryParameters) GetQueryParameter(ctx *fiber.Ctx) error{
	err := ctx.QueryParser(qp)
	if err != nil {
		return err
	}
	return nil
}