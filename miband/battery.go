// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package miband

import (
	"github.com/pkg/errors"
)

func (m *MiBand) GetBatteryPercent() (int, error) {
	data, err := m.client.ReadCharacteristic(m.batteryCharacteristic)
	if err != nil {
		return 0, errors.Wrap(err, "read characteristic")
	}
	return int(data[0]), nil
}
