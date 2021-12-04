package git

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Chiorufarewerin/gitchat/internal/environment"
	"github.com/go-git/go-billy/v5"
	memfs "github.com/go-git/go-billy/v5/memfs"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	memory "github.com/go-git/go-git/v5/storage/memory"
)

var (
	auth       *http.BasicAuth
	repository *git.Repository
	worktree   *git.Worktree

	pushConflict       bool
	previousHashCommit plumbing.Hash

	FileDoesNotExists = os.ErrNotExist
	GitPushConflict   = errors.New("git push conflict")
)

func gitGetOrCreateFileToWrite(filePath string) (billy.File, error) {
	file, err := worktree.Filesystem.OpenFile(filePath, os.O_WRONLY, 0666)
	if err != nil {
		file, err = worktree.Filesystem.Create(filePath)
		if err != nil {
			return nil, err
		}
	}
	return file, nil
}

func GitGetFileSize(filePath string) (int64, error) {
	stat, err := worktree.Filesystem.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

func GitAppendFileData(filePath string, data []byte, offset int64) error {
	file, err := gitGetOrCreateFileToWrite(filePath)
	if err != nil {
		return err
	}

	file.Seek(offset, 2)
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	worktree.Add(filePath)

	return nil
}

func GitWriteFileData(filePath string, data []byte) error {
	file, err := gitGetOrCreateFileToWrite(filePath)
	if err != nil {
		return err
	}

	file.Truncate(0)
	file.Seek(0, 0)
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}
	worktree.Add(filePath)

	return nil
}

func GitReadFileData(filePath string) ([]byte, error) {
	file, err := worktree.Filesystem.Open(filePath)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GitCommitAndPush(message string) error {
	head, err := repository.Head()
	if err != nil {
		return err
	}
	previousHashCommit = head.Hash()

	_, err = worktree.Commit(message, &git.CommitOptions{})
	if err != nil {
		return err
	}

	err = repository.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), git.ErrNonFastForwardUpdate.Error()) {
			pushConflict = true
			err = GitPushConflict
		}
		return err
	}

	return nil
}

func GitPull() error {
	err := worktree.Pull(&git.PullOptions{Force: true})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}
	return nil
}

func GitCheckErrors() error {
	if pushConflict {
		head, err := repository.Head()
		if err != nil {
			return err
		}
		if head.Hash() == previousHashCommit {
			log.Fatalln("With push conflict head and previous commit equals")
		}

		err = worktree.Reset(&git.ResetOptions{Commit: previousHashCommit, Mode: git.HardReset})
		if err != nil {
			return err
		}

		pushConflict = false
	}

	status, err := worktree.Status()
	if err != nil {
		return err
	}

	if len(status) > 0 {
		head, err := repository.Head()
		if err != nil {
			return err
		}

		err = worktree.Reset(&git.ResetOptions{Commit: head.Hash(), Mode: git.HardReset})
		if err != nil {
			return err
		}

		err = worktree.Clean(&git.CleanOptions{Dir: true})
		if err != nil {
			return err
		}
	}
	return nil
}

func InitializeGit() {
	var err error
	auth = &http.BasicAuth{
		Username: environment.GitAuthUsername,
		Password: environment.GitAuthPassword,
	}

	storage := memory.NewStorage()
	gitConfig := config.NewConfig()
	gitConfig.User.Name = environment.GitUserName
	gitConfig.User.Email = environment.GitUserEmail
	err = storage.SetConfig(gitConfig)

	repository, err = git.Clone(storage, memfs.New(), &git.CloneOptions{
		URL: environment.GitRepositoryURL,
	})

	if err != nil {
		log.Fatalf("%v", err)
	}

	worktree, err = repository.Worktree()
	if err != nil {
		log.Fatalf("%v", err)
	}
}
