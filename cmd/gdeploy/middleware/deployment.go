package middleware

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/granted-approvals/internal/build"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/hashicorp/go-version"
	"github.com/urfave/cli/v2"
)

func RequireDeploymentConfig() cli.BeforeFunc {
	return func(c *cli.Context) error {

		f := c.Path("file")
		dc, err := deploy.LoadConfig(f)
		if err == deploy.ErrConfigNotExist {
			return clio.NewCLIError(fmt.Sprintf("Tried to load Granted deployment configuration from %s but the file doesn't exist.", f),
				clio.LogMsg(`
To fix this, take one of the following actions:
  a) run this command from a folder which contains a Granted deployment configuration file (like 'granted-deployment.yml')
  b) run 'gdeploy init' to set up a new deployment configuration file
`),
			)
		}
		if err != nil {
			return fmt.Errorf("failed to load config with error: %s", err)
		}

		if dc.Version != 2 {
			return fmt.Errorf("unexpected deployment config file version found, expected version: 2 found version: %d", dc.Version)
		}

		c.Context = deploy.SetConfigInContext(c.Context, dc)
		return nil
	}
}

func PreventDevUsage() cli.BeforeFunc {
	return func(c *cli.Context) error {
		ctx := c.Context
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}
		if dc.Deployment.Dev != nil && *dc.Deployment.Dev {
			return clio.NewCLIError("Unsupported command used on developement deployment", clio.WarnMsg("It looks like you tried to use an unsupported command on your developement stack: '%s'.", c.Command.Name), clio.InfoMsg("If you were trying to update your stack, use 'mage deploy:dev', if you didn't expect to see this message, check you are in the correct directory!"))
		}
		return nil
	}
}

// BeforeFunc wrapper for IsReleaseVersionDifferent func.
// Prompt user to save `gdeploy` version as release version to `granted-deployment.yml`
// if release version  and gdeploy version is different.
func VerifyGDeployCompatibility() cli.BeforeFunc {
	return func(c *cli.Context) error {
		ctx := c.Context
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		isVersionMismatch, err := IsReleaseVersionDifferent(dc.Deployment, build.Version)
		if err != nil {
			return err
		}

		if isVersionMismatch && c.Bool("ignore-version-mismatch") {
			clio.Warn("Ignoring version mismatch between gdeploy CLI (%s) and deployment release version (%s)", build.Version, dc.Deployment.Release)
			return nil
		}

		if isVersionMismatch {
			var shouldUpdate bool
			clio.Error("Incompatible gdeploy version. Expected %s got %s.", dc.Deployment.Release, build.Version)
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("Would you like to update your 'granted-deployment.yml' to release %s?", build.Version),
			}

			err = survey.AskOne(prompt, &shouldUpdate)
			if err != nil {
				return err
			}

			if shouldUpdate {
				dc.Deployment.Release = build.Version

				f := c.Path("file")

				err := dc.Save(f)
				if err != nil {
					return err
				}

				clio.Success("Release version updated to %s", build.Version)

				return nil

			}

			return clio.NewCLIError("Please update gdeploy version to match your release version in 'granted-deployment.yml'.")
		}

		return nil

	}
}

// Validate if the passed deployment configuration's release value and gdeploy version matches or not.
// Returns true if version same.
// Returns false if we can skip this check.
// Returns error for anything else.
func IsReleaseVersionDifferent(d deploy.Deployment, buildVersion string) (bool, error) {
	// skip compatibility check for dev deployments.
	if d.Dev != nil && *d.Dev {
		return false, nil
	}

	isValidReleaseNumber, err := regexp.MatchString(`v?\d.\d+.\d+`, d.Release)
	if err != nil {
		return false, err
	}

	if isValidReleaseNumber {
		if buildVersion == "dev" {
			clio.Warn("Skipping version compatibility check for dev gdeploy build")
			return false, nil
		}

		formattedBuildVersion, err := version.NewVersion(buildVersion)
		if err != nil {
			return false, err
		}

		formattedReleaseVersion, err := version.NewVersion(d.Release)
		if err != nil {
			return false, err
		}

		if formattedBuildVersion.Equal(formattedReleaseVersion) {
			return false, nil
		}

		return true, nil
	}

	// release value are added as URL for UAT. In such case it should skip this check.
	// if invalid URL, return with error.
	_, err = url.ParseRequestURI(d.Release)
	if err != nil {
		return false, fmt.Errorf("invalid URL. Please update your release version in 'granted-deployment.yml' to %s", buildVersion)
	}

	return false, nil
}
