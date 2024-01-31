package command

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/korfairo/migratory/internal/gomigrator"
	"github.com/korfairo/migratory/internal/sqlmigration"
	"os"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up [--dir <path>] [-d <db-string>] [-s <schema>] [-t <table>]",
	Short: "Up all unapplied migrations",
	Long: `The "up" command applies all unapplied migrations. 
It searches for SQL migration files in the directory that is passed as an argument, 
checks the version of the database - which is the ID of the last applied migration 
in your migrations database table - and applies any missing migrations one-by-one. 
If there are migrations in your directory with ID numbers less than the database version, 
they are considered "dirty migrations". By default, the command will return an error 
in this case, but you can ignore it with the "–force" flag to apply all missing migrations. 
Additionally, the command will create the migrations table if it does not already exist.`,
	Example: `up -c /etc/config.yml
up -d postgresql://role:password@127.0.0.1:5432/database --dir example/migrations/
up -d postgresql://role:password@127.0.0.1:5432/database --dir migrations/ -t my_migrations_table`,
	Run: func(cmd *cobra.Command, args []string) {
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			fmt.Println("failed to get bool --force flag")
			os.Exit(1)
		}

		appliedCount, err := up(config.Directory, config.Schema, config.Table, force)
		if err != nil {
			fmt.Printf("%d migration(s) applied, an error occured: %s\n", appliedCount, err)
			return
		}

		if appliedCount == 0 {
			fmt.Println("the database is up to date, nothing to migrate")
			return
		}

		fmt.Printf("success: applied %d migration(s)\n", appliedCount)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	upCmd.Flags().BoolP("force", "f", false, `ignore "dirty migrations" error`)
}

func up(dir, schema, table string, force bool) (int, error) {
	ctx := context.Background()
	migrator, err := gomigrator.New("postgres", schema, table)
	if err != nil {
		return 0, fmt.Errorf("failed to create migrator: %s", err)
	}

	db, err := sql.Open("postgres", config.DBString)
	if err != nil {
		return 0, fmt.Errorf("could not open database: %s", err)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			fmt.Println("failed to close database connection")
		}
	}()

	migrations, err := sqlmigration.SeekMigrations(dir, nil)
	if err != nil {
		return 0, fmt.Errorf("could not find migrations in directory %s: %s", dir, err)
	}

	appliedCount, err := migrator.Up(ctx, migrations, db, force)
	if err != nil {
		return appliedCount, fmt.Errorf("failed to execute migration: %s", err)
	}

	return appliedCount, nil
}
