package cmd

import (
	"context"
	"github.com/spf13/cobra"
	rpc "go_grpc_demo/server/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
	"time"
)

// 测试客户端

type clientArgs struct {
	addr    string
	token   string
	message string
}

// represents the serve command
func init() {
	var (
		args clientArgs

		clientCmd = &cobra.Command{
			Use:   "client",
			Short: "start client",
			Run: func(_ *cobra.Command, _ []string) {
				startClient(&args)
			},
		}
	)

	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVarP(&args.addr, "server", "s", "localhost:12346", "addr of server")
	clientCmd.Flags().StringVarP(&args.token, "token", "t", "", "token for auth")
	clientCmd.Flags().StringVarP(&args.message, "message", "m", "world", "context for say")
}

func startClient(args *clientArgs) {
	// 连接到服务器
	conn, err := grpc.Dial(args.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	randSleep()

	// 验证登录
	{
		authClient := rpc.NewAuthClient(conn)
		r, err := authClient.Auth(ctx, &rpc.AuthRequest{Token: args.token})
		if err != nil {
			log.Fatalf("could not auth: %v", err)
		}
		log.Printf("auth result: {ok: %d, message: %s}", r.GetOk(), r.GetMessage())
		if r.GetOk() != 1 {
			return
		}
	}

	randSleep()

	// 其他服务
	{
		demoClient := rpc.NewDemoClient(conn)
		r, err := demoClient.Say(ctx, &rpc.SayRequest{Content: args.message})
		if err != nil {
			log.Fatalf("could not say: %v", err)
		}
		log.Printf("say result: {message: %s}", r.GetMessage())
	}
}

func randSleep() {
	rand.Seed(time.Now().Unix())
	d := rand.Int() % 10
	log.Printf("sleeping %d seconds....", d)
	time.Sleep(time.Second * time.Duration(d))
	log.Printf("sleep ended")
}
