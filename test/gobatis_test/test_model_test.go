package test

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xfali/gobatis"
	"github.com/xfali/gobatis/datasource"
	"log"
	"testing"
	"time"
)

var sessionMgr *gobatis.SessionManager

func Init() {
	fac := gobatis.NewFactory(
		gobatis.SetMaxConn(100),
		gobatis.SetMaxIdleConn(50),
		gobatis.SetDataSource(&datasource.MysqlDataSource{
			Host:     "localhost",
			Port:     3306,
			DBName:   "test",
			Username: "root",
			Password: "123456",
			Charset:  "utf8",
		}))
	var testV TestTable
	gobatis.RegisterModel(&testV)
	sessionMgr = gobatis.NewSessionManager(fac)
	gobatis.RegisterMapperFile("./xml/test_table_mapper.xml")
}

func TestTestTable_Insert(t *testing.T) {
	Init()
	session := sessionMgr.NewSession()
	table := &TestTable{
		CreateTime: time.Time{},
		Password:   "123",
		Username:   "123",
	}
	table.Insert(session)
}

func TestTestTable_Select(t *testing.T) {
	Init()
	session := sessionMgr.NewSession()
	table := &TestTable{
		Username: "user1;drop table sentence",
		Id:       52,
	}
	tables, _ := table.Select(session)
	marshal, _ := json.Marshal(tables)
	log.Println(string(marshal))
}
