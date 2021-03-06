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

type Database interface {
	Close()
	Open() (Session, error)
}

type Session interface {
	Close()
	Load(id uint64, data []byte) []byte
	Store(id uint64, data []byte)
}

func NewDatabase(driverFlag, sourceFlag string, locklessFlag bool) Database {
	switch driverFlag {
	case "ram":
		return NewRAMDatabase()
	case "mysql", "sqlite3":
		return NewSQLDatabase(driverFlag, sourceFlag, locklessFlag)
	case "bolt":
		return NewBoltDatabase(sourceFlag)
	default:
		panic("Invalid database driver!")
	}
}
