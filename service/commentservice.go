package service

import "errors"

type Commentinterf interface {
	CekInputComment(message string) error
}

type CommentService struct{}

func CommentServic() Commentinterf {
	return &CommentService{}
}

func (cs *CommentService) CekInputComment(message string) error {
	if message == "" {
		return errors.New("message cannot be empty")
	}
	return nil
}
