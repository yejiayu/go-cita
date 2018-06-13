package main

import (
	"github.com/yejiayu/go-cita/chain"
	chainConfig "github.com/yejiayu/go-cita/config/chain"
	dbConfig "github.com/yejiayu/go-cita/config/db"
	"github.com/yejiayu/go-cita/types"
)

func main() {
	dbConfig := dbConfig.NewDefault()
	dbConfig.Columns = 7
	chainConfig := chainConfig.NewDefault()

	c := chain.New(types.ProofTendermint)

}
