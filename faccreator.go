/*
 * Copyright (c) 2022, OpeningO
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
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
	"github.com/xfali/gobatis/datasource"
	"github.com/xfali/gobatis/factory"
	"github.com/xfali/gobatis/logging"
	"time"
)

type FacOpt func(f *factory.DefaultFactory)

// 创建工厂，传入的是可变参数，而且都是方法参数
// 方法的参数必须是 f *factory.DefaultFactory 类型的参数
func NewFactory(opts ...FacOpt) factory.Factory {
	f, _ := CreateFactory(opts...)
	return f
}

func CreateFactory(opts ...FacOpt) (factory.Factory, error) {

	// 创建默认工厂对象
	f := &factory.DefaultFactory{
		Log: logging.DefaultLogf,
	}

	// 遍历所有的参数，执行所有的参数方法
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(f)
		}
	}

	// 和数据库建立连接
	err := f.Open(f.DataSource)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func SetMaxConn(v int) FacOpt {
	return func(f *factory.DefaultFactory) {
		f.MaxConn = v
	}
}

func SetMaxIdleConn(v int) FacOpt {
	return func(f *factory.DefaultFactory) {
		f.MaxIdleConn = v
	}
}

func SetConnMaxLifetime(v time.Duration) FacOpt {
	return func(f *factory.DefaultFactory) {
		f.ConnMaxLifetime = v
	}
}

func SetLog(v logging.LogFunc) FacOpt {
	return func(f *factory.DefaultFactory) {
		f.Log = v
	}
}

func SetDataSource(v datasource.DataSource) FacOpt {
	return func(f *factory.DefaultFactory) {
		f.WithLock(func(fac *factory.DefaultFactory) {
			fac.DataSource = v
		})
	}
}
