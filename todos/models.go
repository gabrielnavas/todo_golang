package todos

import "time"

type Todo struct {
	ID          int64     `json: "id"`
	Title       string    `json: "title"`
	Description string    `json: "description"`
	CreatedAt   time.Time `json: "createdAt"`
	UpdatedAt   time.Time `json: "updatedAt"`
	StatusID    int64     `json: "statusId"`
}

type StatusTodo struct {
	ID        int64     `json: "id"`
	Name      string    `json: "name"`
	CreatedAt time.Time `json: "createdAt"`
	UpdatedAt time.Time `json: "updateAt"`
}
