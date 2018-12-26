package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Notifier represents an instance that pushes out notifications to
// Google Chat Webhook endpoint.
type Notifier struct {
	root       string
	httpClient *http.Client
}

// NewNotifier initialises a new instance of the Notifier.
func NewNotifier(root string, h http.Client) Notifier {
	return Notifier{
		root:       root,
		httpClient: &h,
	}
}

// PushNotification pushes out a notification to Google Chat Room.
func (n *Notifier) PushNotification(notif ChatNotification) error {
	out, err := json.Marshal(notif)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", n.root, bytes.NewBuffer(out))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// send the request
	resp, err := n.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respMsg, _ := ioutil.ReadAll(resp.Body)
		errLog.Printf("Error sending alert Webhook API error: %s", string(respMsg))
		return errors.New("Error while sending alert")
	}

	return nil
}
