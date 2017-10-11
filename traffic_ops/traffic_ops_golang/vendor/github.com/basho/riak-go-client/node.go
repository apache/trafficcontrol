// Copyright 2015-present Basho Technologies, Inc.
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

package riak

import (
	"fmt"
	"net"
	"time"
)

// Constants identifying Node state
const (
	nodeCreated state = iota
	nodeRunning
	nodeHealthChecking
	nodeShuttingDown
	nodeShutdown
	nodeError
)

// NodeOptions defines the RemoteAddress and operational configuration for connections to a Riak KV
// instance
type NodeOptions struct {
	RemoteAddress       string
	MinConnections      uint16
	MaxConnections      uint16
	TempNetErrorRetries uint16
	IdleTimeout         time.Duration
	ConnectTimeout      time.Duration
	RequestTimeout      time.Duration
	HealthCheckInterval time.Duration
	HealthCheckBuilder  CommandBuilder
	AuthOptions         *AuthOptions
}

// Node is a struct that contains all of the information needed to connect and maintain connections
// with a Riak KV instance
type Node struct {
	addr                *net.TCPAddr
	healthCheckInterval time.Duration
	healthCheckBuilder  CommandBuilder
	stopChan            chan struct{}
	cm                  *connectionManager
	stateData
}

var defaultNodeOptions = &NodeOptions{
	RemoteAddress:       defaultRemoteAddress,
	MinConnections:      defaultMinConnections,
	MaxConnections:      defaultMaxConnections,
	TempNetErrorRetries: defaultTempNetErrorRetries,
	IdleTimeout:         defaultIdleTimeout,
	ConnectTimeout:      defaultConnectTimeout,
	RequestTimeout:      defaultRequestTimeout,
}

// NewNode is a factory function that takes a NodeOptions struct and returns a Node struct
func NewNode(options *NodeOptions) (*Node, error) {
	if options == nil {
		options = defaultNodeOptions
	}
	if options.RemoteAddress == "" {
		options.RemoteAddress = defaultRemoteAddress
	}
	if options.MinConnections == 0 {
		options.MinConnections = defaultMinConnections
	}
	if options.MaxConnections == 0 {
		options.MaxConnections = defaultMaxConnections
	}
	if options.TempNetErrorRetries == 0 {
		options.TempNetErrorRetries = defaultTempNetErrorRetries
	}
	if options.IdleTimeout == 0 {
		options.IdleTimeout = defaultIdleTimeout
	}
	if options.ConnectTimeout == 0 {
		options.ConnectTimeout = defaultConnectTimeout
	}
	if options.RequestTimeout == 0 {
		options.RequestTimeout = defaultRequestTimeout
	}
	if options.HealthCheckInterval == 0 {
		options.HealthCheckInterval = defaultHealthCheckInterval
	}

	var err error
	var resolvedAddress *net.TCPAddr
	resolvedAddress, err = net.ResolveTCPAddr("tcp", options.RemoteAddress)
	if err == nil {
		n := &Node{
			stopChan:            make(chan struct{}),
			addr:                resolvedAddress,
			healthCheckInterval: options.HealthCheckInterval,
			healthCheckBuilder:  options.HealthCheckBuilder,
		}

		connMgrOpts := &connectionManagerOptions{
			addr:                resolvedAddress,
			minConnections:      options.MinConnections,
			maxConnections:      options.MaxConnections,
			tempNetErrorRetries: options.TempNetErrorRetries,
			idleTimeout:         options.IdleTimeout,
			connectTimeout:      options.ConnectTimeout,
			requestTimeout:      options.RequestTimeout,
			authOptions:         options.AuthOptions,
		}

		var cm *connectionManager
		if cm, err = newConnectionManager(connMgrOpts); err == nil {
			n.cm = cm
			n.initStateData("nodeCreated", "nodeRunning", "nodeHealthChecking", "nodeShuttingDown", "nodeShutdown", "nodeError")
			n.setState(nodeCreated)
			return n, nil
		}
	}

	return nil, err
}

// String returns a formatted string including the remoteAddress for the Node and its current
// connection count
func (n *Node) String() string {
	return fmt.Sprintf("%v|%d|%d", n.addr, n.cm.count(), n.cm.q.count())
}

// Start opens a connection with Riak at the configured remoteAddress and adds the connections to the
// active pool
func (n *Node) start() error {
	if err := n.stateCheck(nodeCreated); err != nil {
		return err
	}

	logDebug("[Node]", "(%v) starting", n)
	if err := n.cm.start(); err != nil {
		logErr("[Node]", err)
	}
	n.setState(nodeRunning)
	logDebug("[Node]", "(%v) started", n)

	return nil
}

