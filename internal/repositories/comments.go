package repositories

import "github.com/jackc/pgx/v5/pgxpool"

// Comments table name in db
const COMMENTS_TABLE = "comments"

type CommentsRepository struct {
	_  ContentRepository
	Db *pgxpool.Pool
}

var Comments *CommentsRepository

func (r *CommentsRepository) Create() {

}
