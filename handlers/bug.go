package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshua152/bug-tracker-backend/data"
)

type BugContextKey struct{}
type InputTypeContextKey struct{}

type Bug struct {
	l    *log.Logger
	conn *pgxpool.Pool
}

const (
	ArrayInputType  = "Array"
	ObjectInputType = "Object"
)

func NewBug(p *pgxpool.Pool, l *log.Logger) *Bug {
	return &Bug{
		l:    l,
		conn: p,
	}
}

func (b *Bug) GetBugs(rw http.ResponseWriter, r *http.Request) {
	b.l.Println("AT GET BUGS")

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
	b.l.Println("AT GET BUG")

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

	var err error

	in := r.Context().Value(InputTypeContextKey{}).(string)
	bugCtxVal := r.Context().Value(BugContextKey{})
	switch in {
	case ObjectInputType:
		err = b.addSingularBug(bugCtxVal)
	case ArrayInputType:
		err = b.addMultipleBugs(bugCtxVal)
	}

	if err != nil {
		http.Error(rw, "Unable to query bug databse", http.StatusBadRequest)
		b.l.Printf("Unable to query bug for ADD: %v\n", err)
		return
	}
}

func (b *Bug) addSingularBug(bugCxtVal interface{}) error {
	b.l.Println("Add singular bug")

	bug := bugCxtVal.(data.Bug)
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

	return err
}

func (b *Bug) addMultipleBugs(bugCxtVal interface{}) error {
	bugs := bugCxtVal.([]data.Bug)

	copyCount, err := b.conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"bug"},
		[]string{"bug_id", "title", "description", "time_amt", "complexity", "project_id"},
		pgx.CopyFromSlice(len(bugs), func(i int) ([]interface{}, error) {
			return []interface{}{
				bugs[i].BugID,
				bugs[i].Title,
				bugs[i].Description,
				bugs[i].TimeAmt,
				bugs[i].Complexity,
				bugs[i].ProjectID,
			}, nil
		}),
	)
	if err != nil {
		return err
	}

	if int(copyCount) != len(bugs) {
		return errors.New("error converting bug array to COPY query")
	}

	return nil
}

func (b *Bug) UpdateBug(rw http.ResponseWriter, r *http.Request) {
	b.l.Println("AT UPDATE BUG")

	vars := mux.Vars(r)
	id := vars["id"]

	bug := r.Context().Value(BugContextKey{}).(data.Bug)
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
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			b.l.Println("Unable to read from request body")
			return
		}

		trimmedBytes := bytes.TrimLeft(bodyBytes, " \t\r\n")
		if len(trimmedBytes) == 0 {
			b.l.Printf("Request body required")
			http.Error(rw, "Request body required", http.StatusBadRequest)
			return
		}

		b.l.Println(string(trimmedBytes[0]))

		switch trimmedBytes[0] {
		case '{':
			r, err = b.validateObject(r, trimmedBytes)
		case '[':
			r, err = b.validateArray(r, trimmedBytes)
		default:
			b.l.Printf("Invalid opening delimiter")
			http.Error(rw, "Invalid opening delimiter. May only be '[' or '{'", http.StatusBadRequest)
			return
		}

		if err != nil {
			b.l.Printf("Unable to deserialize bug: %v\n", err)
			http.Error(rw, "Error reading bug", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(rw, r)
	})
}

/*
Validates a singular bug. Will store the context in the returned request
under the BugContextKey{}. The InputTypeContextKey will be set to the
InputTypeObject constant.
*/
func (b Bug) validateObject(r *http.Request, bytes []byte) (*http.Request, error) {
	bug := data.Bug{}

	err := json.Unmarshal(bytes, &bug)
	if err != nil {
		return r, err
	}

	j, _ := json.MarshalIndent(bug, "", "  ")
	b.l.Printf("Bug: %v\n", string(j))

	ctx := context.WithValue(r.Context(), BugContextKey{}, bug)
	ctx = context.WithValue(ctx, InputTypeContextKey{}, ObjectInputType)
	r = r.WithContext(ctx)

	return r, nil
}

/*
Validates an array of bugs. Will store the context in the returned request
under the BugContextKey{}. The InputTypeContextKey will be set to the
InputTypeArray constant.
*/
func (b Bug) validateArray(r *http.Request, bytes []byte) (*http.Request, error) {
	bugs := []data.Bug{}

	err := json.Unmarshal(bytes, &bugs)
	if err != nil {
		return r, err
	}

	j, _ := json.MarshalIndent(bugs, "", "  ")
	b.l.Printf("Bugs: %v\n", string(j))

	ctx := context.WithValue(r.Context(), BugContextKey{}, bugs)
	ctx = context.WithValue(ctx, InputTypeContextKey{}, ArrayInputType)
	r = r.WithContext(ctx)

	return r, nil
}