// Stop closes the connections with Riak at the configured remoteAddress and removes the connections
// from the active pool
func (n *Node) stop() error {
	if err := n.stateCheck(nodeRunning, nodeHealthChecking); err != nil {
		return err
	}

	logDebug("[Node]", "(%v) shutting down.", n)

	n.setState(nodeShuttingDown)
	close(n.stopChan)

	err := n.cm.stop()

	if err == nil {
		n.setState(nodeShutdown)
		logDebug("[Node]", "(%v) shut down.", n)
	} else {
		n.setState(nodeError)
		logErr("[Node]", err)
	}

	return err
}

// Execute retrieves an available connection from the pool and executes the Command operation against
// Riak
func (n *Node) execute(cmd Command) (bool, error) {
	if err := n.stateCheck(nodeRunning, nodeHealthChecking); err != nil {
		return false, err
	}

	if n.isCurrentState(nodeRunning) {
		conn, err := n.cm.get()
		if err != nil {
			logErr("[Node]", err)
			n.doHealthCheck()
			return false, err
		}

		if conn == nil {
			panic(fmt.Sprintf("[Node] (%v) expected non-nil connection", n))
		}

		if rc, ok := cmd.(retryableCommand); ok {
			rc.setLastNode(n)
		}

		logDebug("[Node]", "(%v) - executing command '%v'", n, cmd.Name())
		err = conn.execute(cmd)
		if err == nil {
			// NB: basically the success path of _responseReceived in Node.js client
			if cmErr := n.cm.put(conn); cmErr != nil {
				logErr("[Node]", cmErr)
			}
			return true, nil
		} else {
			// NB: basically, this is _connectionClosed / _responseReceived in Node.js client
			// must differentiate between Riak and non-Riak errors here and within execute() in connection
			switch err.(type) {
			case RiakError, ClientError:
				// Riak and Client errors will not close connection
				if cmErr := n.cm.put(conn); cmErr != nil {
					logErr("[Node]", cmErr)
				}
				return true, err
			default:
				// NB: must be a non-Riak, non-Client error, close the connection
				if cmErr := n.cm.remove(conn); cmErr != nil {
					logErr("[Node]", cmErr)
				}
				if !isTemporaryNetError(err) {
					n.doHealthCheck()
				}
				return true, err
			}
		}
	} else {
		return false, nil
	}
}

func (n *Node) doHealthCheck() {
	// NB: ensure we're not already healthchecking or shutting down
	if n.isStateLessThan(nodeHealthChecking) {
		n.setState(nodeHealthChecking)
		go n.healthCheck()
	} else {
		logDebug("[Node]", "(%v) is already healthchecking or shutting down.", n)
	}
}

func (n *Node) getHealthCheckCommand() (hc Command) {
	// This is necessary to have a unique Command struct as part of each
	// connection so that concurrent calls to check health can all have
	// unique results
	var err error
	if n.healthCheckBuilder != nil {
		hc, err = n.healthCheckBuilder.Build()
	} else {
		hc = &PingCommand{}
	}

	if err != nil {
		logErr("[Node]", err)
		hc = &PingCommand{}
	}

	return
}

func (n *Node) ensureHealthCheckCanContinue() bool {
	// ensure we ARE healthchecking
	if !n.isCurrentState(nodeHealthChecking) {
		logDebug("[Node]", "(%v) expected healthchecking state, got %s", n, n.stateData.String())
		return false
	}
	return true
}

// private goroutine funcs

func (n *Node) healthCheck() {
	logDebug("[Node]", "(%v) starting healthcheck routine", n)

	healthCheckTicker := time.NewTicker(n.healthCheckInterval)
	defer healthCheckTicker.Stop()

	for {
		if !n.ensureHealthCheckCanContinue() {
			return
		}
		select {
		case <-n.stopChan:
			logDebug("[Node]", "(%v) healthcheck quitting", n)
			return
		case t := <-healthCheckTicker.C:
			if !n.ensureHealthCheckCanContinue() {
				return
			}
			logDebug("[Node]", "(%v) running healthcheck at %v", n, t)
			conn, cerr := n.cm.createConnection()
			if cerr != nil {
				conn.close()
				logError("[Node]", "(%v) failed healthcheck in createConnection, err: %v", n, cerr)
			} else {
				if !n.ensureHealthCheckCanContinue() {
					conn.close()
					return
				}
				hcmd := n.getHealthCheckCommand()
				logDebug("[Node]", "(%v) healthcheck executing %v", n, hcmd.Name())
				if hcerr := conn.execute(hcmd); hcerr != nil || !hcmd.Success() {
					conn.close()
					logError("[Node]", "(%v) failed healthcheck, err: %v", n, hcerr)
				} else {
					conn.close()
					logDebug("[Node]", "(%v) healthcheck success, err: %v, success: %v", n, hcerr, hcmd.Success())
					if n.ensureHealthCheckCanContinue() {
						n.setState(nodeRunning)
					}
					return
				}
			}
		}
	}
}
