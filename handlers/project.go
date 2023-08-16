package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshua152/bug-tracker-backend/data"
)

type ContextKey struct{}

type Project struct {
	l    *log.Logger
	conn *pgxpool.Pool
}

func NewProject(p *pgxpool.Pool, l *log.Logger) *Project {
	return &Project{
		l:    l,
		conn: p,
	}
}

func (p *Project) GetProjects(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("AT GET PROJECTS")

	query := "SELECT * FROM project ORDER BY project_id"

	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		http.Error(rw, "Unable to query project database", http.StatusBadRequest)
		log.Printf("GET PROJECTS query failed: %v\n", err)
		return
	}

	var projects []data.Project

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			http.Error(rw, "Unable to process data", http.StatusInternalServerError)
			log.Printf("GET PROJECTS error while iterating dataset: %v\n", err)
			return
		}

		projects = append(projects, data.Project{
			ProjectID: vals[0].(int32),
			Name:      vals[1].(string),
		})
	}

	data.ToJSON(projects, rw)
}

func (p *Project) GetProject(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("AT GET PROJECT")

	vars := mux.Vars(r)
	id := vars["id"]

	var projectID int32
	var name string

	query := "SELECT * FROM project WHERE project_id=$1"
	err := p.conn.QueryRow(context.Background(), query, id).Scan(&projectID, &name)
	if err != nil {
		http.Error(rw, "Unable to query project database", http.StatusBadRequest)
		p.l.Printf("Unable to query project for GET: %v", err)
		return
	}

	project := data.Project{
		ProjectID: projectID,
		Name:      name,
	}

	data.ToJSON(project, rw)
}

func (p *Project) GetBugs(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("AT GET PROJECT BUGS")

	vars := mux.Vars(r)
	id := vars["id"]

	query := `
		SELECT 
			bug.bug_id,
			bug.title,
			bug.description,
			bug.time_amt,
			bug.complexity,
			bug.project_id
		FROM 
			bug
		WHERE bug.project_id = $1
	`
	rows, err := p.conn.Query(context.Background(), query, id)
	if err != nil {
		http.Error(rw, "Unable to query bug database", http.StatusBadRequest)
		p.l.Printf("Unable to query bug database for project bugs GET: %v", err)
		return
	}

	var bugs []data.Bug

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			http.Error(rw, "Unable to process data", http.StatusInternalServerError)
			log.Printf("GET BUGS error while iterating database: %v\n", err)
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

func (p *Project) AddProject(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("AT ADD PROJECT")

	project := r.Context().Value(ContextKey{}).(data.Project)
	query := `INSERT INTO project (project_id, name) VALUES($1, $2)`
	_, err := p.conn.Query(context.Background(), query, project.ProjectID, project.Name)
	if err != nil {
		http.Error(rw, "Unable to query project database", http.StatusBadRequest)
		p.l.Printf("Unable to query project for ADD: %v", err)
		return
	}
}

func (p *Project) UpdateProject(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("AT UPDATE PROJECT")

	vars := mux.Vars(r)
	id := vars["id"]

	project := r.Context().Value(ContextKey{}).(data.Project)
	query := `
		UPDATE project
		SET (project_id, name) = ($1, $2)
		WHERE project_id = $3;
	`
	_, err := p.conn.Query(context.Background(), query, project.ProjectID, project.Name, id)
	if err != nil {
		http.Error(rw, "Unable to query project databse", http.StatusBadRequest)
		p.l.Printf("Unable to query project for UPDATE: %v", err)
		return
	}
}

func (p *Project) DeleteProject(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("AT DELETE PROJECT")

	vars := mux.Vars(r)
	id := vars["id"]

	query := `DELETE FROM project WHERE project_id = $1`
	_, err := p.conn.Query(context.Background(), query, id)
	if err != nil {
		http.Error(rw, "Unable to query project database", http.StatusBadRequest)
		p.l.Printf("Unable to query project for DELETE: %v", err)
		return
	}
}

func (p Project) ValidateProject(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		project := data.Project{}

		err := project.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] DESERIALIZING PRODUCT", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		p.l.Printf("Project: %v\n", project)

		ctx := context.WithValue(r.Context(), ContextKey{}, project)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
