/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package gobatis

import (
	"github.com/xfali/gobatis/reflection"
	"reflect"
)

func init() {
	var typeModelName ModelName
	reflection.SetModelNameType(reflect.TypeOf(typeModelName))
}
