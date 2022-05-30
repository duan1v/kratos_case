package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
	v1 "user/api/user/v1"

	pb "user/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/middleware/tracing"
	grpcx "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"

	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	consulAPI "github.com/hashicorp/consul/api"
)

var userClient v1.UserClient
var conn *grpc.ClientConn

func main() {
	Init()

	// TestCreateUser() // 创建用户

	conn.Close()
}

// Init 初始化 grpc 链接 注意这里链接的 端口
func Init() {
	var err error
	// conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	// if err != nil {
	// 	panic("grpc link err" + err.Error())
	// }
	// userClient = v1.NewUserClient(conn)

	c := consulAPI.DefaultConfig()
	c.Address = "127.0.0.1:8500"
	c.Scheme = "http"
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	conn, err := grpcx.DialInsecure(
		context.Background(),
		grpcx.WithEndpoint("127.0.0.1:50051"),
		// 127.0.0.1:50051/helloworld
		grpcx.WithDiscovery(r),
		grpcx.WithTimeout(2*time.Second),
		grpcx.WithOptions(grpc.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	// userClient = v1.NewUserClient(conn)

	cx := pb.NewGreeterClient(conn)
	for {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		rx, err := cx.SayHello(ctx, &pb.HelloRequest{Name: "efagrteyjr"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", rx.Message)
		time.Sleep(time.Second * 2)
	}
}

func TestCreateUser() {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 10; i++ {
		nn := ""
		for j := 0; j < rand.Intn(20); j++ {
			nn = fmt.Sprintf("%s%s", nn, string(byte(rand.Intn(27)+'a')))
		}
		rsp, err := userClient.CreateUser(context.Background(), &v1.CreateUserInfo{
			Mobile:   fmt.Sprintf("188888%d", rand.Intn(90000)+10000),
			Password: "admin123",
			Nickname: nn,
		})
		if err != nil {
			panic("grpc 创建用户失败" + err.Error())
		}
		fmt.Println(rsp.Id)
	}
	// rsp, err := userClient.GetUserList(context.Background(), &v1.PageInfo{Pn: 1, Psize: 10})
	// if err != nil {
	// 	panic("grpc 创建用户失败" + err.Error())
	// }
	// for k, v := range rsp.Data {
	// 	fmt.Println(k, v)
	// }
}
