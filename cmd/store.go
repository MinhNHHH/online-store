package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/MinhNHHH/online-store/pkg/cfgs"
	"github.com/MinhNHHH/online-store/pkg/databases/repositories/dbrepo"
	"github.com/MinhNHHH/online-store/pkg/store"
)

func initializeApp(cfgs cfgs.Configs) store.OnlineStore {
	app := store.OnlineStore{}

	sqlConn, err := app.ConnectDB(cfgs.DB_CONNECTION_URI)
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.DBRepo{SqlConn: sqlConn}
	app.Cfgs = cfgs
	return app
}

func startServer(app store.OnlineStore) {
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", app.Routes())
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: step = 0 means run all migrations
// step = n means run only the next n migrations
// step = -n means run only the previous n migrations
func handleMigrate(app store.OnlineStore, args []string) {
	step := 0
	if len(args) == 1 {
		var err error
		step, err = strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
	}
	app.Migrate(step, app.Cfgs.DB_CONNECTION_URI)
}

func handleCreate(app store.OnlineStore, args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a name for the migration")
		os.Exit(1)
	}
	app.GenerateMigration(args[0])
}

func handleCommand(app store.OnlineStore, command string, args []string) {
	switch command {
	case "migrate":
		handleMigrate(app, args)
	case "create":
		handleCreate(app, args)
	default:
		startServer(app)
	}
}

func main() {
	cfgs := cfgs.LoadConfigs()
	app := initializeApp(cfgs)

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/store.go [command]")
		fmt.Println("Commands:")
		fmt.Println("  migrate [steps] - Run migrations (optional number of steps)")
		fmt.Println("  create [name]   - Create new migration files")
		fmt.Println("  start           - Run server")
		os.Exit(1)
	}
	command := os.Args[1]
	args := os.Args[2:]

	handleCommand(app, command, args)
}
