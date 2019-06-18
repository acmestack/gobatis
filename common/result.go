/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package common

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
