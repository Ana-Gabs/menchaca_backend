package models

import "time"

type LogEntry struct {
	Email        string    `bson:"email"`
	Action       string    `bson:"action"`
	LogLevel     string    `bson:"logLevel"`
	Timestamp    time.Time `bson:"timestamp"`
	IP           string    `bson:"ip"`
	UserAgent    string    `bson:"userAgent"`
	Referer      string    `bson:"referer"`
	Origin       string    `bson:"origin"`
	Method       string    `bson:"method"`
	URL          string    `bson:"url"`
	Status       int       `bson:"status"`
	ResponseTime float64   `bson:"responseTime"`
	Protocol     string    `bson:"protocol"`
	Hostname     string    `bson:"hostname"`
	Environment  string    `bson:"environment"`
	GoVersion    string    `bson:"goVersion"`
	PID          int       `bson:"pid"`
}
