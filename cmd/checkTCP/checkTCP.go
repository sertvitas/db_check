package checktcpresponse

import (
	"fmt"
	"net"
	"time"
)

// CheckTCPResponse attempts to establish a TCP connection to the specified IP address and port
func CheckTCPResponse(ip string, port int, timeout time.Duration) string {
	// Concatenate IP and port
	address := fmt.Sprintf("%s:%d", ip, port)

	// Attempt to establish a TCP connection
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		// Handle timeout separately
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return "timeout"
		}
		return "failure"
	}

	// Close the connection if successful
	defer conn.Close()

	return "success"
}
