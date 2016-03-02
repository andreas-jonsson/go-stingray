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
	"encoding/binary"

	"github.com/boltdb/bolt"
)

type boltDatabase struct {
	db *bolt.DB
}

type boltSession struct {
	db *bolt.DB
}

func (p *boltSession) Close() {
}

func (p *boltSession) Load(id uint64, data []byte) []byte {
	p.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("cache"))
		var buf [8]byte
		s := buf[:]
		binary.BigEndian.PutUint64(s, id)
		r := b.Get(s)
		ln := len(r)
		memory := &data

		for ca := cap(*memory); ln > ca; {
			ca = ca*2 + 1
			*memory = make([]byte, 0, ca)
		}

		*memory = (*memory)[:ln]
		copy(*memory, r)
		return nil
	})
	return data
}

func (p *boltSession) Store(id uint64, data []byte) {
	metaSize := 0
	dataSize := len(data)

	for metaSize < dataSize {
		if data[metaSize] == 0 {
			break
		}
		metaSize++
	}

	p.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("cache"))
		var buf [8]byte
		s := buf[:]
		binary.BigEndian.PutUint64(s, id)
		return b.Put(s, data[metaSize+1:])
	})
}

func NewBoltDatabase(con string) Database {
	db, err := bolt.Open(con, 0600, nil)
	if err != nil {
		panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("cache"))
		if err != nil {
			panic(err)
		}
		return nil
	})

	return &boltDatabase{db}
}

func (p *boltDatabase) Open() (Session, error) {
	return &boltSession{p.db}, nil
}

func (p *boltDatabase) Close() {
	p.db.Close()
}
