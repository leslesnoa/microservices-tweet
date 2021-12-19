package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/leslesnoa/microservices-tweet/pb"
	"google.golang.org/grpc"
)

const (
	port       = "1234"
	targetPort = "localhost:9090"
)

// Echoのリクエスト/レスポンスボディを出力するミドルウェア
func bodyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	fmt.Printf("Request Body: %v\n", string(reqBody))
	fmt.Printf("Response Body: %v\n", string(resBody))
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.BodyDump(bodyDumpHandler))
	e.GET("/tweets/:tweet_id", GetTweetById)
	e.GET("/tweets", GetAllTweet)
	e.POST("/tweets", CreateTweet)
	e.DELETE("/tweets/:tweet_id", DeleteTweetById)
	e.Logger.Fatal(e.Start("localhost:" + port))
}

// Interceptorの定義
func unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("before call: %s, request: %+v", method, req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("after call: %s, response: %+v", method, reply)
	return err
}

func GetAllTweet(c echo.Context) error {
	log.Println("starting GetAllTweetFunc.")
	conn, err := grpc.Dial(targetPort, grpc.WithInsecure(), grpc.WithUnaryInterceptor(unaryInterceptor))
	if err != nil {
		log.Fatalf("gRPC connection error: %v", err)
	}
	defer conn.Close()

	client := pb.NewTweetServiceClient(conn)
	ctx := context.Background()

	res, err := client.GetAllTweet(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not GetAllTweet: %v", err)
	}
	log.Printf("ResponseValue: %v", res)
	return c.JSON(http.StatusOK, res)
}

func GetTweetById(c echo.Context) error {
	log.Println("Start GetTweetByIdFunc.")
	conn, err := grpc.Dial(targetPort, grpc.WithInsecure(), grpc.WithUnaryInterceptor(unaryInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTweetServiceClient(conn)
	ctx := context.Background()

	// pathパラメータからtweet_idを取得
	tweetId, err := strconv.ParseInt(c.Param("tweet_id"), 10, 64)
	if err != nil {
		fmt.Println("parse error!")
	}

	var t pb.Tweet
	t.Id = tweetId
	res, err := client.GetTweetById(ctx, &pb.TweetByIdRequest{TweetId: tweetId})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("ResponseValue: %v", res)
	return c.JSON(http.StatusOK, res)
}

func CreateTweet(c echo.Context) error {
	log.Println("Start CreateTweets.")
	conn, err := grpc.Dial(targetPort, grpc.WithInsecure(), grpc.WithUnaryInterceptor(unaryInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTweetServiceClient(conn)
	ctx := context.Background()

	var t pb.Tweet
	c.Bind(&t)

	res, err := client.CreateTweet(ctx, &t)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("ResponseValue: %v", res)

	return c.JSON(http.StatusOK, res)
}

func DeleteTweetById(c echo.Context) error {
	log.Println("Start DeleteTweetById.")
	conn, err := grpc.Dial(targetPort, grpc.WithInsecure(), grpc.WithUnaryInterceptor(unaryInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTweetServiceClient(conn)
	ctx := context.Background()

	// pathパラメータからtweet_idを取得
	tweetId, err := strconv.ParseInt(c.Param("tweet_id"), 10, 64)
	if err != nil {
		fmt.Println("parse error!")
	}

	var t pb.Tweet
	t.Id = tweetId

	res, err := client.DeleteTweetById(ctx, &pb.TweetByIdRequest{TweetId: tweetId})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("ResponseValue: %v", res)

	return c.JSON(http.StatusOK, "ok")
}
