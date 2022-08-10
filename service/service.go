package service

import (
	"BookManagementService/model"
	pb "BookManagementService/protoFiles"
	"context"
	"errors"
	"reflect"

	"BookManagementService/database"
)

//Implemetation of proto interfaces
type BookMgmtServiceServer struct {
	pb.UnimplementedBookMgmtServiceServer
}

//Create Book
func (s *BookMgmtServiceServer) CreateBook(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	book := req.Book

	//request Validation
	if book.BookTitle == "" || book.BookAuthor == "" {
		return nil, errors.New("invalid Data")
	}

	bookData := model.Book{
		Title:  book.BookTitle,
		Author: book.BookAuthor,
	}

	id, err := database.CreateBook(ctx, bookData)

	if err != nil {
		return nil, err
	}

	book.Id = uint64(id)

	return &pb.CreateResponse{Book: &pb.Book{Id: book.Id, BookTitle: book.BookTitle, BookAuthor: book.BookAuthor}}, nil
}

//Update Book
func (s *BookMgmtServiceServer) UpdateBook(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	bookUpdates := req.GetBook()
	//oldTitle := req.Title


	//request Validation
	if bookUpdates.BookTitle == "" || bookUpdates.BookAuthor == "" || reflect.ValueOf(bookUpdates.Id) == reflect.Zero(reflect.TypeOf(bookUpdates.Id)) {
		return nil, errors.New("invalid Data")
	}

	bookData := model.Book{
		Title:  bookUpdates.BookTitle,
		Author: bookUpdates.BookAuthor,
	}

	id, err := database.UpdateBook(ctx, bookData, bookUpdates.Id)
	if err != nil {
		return nil, err
	}

	bookUpdates.Id = uint64(id)

	return &pb.UpdateResponse{Book: &pb.Book{Id: bookUpdates.Id, BookTitle: bookUpdates.BookTitle, BookAuthor: bookUpdates.BookAuthor}}, nil
}

//Delete Book
func (s *BookMgmtServiceServer) DeleteBook(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	//Validation
	del := req.BookTitle
	if del == "" {
		return nil, errors.New("invalid Delete Request")
	}

	err := database.DeleteBook(del)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteResponse{Delete: true}, nil
}

//Get All Books
func (s *BookMgmtServiceServer) GetAllBooks(req *pb.GetAllRequest, stream pb.BookMgmtService_GetAllBooksServer) error {

	list, err := database.GetAllBooks()
	if err != nil {
		return err

	}

	for _, book := range list {
		err = stream.Send(&pb.GetAllResponse{Book: &pb.Book{Id: uint64(book.ID), BookTitle: book.Title, BookAuthor: book.Author}})

		if err != nil {
			return err
		}
	}

	return nil

}

//Search Book
func (s *BookMgmtServiceServer) SearchBook(req *pb.SearchRequest, stream pb.BookMgmtService_SearchBookServer) error {

	title := req.GetBookTitle()
	author := req.GetBookAuthor()

	if title == "" && author == "" {
		return errors.New("nothing to search, empty argment")
	}

	search, err := database.SearchBook(title, author)

	if err != nil {
		return err
	}

	for _, book := range search {
		err = stream.Send(&pb.SearchResponse{Book: &pb.Book{Id: uint64(book.ID), BookTitle: book.Title, BookAuthor: book.Author}})

		if err != nil {
			return err
		}
	}

	return nil
}