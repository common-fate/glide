package middleware

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/common-fate/clio"
	"github.com/common-fate/clio/clierr"
	"github.com/urfave/cli/v2"
)

// RequireCleanGitWorktree checks if this is a git repo and if so, checks that the worktree is clean.
// this ensures that users working with a deployment config in a repo always commit their changes prior to deploying.
//
// This method calls out to git if it is installed on the users system.
// Unfortunately, the go library go-git is very slow when checking status.
// https://github.com/go-git/go-git/issues/181
// So this command uses the git cli directly.
// assumption is if a user is using a repository, they will have git installed
func RequireCleanGitWorktree() cli.BeforeFunc {
	return func(c *cli.Context) error {
		if !c.Bool("ignore-git-dirty") {
			_, err := os.Stat(".git")
			if os.IsNotExist(err) {
				// not a git repo, skip check
				return nil
			}
			if err != nil {
				return clierr.New(err.Error(), clierr.Infof("The above error occurred while checking if this is a git repo.\nTo silence this warning, add the 'ignore-git-dirty' flag e.g 'gdeploy --ignore-git-dirty %s'", c.Command.Name))
			}
			_, err = exec.LookPath("git")
			if err != nil {
				// ignore check if git is not installed
				clio.Debugf("could not find 'git' when trying to check if repository is clean. err: %s", err)
				return nil
			}
			cmd := exec.Command("git", "status", "--porcelain")
			var stdout bytes.Buffer
			cmd.Stdout = &stdout
			err = cmd.Run()
			if err != nil {
				return clierr.New(err.Error(), clierr.Infof("The above error occurred while checking if this git repo worktree is clean.\nTo silence this warning, add the 'ignore-git-dirty' flag e.g 'gdeploy --ignore-git-dirty %s'", c.Command.Name))
			}
			if stdout.Len() > 0 {
				return clierr.New("Git worktree is not clean", clierr.Infof("We recommend that you commit all changes before creating or updating your deployment.\nTo silence this warning, add the 'ignore-git-dirty' flag e.g 'gdeploy --ignore-git-dirty %s'", c.Command.Name))
			}
		}
		return nil
	}
}
