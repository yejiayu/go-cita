// Copyright (C) 2018 yejiayu

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package raw

import (
	"context"
)

type Interface interface {
	Put(ctx context.Context, namespace string, key, value []byte) error
	Get(ctx context.Context, namespace string, key []byte) ([]byte, error)
	Delete(ctx context.Context, namespace string, key []byte) error
	BatchGet(ctx context.Context, namespace string, keys [][]byte) ([][]byte, error)
	BatchPut(ctx context.Context, namespace string, keys, values [][]byte) error
	BatchDelete(ctx context.Context, namespace string, keys [][]byte) error

	Scan(ctx context.Context, namespace string, prefix []byte, limit int) ([][]byte, [][]byte, error)

	Close() error
}
