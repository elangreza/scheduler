package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/elangreza/scheduler/config"
	"github.com/elangreza/scheduler/internal/rest"
	"github.com/elangreza/scheduler/internal/service"
	"github.com/elangreza/scheduler/internal/sqliterepo"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// if file db not exist create new one
	if !fileExists(cfg.DBFile) {
		if err := createDBFile(cfg.DBFile); err != nil {
			log.Fatal(err)
		}
	}

	db, err := sqliterepo.NewSql(cfg.DBFile)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := migrateDB(db); err != nil {
		log.Fatal(err)
	}

	taskRepo := sqliterepo.NewTaskRepository(db)
	schedulerService := service.NewTaskService(taskRepo)
	handler := rest.NewHandler(schedulerService)

	http.HandleFunc("/", handler.RootHandler)
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListTaskHandler(w, r)
		case http.MethodDelete:
			handler.DeleteTaskHandler(w, r)
		case http.MethodPut, http.MethodPatch:
			handler.UpdateTaskHandler(w, r)
		case http.MethodPost:
			handler.CreateTask(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server started at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)

	// to := []string{"babehracing14@gmail.com"}
	// cc := []string{}
	// subject := "Test mail"
	// message := "Hello"

	// err = sendMail(cfg, to, cc, subject, message)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func createDBFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func migrateDB(db *sql.DB) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatal(err)
	}
	migrationsPath := "file://./migrations"
	if _, err := os.Stat("./migrations"); os.IsNotExist(err) {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"sqlite3", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
