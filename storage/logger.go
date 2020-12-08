package storage

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type BadgerLogger struct {
	log *zap.Logger
}

func NewBadgerLogger(log *zap.Logger) *BadgerLogger {
	return &BadgerLogger{log: log}
}

func (logger *BadgerLogger) Debugf(msg string, args ...interface{}) {
	logger.log.Debug(fmt.Sprintf(strings.TrimSuffix(msg, "\n"), args...))
}

func (logger *BadgerLogger) Infof(msg string, args ...interface{}) {
	logger.log.Info(fmt.Sprintf(strings.TrimSuffix(msg, "\n"), args...))
}

func (logger *BadgerLogger) Warningf(msg string, args ...interface{}) {
	logger.log.Warn(fmt.Sprintf(strings.TrimSuffix(msg, "\n"), args...))
}

func (logger *BadgerLogger) Errorf(msg string, args ...interface{}) {
	logger.log.Error(fmt.Sprintf(strings.TrimSuffix(msg, "\n"), args...))
}
