package database

import (
	"context"
	"errors"
	"fmt"

	"BookManagementService/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

const DNS = "root:BbSs@563884(localhost:3306)/bookservicedb"

func InitialMigration() {
	DB, err := gorm.Open(mysql.Open(DNS), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	}
	DB.AutoMigrate(&model.Book{})

}

//Create Book
func CreateBook(ctx context.Context, newBook model.Book) (uint, error) {
	user := model.Book{Title: newBook.Title, Author: newBook.Author}
	create := DB.Create(&user)
	if create.Error != nil {
		return 0, create.Error
	}

	return user.ID, nil
}

//Update Book
func UpdateBook(ctx context.Context, updateBook model.Book, id uint64) (uint, error) {

	var book model.Book

	err := DB.Where("ID=?", id).Find(&book).Error
	//(&model.Book{T}).Find(&book).Error
	if err != nil {
		return 0, err
	}

	err = DB.Model(&book).Updates(updateBook).Error
	if err != nil {
		return 0, err
	}

	return book.ID, nil
}

//Delete Book
func DeleteBook(s string) error {
	tx := DB.Delete(s)
	if tx != nil {
		return tx.Error
	}

	return nil
}

//GetAll Books
func GetAllBooks() ([]model.Book, error) {
	var listOfBooks []model.Book
	err := DB.Find(&listOfBooks).Error
	if err != nil {
		return nil, err
	}

	if len(listOfBooks) == 0 {
		return nil, errors.New("empty list")
	}

	return listOfBooks, nil
}

// Search Books
func SearchBook(title string, author string) ([]model.Book, error) {

	if title == "" && author == "" {
		return nil, errors.New("nothing to search, empty argment")
	}

	if title == "" {
		err := DB.Find(&author).Error
		if err != nil {
			return nil, err
		}
	}
	if author == "" {
		err := DB.First(&title).Error
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("Invalid")

}