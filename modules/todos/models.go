package todos

import (
	"bytes"
	"fmt"
	"time"
)

type Todo struct {
	ID          int64
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	StatusID    int64
	Image       bytes.Buffer
}

func (t *Todo) ToDtoHttpResponse() *TodoDtoHttpResponse {
	var imageUrl string
	if t.Image.Len() > 0 {
		imageUrl = fmt.Sprintf("/todos/image/%d", t.ID)
	}
	return &TodoDtoHttpResponse{
		t.ID, t.Title, t.Description, t.CreatedAt, t.UpdatedAt, t.StatusID, imageUrl,
	}
}

type TodoDtoHttpResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	StatusID    int64     `json:"statusId"`
	ImageUrl    string    `json:"imageUrl"`
}

type StatusTodo struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	UserId    int64     `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
