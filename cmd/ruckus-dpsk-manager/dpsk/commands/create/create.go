package create

import (
	"flag"
	"fmt"

	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/errors"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"
)

func Handle(svc *client.DpskService, args []string) error {
	dpskCmd := flag.NewFlagSet("create", flag.ExitOnError)
	wlanID := dpskCmd.Int("wlan-id", -1, "Wlan ID")
	username := dpskCmd.String("username", "", "Username")
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

	dpskData, err := svc.List()
	if err != nil {
		return fmt.Errorf("error getting DPSK list: %v", err)
	}

	dpsk, err := dpskData.FindByWlanUser(*wlanID, *username)
	if err != nil {
		err = svc.Create(*wlanID, *username)
		if err != nil {
			return fmt.Errorf("error creating DPSK user: %v", err)
		}

		dpskData, err = svc.List()
		if err != nil {
			return fmt.Errorf("error getting DPSK list: %v", err)
		}

		dpsk, err = dpskData.FindByWlanUser(*wlanID, *username)
		if err != nil {
			return fmt.Errorf("error finding user: %v", err)
		}
	}

	fmt.Println(dpsk.Passphrase)
	return nil
}
