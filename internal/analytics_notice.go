package internal

import (
	"errors"
	"os"
	"path"

	"github.com/common-fate/clio"
)

const notice = `Attention: Common Fate collects anonymous product analytics about Granted Approvals deployments.
This information is used to shape Common Fate's roadmap and prioritize features.
You can learn more, including how to opt-out and how to provide feedback via our RFD process, by visiting the following URL:
https://docs.commonfate.io/telemetry`

// PrintAnalyticsNotice prints an analytics notice to the user.
// If 'always' is true, it will always print the notice. Otherwise,
// the notice will be displayed once.
//
// More info on our RFD: https://github.com/common-fate/rfds/discussions/8
func PrintAnalyticsNotice(always bool) error {
	if always {
		clio.Info(notice)
		return nil
	}

	cd, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	dir := path.Join(cd, "commonfate")
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	noticefile := path.Join(dir, "analyticsnotice")
	if _, err := os.Stat(noticefile); err == nil {
		// notice already displayed
		return nil
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	file, err := os.Create(noticefile)
	if err != nil {
		return err
	}
	defer file.Close()

	clio.Info(notice)
	return nil
}
