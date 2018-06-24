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

package types

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// scalar Hash
// scalar Address
// scalar Hex
// scalar Uint64
// scalar Uint32
var (
	Hash = graphql.NewScalar(graphql.ScalarConfig{
		Name:        "Hash",
		Description: "Hash represents the 32 byte Keccak256 hash of arbitrary data.",
		Serialize:   coerceString,
		ParseValue:  coerceString,
		ParseLiteral: func(valueAST ast.Value) interface{} {
			switch valueAST := valueAST.(type) {
			case *ast.StringValue:
				return common.HexToHash(valueAST.Value)
			}

			return nil
		},
	})

	Hex = graphql.NewScalar(graphql.ScalarConfig{
		Name:        "Hex",
		Description: "Hex is a hexadecimal string",
		Serialize:   coerceString,
		ParseValue:  coerceString,
		ParseLiteral: func(valueAST ast.Value) interface{} {
			switch valueAST := valueAST.(type) {
			case *ast.StringValue:
				if isHex(valueAST.Value) {
					return valueAST.Value
				}
				return nil
			}

			return nil
		},
	})

	Uint64 = graphql.NewScalar(graphql.ScalarConfig{
		Name:        "Uint64",
		Description: "",
		Serialize:   coerceUint64,
		ParseValue:  coerceUint64,
		ParseLiteral: func(valueAST ast.Value) interface{} {
			switch valueAST := valueAST.(type) {
			case *ast.StringValue:
				v, err := strconv.ParseUint(valueAST.Value, 10, 64)
				if err != nil {
					glog.Error(err)
					return nil
				}
				return v
			case *ast.IntValue:
				v, err := strconv.ParseUint(valueAST.Value, 10, 64)
				if err != nil {
					glog.Error(err)
					return nil
				}
				return v
			}
			return nil
		},
	})

	Uint32 = graphql.NewScalar(graphql.ScalarConfig{
		Name:        "Uint32",
		Description: "",
		Serialize:   coerceUint32,
		ParseValue:  coerceUint32,
		ParseLiteral: func(valueAST ast.Value) interface{} {
			switch valueAST := valueAST.(type) {
			case *ast.StringValue:
				v, err := strconv.ParseUint(valueAST.Value, 10, 32)
				if err != nil {
					glog.Error(err)
					return nil
				}
				return uint32(v)
			case *ast.IntValue:
				v, err := strconv.ParseUint(valueAST.Value, 10, 32)
				if err != nil {
					glog.Error(err)
					return nil
				}
				return uint32(v)
			}
			return nil
		},
	})

	Address = graphql.NewScalar(graphql.ScalarConfig{
		Name:        "Address",
		Description: "Address represents the 20 byte address of an Ethereum account.",
		Serialize:   coerceString,
		ParseValue:  coerceString,
		ParseLiteral: func(valueAST ast.Value) interface{} {
			switch valueAST := valueAST.(type) {
			case *ast.StringValue:
				return common.HexToAddress(valueAST.Value)
			}

			return nil
		},
	})
)

func coerceString(value interface{}) interface{} {
	if v, ok := value.(*string); ok {
		return *v
	}
	return fmt.Sprintf("%v", value)
}

func coerceUint64(value interface{}) interface{} {
	switch value := value.(type) {
	case uint:
		return uint64(value)
	case uint8:
		return uint64(value)
	case uint32:
		return uint64(value)
	case uint64:
		return uint64(value)
	case *uint:
		return uint64(*value)
	case *uint8:
		return uint64(*value)
	case *uint32:
		return uint64(*value)
	case *uint64:
		return uint64(*value)
	case string:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			glog.Error(err)
			return nil
		}
		return v
	case *string:
		v, err := strconv.ParseUint(*value, 10, 64)
		if err != nil {
			glog.Error(err)
			return nil
		}
		return v
	}

	return nil
}

func coerceUint32(value interface{}) interface{} {
	switch value := value.(type) {
	case uint:
		return uint32(value)
	case uint8:
		return uint32(value)
	case uint32:
		return uint32(value)
	case *uint:
		return uint32(*value)
	case *uint8:
		return uint32(*value)
	case *uint32:
		return uint32(*value)
	case string:
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			glog.Error(err)
			return nil
		}
		return uint32(v)
	case *string:
		v, err := strconv.ParseUint(*value, 10, 32)
		if err != nil {
			glog.Error(err)
			return nil
		}
		return uint32(v)
	}

	return nil
}

// isHex validates whether each byte is valid hexadecimal string.
func isHex(str string) bool {
	if !strings.HasPrefix(str, "0x") {
		return false
	}

	str = strings.Replace(str, "0x", "", 1)
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}
