package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	_ "github.com/lib/pq"
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

// type DataServiceCredentials struct {
// 	Host         string
// 	DatabaseName string
// 	Username     string
// 	Password     string
// 	Port         string
// }

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

func connectDatabase() *xorm.Engine {

	// Assumption: We are in a CF like environment. https://docs.cloudfoundry.org/devguide/deploy-apps/environment-variable.html#VCAP-SERVICES
	vcapServices := os.Getenv("VCAP_SERVICES")

	var credentials interface{} // := DataServiceCredentials{}
	err := json.Unmarshal([]byte(vcapServices), &credentials)

	if err != nil {
		app.Logger().Fatalf("Couldn't read database credentials from VCAP_SERVICES ENV variable. Are we in a Cloud Foundry? %v", err)
	}

	app.Logger().Debug("Cred: ", credentials)

	orm, err := xorm.NewEngine("postgres", "host=localhost user=jfischer dbname=gostock sslmode=disable")
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
		app.Logger().Fatalf("Cannot create db schema: ", err)
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
