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
	wlanID := dpskCmd.Int("wlansvc-id", -1, "Ruckus Wlan Service ID")
	username := dpskCmd.String("user", "", "Username")
	dpskCmd.Parse(args)

	if *wlanID < 0 {
		return &errors.CommandError{
			Msg:     fmt.Sprintf("wlan-id is invalid: %d", *wlanID),
			FlagSet: dpskCmd,
		}
	}

	if *username == "" {
		return &errors.CommandError{
			Msg:     fmt.Sprintf("username is invalid: %s", *username),
			FlagSet: dpskCmd,
		}
	}

	wlanIDStr := strconv.Itoa(*wlanID)

	filters := make(map[string]dpsk.Filter)
	filters["WlansvcID"] = &filterExact{value: &wlanIDStr}
	filters["User"] = &filterExact{value: username}

	entries, err := loadEntries(svc, filters)
	if err != nil {
		return fmt.Errorf("error filtering DPSK list: %v", err)
	}

	if len(*entries) > 0 {
		return fmt.Errorf("DPSK already exists for username: %s and wlanID: %d", *username, *wlanID)
	}

	err = svc.Create(*wlanID, *username)
	if err != nil {
		return fmt.Errorf("error creating DPSK user: %v", err)
	}

	entries, err = loadEntries(svc, filters)
	if err != nil {
		return fmt.Errorf("error filtering DPSK list: %v", err)
	}

	for _, entry := range *entries {
		fmt.Println(entry.Passphrase)
	}
	return nil
}

func loadEntries(svc *client.DpskService, filters map[string]dpsk.Filter) (*dpsk.Entries, error) {
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
