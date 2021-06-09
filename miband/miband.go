// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package miband

import (
	"context"
	"crypto/aes"

	"github.com/JuulLabs-OSS/ble"
	"github.com/JuulLabs-OSS/ble/darwin"
	"github.com/pkg/errors"

	"github.com/wuhan005/mebeats/cryptoutil"
)

type MiBand struct {
	client  ble.Client
	authKey string
	authed  chan struct{}

	state MiBandState

	service *ble.Service

	notify             *ble.Service
	authCharacteristic *ble.Characteristic

	battery               *ble.Service
	batteryCharacteristic *ble.Characteristic

	heartRate                      *ble.Service
	heartRateControlCharacteristic *ble.Characteristic
	heartRateMeasureCharacteristic *ble.Characteristic

	currentHeartRate int
	subscriber       []chan struct{}
}

// NewMiBand searches and connect to the Mi band.
func NewMiBand(deviceAddr, authKey string) (*MiBand, error) {
	ctx := context.Background()

	d, err := darwin.NewDevice()
	if err != nil {
		return nil, errors.Wrap(err, "new device")
	}
	ble.SetDefaultDevice(d)

	// Search and connect to the Mi band.
	client, err := ble.Connect(ctx, func(a ble.Advertisement) bool {
		return a.Addr().String() == deviceAddr
	})
	if err != nil {
		return nil, errors.Wrap(err, "dial client")
	}

	// Discover all the device services.
	services, err := client.DiscoverServices(nil)
	if err != nil {
		return nil, errors.Wrap(err, "discover services")
	}

	miband := &MiBand{
		authKey: authKey,
		client:  client,
		authed:  make(chan struct{}),
	}

	for _, service := range services {
		service := service
		switch service.UUID.String() {
		case "fee0": // Service
			miband.service = service
		case "fee1": // Notify
			miband.notify = service
		case "180d": // HeartRate
			miband.heartRate = service
		case "180f": // Battery
			miband.battery = service
		}
	}

	// Notify service
	characteristics, err := client.DiscoverCharacteristics(nil, miband.notify)
	if err != nil {
		return nil, errors.Wrap(err, "discover notify service characteristics")
	}
	for _, characteristic := range characteristics {
		characteristic := characteristic
		if characteristic.UUID.String() == "000000090000351221180009af100700" {
			miband.authCharacteristic = characteristic
		}
	}
	// Subscribe to the auth characteristic
	err = client.Subscribe(miband.authCharacteristic, false, miband.handleAuthNotification)
	if err != nil {
		return nil, errors.Wrap(err, "subscribe auth characteristic")
	}

	// Battery service
	characteristics, err = client.DiscoverCharacteristics(nil, miband.battery)
	if err != nil {
		return nil, errors.Wrap(err, "discover battery service characteristics")
	}
	for _, characteristic := range characteristics {
		characteristic := characteristic
		if characteristic.UUID.String() == "2a19" {
			miband.batteryCharacteristic = characteristic
		}
	}

	// HeartRate service
	characteristics, err = client.DiscoverCharacteristics(nil, miband.heartRate)
	if err != nil {
		return nil, errors.Wrap(err, "discover heart rate service characteristics")
	}
	for _, characteristic := range characteristics {
		characteristic := characteristic
		if characteristic.UUID.String() == "2a37" {
			miband.heartRateMeasureCharacteristic = characteristic
		} else if characteristic.UUID.String() == "2a39" {
			miband.heartRateControlCharacteristic = characteristic
		}
	}

	return miband, nil
}

func (m *MiBand) Initialize() error {
	if err := m.requestRandomNumber(); err != nil {
		return errors.Wrap(err, "request random number")
	}
	return nil
}

func (m *MiBand) Subscribe() chan struct{} {
	ch := make(chan struct{})
	// FIXME: Here is not thread-safe.
	m.subscriber = append(m.subscriber, ch)
	return ch
}

func (m *MiBand) boardcast() {
	for _, ch := range m.subscriber {
		ch := ch
		go func() {
			ch <- struct{}{}
		}()
	}
}

func (m *MiBand) requestRandomNumber() error {
	return m.client.WriteCharacteristic(m.authCharacteristic, []byte("\x02\x00"), false)
}

func (m *MiBand) sendEncryptRandomNumber(data []byte) error {
	encryptData, err := m.encrypt(data)
	if err != nil {
		return errors.Wrap(err, "encrypt")
	}

	message := append([]byte("\x03\x00"), encryptData...)
	if err := m.client.WriteCharacteristic(m.authCharacteristic, message, false); err != nil {
		return errors.Wrap(err, "write characteristic")
	}
	return nil
}

func (m *MiBand) sendKey() error {
	return m.client.WriteCharacteristic(m.authCharacteristic, []byte("\x01\x00"+m.authKey), false)
}

func (m *MiBand) encrypt(message []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(m.authKey))
	if err != nil {
		return nil, errors.Wrap(err, "create new cipher")
	}

	cipherText := make([]byte, len(message))
	mode := cryptoutil.NewECBEncrypter(block)
	mode.CryptBlocks(cipherText, message)
	return cipherText, nil
}
