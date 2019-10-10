module github.com/growerlab/codev-svc

go 1.13

require (
	github.com/gin-gonic/gin v1.4.0
	github.com/graphql-go/graphql v0.7.8
	github.com/graphql-go/handler v0.2.3
	github.com/joho/godotenv v1.3.0
	github.com/rs/zerolog v1.15.0
)

replace github.com/libgit2/git2go => ./vendor/github.com/libgit2/git2go
