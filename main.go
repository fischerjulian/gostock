package main

import (
	"time"

	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Stock represents a single stock title.
type Stock struct {
	ID   int64  // auto-increment by-default by xorm
	Name string `json:"name" validate:"required" xorm:"varchar(200)"`

	// Value of the stock in EUR cent (no decimals)
	Value     uint32    `json:"value" validate:"required"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	// recover panics
	app.Use(recover.New())
	app.Use(logger.New())

	// Establish DB connection
	// orm, err := xorm.NewEngine("sqlite3", "./test.db")
	orm, err := xorm.NewEngine("postgres", "user=jfischer dbname=gostock sslmode=disable")

	if err != nil {
		app.Logger().Fatalf("orm failed to initialized: %v", err)
	}

	// Create schema
	err = orm.Sync2(new(Stock))

	// Close ORM later
	iris.RegisterOnInterrupt(func() {
		orm.Close()
	})

	stock := Stock{Name: "Apple", Value: 17780}

	// Seed
	doesExist, err := orm.Exist(&stock)

	if !doesExist {
		orm.Insert(&stock, &Stock{Name: "Alphabet Inc Class A", Value: 102140})
	}

	// Method: GET
	// Resource http://localhost:8080
	app.Handle("GET", "/stocks", func(ctx iris.Context) {

		stocks := make(map[int64]Stock)
		err := orm.Find(&stocks)

		app.Logger().Debug("Stocks: ", stocks)

		if err != nil {
			app.Logger().Fatalf("orm failed to load stocks: %v", err)
		}

		ctx.JSON(iris.Map{"stocks": stocks})
	})

	// same as app.Handle("GET", "/ping", ..)
	// Method: GET
	// Resource: http://localhost:8080/ping
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString("pong")
	})

	// Method: GET
	// Resource: http://localhost:8080/hello
	app.Get("/hello", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello Iris!"})
	})

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
