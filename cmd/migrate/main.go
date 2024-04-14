package main

import (
	"backend-trainee-assignment-2024/config"
	"backend-trainee-assignment-2024/pkg/postgres"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := postgres.OpenDB(cfg.PG.URL)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	app := &cli.App{
		Name:  "Migration CLI",
		Usage: "CLI tool for managing database migrations",
		Commands: []*cli.Command{
			{
				Name:    "migrate",
				Aliases: []string{"m"},
				Usage:   "Run pending migrations",
				Action:  NewMigratorAction((*Migrator).migrate),
			},
			{
				Name:    "rollback",
				Aliases: []string{"r"},
				Usage:   "Rollback the last group of migrations",
				Action:  NewMigratorAction((*Migrator).rollback),
			},
			{
				Name:    "create-go",
				Aliases: []string{"cg"},
				Usage:   "Create a new Go migration file",
				Action:  NewMigratorAction((*Migrator).createGoMigration),
			},
			{
				Name:    "create-sql",
				Aliases: []string{"cs"},
				Usage:   "Create new SQL migration files",
				Action:  NewMigratorAction((*Migrator).createSQLMigration),
			},
			{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "Check the status of migrations",
				Action:  NewMigratorAction((*Migrator).status),
			},
			{
				Name:    "mark-applied",
				Aliases: []string{"ma"},
				Usage:   "Mark pending migrations as applied",
				Action:  NewMigratorAction((*Migrator).markApplied),
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Println(err)
	}
}

func NewMigratorAction(method func(*Migrator, *cli.Context) error) cli.ActionFunc {
	return func(c *cli.Context) error {
		bunMigrator := &migrate.Migrator{}
		migrator := newMigrator(bunMigrator)
		return method(migrator, c)
	}
}
