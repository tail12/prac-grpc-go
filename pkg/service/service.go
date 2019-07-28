package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tail12/prac-grpc-go/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	db *sql.DB
}

// NewBookServiceServer serverにdbを登録する
func NewBookServiceServer(db *sql.DB) api.BookServiceServer {
	return &server{db: db}
}

func (s *server) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "database connect error "+err.Error())
	}
	return c, nil
}

func (s *server) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "INSERT INTO books(`title`,`author`, `description`, `pages`, `price`) VALUES(?, ?, ?, ?, ?)",
		req.Book.Title, req.Book.Author, req.Book.Description, req.Book.Pages, req.Book.Price)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert "+err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id "+err.Error())
	}
	return &api.CreateResponse{
		Id: id,
	}, nil
}
func (s *server) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT `id`, `title`, `author`, `description`, `pages`, `price` FROM books WHERE `id`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ID='%d' is not found",
			req.Id))
	}

	var bk api.Book
	if err := rows.Scan(&bk.Id, &bk.Title, &bk.Author, &bk.Description, &bk.Pages, &bk.Price); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &api.GetResponse{
		Book: &bk,
	}, nil

}

func (s *server) Update(ctx context.Context, req *api.UpdateRequest) (*api.UpdateResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "UPDATE books SET `title`=?, `author`=?, `description`=?, `pages`=?, `price`=?  WHERE `id`=?",
		req.Book.Title, req.Book.Author, req.Book.Description, req.Book.Pages, req.Book.Price, req.Book.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ID='%d' is not found",
			req.Book.Id))
	}

	return &api.UpdateResponse{
		Updated: rows,
	}, nil

}

func (s *server) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "DELETE FROM books WHERE `id`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ID='%d' is not found",
			req.Id))
	}
	return &api.DeleteResponse{
		Deleted: rows,
	}, nil
}

func (s *server) GetAll(ctx context.Context, req *api.GetAllRequest) (*api.GetAllResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT `id`, `title`, `author`, `description`, `pages`, `price` FROM books")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select "+err.Error())
	}
	defer rows.Close()

	list := []*api.Book{}
	for rows.Next() {
		bk := new(api.Book)
		if err := rows.Scan(&bk.Id, &bk.Title, &bk.Author, &bk.Description, &bk.Pages, &bk.Price); err != nil {
			return nil, status.Error(codes.Unknown, err.Error())
		}
		list = append(list, bk)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &api.GetAllResponse{
		Books: list,
	}, nil
}
