// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package miband

import (
	"time"

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"
)

func (m *MiBand) GetHeartRateOneTime() error {
	<-m.authed

	// Subscribe to the heart rate characteristic
	err := m.client.Subscribe(m.heartRateMeasureCharacteristic, false, m.handleHeartRateNotification)
	if err != nil {
		return errors.Wrap(err, "subscribe heart rate characteristic")
	}

	// Stop continuous.
	err = m.client.WriteCharacteristic(m.heartRateControlCharacteristic, []byte("\x15\x02\x00"), false)
	if err != nil {
		return errors.Wrap(err, "stop continuous")
	}

	// Stop manual.
	err = m.client.WriteCharacteristic(m.heartRateControlCharacteristic, []byte("\x15\x01\x00"), false)
	if err != nil {
		return errors.Wrap(err, "stop manual")
	}

	// Start manual.
	err = m.client.WriteCharacteristic(m.heartRateControlCharacteristic, []byte("\x15\x01\x01"), false)
	if err != nil {
		return errors.Wrap(err, "start manual")
	}

	go func() {
		for {
			time.Sleep(12 * time.Second)
			log.Trace("Send ping...")
			err = m.client.WriteCharacteristic(m.heartRateControlCharacteristic, []byte("\x16"), false)
			if err != nil {
				log.Error("Failed to send ping: %v", err)
			}
		}
	}()
	return nil
}

func (m *MiBand) GetCurrentHeartRate() int {
	return m.currentHeartRate
}
