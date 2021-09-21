package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	"fguzman/grpc-demo/proto-generated/item_proto"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("couldn't connect: %v", err)
	}
	defer conn.Close()

	c := item_proto.NewItemServiceClient(conn)

	fmt.Printf("\n\n###################################\n")
	getItem(c)
	fmt.Printf("\n\n###################################\n")
	getItemBySellerID(c)
	fmt.Printf("\n\n###################################\n")
	createItems(c)
	fmt.Printf("\n\n###################################\n")
	bidi(c)
}

func getItem(c item_proto.ItemServiceClient) {
	fmt.Printf("Starting to do a Unary RPC\n\n")
	itemId := &item_proto.ItemId{Id: "MLA123"}

	res, err := c.Get(context.Background(), itemId)
	if err != nil {
		log.Panicf("couldn't get response: %v", err)
	}
	log.Printf("Item: %v\n", res)
}

func getItemBySellerID(c item_proto.ItemServiceClient) {
	fmt.Printf("Starting to do a Server Streaming RPC\n\n")

	sellerID := &item_proto.SellerId{Id: 11111}
	stream, err := c.GetItemsBySeller(context.Background(), sellerID)
	if err != nil {
		log.Fatalf("couldn't get response from GreatManyTimes Server: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Panicf("error reading stream: %v", err)
		}

		fmt.Printf("Item %v\n", msg)
	}
}

func createItems(c item_proto.ItemServiceClient) {
	fmt.Printf("Starting to do a Client Streaming RPC\n\n")

	items := []*item_proto.Item{
		{
			Id:    "MLA1234",
			Title: "title 1",
		},
		{
			Id:    "MLA5678",
			Title: "title 2",
		},
		{
			Id:    "MLA8888",
			Title: "title 3",
		},
		{
			Id:    "MLA6666",
			Title: "title 4",
		},
	}

	stream, err := c.CreateItems(context.Background())
	if err != nil {
		log.Panicf("error while calling CreateItems: %v", err)
	}

	for _, req := range items {
		fmt.Printf("Sending item:%v\n", req)
		stream.Send(req)
		time.Sleep(time.Second)
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Panicf("error while recieving response from CreateItems: %v", err)
	}

	fmt.Printf("\nItem ids created: %v\n", response)
}

func bidi(c item_proto.ItemServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC")

	requests := []*item_proto.ItemId{
		{
			Id: "MLA1234",
		},
		{
			Id: "MLA3333",
		},
		{
			Id: "MLA4444",
		},
		{
			Id: "MLA5555",
		},
	}

	stream, err := c.BidiItems(context.Background())
	if err != nil {
		log.Panicf("error while calling BiDI stream: %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending item id:%v\n", req)
			stream.Send(req)
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v\n", err)
			}
			fmt.Printf("Received: %v\n", response.GetId())
		}
		close(waitc)
	}()

	<-waitc
}
