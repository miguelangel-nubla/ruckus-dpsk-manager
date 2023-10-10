package create

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/errors"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/data/dpsk"
)

type filterExact struct {
	value *string
	dpsk.Filter
}

func (filter *filterExact) Test(s string) bool {
	return *filter.value == s
}

func Handle(svc *client.DpskService, args []string) error {
	dpskCmd := flag.NewFlagSet("create", flag.ExitOnError)
	wlansvcID := dpskCmd.Int("wlansvc-id", -1, "Ruckus Wlan Service ID")
	user := dpskCmd.String("user", "", "Username")
	dpskLen := dpskCmd.Int("dpsk-len", 8, "DPSK characger length")
	dpskCmd.Parse(args)

	if *wlansvcID < 0 {
		return &errors.CommandError{
			Msg:     fmt.Sprintf("wlan-id is invalid: %d", *wlansvcID),
			FlagSet: dpskCmd,
		}
	}

	if *user == "" {
		return &errors.CommandError{
			Msg:     fmt.Sprintf("username is invalid: %s", *user),
			FlagSet: dpskCmd,
		}
	}

	wlanIDStr := strconv.Itoa(*wlansvcID)

	filters := make(map[string]dpsk.Filter)
	filters["wlansvc-id"] = &filterExact{value: &wlanIDStr}
	filters["user"] = &filterExact{value: user}

	entries, err := loadEntries(svc, filters)
	if err != nil {
		return fmt.Errorf("error filtering DPSK list: %v", err)
	}

	if len(entries) > 0 {
		return fmt.Errorf("DPSK already exists for username: %s and wlanID: %d", *user, *wlansvcID)
	}

	err = svc.Create(*wlansvcID, *user, *dpskLen)
	if err != nil {
		return fmt.Errorf("error creating DPSK user: %v", err)
	}

	entries, err = loadEntries(svc, filters)
	if err != nil {
		return fmt.Errorf("error filtering DPSK list: %v", err)
	}

	for _, entry := range entries {
		fmt.Println(entry.Passphrase)
	}
	return nil
}

func loadEntries(svc *client.DpskService, filters map[string]dpsk.Filter) (dpsk.Entries, error) {
	dpskData, err := svc.List()
	if err != nil {
		return nil, err
	}

	entries, err := dpskData.Filter(filters)
	if err != nil {
		return nil, err
	}

	return entries, nil
}
