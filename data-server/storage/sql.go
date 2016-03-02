// +build sql

/*
Copyright (C) 2015-2016 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package storage

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type sqlDatabaseInterface interface {
	Close()
	Open() (Session, error)
	load(id uint64, data []byte) []byte
	store(id uint64, data []byte)
}

type sqlLocklessDatabase struct {
	db *sql.DB

	stmtSum    *sql.Stmt
	stmtGet    *sql.Stmt
	stmtPut    *sql.Stmt
	stmtUpdate *sql.Stmt
	stmtCount  *sql.Stmt
}

type sqlDatabase struct {
	db    *sqlLocklessDatabase
	mutex sync.Mutex
}

type sqlSession struct {
	db sqlDatabaseInterface
}

func (p *sqlSession) Load(id uint64, data []byte) []byte {
	return p.db.load(id, data)
}

func (p *sqlSession) Store(id uint64, data []byte) {
	p.db.store(id, data)
}

func (p *sqlSession) Close() {
}

func NewSQLDatabase(driver string, con string, lockless bool) Database {
	if lockless {
		return newSQLLocklessDatabase(driver, con)
	} else {
		return &sqlDatabase{db: newSQLLocklessDatabase(driver, con)}
	}
}

func (p *sqlDatabase) Close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.db.Close()
}

func (p *sqlDatabase) Open() (Session, error) {
	return &sqlSession{p.db}, nil
}

func (p *sqlDatabase) load(id uint64, data []byte) []byte {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.db.load(id, data)
}

func (p *sqlDatabase) store(id uint64, data []byte) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.db.store(id, data)
}

func newSQLLocklessDatabase(driver string, con string) *sqlLocklessDatabase {
	db, err := sql.Open(driver, con)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cache (hash INT8 PRIMARY KEY, size INT8, time TEXT, metadata TEXT, data BLOB);")
	if err != nil {
		panic(err)
	}

	p := &sqlLocklessDatabase{db: db}
	p.stmtSum = p.prepStmt("SELECT SUM(size) FROM cache;")
	p.stmtGet = p.prepStmt("SELECT data FROM cache WHERE hash = ?;")
	p.stmtPut = p.prepStmt("INSERT OR REPLACE INTO cache VALUES (?, ?, datetime('now', 'localtime'), ?, ?);")
	p.stmtUpdate = p.prepStmt("UPDATE cache SET time = datetime('now', 'localtime') WHERE hash = ?;")
	p.stmtCount = p.prepStmt("SELECT COUNT(*) FROM cache;")
	return p
}

func (p *sqlLocklessDatabase) Close() {
	p.db.Close()
}

func (p *sqlLocklessDatabase) Open() (Session, error) {
	return &sqlSession{p}, nil
}

func (p *sqlLocklessDatabase) prepStmt(stmt string) *sql.Stmt {
	s, err := p.db.Prepare(stmt)
	if err != nil {
		panic(err)
	}
	return s
}

func (p *sqlLocklessDatabase) load(id uint64, data []byte) []byte {
	err := p.stmtGet.QueryRow(id).Scan(&data)
	if err != nil {
		return data[0:0]
	} else {
		return data
	}
}

func (p *sqlLocklessDatabase) store(id uint64, data []byte) {
	dataSize := len(data)
	metaSize := 0

	for metaSize < dataSize {
		if data[metaSize] == 0 {
			break
		}
		metaSize++
	}

	_, err := p.stmtPut.Exec(int64(id), int64(dataSize-metaSize), string(data[:metaSize]), data[metaSize+1:])
	if err != nil {
		panic(err)
	}
}
