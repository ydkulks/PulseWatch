package notification

import (
	"github.com/gen2brain/beeep"
)

func DesktopNotification(message string) error {
	err := beeep.Notify("PulseWatch", message, "")
	if err != nil {
		return err
	}
	return nil
}
