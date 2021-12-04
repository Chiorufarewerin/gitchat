package environment

import "os"

const (
	Addr = ":3333"

	GitRepositoryURL = "https://github.com/Chiorufarewerin/chiorufarewerin.github.io"

	GitUserName  = "commentator"
	GitUserEmail = "commentator@example.com"

	DateFormat = "2006-01-02 15:04:05"

	CommentsFolderPath       = "data/comments/v1/"
	CommentsConfigFilePath   = CommentsFolderPath + "config.json"
	CommentsBlocksFolderPath = CommentsFolderPath + "blocks/"
	CommentsNextSizeBlock    = 1024 * 128
	CommentsConfigVersion    = "v1"
)

var (
	GitAuthUsername = os.Getenv("GIT_AUTH_USERNAME")
	GitAuthPassword = os.Getenv("GIT_AUTH_PASSWORD")
)
