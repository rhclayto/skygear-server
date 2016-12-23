// Copyright 2015-present Oursky Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package push

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/SkygearIO/buford/push"
	"github.com/skygeario/skygear-server/pkg/server/skydb"
)

type certBaseAPNSPusher struct {
	APNSPusher

	// Function to obtain a skydb connection
	connOpener func() (skydb.Conn, error)

	conn    skydb.Conn
	service pushService

	failed chan failedNotification
}

// NewCertBaseAPNSPusher returns a new APNSPusher from content of certificate
// and private key as string
func NewCertBaseAPNSPusher(
	connOpener func() (skydb.Conn, error),
	gatewayType GatewayType,
	cert string,
	key string,
) (APNSPusher, error) {
	certificate, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return nil, err
	}

	client, err := push.NewClient(certificate)
	if err != nil {
		return nil, err
	}

	service, err := newPushService(client, gatewayType)
	if err != nil {
		return nil, err
	}

	return &certBaseAPNSPusher{
		connOpener: connOpener,
		service:    service,
	}, nil
}

// Start setups the pusher and starts it
func (pusher *certBaseAPNSPusher) Start() {
	conn, err := pusher.connOpener()
	if err != nil {
		log.Errorf("push/apns: failed to open skydb.Conn, abort feedback retrival: %v\n", err)
		return
	}

	pusher.conn = conn
	pusher.failed = make(chan failedNotification)

	go func() {
		checkFailedNotifications(pusher)
	}()
}

// Stop stops and cleans up the pusher
func (pusher *certBaseAPNSPusher) Stop() {
	close(pusher.failed)
	pusher.failed = nil
}

// Send sends a notification to the device identified by the
// specified device
func (pusher *certBaseAPNSPusher) Send(m Mapper, device skydb.Device) error {
	logger := log.WithFields(logrus.Fields{
		"deviceToken": device.Token,
		"deviceID":    device.ID,
		"deviceTopic": device.Topic,
	})

	if m == nil {
		logger.Warn("Cannot send push notification with nil data.")
		return errors.New("push/apns: push notification has no data")
	}

	if device.Topic == "" {
		logger.Warn("Cannot send push notification with empty topic.")
		return errors.New("push/apns: push notification has empty topic")
	}

	apnsMap, ok := m.Map()["apns"].(map[string]interface{})
	if !ok {
		return errors.New("push/apns: payload has no apns dictionary")
	}

	serializedPayload, err := json.Marshal(apnsMap)
	if err != nil {
		return err
	}

	headers := push.Headers{
		Topic: device.Topic,
	}

	apnsid, err := pusher.service.Push(device.Token, &headers, serializedPayload)
	if err != nil {
		if pushError, ok := err.(*push.Error); ok && pushError != nil {
			// We recognize the error, and that error comes from APNS
			logger.WithFields(logrus.Fields{
				"apnsErrorReason":    pushError.Reason,
				"apnsErrorStatus":    pushError.Status,
				"apnsErrorTimestamp": pushError.Timestamp,
			}).Error("push/apns: failed to send push notification")
			queueFailedNotification(pusher, device.Token, *pushError)
			return err
		}

		logger.Errorf("Failed to send push notification: %s", err)
		return err
	}

	logger.WithFields(logrus.Fields{
		"apnsID": apnsid,
	}).Info("push/apns: push notification is sent")

	return nil
}

func (pusher certBaseAPNSPusher) getFailedNotificationChannel() chan failedNotification {
	return pusher.failed
}

func (pusher certBaseAPNSPusher) deleteDeviceToken(token string, beforeTime time.Time) error {
	return pusher.conn.DeleteDevicesByToken(token, beforeTime)
}
