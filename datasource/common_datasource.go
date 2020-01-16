/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package datasource

type CommonDataSource struct {
    Name string
    Info string
}

func (ds *CommonDataSource) DriverName() string {
    return ds.Name
}

func (ds *CommonDataSource) DriverInfo() string {
    return ds.Info
}
