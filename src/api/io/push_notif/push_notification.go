package push_notif

import (
	"be-tpms/src/api/io"
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"log"
	"time"
)

type NotificationSender struct {
	cxt context.Context
	app *firebase.App
}

func NewNotificationSender(ctx context.Context, app *firebase.App) *NotificationSender {
	return &NotificationSender{
		cxt: ctx,
		app: app,
	}
}

//GetClient Obtain a messaging.Client from the App.
func (ns *NotificationSender) GetClient() *messaging.Client {
	client, err := ns.app.Messaging(ns.cxt)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}
	return client
}

func (ns *NotificationSender) SendMessage(token string, data map[string]string) error {
	if token == "" {
		return fmt.Errorf("token is empty string")
	}

	// Send a message to the device corresponding to the provided registration token.
	response, err := ns.GetClient().Send(ns.cxt, buildNotificationMessage(token, data))
	if err != nil {
		return err
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
	return nil
}

func buildNotificationMessage(token string, data map[string]string) *messaging.Message {
	duration := time.Duration(48) * time.Hour
	badge := 42
	return &messaging.Message{
		Notification: &messaging.Notification{
			Title: data[io.TITLE],
			Body:  data[io.BODY],
		},
		Android: &messaging.AndroidConfig{
			TTL: &duration,
			Notification: &messaging.AndroidNotification{
				Icon:  "stock_ticker_update",
				Color: "#f45342",
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Badge: &badge,
				},
			},
		},
		Token: token,
	}
}
