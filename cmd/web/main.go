package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	infoLog, errorLog *log.Logger
}

func main() {
	addr := flag.Int("addr", 4000, "HTTP network address")
	static := flag.String("static-dir", "./ui/static", "Static files directory")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	app := &application{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	db, err := openDB(*dsn)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	defer db.Close()

	srv := &http.Server{
		Addr:     fmt.Sprint(":", *addr),
		ErrorLog: app.errorLog,
		Handler:  app.routes(static),
	}

	app.infoLog.Println("Starting server on port", *addr)
	err = srv.ListenAndServe()
	app.errorLog.Fatal("ERROR", err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
