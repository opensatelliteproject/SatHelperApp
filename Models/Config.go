package Models

// region Config File Structs
type BaseConfig struct {
	SymbolRate        uint32
	RRCAlpha          float32
	Mode              string
	Decimation        uint8
	AGCEnabled        bool
	DeviceType        string
	SendConstellation bool
	PLLAlpha          float32
	DemuxerType       string
}

type CFileSourceConfig struct {
	Filename       string
	FastAsPossible bool
}

type AirspySourceConfig struct {
	MixerGain    uint8
	LNAGain      uint8
	VGAGain      uint8
	BiasTEnabled bool
}

type LimeSourceConfig struct {
	LNAGain uint8
	Antenna string
}

type SpyserverSourceConfig struct {
	Gain     uint8
	Hostname string
	Port     int
}

type SourceConfig struct {
	SampleRate uint32
	Frequency  uint32
}

type DecoderConfig struct {
	Display          bool
	UseLastFrameData bool
}

type TCPServerDemuxerConfig struct {
	Port int
	Host string
}

type AppConfig struct {
	Title            string
	Base             BaseConfig
	Decoder          DecoderConfig
	Source           SourceConfig
	AirspySource     AirspySourceConfig
	LimeSource       LimeSourceConfig
	CFileSource      CFileSourceConfig
	TCPServerDemuxer TCPServerDemuxerConfig
	SpyserverSource  SpyserverSourceConfig
}

// endregion
