package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var jvmPath = "/admin/metrics?group=jvm&wt=json"

//JVMCollector collects JVM type metrics from solr
type JVMCollector struct {
	gcConcurrentMarkSweepCount *prometheus.Desc
	gcConcurrentMarkSweepTime  *prometheus.Desc
	gcParNewCount              *prometheus.Desc
	gcParNewTime               *prometheus.Desc

	memoryHeapCommitted *prometheus.Desc
	memoryHeapInit      *prometheus.Desc
	memoryHeapMax       *prometheus.Desc
	memoryHeapUsage     *prometheus.Desc
	memoryHeapUsed      *prometheus.Desc

	memoryNonHeapCommitted *prometheus.Desc
	memoryNonHeapInit      *prometheus.Desc
	memoryNonHeapMax       *prometheus.Desc
	memoryNonHeapUsage     *prometheus.Desc
	memoryNonHeapUsed      *prometheus.Desc

	memoryTotalCommitted *prometheus.Desc
	memoryTotalInit      *prometheus.Desc
	memoryTotalMax       *prometheus.Desc
	memoryTotalUsed      *prometheus.Desc

	osAvailableProcessors        *prometheus.Desc
	osCommittedVirtualMemorySize *prometheus.Desc
	osFreePhysicalMemorySize     *prometheus.Desc
	osFreeSwapSpaceSize          *prometheus.Desc
	osMaxFileDescriptorCount     *prometheus.Desc
	osOpenFileDescriptorCount    *prometheus.Desc
	osProcessCPUTime             *prometheus.Desc
	osSystemLoadAverage          *prometheus.Desc
	osTotalPhysicalMemorySize    *prometheus.Desc
	osTotalSwapSapceSize         *prometheus.Desc

	threadsBlockedCount      *prometheus.Desc
	threadsDaemonCount       *prometheus.Desc
	threadsDeadlockCount     *prometheus.Desc
	threadsNewCount          *prometheus.Desc
	threadsRunnableCount     *prometheus.Desc
	threadsTerminatedCount   *prometheus.Desc
	threadsTimedWaitingCount *prometheus.Desc
	threadsWaitingCount      *prometheus.Desc

	client http.Client
	jvmURL string
}

// NewJVMCollector returns a new Collector exposing solr jvm statistics.
func NewJVMCollector(client http.Client, solrBaseURL string) (*JVMCollector, error) {
	jvmURL := fmt.Sprintf("%s%s", solrBaseURL, jvmPath)
	return &JVMCollector{
		gcConcurrentMarkSweepCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "gc_concurrentmarksweep_count"),
			"Garbage collector concurrent mark sweep count.",
			[]string{}, nil,
		),
		gcConcurrentMarkSweepTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "gc_concurrentmarksweep_time"),
			"Garbage collector concurrent mark sweep time in miliseconds.",
			[]string{},
			nil,
		),
		gcParNewCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "gc_parnew_count"),
			"Garbage collector parnew count.",
			[]string{},
			nil,
		),
		gcParNewTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "gc_parnew_time"),
			"Garbage collector parnew time in miliseconds.",
			[]string{},
			nil,
		),

		memoryHeapCommitted: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "memory_heap_committed"),
			"JVM memory heap committed bytes.",
			[]string{},
			nil,
		),
		memoryHeapInit: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "memory_heap_init"),
			"JVM memory heap initial bytes.",
			[]string{},
			nil,
		),
		memoryHeapMax: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "memory_heap_max"),
			"JVM memory heap max bytes.",
			[]string{},
			nil,
		),
		memoryHeapUsage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "memory_heap_usage"),
			"JVM memory heap percentage usage.",
			[]string{},
			nil,
		),
		memoryHeapUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "memory_heap_used"),
			"JVM memory heap used bytes.",
			[]string{},
			nil,
		),

		memoryNonHeapCommitted: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "memory_nonheap_committed"),
			"JVM memory non heap committed bytes.",
			[]string{},
			nil,
		),
		memoryNonHeapInit: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "memory_nonheap_init"),
			"JVM memory non heap initial bytes.",
			[]string{},
			nil,
		),
		memoryNonHeapMax: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "memory_nonheap_max"),
			"JVM memory non heap max bytes.",
			[]string{},
			nil,
		),
		memoryNonHeapUsage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "memory_nonheap_usage"),
			"JVM memory non heap percentage usage.",
			[]string{},
			nil,
		),
		memoryNonHeapUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "memory_nonheap_used"),
			"JVM memory non heap used bytes.",
			[]string{},
			nil,
		),

		memoryTotalCommitted: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "memory_total_committed"),
			"JVM memory total committed bytes.",
			[]string{},
			nil,
		),
		memoryTotalInit: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "memory_total_init"),
			"JVM memory total inital bytes.",
			[]string{},
			nil,
		),
		memoryTotalMax: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "memory_total_max"),
			"JVM memory total max bytes.",
			[]string{},
			nil,
		),
		memoryTotalUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "memory_total_used"),
			"JVM memory total used bytes.",
			[]string{},
			nil,
		),

		osAvailableProcessors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "os_availableprocessors"),
			"Avaialable number of processors.",
			[]string{},
			nil,
		),
		osCommittedVirtualMemorySize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "os_committedvirtualmemorysize"),
			"Operating system commited virtual memory size in bytes.",
			[]string{},
			nil,
		),
		osFreePhysicalMemorySize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "os_freephysicalmemorysize"),
			"Operating system free physical memory in bytes.",
			[]string{},
			nil,
		),
		osFreeSwapSpaceSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "os_freeswapspacesize"),
			"Operating system free swap memory in bytes.",
			[]string{},
			nil,
		),
		osMaxFileDescriptorCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "os_maxfiledescriptorcount"),
			"Operating system maximum number of open file descriptors.",
			[]string{},
			nil,
		),
		osOpenFileDescriptorCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "os_openfiledescriptorcount"),
			"Operating system current number of open file descriptors.",
			[]string{},
			nil,
		),
		osProcessCPUTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "os_processcputime"),
			"Time process was running on the cpu in milliseconds.",
			[]string{},
			nil,
		),
		osSystemLoadAverage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "os_systemloadaverage"),
			"Operating system load average.",
			[]string{},
			nil,
		),
		osTotalPhysicalMemorySize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "os_totalphysicalmemorysize"),
			"Operating System total physical memory size in bytes.",
			[]string{},
			nil,
		),
		osTotalSwapSapceSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "os_totalswapspacesize"),
			"Operating System totale swap memory size in bytes.",
			[]string{},
			nil,
		),

		threadsBlockedCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "threads_blocked_count"),
			"Count of blocked threads.",
			[]string{},
			nil,
		),
		threadsDaemonCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "threads_daemon_conut"),
			"Count of daemon threads.",
			[]string{},
			nil,
		),
		threadsDeadlockCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "threads_deadlock_count"),
			"Count of deadlocked threads.",
			[]string{},
			nil,
		),
		threadsNewCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "threads_new_count"),
			"Count of new threads.",
			[]string{},
			nil,
		),
		threadsRunnableCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "threads_runnable_count"),
			"Count of runnable threads.",
			[]string{},
			nil,
		),
		threadsTerminatedCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "threads_terminated_count"),
			"Count of terminated threads.",
			[]string{},
			nil,
		),
		threadsTimedWaitingCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "threads_timedwaiting_count"),
			"Count of threads in timed_waiting state.",
			[]string{},
			nil,
		),
		threadsWaitingCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jvm", "threads_waiting_count"),
			"Count of waiting threads.",
			[]string{},
			nil,
		),
		client: client,
		jvmURL: jvmURL,
	}, nil
}

