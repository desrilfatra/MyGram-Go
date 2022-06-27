package middleware

import (
	"context"
	"fmt"
	"mygram-go/entity"
	"mygram-go/service"
	"net/http"
	"strings"
)

var tempKey = &tempcontext{"user"}

type tempcontext struct {
	data string
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		servic := service.NewUserService()
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]
		// fmt.Println(reqToken)
		fmt.Println("Auth token")
		temp_id, err := servic.VerifToken(reqToken)
		if err != nil {
			w.Write([]byte(fmt.Sprint(err)))
		}
		fmt.Print(" id: ")
		fmt.Println(temp_id)
		user := entity.User{Id: int(temp_id)}

		ctx := context.WithValue(r.Context(), tempKey, &user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})
}

func RunUser(ctx context.Context) *entity.User {
	rawData, _ := ctx.Value(tempKey).(*entity.User)
	return rawData
}
