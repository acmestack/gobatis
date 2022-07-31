package impl

import (
	"github.com/acmestack/gobatis/plus/mapper"
	_ "github.com/go-sql-driver/mysql"
)

type UserMapperImpl[T any] struct {
	mapper.BaseMapper[T]
}