// Update exposes jvm related metrics from solr.
func (c *JVMCollector) Update(ch chan<- prometheus.Metric) error {
	resp, err := c.client.Get(c.jvmURL)
	if err != nil {
		return fmt.Errorf("Error while querying Solr for jvm stats: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read jvm stats response body: %v", err)

	}

	jvmStatus := &JVMStatus{}
	err = json.Unmarshal(body, jvmStatus)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal solr jvm JSON into struct: %v", err)

	}

	ch <- prometheus.MustNewConstMetric(c.gcConcurrentMarkSweepCount, prometheus.CounterValue, float64(jvmStatus.Metrics.JVM.GCConcurrentMarkSweepCount))
	ch <- prometheus.MustNewConstMetric(c.gcConcurrentMarkSweepTime, prometheus.CounterValue, float64(jvmStatus.Metrics.JVM.GCConcurrentMarkSweepTime))
	ch <- prometheus.MustNewConstMetric(c.gcParNewCount, prometheus.CounterValue, float64(jvmStatus.Metrics.JVM.GCParNewCount))
	ch <- prometheus.MustNewConstMetric(c.gcParNewTime, prometheus.CounterValue, float64(jvmStatus.Metrics.JVM.GCParNewTime))

	ch <- prometheus.MustNewConstMetric(c.memoryHeapCommitted, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.MemoryHeapCommitted))
	ch <- prometheus.MustNewConstMetric(c.memoryHeapInit, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.MemoryHeapInit))
	ch <- prometheus.MustNewConstMetric(c.memoryHeapMax, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.MemoryHeapMax))
	ch <- prometheus.MustNewConstMetric(c.memoryHeapUsage, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.MemoryHeapUsage))
	ch <- prometheus.MustNewConstMetric(c.memoryHeapUsed, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.MemoryHeapUsed))

	ch <- prometheus.MustNewConstMetric(c.memoryNonHeapCommitted, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.MemoryNonHeapCommitted))
	ch <- prometheus.MustNewConstMetric(c.memoryNonHeapInit, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.MemoryNonHeapInit))
	ch <- prometheus.MustNewConstMetric(c.memoryNonHeapMax, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.MemoryNonHeapMax))
	ch <- prometheus.MustNewConstMetric(c.memoryNonHeapUsage, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.MemoryNonHeapUsage))
	ch <- prometheus.MustNewConstMetric(c.memoryNonHeapUsed, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.MemoryNonHeapUsed))

	ch <- prometheus.MustNewConstMetric(c.memoryTotalCommitted, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.MemoryTotalCommitted))
	ch <- prometheus.MustNewConstMetric(c.memoryTotalInit, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.MemoryTotalInit))
	ch <- prometheus.MustNewConstMetric(c.memoryTotalMax, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.MemoryTotalMax))
	ch <- prometheus.MustNewConstMetric(c.memoryTotalUsed, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.MemoryTotalUsed))

	ch <- prometheus.MustNewConstMetric(c.osAvailableProcessors, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.OSAvailableProcessors))
	ch <- prometheus.MustNewConstMetric(c.osCommittedVirtualMemorySize, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.OSCommittedVirtualMemorySize))
	ch <- prometheus.MustNewConstMetric(c.osFreePhysicalMemorySize, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.OSFreePhysicalMemorySize))
	ch <- prometheus.MustNewConstMetric(c.osFreeSwapSpaceSize, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.OSFreeSwapSpaceSize))
	ch <- prometheus.MustNewConstMetric(c.osMaxFileDescriptorCount, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.OSMaxFileDescriptorCount))
	ch <- prometheus.MustNewConstMetric(c.osOpenFileDescriptorCount, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.OSOpenFileDescriptorCount))
	ch <- prometheus.MustNewConstMetric(c.osProcessCPUTime, prometheus.CounterValue, float64(jvmStatus.Metrics.JVM.OSProcessCPUTime))
	ch <- prometheus.MustNewConstMetric(c.osSystemLoadAverage, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.OSSystemLoadAverage))
	ch <- prometheus.MustNewConstMetric(c.osTotalPhysicalMemorySize, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.OSTotalPhysicalMemorySize))
	ch <- prometheus.MustNewConstMetric(c.osTotalSwapSapceSize, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.OSTotalSwapSapceSize))

	ch <- prometheus.MustNewConstMetric(c.threadsBlockedCount, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.ThreadsBlockedCount))
	ch <- prometheus.MustNewConstMetric(c.threadsDaemonCount, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.ThreadsDaemonCount))
	ch <- prometheus.MustNewConstMetric(c.threadsDeadlockCount, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.ThreadsDeadlockCount))
	ch <- prometheus.MustNewConstMetric(c.threadsNewCount, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.ThreadsNewCount))
	ch <- prometheus.MustNewConstMetric(c.threadsRunnableCount, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.ThreadsRunnableCount))
	ch <- prometheus.MustNewConstMetric(c.threadsTerminatedCount, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.ThreadsTerminatedCount))
	ch <- prometheus.MustNewConstMetric(c.threadsTimedWaitingCount, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.ThreadsTimedWaitingCount))
	ch <- prometheus.MustNewConstMetric(c.threadsWaitingCount, prometheus.GaugeValue, float64(jvmStatus.Metrics.JVM.ThreadsWaitingCount))

	return nil
}

