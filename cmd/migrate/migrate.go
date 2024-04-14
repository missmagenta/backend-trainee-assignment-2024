package main

import (
	"context"
	"fmt"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"strings"
)

type Migrator struct {
	bunMigrator *migrate.Migrator
}

func newMigrator(bunMigrator *migrate.Migrator) *Migrator {
	return &Migrator{bunMigrator: bunMigrator}
}

func (m *Migrator) migrate(c *cli.Context) error {
	ctx := c.Context
	if err := m.lock(ctx); err != nil {
		return err
	}
	defer m.unlock(ctx)

	group, err := m.bunMigrator.Migrate(ctx)
	if err != nil {
		return err
	}
	if group.IsZero() {
		fmt.Println("No new migrations to run (database is up to date)")
		return nil
	}
	fmt.Printf("Migrated to %s\n", group)
	return nil
}

func (m *Migrator) lock(ctx context.Context) error {
	return m.bunMigrator.Lock(ctx)
}

func (m *Migrator) unlock(ctx context.Context) {
	_ = m.bunMigrator.Unlock(ctx)
}

func (m *Migrator) rollback(c *cli.Context) error {
	ctx := c.Context
	if err := m.lock(ctx); err != nil {
		return err
	}
	defer m.unlock(ctx)

	group, err := m.bunMigrator.Rollback(ctx)
	if err != nil {
		return err
	}
	if group.IsZero() {
		fmt.Println("No groups to roll back")
		return nil
	}
	fmt.Printf("Rolled back %s\n", group)
	return nil
}

type createMigrationFunc func(ctx context.Context, name string, opts ...interface{}) (*migrate.MigrationFile, error)

func (m *Migrator) createMigration(c *cli.Context, createFunc createMigrationFunc) error {
	name := strings.Join(c.Args().Slice(), "_")
	mf, err := createFunc(c.Context, name)
	if err != nil {
		return err
	}
	fmt.Printf("Created migration %s (%s)\n", mf.Name, mf.Path)
	return nil
}

func (m *Migrator) createGoMigration(c *cli.Context) error {
	createFunc := func(ctx context.Context, name string, opts ...interface{}) (*migrate.MigrationFile, error) {
		var options []migrate.GoMigrationOption
		for _, opt := range opts {
			if o, ok := opt.(migrate.GoMigrationOption); ok {
				options = append(options, o)
			} else {
				return nil, fmt.Errorf("invalid migration option type")
			}
		}
		return m.bunMigrator.CreateGoMigration(ctx, name, options...)
	}
	return m.createMigration(c, createFunc)
}

func (m *Migrator) createSQLMigration(c *cli.Context) error {
	createFunc := func(ctx context.Context, name string, opts ...interface{}) (*migrate.MigrationFile, error) {
		files, err := m.bunMigrator.CreateSQLMigrations(ctx, name)
		if err != nil {
			return nil, err
		}
		if len(files) == 0 {
			return nil, fmt.Errorf("no SQL migration files created")
		}
		return files[0], nil
	}
	return m.createMigration(c, createFunc)
}

func (m *Migrator) status(c *cli.Context) error {
	ms, err := m.bunMigrator.MigrationsWithStatus(c.Context)
	if err != nil {
		return err
	}
	fmt.Printf("Migrations: %s\n", ms)
	fmt.Printf("Unapplied migrations: %s\n", ms.Unapplied())
	fmt.Printf("Last migration group: %s\n", ms.LastGroup())
	return nil
}

func (m *Migrator) markApplied(c *cli.Context) error {
	group, err := m.bunMigrator.Migrate(c.Context, migrate.WithNopMigration())
	if err != nil {
		return err
	}
	if group.IsZero() {
		fmt.Println("No new migrations to mark as applied")
		return nil
	}
	fmt.Printf("Marked as applied %s\n", group)
	return nil
}
