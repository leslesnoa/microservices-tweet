package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	pb "github.com/leslesnoa/microservices-tweet/pb"
)

const (
	SqlNoRows = "no rows in result set"
)

// type Post struct {
// 	Id   int
// 	Name string
// 	Text string
// }

func Connect() *sql.DB {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	log.Println("start dbConnectFunc.")
	user := "test"
	password := "test"
	host := "localhost"
	port := "3306"
	dbName := "testdb"

	conn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=utf8mb4"

	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Println("DB接続時エラー")
		panic(err.Error())
	}
	return db
}

func GetRows(db *sql.DB) []*pb.Tweet {
	log.Println("start GetRowsFunc.")
	cmd := "SELECT * FROM tweets;"
	rows, err := db.Query(cmd)
	if err != nil {
		log.Println("Query実行時エラー")
		log.Println(err.Error())
	}
	defer rows.Close()
	log.Println("get rows query success.")

	var result []*pb.Tweet
	for rows.Next() {
		var p *pb.Tweet
		p = new(pb.Tweet)
		err := rows.Scan(&p.Id, &p.UserId, &p.Content)
		if err != nil {
			log.Println("クエリ実行時エラー")
			log.Fatal(err.Error())
		}
		result = append(result, p)
	}
	return result
}

func CreateRow(db *sql.DB, r *pb.Tweet) {
	log.Println("Starting CreateRowFunc.")
	// db.Prepare()
	stmtInsert, err := db.Prepare("INSERT INTO tweets(user_id, content) VALUES(?, ?);")
	if err != nil {
		panic(err.Error())
	}
	defer stmtInsert.Close()

	t := pb.Tweet{UserId: r.UserId, Content: r.Content}
	result, err := stmtInsert.Exec(t.UserId, t.Content)
	if err != nil {
		panic(err.Error())
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(lastInsertID)

	// return &pb.User{Id: int32(lastInsertID), Name: r.Name, Email: r.Email}, nil
}

func GetRowByTweetId(db *sql.DB, r *pb.TweetByIdRequest) []*pb.Tweet {
	log.Println("start GetRowByTweetIdFunc.")
	log.Println(r.GetTweetId())
	var t *pb.Tweet
	t = new(pb.Tweet)
	err := db.QueryRow("SELECT * FROM tweets WHERE id=?;", r.TweetId).Scan(&t.Id, &t.UserId, &t.Content)
	if err != nil {
		panic(err)
	}
	log.Println("get rows query success.")

	var result []*pb.Tweet
	result = append(result, t)

	return result
}

func DeleteRowByTweetId(db *sql.DB, r *pb.TweetByIdRequest) {
	log.Println("start DeleteRowByTweetIdFunc.")
	log.Println(r.GetTweetId())
	stmtDelete, err := db.Prepare("DELETE FROM tweets WHERE id=?;")
	if err != nil {
		panic(err)
	}
	defer stmtDelete.Close()

	log.Println(r.TweetId)
	result, err := stmtDelete.Exec(r.TweetId)
	if err != nil {
		panic(err.Error())
	}
	rowsAffect, err := result.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(rowsAffect)

	log.Println("delete rows query success.")
}
