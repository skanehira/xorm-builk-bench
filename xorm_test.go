package main

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

type Todo struct {
	ID   int    `xorm:"id"`
	Todo string `xorm:"todo"`
}

func (Todo) TableName() string {
	return "todos"
}

func setup() *xorm.Engine {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Asia%%2FTokyo",
		"test", "test", "localhost", "3306", "test")

	db, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func Benchmark_builkInsert(b *testing.B) {
	db := setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var todos []Todo
		for j := 0; j < 5000; j++ {
			todos = append(todos, Todo{
				ID:   j,
				Todo: "gorilla" + strconv.Itoa(j),
			})
		}
		if _, err := db.Insert(&todos); err != nil {
			log.Fatal(err)
		}
	}
}

func Benchmark_insert(b *testing.B) {
	db := setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := db.NewSession()
		for j := 0; j < 5000; j++ {
			if _, err := x.InsertOne(&Todo{
				ID:   j,
				Todo: "gorilla" + strconv.Itoa(j),
			}); err != nil {
				x.Rollback()
				log.Fatal(err)
			}
		}
		x.Commit()
	}
}
