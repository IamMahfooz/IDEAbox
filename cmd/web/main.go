// DO SWIPE OVER THE PACKAGES TO READ THE DOCUMENTATION AND UNDERSTAND THE ARGS.
package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/leandroveronezi/go-recognizer"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetbox/pkg/models"
	"snippetbox/pkg/models/mysql"
	"time"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

// dependencies created here
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets interface {
		Insert(string, string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}
	templateCache map[string]*template.Template
	session       *sessions.Session
	users         interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
	rec        recognizer.Recognizer
	studentMap map[string]map[string]interface{}
	name       string
}

func main() {
	rec := recognizer.Recognizer{}
	studentMap := make(map[string]map[string]interface{})

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	dsn := flag.String("dsn", "web:mahfooz@/snippetbox?parseTime=true", "MySQL data source name")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	//entering  of dependencies !
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		session:       session,
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
		rec:           rec,
		studentMap:    studentMap,
	}
	app.dataTrain()
	app.jsonToMap()

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	//err = srv.ListenAndServe()
	//it now processes middleware(mux) not just mux everytime
	errorLog.Fatal(err)
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
