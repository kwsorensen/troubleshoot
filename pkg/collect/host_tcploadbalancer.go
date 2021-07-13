package collect

import (
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/pkg/errors"
	troubleshootv1beta2 "github.com/replicatedhq/troubleshoot/pkg/apis/troubleshoot/v1beta2"
)

type CollectHostTCPLoadBalancer struct {
	hostCollector *troubleshootv1beta2.TCPLoadBalancer
}

func (c *CollectHostTCPLoadBalancer) Title() string {
	return hostCollectorTitleOrDefault(c.hostCollector.HostCollectorMeta, "TCP Load Balancer")
}

func (c *CollectHostTCPLoadBalancer) IsExcluded() (bool, error) {
	return isExcluded(c.hostCollector.Exclude)
}

func (c *CollectHostTCPLoadBalancer) Collect(progressChan chan<- interface{}) (map[string][]byte, error) {
	listenAddress := fmt.Sprintf("0.0.0.0:%d", c.hostCollector.Port)
	dialAddress := c.hostCollector.Address

	if !isValidLoadBalancerAddress(dialAddress) {
		// create a structure and return it with error
		println("Error in validating LB address.")
		result := NetworkStatusResult{
			Status: NetworkStatusInvalidAddress,
			Error:  "Invalid Load Balancer Address",
		}

		b, err := json.Marshal(result)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal result")
		}
		name := path.Join("tcpLoadBalancer", "tcpLoadBalancer.json")
		name = path.Join("tcpLoadBalancer", fmt.Sprintf("%s.json", c.hostCollector.CollectorName))
		return map[string][]byte{
			name: b,
		}, errors.New("Errors")

	}

	timeout := 60 * time.Minute
	if c.hostCollector.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(c.hostCollector.Timeout)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse duration")
		}
	}
	println("Calling checkTCPConnection(progressChan, listenAddress, dialAddress, timeout)")
	networkStatus, err := checkTCPConnection(progressChan, listenAddress, dialAddress, timeout)
	if err != nil {
		println("Inside error line 40")
		return nil, err
	}
	println("Line 43")
	result := NetworkStatusResult{
		Status: networkStatus,
	}

	b, err := json.Marshal(result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal result")
	}

	name := path.Join("tcpLoadBalancer", "tcpLoadBalancer.json")
	if c.hostCollector.CollectorName != "" {
		name = path.Join("tcpLoadBalancer", fmt.Sprintf("%s.json", c.hostCollector.CollectorName))
	}

	return map[string][]byte{
		name: b,
	}, nil
}
