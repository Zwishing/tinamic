package query

import "github.com/gofiber/fiber/v2"

type QueryParameters struct {
	Limit      int  //-1 无限制
	Properties []string
	Resolution int
	Buffer     int
}

func NewQueryParameters() *QueryParameters {
	return &QueryParameters{
		Limit: -1,
		Resolution: 4096,
		Buffer: 256,
	}
}

func (qp *QueryParameters) GetQueryParameter(ctx *fiber.Ctx) error{
	err := ctx.QueryParser(qp)
	if err != nil {
		return err
	}
	return nil
}