// Collect implements the prometheus.Collector interface.
func (c *JVMCollector) Collect(ch chan<- prometheus.Metric) {
	if err := c.Update(ch); err != nil {
		log.Errorf("Failed to collect metrics: %v", err)
	}
}

// Describe implements the prometheus.Collector interface.
func (c *JVMCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.gcConcurrentMarkSweepCount
	ch <- c.gcConcurrentMarkSweepTime
	ch <- c.gcParNewCount
	ch <- c.gcParNewTime

	ch <- c.memoryHeapCommitted
	ch <- c.memoryHeapInit
	ch <- c.memoryHeapMax
	ch <- c.memoryHeapUsage
	ch <- c.memoryHeapUsed

	ch <- c.memoryNonHeapCommitted
	ch <- c.memoryNonHeapInit
	ch <- c.memoryNonHeapMax
	ch <- c.memoryNonHeapUsage
	ch <- c.memoryNonHeapUsed

	ch <- c.memoryTotalCommitted
	ch <- c.memoryTotalInit
	ch <- c.memoryTotalMax
	ch <- c.memoryTotalUsed

	ch <- c.osAvailableProcessors
	ch <- c.osCommittedVirtualMemorySize
	ch <- c.osFreePhysicalMemorySize
	ch <- c.osFreeSwapSpaceSize
	ch <- c.osMaxFileDescriptorCount
	ch <- c.osOpenFileDescriptorCount
	ch <- c.osProcessCPUTime
	ch <- c.osSystemLoadAverage
	ch <- c.osTotalPhysicalMemorySize
	ch <- c.osTotalSwapSapceSize

	ch <- c.threadsBlockedCount
	ch <- c.threadsDaemonCount
	ch <- c.threadsDeadlockCount
	ch <- c.threadsNewCount
	ch <- c.threadsRunnableCount
	ch <- c.threadsTerminatedCount
	ch <- c.threadsTimedWaitingCount
	ch <- c.threadsWaitingCount
}
