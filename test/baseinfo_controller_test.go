/*
 * @Author: your name
 * @Date: 2021-11-28 10:17:19
 * @LastEditTime: 2021-11-28 11:23:34
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \tinamic\test\baseinfo_controller_test.go
 */
package test

import (
	"testing"
	"tinamic/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func TestGetLyrBaseInfo(t *testing.T) {
	c := fiber.Ctx{}
	controllers.GetLyrBaseInfo(&c)
}
