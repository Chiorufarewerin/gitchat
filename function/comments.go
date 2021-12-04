package function

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Chiorufarewerin/gitchat/internal/environment"
	"github.com/Chiorufarewerin/gitchat/internal/git"
)

func getCommentConfig() (*CommentConfig, error) {
	fileData, err := git.GitReadFileData(environment.CommentsConfigFilePath)
	if err != nil {
		if err == git.FileDoesNotExists {
			config := &CommentConfig{
				Version:    environment.CommentsConfigVersion,
				FirstBlock: 1,
				LastBlock:  1,
			}
			if err := writeCommentConfig(config); err != nil {
				return nil, err
			}
			return config, nil
		}
		return nil, err
	}

	commentConfig := &CommentConfig{}
	if err := json.Unmarshal(fileData, commentConfig); err != nil {
		return nil, err
	}

	return commentConfig, nil
}

func writeCommentConfig(config *CommentConfig) error {
	version := MakeJsonField("version", config.Version)
	firstBlock := MakeJsonField("firstBlock", config.FirstBlock)
	lastBlock := MakeJsonField("lastBlock", config.LastBlock)

	strFmt := fmt.Sprintf("{\n%s%s}\n", strings.Repeat("  %s,\n", 2), "  %s\n")
	data := []byte(fmt.Sprintf(strFmt, version, firstBlock, lastBlock))
	return git.GitWriteFileData(environment.CommentsConfigFilePath, data)
}

func generateCommentID(comment *Comment, config *CommentConfig) string {
	return fmt.Sprintf("%s_%d_%s", config.Version, config.LastBlock, GenerateUniqueId())
}

func writeComment(comment *Comment, config *CommentConfig, append bool) error {
	id := MakeJsonField("id", comment.ID)
	author := MakeJsonField("author", comment.Author)
	text := MakeJsonField("text", comment.Text)
	color := MakeJsonField("color", comment.Color)
	date := MakeJsonField("date", comment.Date)
	reply := MakeJsonField("reply", comment.Reply)

	strFmt := fmt.Sprintf(",\n  {\n%s%s  }\n]\n", strings.Repeat("    %s,\n", 5), "    %s\n")
	offset := int64(-3)
	if !append {
		strFmt = "[" + strFmt[1:]
		offset = 0
	}
	data := []byte(fmt.Sprintf(strFmt, id, author, text, color, date, reply))

	lastBlockPath := fmt.Sprintf("%s%d.json", environment.CommentsBlocksFolderPath, config.LastBlock)
	return git.GitAppendFileData(lastBlockPath, data, offset)
}

func addCommentAttempts(comment *Comment, attempt int) (*Comment, error) {
	var err error

	defer func() {
		git.GitCheckErrors()
		if err == git.GitPushConflict && attempt < 1 {
			addCommentAttempts(comment, attempt+1)
		}
	}()

	err = git.GitPull()
	if err != nil {
		return nil, err
	}

	config, err := getCommentConfig()
	if err != nil {
		return nil, err
	}

	lastBlockPath := fmt.Sprintf("%s%d.json", environment.CommentsBlocksFolderPath, config.LastBlock)
	lastBlockFileSize, err := git.GitGetFileSize(lastBlockPath)
	if err != nil && err != git.FileDoesNotExists {
		return nil, err
	}

	if lastBlockFileSize > environment.CommentsNextSizeBlock {
		config.LastBlock += 1
		lastBlockFileSize = 0
		err = writeCommentConfig(config)
		if err != nil {
			return nil, err
		}
	}

	comment.ID = generateCommentID(comment, config)
	comment.Date = GetCurrentDateUTCString()
	err = writeComment(comment, config, lastBlockFileSize != 0)
	if err != nil {
		return nil, err
	}

	err = git.GitCommitAndPush("Add comment: " + comment.ID)

	return comment, err
}

func AddComment(comment *Comment) (*Comment, error) {
	return addCommentAttempts(comment, 0)
}
