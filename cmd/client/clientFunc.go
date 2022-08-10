package main

import (
	"context"
	"errors"
	"fmt"
	pb "BookManagementService/protoFiles"
	"io"
	"log"
	"strconv"
	"strings"
)

//Create Book
func BookCreate(client pb.BookMgmtServiceClient, ctx context.Context) error {

	//Input
	fmt.Println("Enter Details:  ")

	fmt.Print("Book Title:  ")
	title, err := reader.ReadString('\n')
	if err != nil {
		return errors.New(fmt.Sprint("Invalid Title", err))
	}
	title = strings.TrimSpace(title)

	fmt.Print("Book Author:  ")
	author, err := reader.ReadString('\n')
	if err != nil {
		return errors.New(fmt.Sprint("Invalid Author", err))
	}
	author = strings.TrimSpace(author)

	if title == "" || author == "" {
		return errors.New(("Empty..."))
	}

	//Creating book to send as request
	NewBook := &pb.Book{
		BookTitle:  title,
		BookAuthor: author,
	}
	//Call CreateBook that returns a book as response
	response, err := client.CreateBook(ctx, &pb.CreateRequest{Book: NewBook})
	if err != nil {
		return errors.New(fmt.Sprint("Book Creation Failed... \n", err))
	}
	//print
	log.Printf(`New Book Uploaded:
	Book Id: %d
	Title: %s
	Author: %s`, response.Book.Id, response.Book.BookTitle, response.Book.BookAuthor)

	return nil
}

//Get All Books
func GetBooks(client pb.BookMgmtServiceClient, ctx context.Context) error {

	fmt.Print("Enter page no.: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	input = strings.TrimSpace(input)
	pageNo, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		return errors.New("Invalid page no.")
	}

	fmt.Print("No. of books in page:  ")
	input, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	input = strings.TrimSpace(input)
	NoOfBooks, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		return errors.New("Invalid input...")
	}

	print := (pageNo - 1) * NoOfBooks

	// Call ListAllBook that returns a stream
	stream, err := client.GetAllBooks(ctx, &pb.GetAllRequest{})
	// Check for errors
	if err != nil {
		return errors.New(fmt.Sprint("GetBooks failed to stream ", err))
	}

	fmt.Println("Books: ")

	for i := 0; i < int(pageNo*NoOfBooks); i++ {

		// stream.Recv returns a pointer to a book in a current iteration
		response, err := stream.Recv()
		// If end of stream, break the loop
		if err == io.EOF {
			break
		}

		if err != nil {
			return errors.New(fmt.Sprint("Stream error: ", err))
		}

		if i >= int(print) {

			fmt.Printf("%d %v\n", i+1, response.GetBook())

		}

	}

	return nil
}

//Update Book
func BookUpdate(client pb.BookMgmtServiceClient, ctx context.Context) error {

	fmt.Println("Update--- ")

	fmt.Print(" Book Id: ")
	id, err := reader.ReadString('\n')
	if err != nil {
		return errors.New(fmt.Sprint("Invalid ID..", err))
	}
	id = strings.TrimSpace(id)

	fmt.Print("Book Title: ")
	title, err := reader.ReadString('\n')
	if err != nil {
		return errors.New(fmt.Sprint("Invalid Book Title..", err))
	}
	title = strings.TrimSpace(title)

	fmt.Print("Book Author: ")
	author, err := reader.ReadString('\n')
	if err != nil {
		return errors.New(fmt.Sprint("Invalid Book Author..", err))
	}
	author = strings.TrimSpace(author)

	if id == "" || title == "" || author == "" {
		return errors.New(fmt.Sprint("Empty..."))
	}

	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	//Creating Request
	updateBook := &pb.Book{
		Id:         uint64(num),
		BookTitle:  title,
		BookAuthor: author,
	}
	//Call UpdateBook that returns a Book as response
	response, err := client.UpdateBook(ctx, &pb.UpdateRequest{Book: updateBook})
	if err != nil {
		return errors.New(fmt.Sprint("Could not update book: \n", err))
	}
	//print
	log.Printf(`Book Updated:
	Book Id: %d
	Title: %s
	Author: %s`, response.Book.Id, response.Book.BookTitle, response.Book.BookAuthor)

	return nil
}

//Delete Book
func BookDelete(client pb.BookMgmtServiceClient, ctx context.Context) error {

	fmt.Printf("Enter Book Title u want to delete:  ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	title := strings.TrimSpace(input)

	//Call DeleteBook
	_, err = client.DeleteBook(ctx, &pb.DeleteRequest{BookTitle: title})
	if err != nil {
		return err
	}

	//Print Result
	fmt.Print("\nDeleted book with Book Title: ", title)
	return nil
}

//Search Books
func BookSearch(client pb.BookMgmtServiceClient, ctx context.Context) error {

	var stream pb.BookMgmtService_SearchBookClient

	fmt.Println(`Choose an option...
	
	1. Search by Book Title
	2. Search by Book Author
	
	Choose: `)

	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	input = strings.TrimSpace(input)
	choice, err := strconv.ParseInt(input, 10, 64)

	switch choice {
	case 1:
		fmt.Println("Enter Book Title: ")
		title, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		title = strings.TrimSpace(title)

		if title == "" {
			return errors.New("empty title")
		}

		//Streaming
		stream, err = client.SearchBook(ctx, &pb.SearchRequest{
			Search: &pb.SearchRequest_BookTitle{BookTitle: title},
		})

		if err != nil {
			return err
		}

	case 2:
		fmt.Println("Enter Book Author: ")
		author, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		author = strings.TrimSpace(author)

		if author == "" {
			return errors.New("empty author")
		}

		//Streaming
		stream, err = client.SearchBook(ctx, &pb.SearchRequest{
			Search: &pb.SearchRequest_BookAuthor{BookAuthor: author},
		})

		if err != nil {
			return err
		}

	default:
		return errors.New("search book failed")
	}
	fmt.Println("Books we have found...")

	// Start iterating
	for i := 0; i < 3; i++ {

		// stream.Recv returns a pointer to a book in a current iteration
		responseStream, err := stream.Recv()
		// If end of stream, break the loop
		if err == io.EOF {
			break
		}
		// if err, print error
		if err != nil {
			return errors.New(fmt.Sprint("Stream error: ", err))
		}

		// If everything went well use the generated getter to print the Book Details
		fmt.Println(responseStream.GetBook())

	}
	return nil
}