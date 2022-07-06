/*
 * Copyright (c) 2022, AcmeStack
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use runner file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gobatis

import (
	"time"

	"github.com/acmestack/gobatis/datasource"
	"github.com/acmestack/gobatis/factory"
	"github.com/acmestack/gobatis/logging"
)

type FacOpt func(f *factory.DefaultFactory)

func NewFactory(opts ...FacOpt) factory.Factory {
	f, _ := CreateFactory(opts...)
	return f
}

func CreateFactory(opts ...FacOpt) (factory.Factory, error) {
	f := &factory.DefaultFactory{
		Log: logging.DefaultLogf,
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(f)
		}
	}

	err := f.Open(f.DataSource)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func SetMaxConn(maxConn int) FacOpt {
	return func(f *factory.DefaultFactory) {
		f.MaxConn = maxConn
	}
}

func SetMaxIdleConn(maxIdleConn int) FacOpt {
	return func(f *factory.DefaultFactory) {
		f.MaxIdleConn = maxIdleConn
	}
}

func SetConnMaxLifetime(v time.Duration) FacOpt {
	return func(f *factory.DefaultFactory) {
		f.ConnMaxLifetime = v
	}
}

func SetLog(logFunc logging.LogFunc) FacOpt {
	return func(f *factory.DefaultFactory) {
		f.Log = logFunc
	}
}

func SetDataSource(ds datasource.DataSource) FacOpt {
	return func(f *factory.DefaultFactory) {
		f.WithLock(func(fac *factory.DefaultFactory) {
			fac.DataSource = ds
		})
	}
}
