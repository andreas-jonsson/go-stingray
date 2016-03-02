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

import "sync"

const defaultSize = 1000

type ramDatabase struct {
	db    map[uint64][]byte
	mutex sync.Mutex
}

type ramSession struct {
	db *ramDatabase
}

func (p *ramSession) Close() {
}

func (p *ramSession) Load(id uint64, data []byte) []byte {
	return p.db.load(id, data)
}

func (p *ramSession) Store(id uint64, data []byte) {
	p.db.store(id, data)
}

func NewRAMDatabase() Database {
	return &ramDatabase{db: make(map[uint64][]byte, defaultSize)}
}

func (p *ramDatabase) Open() (Session, error) {
	return &ramSession{p}, nil
}

func (p *ramDatabase) Close() {
}

func (p *ramDatabase) store(id uint64, data []byte) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	size := len(data)
	blob, ok := p.db[id]

	if !ok || size > cap(blob) {
		blob = make([]byte, 0, size)
	}

	copy(blob, data)
	p.db[id] = blob
}

func (p *ramDatabase) load(id uint64, data []byte) []byte {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	blob, ok := p.db[id]
	if !ok {
		return data[:0]
	}

	size := len(blob)
	if size > cap(data) {
		data = make([]byte, 0, size)
	}

	copy(data, blob)
	return data
}
