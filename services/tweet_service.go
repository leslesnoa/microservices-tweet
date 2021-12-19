package main

import (
	"context"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
	"github.com/leslesnoa/microservices-tweet/db"
	pb "github.com/leslesnoa/microservices-tweet/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// gRPCサーバ
type server struct {
}

const (
	port = ":9090"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTweetServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %c", err)
	}
	log.Println("starting on gRPC server on " + port)
}

func (s *server) CreateTweet(ctx context.Context, r *pb.Tweet) (*pb.Empty, error) {
	log.Printf("Recieved CreateTweet : %s", r)
	// DB接続
	conn := db.Connect()
	db.CreateRow(conn, r)
	return &pb.Empty{}, nil
}

func (s *server) GetAllTweet(ctx context.Context, r *pb.Empty) (*pb.Tweets, error) {
	log.Printf("Recieved GetAllTweet : %s", r)

	// DB接続
	conn := db.Connect()

	// 接続確認
	err := conn.Ping()
	if err != nil {
		log.Println("connection failed.")
		panic(err)
	} else {
		log.Println("connection success.")
	}

	// 行データ取得
	rows := db.GetRows(conn)
	log.Println(rows)
	var result []*pb.Tweet
	for _, p := range rows {
		result = append(result, p)
	}
	defer conn.Close()
	return &pb.Tweets{Tweets: rows}, nil

	// return &pb.Posts{Posts: []*pb.Post{
	// 	{
	// 		Id:    1,
	// 		Title: "testTitle",
	// 		Text:  "testText",
	// 	},
	// 	{
	// 		Id:    2,
	// 		Title: "testTitle",
	// 		Text:  "testText",
	// 	},
	// }}, nil
}

func (s *server) GetTweetById(ctx context.Context, r *pb.TweetByIdRequest) (*pb.Tweets, error) {
	log.Printf("Recieved GetTweetById : %s", r)

	// DB接続
	conn := db.Connect()
	defer conn.Close()

	// 接続確認
	err := conn.Ping()
	if err != nil {
		log.Println("connection failed.")
		panic(err)
	} else {
		log.Println("connection success.")
	}

	log.Println(r.GetTweetId())

	// 行データ取得
	row := db.GetRowByTweetId(conn, r)
	log.Println(row)
	// var result []*pb.Tweet
	// result = append(result, row)

	// return &pb.Tweets{Tweets: []*pb.Tweet{}}, nil
	return &pb.Tweets{Tweets: row}, nil
}

func (s *server) DeleteTweetById(ctx context.Context, r *pb.TweetByIdRequest) (*pb.Empty, error) {
	log.Printf("Recieved DeleteTweetById : %s", r)

	// DB接続
	conn := db.Connect()
	defer conn.Close()

	// 接続確認
	err := conn.Ping()
	if err != nil {
		log.Println("connection failed.")
		panic(err)
	} else {
		log.Println("connection success.")
	}
	log.Println(r.TweetId)

	// 行データ削除
	db.DeleteRowByTweetId(conn, r)

	return &pb.Empty{}, nil
}
