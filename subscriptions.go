package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// A Subscription is something that should be periodically checked for
// updated.
type Subscription struct {
	Title, URL  string
	LastUpdated time.Time
}

// A Subscriptions should provide a means to List, Add and Remove
// subscriptions given their URL.
type Subscriptions interface {
	List() ([]*Subscription, error)
	Add(string) error
	Remove(string) error
}

type inMemorySubscriptions struct{}

func (s *inMemorySubscriptions) List() ([]*Subscription, error) {
	return []*Subscription{}, nil
}

func (s *inMemorySubscriptions) Add(url string) error {
	return nil
}

func (s *inMemorySubscriptions) Remove(url string) error {
	return nil
}

var subscriptions Subscriptions
var defaultSubscriptions = &inMemorySubscriptions{}

var subscriptionsCmd = &cobra.Command{
	Use:     "subscriptions",
	Aliases: []string{"sub", "subs"},
	Short:   "Manage podcasts subscriptions",
	Long:    `Add or remove podcast subscriptions.`,
	Run: func(c *cobra.Command, args []string) {
		subs, err := subscriptions.List()
		if err != nil {
			fmt.Fprintf(c.OutOrStderr(), "%q\n", err)
		}
		for _, sub := range subs {
			fmt.Fprintf(c.OutOrStdout(), "%s [%s]\n", sub.Title, sub.URL)
		}
	},
}

func init() {
	cmd.AddCommand(subscriptionsCmd)
	subscriptions = defaultSubscriptions
}
