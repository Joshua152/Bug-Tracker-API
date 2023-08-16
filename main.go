package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/joshua152/bug-tracker-backend/handlers"
)

func main() {
	l := log.New(os.Stdout, "bug-tracker", log.LstdFlags)

	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Unable to load environment variables: %v\n", err)
	}

	// "postgres://joshuaau:postgres@localhost:5432/bugtracker"
	dbURL := os.Getenv("DB_URL")
	conn, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	sm := mux.NewRouter()

	ph := handlers.NewProject(conn, l)
	bh := handlers.NewBug(conn, l)

	// PROJECTS SUB ROUTER

	projects := sm.PathPrefix("/projects").Subrouter()

	// GET
	pGetRouter := projects.Methods(http.MethodGet).Subrouter()
	pGetRouter.HandleFunc("", ph.GetProjects)
	pGetRouter.HandleFunc("/{id:[0-9a-zA-Z]+}", ph.GetProject)
	pGetRouter.HandleFunc("/{id:[0-9a-zA-Z]+}/bugs", ph.GetBugs)

	// PUT
	pPutRouter := projects.Methods(http.MethodPut).Subrouter()
	pPutRouter.HandleFunc("/{id:[0-9a-zA-Z]+}", ph.UpdateProject)
	pPutRouter.Use(ph.ValidateProject)

	// POST
	pPostRouter := projects.Methods(http.MethodPost).Subrouter()
	pPostRouter.HandleFunc("", ph.AddProject)
	pPostRouter.Use(ph.ValidateProject)

	// DELETE
	pDeleteRouter := projects.Methods(http.MethodDelete).Subrouter()
	pDeleteRouter.HandleFunc("/{id:[0-9a-zA-Z]+}", ph.DeleteProject)

	// BUGS SUB ROUTER

	bugs := sm.PathPrefix("/bugs").Subrouter()

	// GET
	bGetRouter := bugs.Methods(http.MethodGet).Subrouter()
	bGetRouter.HandleFunc("", bh.GetBugs)
	bGetRouter.HandleFunc("/{id:[0-9a-zA-z]+}", bh.GetBug)

	// PUT
	bPutRouter := bugs.Methods(http.MethodPut).Subrouter()
	bPutRouter.HandleFunc("/{id:[0-9a-zA-z]+}", bh.UpdateBug)
	bPutRouter.Use(bh.ValidateBug)

	// POST
	bPostRouter := bugs.Methods(http.MethodPost).Subrouter()
	bPostRouter.HandleFunc("", bh.AddBug)
	bPostRouter.Use(bh.ValidateBug)

	// DELETE
	bDeleteRouter := bugs.Methods(http.MethodDelete).Subrouter()
	bDeleteRouter.HandleFunc("/{id:[0-9a-zA-z]+}", bh.DeleteBug)

	s := http.Server{
		Addr:         ":9090",
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Printf("Starting server on port %s\n", s.Addr)

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	sig := <-c
	log.Println("got signal:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(ctx)
}

/*
func main() {
	dbURL := "postgres://joshuaau:postgres@localhost:5432/bugtracker"
	conn, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	query := `
		SELECT
			bug.*,
			project.name
		FROM
			bug
		INNER JOIN project
			ON bug.project_id = project.project_id
		WHERE
			bug.project_id = 1;
	`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}

	fmt.Printf("here\n")

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			log.Fatal("Error while iterating dataset")
		}

		fmt.Printf("(%d) Process row: %v\n", len(vals), vals)

		bugID := vals[0].(int32)
		title := vals[1].(string)
		description := vals[2].(string)
		timeAmt := vals[3].(float64)
		complexity := vals[4].(float64)
		projectID := vals[5].(int32)
		projectName := vals[6].(string)

		fmt.Printf("done")

		_, err = fmt.Printf(`bugID: %v, title: %s, description: %s, timeAmt: %s, complexity: %v,
		projectID: %v, projectName: %s\n`, bugID, title, description, timeAmt, complexity, projectID, projectName)
		if err != nil {
			log.Fatalf("Printf error: %v\n", err)
		}
	}
}
*/
