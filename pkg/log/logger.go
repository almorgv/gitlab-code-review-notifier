package log

import "github.com/rs/zerolog"

type Logger struct {
	logger zerolog.Logger
}

func (log *Logger) Trace(msg string) {
	log.logger.Trace().Msg(msg)
}

func (log *Logger) Debug(msg string) {
	log.logger.Debug().Msg(msg)
}

func (log *Logger) Info(msg string) {
	log.logger.Info().Msg(msg)
}

func (log *Logger) Warn(msg string) {
	log.logger.Warn().Msg(msg)
}

func (log *Logger) Error(msg string) {
	log.logger.Error().Msg(msg)
}

func (log *Logger) Tracef(msg string, v ...interface{}) {
	log.logger.Trace().Msgf(msg, v...)
}

func (log *Logger) Debugf(msg string, v ...interface{}) {
	log.logger.Debug().Msgf(msg, v...)
}

func (log *Logger) Infof(msg string, v ...interface{}) {
	log.logger.Info().Msgf(msg, v...)
}

func (log *Logger) Warnf(msg string, v ...interface{}) {
	log.logger.Warn().Msgf(msg, v...)
}

func (log *Logger) Errorf(msg string, v ...interface{}) {
	log.logger.Error().Msgf(msg, v...)
}

func (log *Logger) Fatalf(msg string, v ...interface{}) {
	log.logger.Fatal().Msgf(msg, v...)
}

func (log *Logger) TraceWithErr(err error, msg string) {
	log.logger.Trace().Err(err).Msg(msg)
}

func (log *Logger) DebugWithErr(err error, msg string) {
	log.logger.Debug().Err(err).Msg(msg)
}

func (log *Logger) InfoWithErr(err error, msg string) {
	log.logger.Info().Err(err).Msg(msg)
}

func (log *Logger) WarnWithErr(err error, msg string) {
	log.logger.Warn().Err(err).Msg(msg)
}

func (log *Logger) ErrorWithErr(err error, msg string) {
	log.logger.Error().Err(err).Msg(msg)
}

func (log *Logger) TracefWithErr(err error, msg string, v ...interface{}) {
	log.logger.Trace().Err(err).Msgf(msg, v...)
}

func (log *Logger) DebugfWithErr(err error, msg string, v ...interface{}) {
	log.logger.Debug().Err(err).Msgf(msg, v...)
}

func (log *Logger) InfofWithErr(err error, msg string, v ...interface{}) {
	log.logger.Info().Err(err).Msgf(msg, v...)
}

func (log *Logger) WarnfWithErr(err error, msg string, v ...interface{}) {
	log.logger.Warn().Err(err).Msgf(msg, v...)
}

func (log *Logger) ErrorfWithErr(err error, msg string, v ...interface{}) {
	log.logger.Error().Err(err).Msgf(msg, v...)
}
