package main

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const (
	host     = "192.168.85.36"
	port     = 30326
	user     = "divvet"
	password = "72139aa53dc9c91ed29cffeedc9dce8ad1ff71485786d38349dbd74bbb29b714bfcde645beb8"
	dbname   = "productsdb"
)

func main() {
	passphrase := os.Args[1]
	log.Println("Passphrase: ", passphrase)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, decrypt(password, passphrase), dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	// Create an empty user and make the sql query (using $1 for the parameter)
	//var product Product
	//sql := "SELECT id, brand_name, description, product_date FROM product"
	sql := "SELECT * FROM product"

	//var id int32

	rows, err := db.Query(sql)
	defer rows.Close()
	//err = db.QueryRow(sql, 1).Scan(&product.id, &product.brand_name, &product.description, &product.product_date)
	if err != nil {
		log.Fatal("Failed to query product table: ", err)
	}

	for rows.Next() {
		product := new(Product)
		//if err := rows.Scan(&product.id, &product.brand_name, &product.description, &product.product_date); err != nil {
		if err := rows.Scan(&product.id, &product.brand_name, &product.description, &product.internalName, &product.primaryCategoryId, &product.product_date, &product.productLabel, &product.productName, &product.productPrice, &product.productType, &product.variantProduct, &product.virtualProduct, &product.virtualVariantMethod); err != nil {
			panic(err)
		}
		if product.productName.Valid {
			fmt.Printf("Product name: ", product.description)
		}
		fmt.Printf("\n")
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
}

func decrypt(encryptedString string, keyString string) (decryptedString string) {

	//key, _ := hex.DecodeString(keyString)
	key := []byte(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}
