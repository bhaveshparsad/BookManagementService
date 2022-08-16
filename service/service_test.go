package service

import (
	"context"
	pb "BookManagementService/protoFiles"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {

	lis = bufconn.Listen(bufSize)
	server := grpc.NewServer()
	pb.RegisterBookMgmtServiceServer(server, &BookMgmtServiceServer{})
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestService(t *testing.T) {

	//Dial a connection to grpc Server
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	//Create new Client
	client := pb.NewBookMgmtServiceClient(conn)

	//Test CreateBook
	newBook := &pb.Book{
		BookTitle:  "test title",
		BookAuthor: "test author",
	}

	response, err := client.CreateBook(ctx, &pb.CreateRequest{Book: newBook})
	if err != nil {
		t.Error("Test CreateBook FAILED!\nerr: ", err)
	}
	t.Log("Test CreateBook PASSED.")

	//get id
	id := response.Book.Id

	_, err = client.GetAllBooks(ctx, &pb.GetAllRequest{})
	if err != nil {
		t.Error("Test ListAllBooks FAILED!\nerr: ", err)
	}
	t.Log("Test ListAllBooks PASSED.")

	//Test SearchBooks
	_, err = client.SearchBook(ctx, &pb.SearchRequest{
		Search: &pb.SearchRequest_BookTitle{BookTitle: newBook.BookTitle},
	})
	if err != nil {
		t.Error("Test SearchBooks FAILED!\nerr: ", err)
	}

	t.Log("Test SearchBooks PASSED.")

	//Test Updatebook
	updateBook := &pb.Book{
		Id:         uint64(id),
		BookTitle:  "updated test title",
		BookAuthor: "updated test title",
	}
	_, err = client.UpdateBook(ctx, &pb.UpdateRequest{Title: updateBook.BookTitle})
	if err != nil {
		t.Error("Test UpdateBooks FAILED!\nerr: ", err)
	}

	t.Log("Test UpdateBooks PASSED.")

	//Test DeleteBook
	_, err = client.DeleteBook(ctx, &pb.DeleteRequest{BookTitle: updateBook.BookTitle})
	if err != nil {
		t.Error("Test DeleteBooks FAILED!\nerr: ", err)
	}

	t.Log("Test DeleteBooks PASSED.")

}