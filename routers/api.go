package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"

	"tinamic/app/controllers"
)

func RegisterAPI(api fiber.Router,db *pgxpool.Pool) {
	registerLayer(api,db)
	//registerRoles(api)
	//registerUsers(api)
}

func registerLayer(api fiber.Router,db *pgxpool.Pool) {
	layer := api.Group("/layer")

	layer.Get("/table-layers", controllers.GetAllTableLayers(db))
	layer.Get("/tile/:name/:z/:x/:y.pbf",controllers.GetTableLayerTile(db))
	//layer.Post("/")
	//layer.Put("/:id")
	//layer.Delete("/:id")
}

func registerRoles(api fiber.Router) {
	roles := api.Group("/roles")

	roles.Get("/", )
	roles.Get("/:id")
	roles.Post("/")
	roles.Put("/:id")
	roles.Delete("/:id")
}

func registerUsers(api fiber.Router) {
	users := api.Group("/users")

	users.Get("/")
	users.Get("/:id")
	users.Post("/")
	users.Put("/:id")
	users.Delete("/:id")
}