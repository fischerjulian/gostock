package main

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	_ "github.com/lib/pq"

	"github.com/cloudfoundry-community/go-cfenv"
)

// Stock represents a single stock title.
type Stock struct {
	ID   int64  // auto-increment by-default by xorm
	Name string `form:"name" json:"name" validate:"required" xorm:"varchar(200)"`

	// Value of the stock in EUR cent (no decimals)
	Value     uint32    `form:"value" json:"value" validate:"required"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

var app *iris.Application
var orm *xorm.Engine

func main() {
	app = iris.New()
	app.Logger().SetLevel("debug")

	// recover panics
	app.Use(recover.New())
	app.Use(logger.New())

	// Establish DB connection
	orm = connectDatabase()

	createDbSchema()

	seedData()

	// Method: GET
	// Resource http://localhost:8080
	app.Handle("GET", "/stocks", listStocks)

	app.Post("/stock", postStock)

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

func databaseConnectionString() string {
	appEnv, _ := cfenv.Current()
	fmt.Println("Services:", appEnv.Services)

	var a9spg cfenv.Service

	if appEnv.Services["a9s-postgresql10"] != nil {
		a9spg = appEnv.Services["a9s-postgresql10"][0]
	} else if appEnv.Services["elephantsql"] != nil {
		a9spg = appEnv.Services["elephantsql"][0]
	} else {
		app.Logger().Fatalf("Neither a9s PostgreSQL nor ElephantSQL has been found. I cannot do without. Please bind a service instance to the app.")
	}

	dbConnectionString := fmt.Sprintf("%s?sslmode=disable", a9spg.Credentials["uri"])

	// Note that this demo app logs the connection string including passwords to facilitate debugging during deployment.
	app.Logger().Debug("DB Connection String: ", dbConnectionString)

	return dbConnectionString
}

func connectDatabase() *xorm.Engine {

	connString := databaseConnectionString()

	orm, err := xorm.NewEngine("postgres", connString)
	if err != nil {
		app.Logger().Fatalf("ORM failed to initialized: %v", err)
	}

	// Close ORM later
	iris.RegisterOnInterrupt(func() {
		orm.Close()
	})
	return orm
}

func createDbSchema() {
	// Create schema
	err := orm.Sync2(new(Stock))
	if err != nil {
		app.Logger().Fatalf("Cannot create db schema. Maybe I was unable to connect to the DB: ", err)
	}
}

func seedData() {

	// Seed
	count, err := orm.Count(new(Stock))

	if err != nil {
		app.Logger().Fatalf("Cannot retrieve data from db during seed check: ", err)
	}

	if count == 0 {
		orm.Insert(&Stock{Name: "Apple", Value: 17780}, &Stock{Name: "Alphabet Inc Class A", Value: 102140})
	}
}

func listStocks(ctx iris.Context) {
	stocks := make(map[int64]Stock)
	err := orm.Find(&stocks)

	app.Logger().Debug("Stocks: ", stocks)

	if err != nil {
		app.Logger().Fatalf("orm failed to load stocks: %v", err)
	}

	ctx.JSON(iris.Map{"stocks": stocks})
}

func postStock(ctx iris.Context) {
	stock := Stock{}
	err := ctx.ReadForm(&stock)

	if err != nil {
		app.Logger().Errorf("Couldn't read form input in postStock: %v", err)
	}

	app.Logger().Debug("Submitted Stock:", stock)

	orm.Insert(&stock)

	// read stock from params
	ctx.JSON(iris.Map{"success": true})
}
