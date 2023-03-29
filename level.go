package loggy

type Level = int

const (
	// LevelStd indicates standard log output. Always shown.
	LevelStd Level = iota
	// LevelCritical indicates a fatal issue.
	LevelCritical
	// LevelError indicates a general issue that was recovered.
	LevelError
	// LevelWarning indicates an issue that may require intervention.
	LevelWarning
	// LevelInfo indicates generic runtime information.
	LevelInfo
	// LevelDebug indicates debug output.
	LevelDebug
)

// LevelNames describe the alphabetical types to label each Level* with in stdout/stderr.
var LevelNames = map[Level]string{
	LevelCritical: "CRIT",
	LevelDebug:    "DEBUG",
	LevelError:    "ERROR",
	LevelInfo:     "INFO",
	LevelWarning:  "WARN",
	LevelStd:      "OUT",
}
