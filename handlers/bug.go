package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshua152/bug-tracker-backend/data"
)

var ContetKey struct{}

type Bug struct {
	l    *log.Logger
	conn *pgxpool.Pool
}

func NewBug(p *pgxpool.Pool, l *log.Logger) *Bug {
	return &Bug{
		l:    l,
		conn: p,
	}
}

func (b *Bug) GetBugs(rw http.ResponseWriter, r *http.Request) {
	b.l.Println("AT GET PROJECTS")

	query := `SELECT * FROM bug ORDER BY project_id, bug_id`
	rows, err := b.conn.Query(context.Background(), query)
	if err != nil {
		http.Error(rw, "Unable to query bug database", http.StatusBadRequest)
		b.l.Printf("Unable to query bugs for GET: %v\n", err)
		return
	}

	var bugs []data.Bug

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			http.Error(rw, "Unable to process data", http.StatusInternalServerError)
			log.Printf("Erorr while iterating database for GET BUGS: %v", err)
			return
		}

		bugs = append(bugs, data.Bug{
			BugID:       vals[0].(int32),
			Title:       vals[1].(string),
			Description: vals[2].(string),
			TimeAmt:     vals[3].(float64),
			Complexity:  vals[4].(float64),
			ProjectID:   vals[5].(int32),
		})
	}

	data.ToJSON(bugs, rw)
}

func (b *Bug) GetBug(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var bug data.Bug

	query := `SELECT * FROM bug WHERE bug_id = $1`
	err := b.conn.QueryRow(context.Background(), query, id).Scan(
		&bug.BugID,
		&bug.Title,
		&bug.Description,
		&bug.TimeAmt,
		&bug.Complexity,
		&bug.ProjectID,
	)
	if err != nil {
		http.Error(rw, "Unable to query bug database", http.StatusBadRequest)
		b.l.Printf("Unable to query bug database for get: %v\n", err)
		return
	}

	data.ToJSON(bug, rw)
}

func (b *Bug) AddBug(rw http.ResponseWriter, r *http.Request) {
	b.l.Println("AT ADD BUG")

	bug := r.Context().Value(ContextKey{}).(data.Bug)
	query := `
		INSERT INTO bug (
			bug_id,
			title,
			description,
			time_amt,
			complexity,
			project_id
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := b.conn.Query(context.Background(), query, bug.BugID, bug.Title,
		bug.Description, bug.TimeAmt, bug.Complexity, bug.ProjectID)
	if err != nil {
		http.Error(rw, "Unable to query bug databse", http.StatusBadRequest)
		b.l.Printf("Unable to query bug for ADD: %v\n", err)
		return
	}
}

func (b *Bug) UpdateBug(rw http.ResponseWriter, r *http.Request) {
	b.l.Println("AT UPDATE BUG")

	vars := mux.Vars(r)
	id := vars["id"]

	bug := r.Context().Value(ContextKey{}).(data.Bug)
	query := `
		UPDATE bug 
		SET (
			bug_id,
			title,
			description,
			time_amt,
			complexity,
			project_id
		) = ($1, $2, $3, $4, $5, $6)
		WHERE bug_id = $7
	`
	_, err := b.conn.Query(context.Background(), query, bug.BugID, bug.Title,
		bug.Description, bug.TimeAmt, bug.Complexity, bug.ProjectID, id)
	if err != nil {
		http.Error(rw, "Unable to query bug database", http.StatusBadRequest)
		b.l.Printf("Unable to query bug for UPDATE: %v\n", err)
		return
	}
}

func (b *Bug) DeleteBug(rw http.ResponseWriter, r *http.Request) {
	b.l.Println("AT DELETE BUG")

	vars := mux.Vars(r)
	id := vars["id"]

	query := `DELETE FROM bug WHERE bug_id = $1`
	_, err := b.conn.Query(context.Background(), query, id)
	if err != nil {
		http.Error(rw, "Unable to query bug databases", http.StatusBadRequest)
		b.l.Printf("Unable to query bug database for DELETE: %v", err)
		return
	}
}

func (b Bug) ValidateBug(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		bug := data.Bug{}

		err := bug.FromJSON(r.Body)
		if err != nil {
			b.l.Printf("Unable to deserialize bug: %v\n", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		b.l.Printf("Bug: %v", bug)

		ctx := context.WithValue(r.Context(), ContextKey{}, bug)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
