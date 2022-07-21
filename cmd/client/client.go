package main

import (
	"bufio"
	"context"
	"fmt"
	pb "BookManagementService/protoFiles"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
)

const address = ":54321"

var reader *bufio.Reader = bufio.NewReader(os.Stdin)

func main() {

	//Dial a connection to grpc Server
	connection, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
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

			err := BookCreate(client, ctx)
			if err != nil {
				fmt.Println(err)
			}

			continue
		case 2:

			err := GetBooks(client, ctx)
			if err != nil {
				fmt.Println(err)
			}

			continue
		case 3:

			err := BookSearch(client, ctx)
			if err != nil {
				fmt.Println(err)
			}

			continue
		case 4:

			err := BookUpdate(client, ctx)
			if err != nil {
				fmt.Println(err)
			}

			continue
		case 5:

			err := BookDelete(client, ctx)
			if err != nil {
				fmt.Println(err)
			}

			continue
		default:
			os.Exit(0)
		}
	}

}
