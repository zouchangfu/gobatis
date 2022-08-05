/*
 * Licensed to the AcmeStack under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for
 * additional information regarding copyright ownership.
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
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"time"

	"github.com/acmestack/gobatis/datasource"
	"github.com/acmestack/gobatis/factory"
	"github.com/acmestack/gobatis/logging"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
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

type Gobatis struct {
	Config Config `yaml:"gobatis"`
}

type Config struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	DBName      string `yaml:"db-name"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Charset     string `yaml:"charset"`
	MaxConn     int    `yaml:"max-conn"`
	MaxIdleConn int    `yaml:"max-idle-conn"`
	Type        string `yaml:"type"`
	SslMode     string `yaml:"ssl-mode"`
	SqlLitePath string `yaml:"sql-lite-path"`
}

func Connect(path string) (factory.Factory, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	gobatis := &Gobatis{}
	if err := yaml.Unmarshal([]byte(data), &gobatis); err != nil {
		return nil, err
	}
	conf := gobatis.Config
	if conf.Type == "" {
		return nil, errors.New("The database type is not configured")
	}
	if strings.ToLower(conf.Type) == MYSQL {
		return NewFactory(
			SetMaxConn(conf.MaxConn),
			SetMaxIdleConn(conf.MaxIdleConn),
			SetDataSource(&datasource.MysqlDataSource{
				Host:     conf.Host,
				Port:     conf.Port,
				DBName:   conf.DBName,
				Username: conf.Username,
				Password: conf.Password,
				Charset:  conf.Charset,
			})), nil
	} else if strings.ToLower(conf.Type) == SQLLITE {
		return NewFactory(
			SetMaxConn(conf.MaxConn),
			SetMaxIdleConn(conf.MaxIdleConn),
			SetDataSource(&datasource.PostgreDataSource{
				Host:     conf.Host,
				Port:     conf.Port,
				DBName:   conf.DBName,
				Username: conf.Username,
				Password: conf.Password,
				SslMode:  conf.SslMode,
			})), nil
	} else if strings.ToLower(conf.Type) == POSTGRE {
		return NewFactory(
			SetMaxConn(conf.MaxConn),
			SetMaxIdleConn(conf.MaxIdleConn),
			SetDataSource(&datasource.SqliteDataSource{
				Path: conf.SqlLitePath,
			})), nil
	}

	return nil, errors.New("Only support mysql,sqllite,portgre ")

}
