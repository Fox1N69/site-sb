package manager

import (
	"sync"

	"site-sb/infra"
	"site-sb/internal/repo"
)

type RepoManager interface {
	AuthRepo() repo.AuthRepo
	PostRepo() repo.PostRepo
}

type repoManager struct {
	infra infra.Infra
}

func NewRepoManager(infra infra.Infra) RepoManager {
	return &repoManager{infra: infra}
}

var (
	authRepoOnce sync.Once
	authRepo     repo.AuthRepo
)

func (rm *repoManager) AuthRepo() repo.AuthRepo {
	authRepoOnce.Do(func() {
		authRepo = repo.NewAuthRepo(rm.infra.GormDB())
	})

	return authRepo
}

var (
	postRepoOnce sync.Once
	postRepo     repo.PostRepo
)

func (rm *repoManager) PostRepo() repo.PostRepo {
	postRepoOnce.Do(func() {
		postRepo = repo.NewPostRepo(rm.infra.GormDB())
	})

	return postRepo
}
