package lamportserver

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "chiman"
	password = "ramila"
	dbname   = "db_test"
)

type skierStat struct {
	resortID  int `db:"resort_id"`
	dayNum    int `db:"day_num"`
	skierID   int `db:"skier_id"`
	liftID    int `db:"lift_id"`
	timeStamp int `db:"time_stamp"`
	vertical  int `db:"verticals"`
}

var schema = `
	CREATE TABLE skier_stats (
		resort_id int,
		day_num int,
		skier_id int,
		lift_id int,
		time_stamp int,
		verticals int
	)`

// func init() {
var postgresURL = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbname)

// fmt.Println("Successfully connected!")
// }

func writeToDB() {
	db, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	db.Exec("DROP TABLE IF EXISTS skier_stats")
	db.MustExec(schema)

	tx := db.MustBegin()
	for _, stat := range statCache {
		tx.NamedExec("INSERT INTO skier_stats (resort_id, day_num, skier_id, lift_id, time_stamp, verticals ) VALUES (:resort_id, :day_num, :skier_id, :lift_id, :time_stamp, :verticals)", stat)
	}
	tx.Commit()
	fmt.Println("single threaded db write took", time.Since(start))
}
