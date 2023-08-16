package handlers

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

func (p *Bug) GetBugs(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("AT GET BUGS")

}

func (p *Bug) GetBug(rw http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	// sid := vars["id"]
}

func (p *Bug) AddBug(rw http.ResponseWriter, r *http.Request) {

}

func (p *Bug) UpdateBug(rw http.ResponseWriter, r *http.Request) {

}

func (p *Bug) DeleteBug(rw http.ResponseWriter, r *http.Request) {
}
