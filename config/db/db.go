package config

// Config is DB configuration
type Config struct {
	// Max number of open files.
	MaxOpenFile uint32
	// Cache sizes (in MiB) for specific columns.
	// cache size
	// Set number of columns
	Columns uint32
	// Should we keep WAL enabled?
	WAL bool

	// Compaction profile
	Compaction CompactionProfile
}

func NewDefault() Config {
	return Config{
		MaxOpenFile: 512,
		Columns:     0,
		WAL:         true,
		Compaction:  NewCompactionProfileOfSSD(),
	}
}

// CompactionProfile for the database settings
type CompactionProfile struct {
	// L0-L1 target file size
	InitFileSize uint64
	// L2-LN target file size multiplier
	FileSizeMultiplier uint32
	// rate limiter for background flushes and compactions, bytes/sec, if any
	WriteRateLimit uint64
}

// NewCompactionProfileOfHDD slow HDD.
func NewCompactionProfileOfHDD() CompactionProfile {
	return CompactionProfile{
		InitFileSize:       192 * 1024 * 1024,
		FileSizeMultiplier: 1,
		WriteRateLimit:     8 * 1024 * 1024,
	}
}

// NewCompactionProfileOfSSD suitable for SSD storage.
func NewCompactionProfileOfSSD() CompactionProfile {
	return CompactionProfile{
		InitFileSize:       32 * 1024 * 1024,
		FileSizeMultiplier: 2,
		WriteRateLimit:     0,
	}
}
