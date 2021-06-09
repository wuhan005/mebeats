// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package miband

import (
	log "unknwon.dev/clog/v2"
)

func (m *MiBand) handleAuthNotification(data []byte) {
	switch string(data[:3]) {
	case "\x10\x01\x01":
		log.Trace("[Auth] Start to request random number...")
		if err := m.requestRandomNumber(); err != nil {
			log.Error("[Auth] Failed to request random number: %v", err)
		}

	case "\x10\x01\x04":
		m.state = AuthKeySendingFailed
		log.Error("[Auth] Failed to send key.")

	case "\x10\x02\x01":
		log.Trace("[Auth] Start to send encrypt random number...")
		randomNumber := data[3:]
		if err := m.sendEncryptRandomNumber(randomNumber); err != nil {
			log.Error("[Auth] Failed to send encrypt random number: %v", err)
		}

	case "\x10\x02\x04":
		m.state = AuthRequestRandomNumberError
		log.Error("[Auth] Failed to request random number.")

	case "\x10\x03\x01":
		m.state = AuthSuccess
		log.Trace("[Auth] Success!")
		close(m.authed)

	case "\x10\x03\x04":
		m.state = AuthEncryptionKeyFailed
		log.Error("[Auth] Encryption key auth fail, sending new key...")
		err := m.sendKey()
		if err != nil {
			log.Error("[Auth] Failed to send new key: %v", err)
		}

	default:
		m.state = AuthFailed
		log.Error("Auth failed: %v", data[:3])
	}
}

func (m *MiBand) handleHeartRateNotification(data []byte) {
	m.currentHeartRate = int(data[1])
	m.boardcast()
	log.Trace("Heart rate: %d", m.currentHeartRate)
}
