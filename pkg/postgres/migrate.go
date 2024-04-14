package postgres

import (
	"backend-trainee-assignment-2024/migrations"
	"context"
	"fmt"
	"github.com/uptrace/bun"
	bunmigrate "github.com/uptrace/bun/migrate"
)

func migrate(db *bun.DB) error {
	migrator := bunmigrate.NewMigrator(db, migrations.Migrations)
	ctx := context.Background()

	err := migrator.Init(ctx)
	if err != nil {
		return err
	}

	if err = migrator.Lock(ctx); err != nil {
		return err
	}
	defer migrator.Unlock(ctx)

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return err
	}
	if group.IsZero() {
		fmt.Printf("no new migrations\n")
		return nil
	}

	fmt.Printf("migrated to %s\n", group)
	return nil
}
