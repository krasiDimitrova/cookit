package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/krasimiraMilkova/cookit/internal/appconfig"
	"log"
	"time"
)

var db *sql.DB

// Get is used to provide connection to the db following the singleton pattern
func Get() *sql.DB {
	if db == nil {
		initDB()
	}
	return db
}

func initDB() {
	db = connect()

	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(4)
	db.SetConnMaxLifetime(time.Second * 15)
}

func connect() *sql.DB {
	config := appconfig.Get()
	username, password, databaseName, databaseHost := config.GetDBConfig()

	dbURI := fmt.Sprintf("%s:%s@(%s)/", username, password, databaseHost)
	db, err := createAndOpen(databaseName, dbURI)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return db
}

func createAndOpen(name string, dbURI string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		return nil, err
	}

	db, err = sql.Open("mysql", dbURI+name)
	if err != nil {
		return nil, err
	}

	err = createUsersTable(db)
	if err != nil {
		return nil, err
	}

	err = createRecipesTable(db)
	if err != nil {
		return nil, err
	}

	err = createIngredientsTable(db)
	if err != nil {
		return nil, err
	}

	err = createIngredientsForRecipeTable(db)
	if err != nil {
		return nil, err
	}

	err = createCommentsTable(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createUsersTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
						id int NOT NULL AUTO_INCREMENT,
						name varchar(100) NOT NULL,
						email varchar(100) NOT NULL,
						password varchar(100) NOT NULL,
						PRIMARY KEY (id)
					);`)
	return err
}

func createRecipesTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS recipes (
						id int NOT NULL AUTO_INCREMENT,
						title varchar(100) NOT NULL,
						directions text NOT NULL,
						PRIMARY KEY (id)
					);`)
	return err
}

func createIngredientsTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS ingredients (
						id int NOT NULL AUTO_INCREMENT,
						name varchar(100) NOT NULL,
						PRIMARY KEY (id),
						UNIQUE (name)
					);`)
	return err
}

func createIngredientsForRecipeTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS recipe_ingredients (
						recipe_id int NOT NULL,
						ingredient_id int NOT NULL,
						quantity int NOT NULL,
						measurement varchar(10),
						PRIMARY KEY (recipe_id, ingredient_id),
						FOREIGN KEY (recipe_id) 
						    REFERENCES recipes(id)
                        	ON DELETE CASCADE
                            ON UPDATE CASCADE,
						FOREIGN KEY (ingredient_id) 
							REFERENCES ingredients(id)
							ON DELETE CASCADE
							ON UPDATE CASCADE
					);`)
	return err
}

func createCommentsTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS comments (
						id int NOT NULL AUTO_INCREMENT,
						recipe_id int NOT NULL,
						comment text NOT NULL,
						PRIMARY KEY (id),
						FOREIGN KEY (recipe_id) 
							REFERENCES recipes(id)
							ON DELETE CASCADE
							ON UPDATE CASCADE       
					);`)
	return err
}
