/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package reflection

import "reflect"

type Object interface {
    NewValue() reflect.Value
    AddValue(reflect.Value)
    GetClassName() string
}
