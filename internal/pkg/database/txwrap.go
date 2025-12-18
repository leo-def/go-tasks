package database

import "gorm.io/gorm"

type Transaction interface {
	Connection
	Commit() error
	Rollback()
}

type Connection interface {
	BeginTransaction() Transaction
}

type Repository interface {
	GetConnection() Connection
}

type GormConnection struct {
	impl *gorm.DB
}

type GormTransaction struct {
	tx *gorm.DB
}

func WrapConnection(db *gorm.DB) Connection { return &GormConnection{impl: db} }

func BeginTransaction(db Connection) Transaction {
	if gc, ok := db.(*GormConnection); ok {
		return &GormTransaction{tx: gc.impl.Begin()}
	}
	if gc, ok := db.(*GormTransaction); ok {
		return gc
	}
	return &GormTransaction{tx: nil}
}

func (d *GormConnection) BeginTransaction() Transaction { return BeginTransaction(d) }

func (d *GormTransaction) BeginTransaction() Transaction { return BeginTransaction(d) }

func (t *GormTransaction) Rollback() { t.tx.Rollback() }

func (t *GormTransaction) Commit() error { return t.tx.Commit().Error }

// Adapter helpers for concrete access (only inside repository implementations)
func AsGormConn(conn Connection) (*gorm.DB, bool) {
	gc, ok := conn.(*GormConnection)
	if !ok || gc == nil {
		return nil, false
	}
	return gc.impl, true
}

func AsGormTx(tx Transaction) (*gorm.DB, bool) {
	gt, ok := tx.(*GormTransaction)
	if !ok || gt == nil {
		return nil, false
	}
	return gt.tx, true
}
