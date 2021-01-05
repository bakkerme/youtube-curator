package collection

import (
	"testing"
)

func TestYTChannelData(t *testing.T) {
	t.Run("Creates valid YTChannelData", func(t *testing.T) {
		ytc := YTChannelData{
			IName:         "a name",
			IID:           "123abc",
			IRSSURL:       "http://testrss.example",
			IChannelURL:   "http://testchannel.example",
			IArchivalMode: ArchivalModeArchive,
		}

		matches := []bool{
			ytc.Name() == ytc.IName,
			ytc.ID() == ytc.IID,
			ytc.RSSURL() == ytc.IRSSURL,
			ytc.ChannelURL() == ytc.IChannelURL,
			ytc.ArchivalMode() == ytc.IArchivalMode,
		}

		for _, match := range matches {
			if !match {
				t.Errorf("Value did not match.")
			}
		}
	})
}
