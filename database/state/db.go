package state

import (
	"github.com/yejiayu/go-cita/database"
)

type Interface interface {
}

type stateDB struct {
	rawDB database.Raw
}
