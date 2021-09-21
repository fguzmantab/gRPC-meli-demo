package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"fguzman/grpc-demo/proto-generated/item_proto"
)

type server struct {
	item_proto.UnimplementedItemServiceServer
}

// Get Unary
func (s *server) Get(ctx context.Context, itemId *item_proto.ItemId) (*item_proto.Item, error) {
	fmt.Printf("Getting item %s", itemId.GetId())
	item := &item_proto.Item{
		Id:             itemId.GetId(),
		SiteId:         "MLA",
		Title:          "Item de test",
		SellerId:       1234566,
		CatalogListing: false,
		Permalink:      "www.test.com/item",
		Attributes: []*item_proto.Item_Attribute{
			{
				Id:        "001",
				Name:      "fipe",
				ValueId:   "123",
				ValueName: "0123123-1",
				Values: []*item_proto.Item_Attribute_Values{
					{
						Id:   "2222",
						Name: "fipe",
					},
				},
			},
		},
	}
	return item, nil
}

// GetItemsBySeller Server streaming
func (s server) GetItemsBySeller(sellerID *item_proto.SellerId, stream item_proto.ItemService_GetItemsBySellerServer) error {
	fmt.Printf("Calling stream server with %s\n", sellerID.GetId())

	for i := 0; i < 3; i++ {
		item := &item_proto.Item{
			Id:             fmt.Sprintf("MLA000%v", i),
			SiteId:         "MLA",
			Title:          "Item de test",
			SellerId:       sellerID.GetId(),
			CatalogListing: false,
			Permalink:      "www.test.com/item",
		}

		stream.Send(item)
		time.Sleep(time.Second)
	}

	return nil
}

// CreateItems Client streaming
func (s server) CreateItems(stream item_proto.ItemService_CreateItemsServer) error {
	fmt.Printf("Calling stream client\n")
	var ids []string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			var itemIDs []*item_proto.ItemId
			for _, id := range ids {
				itemIDs = append(itemIDs, &item_proto.ItemId{Id: id})
			}
			return stream.SendAndClose(&item_proto.ItemIds{Ids: itemIDs})
		}

		if err != nil {
			log.Panicf("Error while reading client stream: %v", err)
		}

		ids = append(ids, req.GetId())
	}

	return nil
}

// BidiItems Bidirectional streaming
func (s server) BidiItems(stream item_proto.ItemService_BidiItemsServer) error {
	fmt.Printf("BidiItems function was invoked with a bi streaming request\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Panicf("Error while reading client stream: %v", err)
			return err
		}

		req.GetId()
		err = stream.Send(&item_proto.Item{Id: req.GetId(), Title: "Title: " + req.GetId()})
		if err != nil {
			log.Fatalf("Error while sending data to client: %v", err)
			return err
		}
	}
}

func main() {
	list, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	item_proto.RegisterItemServiceServer(s, &server{})

	if err := s.Serve(list); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
