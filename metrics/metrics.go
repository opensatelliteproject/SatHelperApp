package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	droppedPackets = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sathelperapp_droppedpackets",
		Help: "Number of dropped packets",
	})
	decoderFifoUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "sathelperapp_decoderfifousage",
		Help: "Decoder FIFO usage in Percent",
	})
	demodulatorFifoUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "sathelperapp_decoderfifousage",
		Help: "Decoder FIFO usage in Percent",
	})
	viterbi = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "sathelperapp_viterbi",
		Help: "The total number of corrected bits by Viterbi",
	})
	viterbiHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "sathelperapp_viterbihist",
		Help: "The Viterbi corrected bits",
	})
	signalQuality = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "sathelperapp_signalquality",
		Help: "The signal quality in percent",
	})
	signalStatus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "sathelperapp_signalstatus",
		Help: "The signal status: 1 for locked, 0 for not locked",
	})
	syncCorrelation = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "sathelperapp_syncorrelation",
		Help: "The number of matched bits in sync correlation",
	})
	reedSolomon = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "sathelperapp_reedsolomon",
		Help: "The number of corrected bytes by ReedSolomon",
	})
	totalFiles = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sathelperapp_totalfiles",
		Help: "The number of total files received",
	})
	totalPackets = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sathelperapp_totalpackets",
		Help: "The number of total packets received",
	})

	prometheusEnabled = false
)

// DroppedPacket add specified number of packets to dropped packets
func DroppedPackets(packets int) {
	if prometheusEnabled {
		droppedPackets.Add(float64(packets))
	}
}

// DecoderFifoUsage sets the Decoder fifo usage
func DecoderFifoUsage(usage float64) {
	if prometheusEnabled {
		decoderFifoUsage.Set(usage)
	}
}

func DemodulatorFifoUsage(usage float64) {
	if prometheusEnabled {
		demodulatorFifoUsage.Set(usage)
	}
}

func Viterbi(vit int) {
	if prometheusEnabled {
		viterbi.Set(float64(vit))
		viterbiHistogram.Observe(float64(vit))
	}
}

func SignalQuality(quality int) {
	if prometheusEnabled {
		signalQuality.Set(float64(quality))
	}
}

func SignalStatus(locked bool) {
	if prometheusEnabled {
		if locked {
			signalStatus.Set(1)
		} else {
			signalStatus.Set(0)
		}
	}
}

func SyncCorrelation(correlation int) {
	if prometheusEnabled {
		syncCorrelation.Set(float64(correlation))
	}
}

func ReedSolomon(rs int) {
	if prometheusEnabled {
		reedSolomon.Set(float64(rs))
	}
}

func NewFile() {
	if prometheusEnabled {
		totalFiles.Inc()
	}
}

func NewPacket() {
	if prometheusEnabled {
		totalPackets.Inc()
	}
}

var registry *prometheus.Registry

// EnablePrometheus enable metric collection in prometheus
func EnablePrometheus() {
	if !prometheusEnabled {
		prometheusEnabled = true

		registry = prometheus.NewRegistry()
		registry.MustRegister(droppedPackets)
		registry.MustRegister(decoderFifoUsage)
		registry.MustRegister(demodulatorFifoUsage)
		registry.MustRegister(viterbi)
		registry.MustRegister(viterbiHistogram)
		registry.MustRegister(signalQuality)
		registry.MustRegister(signalStatus)
		registry.MustRegister(syncCorrelation)
		registry.MustRegister(reedSolomon)
		registry.MustRegister(totalPackets)
		registry.MustRegister(totalFiles)
	}
}

func GetHandler() http.Handler {
	if prometheusEnabled {
		return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	}
	return nil
}
