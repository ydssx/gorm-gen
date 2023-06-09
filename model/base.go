package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseModel struct {
	ID        uint  `json:"id" gorm:"primaryKey"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	DeletedAt int64 `json:"deleted_at"`
}

type DB struct {
	db *gorm.DB
}

func cloneDB(db *gorm.DB) *gorm.DB {
	tx := &gorm.DB{Config: db.Config}

	stmt := db.Statement
	newStmt := &gorm.Statement{
		TableExpr:            stmt.TableExpr,
		Table:                stmt.Table,
		Model:                stmt.Model,
		Unscoped:             stmt.Unscoped,
		Dest:                 stmt.Dest,
		ReflectValue:         stmt.ReflectValue,
		Clauses:              map[string]clause.Clause{},
		Distinct:             stmt.Distinct,
		Selects:              stmt.Selects,
		Omits:                stmt.Omits,
		Preloads:             map[string][]interface{}{},
		ConnPool:             stmt.ConnPool,
		Schema:               stmt.Schema,
		Context:              stmt.Context,
		RaiseErrorOnNotFound: stmt.RaiseErrorOnNotFound,
		SkipHooks:            stmt.SkipHooks,
	}

	if stmt.SQL.Len() > 0 {
		newStmt.SQL.WriteString(stmt.SQL.String())
		newStmt.Vars = make([]interface{}, 0, len(stmt.Vars))
		newStmt.Vars = append(newStmt.Vars, stmt.Vars...)
	}

	for k, c := range stmt.Clauses {
		newStmt.Clauses[k] = c
	}

	for k, p := range stmt.Preloads {
		newStmt.Preloads[k] = p
	}

	// if len(stmt.Joins) > 0 {
	// 	newStmt.Joins = make([]join, len(stmt.Joins))
	// 	copy(newStmt.Joins, stmt.Joins)
	// }

	// if len(stmt.scopes) > 0 {
	// 	newStmt.scopes = make([]func(*DB) *DB, len(stmt.scopes))
	// 	copy(newStmt.scopes, stmt.scopes)
	// }

	stmt.Settings.Range(func(k, v interface{}) bool {
		newStmt.Settings.Store(k, v)
		return true
	})
	tx.Statement = newStmt
	tx.Statement.DB = tx
	return tx
}
