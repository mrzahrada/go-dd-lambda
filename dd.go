package dd

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type metricType string

const (
	gauge     = metricType("gauge")
	count     = metricType("count")
	histogram = metricType("histogram")
	check     = metricType("check")
)

// A Client is a handle for sending messages to datadog. It is safe to use one Client from multiple goroutines simultaneously.
type Client struct {
	// Namespace to prepend to all statsd calls
	Namespace string
	// Tags are global tags to be added to every statsd call
	Tags []string
}

// New DataDog client
func New() (*Client, error) {
	return &Client{}, nil
}

// send handles sampling and sends the message to stdout. It also adds global namespace prefixes and tags.
// format: MONITORING|<unix_epoch_timestamp>|<value>|<metric_type>|<metric_name>|#<tag_list>
func (c *Client) sendFloat(name string, value float64, metric metricType, tags []string, rate float64) error {
	if c == nil {
		return nil
	}
	if rate < 1 && rand.Float64() > rate {
		return nil
	}

	_, err := fmt.Printf("MONITORING|%d|%f|%s|%s.%s|%s",
		time.Now().UTC().Unix(),
		value, metric,
		c.Namespace,
		name,
		strings.Join(tags, ","))

	return err
}

// send handles sampling and sends the message to stdout. It also adds global namespace prefixes and tags.
// format: MONITORING|<unix_epoch_timestamp>|<value>|<metric_type>|<metric_name>|#<tag_list>
func (c *Client) sendInt(name string, value int64, metric metricType, tags []string, rate float64) error {
	if c == nil {
		return nil
	}
	if rate < 1 && rand.Float64() > rate {
		return nil
	}

	_, err := fmt.Printf("MONITORING|%d|%d|%s|%s.%s|%s",
		time.Now().UTC().Unix(),
		value, metric,
		c.Namespace,
		name,
		strings.Join(tags, ","))

	return err
}

// Count tracks how many times something happened per second.
func (c *Client) Count(name string, value int64, tags []string, rate float64) error {
	return c.sendInt(name, value, gauge, tags, rate)
}

// Decr is just Count of -1
func (c *Client) Decr(name string, tags []string, rate float64) error {
	return c.sendInt(name, -1, count, tags, rate)
}

// Gauge measures the value of a metric at a particular time.
func (c *Client) Gauge(name string, value float64, tags []string, rate float64) error {
	return c.sendFloat(name, value, gauge, tags, rate)
}

// Histogram tracks the statistical distribution of a set of values on each host.
func (c *Client) Histogram(name string, value float64, tags []string, rate float64) error {
	return c.sendFloat(name, value, histogram, tags, rate)
}

// Incr is just Count of 1
func (c *Client) Incr(name string, tags []string, rate float64) error {
	return c.sendInt(name, 1, count, tags, rate)
}
