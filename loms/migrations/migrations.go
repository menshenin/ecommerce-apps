// Package migrations миграции
package migrations

import "embed"

// Migrations Миграции
//
//go:embed *.sql
var Migrations embed.FS
