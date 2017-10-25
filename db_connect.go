package lamportserver

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "chiman"
	password = "ramila"
	dbname   = "db_test"
)

type skierStat struct {
	ResortID  int `db:"resort_id"`
	DayNum    int `db:"day_num"`
	SkierID   int `db:"skier_id"`
	LiftID    int `db:"lift_id"`
	TimeStamp int `db:"time_stamp"`
	Vertical  int `db:"verticals"`
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

func init() {
	db, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("DROP TABLE IF EXISTS skier_stats")
	db.MustExec(schema)

}

var postgresURL = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbname)

// fmt.Println("Successfully connected!")
// }

func queryDB(skierID, dayNum int) (int, int) {
	db, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		log.Fatal(err)
	}
	var verticals, lifts int
	for trigger := <-getTrigger; trigger; trigger = <-receiveTrigger {
		err := db.Get(&verticals, "SELECT SUM(verticals) FROM skier_stats WHERE skier_id=$1 AND day_num=$2", skierID, dayNum)
		if err != nil {
			log.Fatal(err)
		}
		err = db.Get(&lifts, "SELECT COUNT(verticals) FROM skier_stats WHERE skier_id=$1 AND day_num=$2", skierID, dayNum)
		if err != nil {
			log.Fatal(err)
		}
	}
	return verticals, lifts
}

func writeToDB() {
	db, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		log.Fatal(err)
	}
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// for stat := range statChan {
	// 	tx.NamedExec("INSERT INTO skier_stats (resort_id, day_num, skier_id, lift_id, time_stamp, verticals ) VALUES (:resort_id, :day_num, :skier_id, :lift_id, :time_stamp, :verticals)", stat)
	// }
	for trigger := <-receiveTrigger; trigger; trigger = <-receiveTrigger {
		select {
		case stat := <-statChan:
			tx := db.MustBegin()
			_, err := tx.NamedExec("INSERT INTO skier_stats (resort_id, day_num, skier_id, lift_id, time_stamp, verticals ) VALUES (:resort_id, :day_num, :skier_id, :lift_id, :time_stamp, :verticals)", stat)
			if err != nil {
				log.Fatal(err)
			}
			tx.Commit()
		}
	}
	// fmt.Println("single threaded db write took", time.Since(start))
}
