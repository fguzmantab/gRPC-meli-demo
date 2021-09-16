package main

import (
	"context"
	"fmt"
	"io"
	"log"

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

	getItem(c)
	getItemBySellerId(c)
}

func getItem(c item_proto.ItemServiceClient) {
	itemId := &item_proto.ItemId{Id: "MLA123"}

	res, err := c.Get(context.Background(), itemId)
	if err != nil {
		log.Panicf("couldn't get response: %v", err)
	}
	log.Printf("Item: %v\n", res)
}

func getItemBySellerId(c item_proto.ItemServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC")

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
