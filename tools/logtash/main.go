package main

import (
	logstash_logger "github.com/KaranJagtiani/go-logstash"
)

func main() {
	logger := logstash_logger.Init("localhost", 5032, "tcp", 5)

	payload := map[string]interface{}{
		"message": "TEST_MSG",
		"error":   false,
	}

	logger.Log(payload)   // Generic log
	logger.Info(payload)  // Adds "severity": "INFO"
	logger.Debug(payload) // Adds "severity": "DEBUG"
	logger.Warn(payload)  // Adds "severity": "WARN"
	logger.Error(payload) // Adds "severity": "ERROR"
}
