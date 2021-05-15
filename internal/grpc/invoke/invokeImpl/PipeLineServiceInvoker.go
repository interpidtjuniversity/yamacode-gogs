package invoker

import (
	"context"
	"gogs.io/gogs/internal/grpc/invoke"
	"google.golang.org/grpc"
	"log"
	"time"
)

func InvokePipeLineService() *invoke.StartYaMaPipeLineResponse {
	//连接到gRPC服务端
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//建立客户端
	c := invoke.NewYaMaPipeLineServiceClient(conn)

	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 调用方法
	r, err := c.StartYaMaPipeLine(context.Background(), &invoke.StartYaMaPipeLineRequest{UserId: 1, UserName: "interpidtjuniversity", Branch: "master", Repository: "init"})
	if err != nil {
		log.Fatalf("couldn not greet: %v", err)
		return nil
	}
	return r
}

func InvokePassMergerRequestCodeReview(actionId, stageId, stepId int64) (bool, error) {
	conn := invoke.GetConnection()
	defer invoke.Return(conn)
	client := invoke.NewYaMaPipeLineServiceClient(conn)
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := client.PassMergeRequestCodeReview(context.Background(), &invoke.PassMergeRequestCodeReviewRequest{
		ActionId: actionId,
		StageId: stageId,
		StepId: stepId,
	})
	if response == nil {
		return false, err
	}
	return response.Success, err
}
