package git

import (
	"guard_rails/client"
	"guard_rails/config"
	"guard_rails/model"
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type gitClientProvider struct {
	config *config.GitConfig
}

func NewGitCleintProvider(config *config.GitConfig) client.GitClientProvider {
	return &gitClientProvider{
		config: config,
	}
}

func (gcp *gitClientProvider) NewGitClient() client.GitClient {
	return &gitClient{
		gitClientProvider: gcp,
	}
}

type gitClient struct {
	*gitClientProvider
	repository *git.Repository
	treeIter   *object.TreeIter
	tree       *object.Tree
	filesIter  *object.FileIter
}

func (gc *gitClient) Clone(repository *model.Repository) error {
	var err error

	path := gc.config.CloneLocation + "/" + repository.Name

	// clone repository
	gc.repository, err = git.PlainClone(path, false, &git.CloneOptions{
		URL: repository.Url,
	})
	if err != nil {
		if err != git.ErrRepositoryAlreadyExists {
			return err
		}

		// if repository already cloned, open locally
		gc.repository, err = git.PlainOpen(path)
		if err != nil {
			return err
		}
	}

	err = gc.initalizeRepositoryIteration()
	if err != nil {
		return err
	}

	return nil
}

func (gc *gitClient) initalizeRepositoryIteration() error {
	var err error

	gc.treeIter, err = gc.repository.TreeObjects()
	if err != nil {
		return err
	}

	gc.tree, err = gc.treeIter.Next()
	if err != nil {
		return err
	}

	gc.filesIter = gc.tree.Files()

	return nil
}

func (gc *gitClient) GetNextFile() (client.File, error) {

	file, err := gc.filesIter.Next()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	// final file in tree
	if file == nil {

		// go to next tree
		gc.treeIter.Close()
		gc.tree, err = gc.treeIter.Next()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}

		// end of the tree
		if gc.tree == nil {
			return nil, nil
		}

		// get next file
		gc.filesIter.Close()
		gc.filesIter = gc.tree.Files()

		file, err = gc.filesIter.Next()
		if err != nil {
			return nil, err
		}

	}

	return newFile(file), nil
}
