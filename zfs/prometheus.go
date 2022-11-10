package zfs

import "github.com/prometheus/client_golang/prometheus"

var prom struct {
	ZFSListFilesystemVersionDuration *prometheus.HistogramVec
	ZFSSnapshotDuration              *prometheus.HistogramVec
	ZFSBookmarkDuration              *prometheus.HistogramVec
	ZFSDestroyDuration               *prometheus.HistogramVec
	ZFSMissingDatasets               *prometheus.GaugeVec

}

func init() {
	prom.ZFSListFilesystemVersionDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "zrepl",
		Subsystem: "zfs",
		Name:      "list_filesystem_versions_duration",
		Help:      "Seconds it took for listing the versions of a given filesystem",
	}, []string{"filesystem"})
	prom.ZFSSnapshotDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "zrepl",
		Subsystem: "zfs",
		Name:      "snapshot_duration",
		Help:      "Seconds it took to create a snapshot a given filesystem",
	}, []string{"filesystem"})
	prom.ZFSBookmarkDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "zrepl",
		Subsystem: "zfs",
		Name:      "bookmark_duration",
		Help:      "Duration it took to bookmark a given snapshot",
	}, []string{"filesystem"})
	prom.ZFSDestroyDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "zrepl",
		Subsystem: "zfs",
		Name:      "destroy_duration",
		Help:      "Duration it took to destroy a dataset",
	}, []string{"dataset_type", "filesystem"})
	prom.ZFSMissingDatasets = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "zrepl",
		Subsystem: "zfs",
		Name:      "unmatched_dataset_filter_rule_count",
		Help:      "number of configured datasets that don't exist",
	}, []string{"zrepl_job"})
}

func PrometheusRegister(registry prometheus.Registerer) error {
	if err := registry.Register(prom.ZFSListFilesystemVersionDuration); err != nil {
		return err
	}
	if err := registry.Register(prom.ZFSBookmarkDuration); err != nil {
		return err
	}
	if err := registry.Register(prom.ZFSSnapshotDuration); err != nil {
		return err
	}
	if err := registry.Register(prom.ZFSDestroyDuration); err != nil {
		return err
	}
	if err := registry.Register(prom.ZFSMissingDatasets); err != nil {
		return err
	}
	return nil
}

func SetMissingCount(job string, count float64) {
  prom.ZFSMissingDatasets.WithLabelValues(job).Set(count)
}
