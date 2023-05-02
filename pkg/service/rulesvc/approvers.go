package rulesvc

import (
	"context"
	"sort"
	"sync"

	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"golang.org/x/sync/errgroup"
)

// userMap holds a set of de-duplicated users.
// It is safe for concurrent use.
type userMap struct {
	mu    sync.RWMutex
	users map[string]bool
}

func newUserMap() *userMap {
	return &userMap{
		users: make(map[string]bool),
	}
}

func (u *userMap) Add(user string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.users[user] = true
}

// All returns a list of all users.
// The list is sorted.
func (u *userMap) All() []string {
	u.mu.Lock()
	defer u.mu.Unlock()

	res := []string{}
	for user := range u.users {
		res = append(res, user)
	}
	sort.Strings(res)
	return res
}

// GetApprovers gets all the approvers for a rule, both those assigned as individuals and those
// assigned via a group. It de-duplicates users, so if a user is assigned as an approver through
// multiple groups they'll only be returned once.
func (s *Service) GetApprovers(ctx context.Context, rule rule.AccessRule) ([]string, error) {
	users := newUserMap()

	for _, u := range rule.Approval.Users {
		users.Add(u)
	}

	wg, gctx := errgroup.WithContext(ctx)
	for _, g := range rule.Approval.Groups {
		id := g
		wg.Go(func() error {
			q := &storage.GetGroup{ID: id}
			_, err := s.DB.Query(gctx, q)
			if err != nil {
				return err
			}
			for _, u := range q.Result.Users {
				users.Add(u)
			}
			return nil
		})
	}
	err := wg.Wait()
	if err != nil {
		return nil, err
	}

	res := users.All()
	return res, nil
}
