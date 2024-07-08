package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"
func main() {
	//connect to database
	db := initDB()
	db.Ping()

	//create sessions
	session := initSession()

	//create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate | log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile)

	//create channels

	//create wait groups
	wg := sync.WaitGroup{}

	//setup application config
	app := Config{
		Session: session,
		DB: db,
		Wait: &wg,
		InfoLog: infoLog,
		ErrorLog: errorLog,

	}

	//set up mail

	//listen for web connections
	app.serve()
}

func (app *Config) serve() {
	//start http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	app.InfoLog.Println("Starting web server...")
	err := srv.ListenAndServe()
	
	if err != nil {
		log.Panic("Failed to start a server ", err)
	}

}


func initDB () *sql.DB {
	conn := connectToDB()

	if conn == nil {
		log.Panic("Cannot connect to database")
	}
	return conn
}

func connectToDB () *sql.DB {
	retryCount := 0

	//connection string
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready yet....")
			retryCount++
		} else {
			log.Println("Connected to database")
			return connection;
		}
		if retryCount > 10 {
			return nil
		}
		log.Println("Backing up for 1 second")
		time.Sleep(1 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func initSession() *scs.SessionManager {
	 session := scs.New()
	 session.Store = redisstore.New(initRedis())

	 session.Lifetime = 24 * time.Hour
	 session.Cookie.Persist = true
	 session.Cookie.SameSite = http.SameSiteLaxMode
	 session.Cookie.Secure = true

	 return session
}

func initRedis() *redis.Pool {
	//pool of redis connections
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial : func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}
	return redisPool
}
