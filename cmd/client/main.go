package main

import (
	pb "BookManagementService/protoFiles"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const address = ":55001"

var reader *bufio.Reader = bufio.NewReader(os.Stdin)

func main() {

	//Dial a connection to grpc Server
	//connection, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	connection, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed Dial... ", err)
	}
	defer connection.Close()

	//Create new Client
	client := pb.NewBookMgmtServiceClient(connection)

	//Context Initialisation
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	//Input
	for {
		fmt.Println(`Welcome...
	 
		**MENU**
		1. Upload a New Book
		2. Get All books 
		3. Search a Book
		4. Update Book
		5. Delete Book 
		Choose other number to exit!`)

		fmt.Printf("Choose Option: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Input Failed...", err)
		}
		option, err := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
		if err != nil {
			log.Fatal("Failed Coversion String to int")
		}

		switch option {
		case 1:

			err := CreateBook(client, ctx)
			if err != nil {
				fmt.Println(err)
			}

			continue
		case 2:

			err := GetAllBooks(client, ctx)
			if err != nil {
				fmt.Println(err)
			}

			continue
		case 3:

			err := SearchBook(client, ctx)
			if err != nil {
				fmt.Println(err)
			}

			continue
		case 4:

			err := UpdateBook(client, ctx)
			if err != nil {
				fmt.Println(err)
			}

			continue
		case 5:

			err := DeleteBook(client, ctx)
			if err != nil {
				fmt.Println(err)
			}

			continue
		default:
			os.Exit(0)
		}
	}

}
