package gomigrator

import (
	"github.com/korfairo/migratory/internal/require"
	"testing"
)

func TestFindMissingMigrations(t *testing.T) {
	type args struct {
		migrations Migrations
		results    []MigrationResult
	}
	tests := map[string]struct {
		args        args
		wantMissing Migrations
		wantDirty   bool
	}{
		"all migrations missing": {
			args: args{
				migrations: allMigrations,
				results:    []MigrationResult{},
			},
			wantMissing: allMigrations,
			wantDirty:   false,
		},
		"no migrations missing": {
			args: args{
				migrations: allMigrations,
				results:    allResults,
			},
			wantMissing: nil,
			wantDirty:   false,
		},
		"half migrations missing": {
			args: args{
				migrations: allMigrations,
				results: []MigrationResult{
					{Id: 1},
					{Id: 2},
					{Id: 3},
					{Id: 4},
					{Id: 5},
				},
			},
			wantMissing: Migrations{
				Migration{id: 6},
				Migration{id: 7},
				Migration{id: 8},
				Migration{id: 9},
				Migration{id: 10},
			},
			wantDirty: false,
		},
		"last migration missing": {
			args: args{
				migrations: allMigrations,
				results: []MigrationResult{
					{Id: 1},
					{Id: 2},
					{Id: 3},
					{Id: 4},
					{Id: 5},
					{Id: 6},
					{Id: 7},
					{Id: 8},
					{Id: 9},
				},
			},
			wantMissing: Migrations{
				Migration{id: 10},
			},
			wantDirty: false,
		},
		"one dirty migration": {
			args: args{
				migrations: allMigrations,
				results: []MigrationResult{
					{Id: 2},
					{Id: 3},
					{Id: 4},
					{Id: 5},
					{Id: 6},
					{Id: 7},
					{Id: 8},
					{Id: 9},
					{Id: 10},
				},
			},
			wantMissing: Migrations{
				Migration{id: 1},
			},
			wantDirty: true,
		},
		"dirty migrations": {
			args: args{
				migrations: allMigrations,
				results: []MigrationResult{
					{Id: 2},
					{Id: 3},
					{Id: 5},
					{Id: 8},
				},
			},
			wantMissing: Migrations{
				Migration{id: 1},
				Migration{id: 4},
				Migration{id: 6},
				Migration{id: 7},
				Migration{id: 9},
				Migration{id: 10},
			},
			wantDirty: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotMissing, gotDirty := findMissingMigrations(test.args.migrations, test.args.results)
			require.Equal(t, gotMissing, test.wantMissing, "findMissingMigrations(...) missing migrations")
			require.Bool(t, gotDirty, test.wantDirty, "findMissingMigrations(...) dirty migrations found")
		})
	}
}

var allMigrations = Migrations{
	Migration{id: 1},
	Migration{id: 2},
	Migration{id: 3},
	Migration{id: 4},
	Migration{id: 5},
	Migration{id: 6},
	Migration{id: 7},
	Migration{id: 8},
	Migration{id: 9},
	Migration{id: 10},
}

var allResults = []MigrationResult{
	{Id: 1},
	{Id: 2},
	{Id: 3},
	{Id: 4},
	{Id: 5},
	{Id: 6},
	{Id: 7},
	{Id: 8},
	{Id: 9},
	{Id: 10},
}
