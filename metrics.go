package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/tidwall/gjson"
)

var adminMetricsPath = "/admin/metrics?group=all&type=all&prefix=&property="

//MetricsCollector collects all metrics from solr
type MetricsCollector struct {
	jettyResponseTotal   *prometheus.Desc
	jettyRequestsTotal   *prometheus.Desc
	jettyDispatchesTotal *prometheus.Desc

	coreClientErrorsTotal                   *prometheus.Desc
	coreErrorsTotal                         *prometheus.Desc
	coreRequestsTotal                       *prometheus.Desc
	coreServerErrorsTotal                   *prometheus.Desc
	coreTimeoutsTotal                       *prometheus.Desc
	coreTimeSecondsTotal                    *prometheus.Desc
	coreFieldCacheTotal                     *prometheus.Desc
	coreHighlighterRequestTotal             *prometheus.Desc
	coreIndexSizeBytes                      *prometheus.Desc
	coreReplicationMaster                   *prometheus.Desc
	coreReplicationSlave                    *prometheus.Desc
	coreReplicationLastSuccess              *prometheus.Desc
	coreReplicationLastFail                 *prometheus.Desc
	coreReplicationSuccessCount             *prometheus.Desc
	coreReplicationFailCount                *prometheus.Desc
	coreReplicationReplicating              *prometheus.Desc
	coreReplicationLastCycleDownloadedBytes *prometheus.Desc
	coreSearcherDocuments                   *prometheus.Desc
	coreUpdateHandlerAdds                   *prometheus.Desc
	coreUpdateHandlerAddsTotal              *prometheus.Desc
	coreUpdateHandlerAutoCommitsTotal       *prometheus.Desc
	coreUpdateHandlerCommitsTotal           *prometheus.Desc
	coreUpdateHandlerDeletesByID            *prometheus.Desc
	coreUpdateHandlerDeletesByIDTotal       *prometheus.Desc
	coreUpdateHandlerDeletesByQuery         *prometheus.Desc
	coreUpdateHandlerDeletesByQueryTotal    *prometheus.Desc
	coreUpdateHandlerErrors                 *prometheus.Desc
	coreUpdateHandlerErrorsTotal            *prometheus.Desc
	coreUpdateHandlerExpungeDeletesTotal    *prometheus.Desc
	coreUpdateHandlerMergesTotal            *prometheus.Desc
	coreUpdateHandlerOptimizesTotal         *prometheus.Desc
	coreUpdateHandlerPendingDocs            *prometheus.Desc
	coreUpdateHandlerRollbacksTotal         *prometheus.Desc
	coreUpdateHandlerSoftAutoCommitsTotal   *prometheus.Desc
	coreUpdateHandlerSplitsTotal            *prometheus.Desc

	coreFSBytes *prometheus.Desc

	coreSearcherCache                *prometheus.Desc
	coreSearcherCacheRatio           *prometheus.Desc
	coreSearcherWarmupTimeSeconds    *prometheus.Desc
	coreSearcherCumulativeCacheTotal *prometheus.Desc
	coreSearcherCumulativeCacheRatio *prometheus.Desc

	coreRequestsP75ms    *prometheus.Desc
	coreRequestsP95ms    *prometheus.Desc
	coreRequestsP99ms    *prometheus.Desc
	coreRequestsStddevMs *prometheus.Desc
	coreRequestsMeanMs   *prometheus.Desc
	coreRequestsMedianMs *prometheus.Desc

	jvmBuffers            *prometheus.Desc
	jvmBuffersBytes       *prometheus.Desc
	jvmGCTotal            *prometheus.Desc
	jvmGCSecondsTotal     *prometheus.Desc
	jvmMemoryHeapBytes    *prometheus.Desc
	jvmMemoryNonHeapBytes *prometheus.Desc
	jvmMemoryPoolsBytes   *prometheus.Desc
	jvmMemoryBytes        *prometheus.Desc
	jvmOsMemoryBytes      *prometheus.Desc
	jvmOsFileDescriptors  *prometheus.Desc
	jvmOsCPULoad          *prometheus.Desc
	jvmOsCPUTimeSeconds   *prometheus.Desc
	jvmOsLoadAverage      *prometheus.Desc
	jvmThreads            *prometheus.Desc

	nodeClientErrorsTotal        *prometheus.Desc
	nodeErrorsTotal              *prometheus.Desc
	nodeRequestTotal             *prometheus.Desc
	nodeServerErrors             *prometheus.Desc
	nodeTimeoutsTotal            *prometheus.Desc
	nodeTimeSecondsTotal         *prometheus.Desc
	nodeCores                    *prometheus.Desc
	nodeCoreRootFsBytes          *prometheus.Desc
	nodeThreadPoolCompletedTotal *prometheus.Desc
	nodeThreadPoolRunning        *prometheus.Desc
	nodeThreadPoolSubmittedTotal *prometheus.Desc
	nodeConnections              *prometheus.Desc

	client          http.Client
	adminMetricsURL string
	solrBaseURL     string
}

// NewMetricsCollector returns a new Collector exposing solr metrics statistics.
func NewMetricsCollector(client http.Client, solrBaseURL string) (*MetricsCollector, error) {
	adminMetricsURL := fmt.Sprintf("%s%s", solrBaseURL, adminMetricsPath)
	return &MetricsCollector{
		jettyResponseTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jetty_response_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"status"}, nil,
		),
		jettyRequestsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jetty_requests_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"method"}, nil,
		),
		jettyDispatchesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jetty_dispatches_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{}, nil,
		),
		coreClientErrorsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_client_errors_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreErrorsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_errors_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreRequestsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_requests_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreServerErrorsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_server_errors_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreTimeoutsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_timeouts_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreTimeSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_time_seconds_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreFieldCacheTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_field_cache_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "core", "collection", "shard", "replica"}, nil,
		),
		coreRequestsP75ms: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_requests_p75_ms"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreRequestsP95ms: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_requests_p95_ms"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreRequestsP99ms: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_requests_p99_ms"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreRequestsStddevMs: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_requests_stddev_ms"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreRequestsMedianMs: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_requests_median_ms"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreRequestsMeanMs: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_requests_mean_ms"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreSearcherCacheRatio: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_searcher_cache_ratio"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "core", "type", "item", "collection", "shard", "replica"}, nil,
		),
		coreSearcherCache: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_searcher_cache"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "core", "type", "item", "collection", "shard", "replica"}, nil,
		),
		coreSearcherWarmupTimeSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_searcher_warmup_time_seconds"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "core", "type", "item", "collection", "shard", "replica"}, nil,
		),
		coreSearcherCumulativeCacheTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_searcher_cumulative_cache_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "core", "type", "item", "collection", "shard", "replica"}, nil,
		),
		coreSearcherCumulativeCacheRatio: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_searcher_cumulative_cache_ratio"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "core", "type", "item", "collection", "shard", "replica"}, nil,
		),
		coreFSBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_fs_bytes"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "core", "item", "collection", "shard", "replica"}, nil,
		),
		coreHighlighterRequestTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_highlighter_request_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "core", "item", "name", "collection", "shard", "replica"}, nil,
		),
		coreIndexSizeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_index_size_bytes"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "core", "collection", "shard", "replica"}, nil,
		),
		coreReplicationMaster: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_replication_master"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreReplicationSlave: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_replication_slave"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreReplicationLastSuccess: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_replication_last_success"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreReplicationLastFail: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_replication_last_fail"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreReplicationSuccessCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_replication_success_count"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreReplicationFailCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_replication_fail_count"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreReplicationReplicating: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_replication_replicating"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreReplicationLastCycleDownloadedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_replication_last_cycle_downloaded_bytes"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreSearcherDocuments: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_searcher_documents"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "core", "item", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerAdds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_adds"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerAutoCommitsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_auto_commits_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerCommitsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_commits_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerAddsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_adds_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerDeletesByIDTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_deletes_by_id_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerDeletesByQueryTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_deletes_by_query_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerErrorsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_errors_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerDeletesByID: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_deletes_by_id"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerDeletesByQuery: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_deletes_by_query"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerPendingDocs: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_pending_docs"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerErrors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_errors"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerExpungeDeletesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_expunge_deletes_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerMergesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_merges_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerOptimizesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_optimizes_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerRollbacksTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_rollbacks_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerSoftAutoCommitsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_soft_auto_commits_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),
		coreUpdateHandlerSplitsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "core_update_handler_splits_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "core", "collection", "shard", "replica"}, nil,
		),

		jvmBuffers: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jvm_buffers"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"pool"}, nil,
		),
		jvmBuffersBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jvm_buffers_bytes"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"pool", "item"}, nil,
		),
		jvmGCTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jvm_gc_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"item"}, nil,
		),
		jvmGCSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jvm_gc_seconds_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"item"}, nil,
		),
		jvmMemoryHeapBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jvm_memory_heap_bytes"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"item"}, nil,
		),
		jvmMemoryNonHeapBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jvm_memory_non_heap_bytes"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"item"}, nil,
		),
		jvmMemoryPoolsBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jvm_memory_pools_bytes"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"space", "item"}, nil,
		),
		jvmMemoryBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jvm_memory_bytes"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"item"}, nil,
		),
		jvmOsMemoryBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jvm_os_memory_bytes"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"item"}, nil,
		),
		jvmOsFileDescriptors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jvm_os_file_descriptors"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"item"}, nil,
		),
		jvmOsCPULoad: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jvm_os_cpu_load"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"item"}, nil,
		),
		jvmOsCPUTimeSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jvm_os_cpu_time_seconds"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"item"}, nil,
		),
		jvmOsLoadAverage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jvm_os_load_average"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"item"}, nil,
		),
		jvmThreads: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "jvm_threads"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"item"}, nil,
		),

		nodeClientErrorsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "node_client_errors_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler"}, nil,
		),
		nodeErrorsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "node_errors_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler"}, nil,
		),
		nodeTimeoutsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "node_timeouts_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler"}, nil,
		),
		nodeRequestTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "node_requests_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler"}, nil,
		),
		nodeServerErrors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "node_server_errors_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler"}, nil,
		),
		nodeTimeSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "node_time_seconds_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler"}, nil,
		),
		nodeCores: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "node_cores"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "item"}, nil,
		),
		nodeCoreRootFsBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "node_core_root_fs_bytes"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "item"}, nil,
		),
		nodeThreadPoolCompletedTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "node_thread_pool_completed_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "executor"}, nil,
		),
		nodeThreadPoolRunning: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "node_thread_pool_running"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "executor"}, nil,
		),
		nodeThreadPoolSubmittedTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "node_thread_pool_submitted_total"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "executor"}, nil,
		),
		nodeConnections: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "metrics", "node_connections"),
			"See following URL: https://lucene.apache.org/solr/guide/metrics-reporting.html",
			[]string{"category", "handler", "item"}, nil,
		),
		client:          client,
		adminMetricsURL: adminMetricsURL,
		solrBaseURL:     solrBaseURL,
	}, nil
}

// Update exposes all related metrics from solr.
func (c *MetricsCollector) Update(ch chan<- prometheus.Metric) error {
	metricsIO, err := fetchHTTP(c.client, c.adminMetricsURL)
	if err != nil {
		return err
	}

	metrics := gjson.Parse(string(metricsIO))

	// jetty metrics
	for metricName, value := range metrics.Get("metrics.solr*jetty").Map() {
		if strings.HasPrefix(metricName, "org.eclipse.jetty.server.handler.DefaultHandler") && strings.HasSuffix(metricName, "xx-responses") {
			metricNameSplitted := strings.Split(metricName, ".")
			status := strings.Split(metricNameSplitted[len(metricNameSplitted)-1], "-")
			if !value.Get("count").Exists() {
				return fmt.Errorf("Fail to find count value in jetty responses metrics JSON : %s", value)
			}
			ch <- prometheus.MustNewConstMetric(c.jettyResponseTotal, prometheus.CounterValue, float64(value.Get("count").Int()), status[0])
		}
		if strings.HasPrefix(metricName, "org.eclipse.jetty.server.handler.DefaultHandler.") && strings.HasSuffix(metricName, "-requests") && value.Get("count").Exists() {
			metricNameSplitted := strings.Split(metricName, ".")
			method := strings.Split(metricNameSplitted[len(metricNameSplitted)-1], "-")
			ch <- prometheus.MustNewConstMetric(c.jettyRequestsTotal, prometheus.CounterValue, float64(value.Get("count").Int()), method[0])
		}
		if strings.HasPrefix(metricName, "org.eclipse.jetty.server.handler.DefaultHandler.dispatches") {
			if !value.Get("count").Exists() {
				return fmt.Errorf("Fail to find count value in jetty dispatches metrics JSON : %s", value)
			}
			ch <- prometheus.MustNewConstMetric(c.jettyDispatchesTotal, prometheus.CounterValue, float64(value.Get("count").Int()))
		}
	}

	// metrics
	metrics.Get("metrics").ForEach(func(metricName, value gjson.Result) bool {
		// jvm metrics
		if strings.HasPrefix(metricName.String(), "solr.jvm") {
			value.ForEach(func(key, value gjson.Result) bool {
				splittedKey := strings.Split(key.String(), ".")
				if strings.HasPrefix(key.String(), "buffers.") && strings.HasSuffix(key.String(), ".Count") {
					ch <- prometheus.MustNewConstMetric(c.jvmBuffers, prometheus.GaugeValue, value.Float(), splittedKey[1])
				}
				if strings.HasPrefix(key.String(), "buffers.") && (strings.HasSuffix(key.String(), ".MemoryUsed") || strings.HasSuffix(key.String(), ".TotalCapacity")) {
					ch <- prometheus.MustNewConstMetric(c.jvmBuffersBytes, prometheus.GaugeValue, value.Float(), splittedKey[1], splittedKey[len(splittedKey)-1])
				}
				if strings.HasPrefix(key.String(), "gc.") && strings.HasSuffix(key.String(), ".count") {
					ch <- prometheus.MustNewConstMetric(c.jvmGCTotal, prometheus.CounterValue, value.Float(), splittedKey[1])
				}
				if strings.HasPrefix(key.String(), "gc.") && strings.HasSuffix(key.String(), ".time") {
					ch <- prometheus.MustNewConstMetric(c.jvmGCSecondsTotal, prometheus.CounterValue, value.Float()/1000, splittedKey[1])
				}
				if strings.HasPrefix(key.String(), "memory.heap.") && !strings.HasSuffix(key.String(), ".usage") {
					ch <- prometheus.MustNewConstMetric(c.jvmMemoryHeapBytes, prometheus.GaugeValue, value.Float(), splittedKey[len(splittedKey)-1])
				}
				if strings.HasPrefix(key.String(), "memory.non-heap.") && !strings.HasSuffix(key.String(), ".usage") {
					ch <- prometheus.MustNewConstMetric(c.jvmMemoryNonHeapBytes, prometheus.GaugeValue, value.Float(), splittedKey[len(splittedKey)-1])
				}
				if strings.HasPrefix(key.String(), "memory.pools.") && !strings.HasSuffix(key.String(), ".usage") {
					ch <- prometheus.MustNewConstMetric(c.jvmMemoryPoolsBytes, prometheus.GaugeValue, value.Float(), splittedKey[2], splittedKey[len(splittedKey)-1])
				}
				if strings.HasPrefix(key.String(), "memory.total.") {
					ch <- prometheus.MustNewConstMetric(c.jvmMemoryBytes, prometheus.GaugeValue, value.Float(), splittedKey[len(splittedKey)-1])
				}
				if key.String() == "os.committedVirtualMemorySize" || key.String() == "os.freePhysicalMemorySize" || key.String() == "os.freeSwapSpaceSize" || key.String() == "os.totalPhysicalMemorySize" || key.String() == "os.totalSwapSpaceSize" {
					ch <- prometheus.MustNewConstMetric(c.jvmOsMemoryBytes, prometheus.GaugeValue, value.Float(), splittedKey[len(splittedKey)-1])
				}
				if key.String() == "os.maxFileDescriptorCount" || key.String() == "os.openFileDescriptorCount" {
					ch <- prometheus.MustNewConstMetric(c.jvmOsFileDescriptors, prometheus.GaugeValue, value.Float(), splittedKey[len(splittedKey)-1])
				}
				if key.String() == "os.processCpuLoad" || key.String() == "os.systemCpuLoad" {
					ch <- prometheus.MustNewConstMetric(c.jvmOsCPULoad, prometheus.GaugeValue, value.Float(), splittedKey[len(splittedKey)-1])
				}
				if key.String() == "os.processCpuTime" {
					ch <- prometheus.MustNewConstMetric(c.jvmOsCPUTimeSeconds, prometheus.CounterValue, value.Float()/1000.0, "processCpuTime")
				}
				if key.String() == "os.systemLoadAverage" {
					ch <- prometheus.MustNewConstMetric(c.jvmOsLoadAverage, prometheus.GaugeValue, value.Float()/1000.0, "systemLoadAverage")
				}
				if strings.HasPrefix(key.String(), "threads.") && strings.HasSuffix(key.String(), ".count") {
					ch <- prometheus.MustNewConstMetric(c.jvmThreads, prometheus.GaugeValue, value.Float(), splittedKey[1])
				}
				return true
			})
		}
		// node metrics
		if strings.HasPrefix(metricName.String(), "solr.node") {
			value.ForEach(func(key, value gjson.Result) bool {
				splittedKey := strings.Split(key.String(), ".")
				category := splittedKey[0]
				handler := splittedKey[1]
				if strings.HasSuffix(key.String(), ".clientErrors") {
					ch <- prometheus.MustNewConstMetric(c.nodeClientErrorsTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler)
					ch <- prometheus.MustNewConstMetric(c.nodeErrorsTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler)
				}
				if strings.HasSuffix(key.String(), ".timeouts") {
					ch <- prometheus.MustNewConstMetric(c.nodeTimeoutsTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler)
				}
				if strings.HasSuffix(key.String(), ".serverErrors") {
					ch <- prometheus.MustNewConstMetric(c.nodeServerErrors, prometheus.CounterValue, value.Get("count").Float(), category, handler)
				}
				if strings.HasSuffix(key.String(), ".requestTimes") {
					ch <- prometheus.MustNewConstMetric(c.nodeRequestTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler)
				}
				if strings.HasSuffix(key.String(), ".totalTime") {
					ch <- prometheus.MustNewConstMetric(c.nodeTimeSecondsTotal, prometheus.CounterValue, value.Float()/1000.0, category, handler)
				}
				if strings.HasPrefix(key.String(), "CONTAINER.cores.") {
					ch <- prometheus.MustNewConstMetric(c.nodeCores, prometheus.GaugeValue, value.Float(), category, splittedKey[2])
				}
				if strings.HasPrefix(key.String(), "CONTAINER.fs.coreRoot.") && (strings.HasSuffix(key.String(), ".totalSpace") || strings.HasSuffix(key.String(), ".usableSpace")) {
					ch <- prometheus.MustNewConstMetric(c.nodeCoreRootFsBytes, prometheus.GaugeValue, value.Float(), category, splittedKey[3])
				}
				if strings.Contains(key.String(), ".threadPool.") && strings.HasSuffix(key.String(), ".completed") {
					executor := ""
					if len(splittedKey) >= 5 {
						executor = splittedKey[3]
					} else {
						executor = splittedKey[2]
					}
					hand := ""
					if len(splittedKey) >= 5 {
						hand = handler
					}
					ch <- prometheus.MustNewConstMetric(c.nodeThreadPoolCompletedTotal, prometheus.CounterValue, value.Float(), category, hand, executor)
				}
				if strings.Contains(key.String(), ".threadPool.") && strings.HasSuffix(key.String(), ".running") {
					executor := ""
					if len(splittedKey) >= 5 {
						executor = splittedKey[3]
					} else {
						executor = splittedKey[2]
					}
					hand := ""
					if len(splittedKey) >= 5 {
						hand = handler
					}
					ch <- prometheus.MustNewConstMetric(c.nodeThreadPoolRunning, prometheus.CounterValue, value.Float(), category, hand, executor)
				}
				if strings.Contains(key.String(), ".threadPool.") && strings.HasSuffix(key.String(), ".submitted") {
					executor := ""
					if len(splittedKey) >= 5 {
						executor = splittedKey[3]
					} else {
						executor = splittedKey[2]
					}
					hand := ""
					if len(splittedKey) >= 5 {
						hand = handler
					}
					ch <- prometheus.MustNewConstMetric(c.nodeThreadPoolSubmittedTotal, prometheus.CounterValue, value.Float(), category, hand, executor)
				}
				if strings.HasSuffix(key.String(), "Connections") {
					ch <- prometheus.MustNewConstMetric(c.nodeConnections, prometheus.GaugeValue, value.Float(), category, handler, splittedKey[2])
				}
				return true
			})
		}
		// core metrics
		if strings.HasPrefix(metricName.String(), "solr.core.") {
			value.ForEach(func(path, value gjson.Result) bool {
				splittedPath := strings.Split(path.String(), ".")
				splittedMetricName := strings.Split(metricName.String(), ".")
				core := splittedMetricName[2]
				category := splittedPath[0]
				handler := splittedPath[1]
				if strings.HasSuffix(path.String(), ".requestTimes") {
					if strings.HasPrefix(handler, "/") {
						ch <- prometheus.MustNewConstMetric(c.coreRequestsTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler, core, "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreRequestsP75ms, prometheus.GaugeValue, value.Get("p75_ms").Float(), category, handler, core, "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreRequestsP95ms, prometheus.GaugeValue, value.Get("p95_ms").Float(), category, handler, core, "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreRequestsP99ms, prometheus.GaugeValue, value.Get("p99_ms").Float(), category, handler, core, "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreRequestsMeanMs, prometheus.GaugeValue, value.Get("mean_ms").Float(), category, handler, core, "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreRequestsMedianMs, prometheus.GaugeValue, value.Get("median_ms").Float(), category, handler, core, "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreRequestsStddevMs, prometheus.GaugeValue, value.Get("stddev_ms").Float(), category, handler, core, "", "", "")
					}
				}
				if strings.HasSuffix(path.String(), ".clientErrors") {
					if strings.HasPrefix(handler, "/") {
						ch <- prometheus.MustNewConstMetric(c.coreClientErrorsTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler, core, "", "", "")
					}
				}
				if strings.HasSuffix(path.String(), ".errors") {
					if strings.HasPrefix(handler, "/") {
						ch <- prometheus.MustNewConstMetric(c.coreErrorsTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler, core, "", "", "")
					}
				}
				if strings.HasSuffix(path.String(), ".serverErrors") {
					if strings.HasPrefix(handler, "/") {
						ch <- prometheus.MustNewConstMetric(c.coreServerErrorsTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler, core, "", "", "")
					}
				}
				if strings.HasSuffix(path.String(), ".timeouts") {
					if strings.HasPrefix(handler, "/") {
						ch <- prometheus.MustNewConstMetric(c.coreTimeoutsTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler, core, "", "", "")
					}
				}
				if strings.HasSuffix(path.String(), ".totalTime") {
					if strings.HasPrefix(handler, "/") {
						ch <- prometheus.MustNewConstMetric(c.coreTimeSecondsTotal, prometheus.CounterValue, value.Float()/1000, category, handler, core, "", "", "")
					}
				}
				if path.String() == "CACHE.core.fieldCache" {
					ch <- prometheus.MustNewConstMetric(c.coreFieldCacheTotal, prometheus.CounterValue, value.Get("entries_count").Float(), category, core, "", "", "")
				}
				if strings.HasPrefix(path.String(), "CACHE.searcher") {
					if strings.HasSuffix(path.String(), "documentCache") || strings.HasSuffix(path.String(), "fieldValueCache") || strings.HasSuffix(path.String(), "filterCache") || strings.HasSuffix(path.String(), "perSegFilter") || strings.HasSuffix(path.String(), "queryResultCache") {
						ch <- prometheus.MustNewConstMetric(c.coreSearcherCacheRatio, prometheus.GaugeValue, value.Get("hitratio").Float(), category, core, splittedPath[2], "hitratio", "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreSearcherCache, prometheus.GaugeValue, value.Get("lookups").Float(), category, core, splittedPath[2], "lookups", "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreSearcherCache, prometheus.GaugeValue, value.Get("hits").Float(), category, core, splittedPath[2], "hits", "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreSearcherCache, prometheus.GaugeValue, value.Get("size").Float(), category, core, splittedPath[2], "size", "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreSearcherCache, prometheus.GaugeValue, value.Get("evictions").Float(), category, core, splittedPath[2], "evictions", "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreSearcherCache, prometheus.GaugeValue, value.Get("inserts").Float(), category, core, splittedPath[2], "inserts", "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreSearcherWarmupTimeSeconds, prometheus.GaugeValue, value.Get("warmupTime").Float()/1000, category, core, splittedPath[2], "warmupTime", "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreSearcherCumulativeCacheTotal, prometheus.CounterValue, value.Get("cumulative_lookups").Float(), category, core, splittedPath[2], "cumulative_lookups", "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreSearcherCumulativeCacheTotal, prometheus.CounterValue, value.Get("cumulative_hits").Float(), category, core, splittedPath[2], "cumulative_hits", "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreSearcherCumulativeCacheTotal, prometheus.CounterValue, value.Get("cumulative_evictions").Float(), category, core, splittedPath[2], "cumulative_evictions", "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreSearcherCumulativeCacheTotal, prometheus.CounterValue, value.Get("cumulative_inserts").Float(), category, core, splittedPath[2], "cumulative_inserts", "", "", "")
						ch <- prometheus.MustNewConstMetric(c.coreSearcherCumulativeCacheRatio, prometheus.GaugeValue, value.Get("cumulative_hitratio").Float(), category, core, splittedPath[2], "cumulative_hitratio", "", "", "")
					}
				}
				if strings.HasPrefix(path.String(), "CORE.fs") {
					if strings.HasSuffix(path.String(), ".totalSpace") || strings.HasSuffix(path.String(), ".usableSpace") {
						ch <- prometheus.MustNewConstMetric(c.coreFSBytes, prometheus.CounterValue, value.Float(), category, core, splittedPath[2], "", "", "")
					}
				}
				if strings.HasPrefix(path.String(), "HIGHLIGHTER.") && strings.HasSuffix(path.String(), ".requests") {
					ch <- prometheus.MustNewConstMetric(c.coreHighlighterRequestTotal, prometheus.CounterValue, value.Float(), category, core, splittedPath[1], splittedPath[2], "", "", "")
				}
				if path.String() == "INDEX.sizeInBytes" {
					ch <- prometheus.MustNewConstMetric(c.coreIndexSizeBytes, prometheus.GaugeValue, value.Float(), category, core, "", "", "")
				}
				if path.String() == "REPLICATION./replication.isMaster" {
					ch <- prometheus.MustNewConstMetric(c.coreReplicationMaster, prometheus.GaugeValue, value.Float(), category, handler, core, "", "", "")
				}
				if path.String() == "REPLICATION./replication.isSlave" {
					ch <- prometheus.MustNewConstMetric(c.coreReplicationSlave, prometheus.GaugeValue, value.Float(), category, handler, core, "", "", "")
				}
				if path.String() == "REPLICATION./replication.fetcher" {
					successReplicationStr := value.Get("indexReplicatedAt").String()
					successReplicationTime, err := time.Parse(time.UnixDate, successReplicationStr)
					successReplication := 0.0
					if err == nil {
						successReplication = float64(successReplicationTime.Unix())
					}
					failReplicationStr := value.Get("replicationFailedAt").String()
					failReplicationTime, err := time.Parse(time.UnixDate, failReplicationStr)
					failReplication := 0.0
					if err == nil {
						failReplication = float64(failReplicationTime.Unix())
					}

					ch <- prometheus.MustNewConstMetric(c.coreReplicationLastSuccess, prometheus.GaugeValue, successReplication, category, handler, core, "", "", "")
					ch <- prometheus.MustNewConstMetric(c.coreReplicationLastFail, prometheus.GaugeValue, failReplication, category, handler, core, "", "", "")
					ch <- prometheus.MustNewConstMetric(c.coreReplicationSuccessCount, prometheus.GaugeValue, value.Get("timesIndexReplicated").Float(), category, handler, core, "", "", "")
					ch <- prometheus.MustNewConstMetric(c.coreReplicationFailCount, prometheus.GaugeValue, value.Get("timesFailed").Float(), category, handler, core, "", "", "")
					ch <- prometheus.MustNewConstMetric(c.coreReplicationReplicating, prometheus.GaugeValue, value.Get("isReplicating").Float(), category, handler, core, "", "", "")
					ch <- prometheus.MustNewConstMetric(c.coreReplicationLastCycleDownloadedBytes, prometheus.GaugeValue, value.Get("lastCycleBytesDownloaded").Float(), category, handler, core, "", "", "")
				}
				if path.String() == "SEARCHER.searcher.deletedDocs" || path.String() == "SEARCHER.searcher.maxDoc" || path.String() == "SEARCHER.searcher.numDocs" {
					ch <- prometheus.MustNewConstMetric(c.coreSearcherDocuments, prometheus.GaugeValue, value.Float(), category, core, splittedPath[2], "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.adds" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerAdds, prometheus.GaugeValue, value.Float(), category, handler, core, "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.autoCommits" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerAutoCommitsTotal, prometheus.CounterValue, value.Float(), category, handler, core, "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.commits" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerCommitsTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler, core, "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.cumulativeAdds" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerAddsTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler, core, "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.cumulativeDeletesById" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerDeletesByIDTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler, core, "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.cumulativeDeletesByQuery" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerDeletesByQueryTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler, core, "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.cumulativeErrors" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerErrorsTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler, core, "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.deletesById" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerDeletesByID, prometheus.GaugeValue, value.Float(), category, handler, core, "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.deletesByQuery" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerDeletesByQuery, prometheus.GaugeValue, value.Float(), category, handler, core, "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.docsPending" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerPendingDocs, prometheus.GaugeValue, value.Float(), category, handler, core, "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.errors" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerErrors, prometheus.GaugeValue, value.Float(), category, handler, core, "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.expungeDeletes" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerExpungeDeletesTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler, core, "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.merges" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerMergesTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler, core, "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.optimizes" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerOptimizesTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler, core, "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.rollbacks" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerRollbacksTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler, core, "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.softAutoCommits" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerSoftAutoCommitsTotal, prometheus.CounterValue, value.Float(), category, handler, core, "", "", "")
				}
				if path.String() == "UPDATE.updateHandler.splits" {
					ch <- prometheus.MustNewConstMetric(c.coreUpdateHandlerSplitsTotal, prometheus.CounterValue, value.Get("count").Float(), category, handler, core, "", "", "")
				}

				return true
			})
		}
		return true
	})

	return nil
}

// Collect implements the prometheus.Collector interface.
func (c *MetricsCollector) Collect(ch chan<- prometheus.Metric) {
	if err := c.Update(ch); err != nil {
		log.Errorf("Failed to collect metrics: %v", err)
	}
}

// Describe implements the prometheus.Collector interface.
func (c *MetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.jettyRequestsTotal
	ch <- c.jettyResponseTotal
	ch <- c.jettyDispatchesTotal

	ch <- c.coreClientErrorsTotal
	ch <- c.coreErrorsTotal
	ch <- c.coreRequestsTotal
	ch <- c.coreServerErrorsTotal
	ch <- c.coreTimeoutsTotal
	ch <- c.coreTimeSecondsTotal
	ch <- c.coreFieldCacheTotal
	ch <- c.coreHighlighterRequestTotal
	ch <- c.coreIndexSizeBytes
	ch <- c.coreReplicationMaster
	ch <- c.coreReplicationSlave
	ch <- c.coreReplicationLastSuccess
	ch <- c.coreReplicationLastFail
	ch <- c.coreReplicationSuccessCount
	ch <- c.coreReplicationFailCount
	ch <- c.coreReplicationReplicating
	ch <- c.coreReplicationLastCycleDownloadedBytes
	ch <- c.coreSearcherDocuments
	ch <- c.coreUpdateHandlerAdds
	ch <- c.coreUpdateHandlerAddsTotal
	ch <- c.coreUpdateHandlerAutoCommitsTotal
	ch <- c.coreUpdateHandlerCommitsTotal
	ch <- c.coreUpdateHandlerDeletesByID
	ch <- c.coreUpdateHandlerDeletesByIDTotal
	ch <- c.coreUpdateHandlerDeletesByQuery
	ch <- c.coreUpdateHandlerDeletesByQueryTotal
	ch <- c.coreUpdateHandlerErrors
	ch <- c.coreUpdateHandlerErrorsTotal
	ch <- c.coreUpdateHandlerExpungeDeletesTotal
	ch <- c.coreUpdateHandlerMergesTotal
	ch <- c.coreUpdateHandlerOptimizesTotal
	ch <- c.coreUpdateHandlerPendingDocs
	ch <- c.coreUpdateHandlerRollbacksTotal
	ch <- c.coreUpdateHandlerSoftAutoCommitsTotal
	ch <- c.coreUpdateHandlerSplitsTotal

	ch <- c.coreSearcherCache
	ch <- c.coreSearcherCacheRatio
	ch <- c.coreSearcherWarmupTimeSeconds
	ch <- c.coreSearcherCumulativeCacheTotal
	ch <- c.coreSearcherCumulativeCacheRatio

	ch <- c.coreFSBytes

	ch <- c.coreRequestsP75ms
	ch <- c.coreRequestsP95ms
	ch <- c.coreRequestsP99ms
	ch <- c.coreRequestsMeanMs
	ch <- c.coreRequestsMedianMs
	ch <- c.coreRequestsStddevMs

	ch <- c.jvmBuffers
	ch <- c.jvmBuffersBytes
	ch <- c.jvmGCTotal
	ch <- c.jvmGCSecondsTotal
	ch <- c.jvmMemoryHeapBytes
	ch <- c.jvmMemoryNonHeapBytes
	ch <- c.jvmMemoryPoolsBytes
	ch <- c.jvmMemoryBytes
	ch <- c.jvmOsMemoryBytes
	ch <- c.jvmOsFileDescriptors
	ch <- c.jvmOsCPULoad
	ch <- c.jvmOsCPUTimeSeconds
	ch <- c.jvmOsLoadAverage
	ch <- c.jvmThreads

	ch <- c.nodeClientErrorsTotal
	ch <- c.nodeErrorsTotal
	ch <- c.nodeRequestTotal
	ch <- c.nodeServerErrors
	ch <- c.nodeTimeoutsTotal
	ch <- c.nodeTimeSecondsTotal
	ch <- c.nodeCores
	ch <- c.nodeCoreRootFsBytes
	ch <- c.nodeThreadPoolCompletedTotal
	ch <- c.nodeThreadPoolRunning
	ch <- c.nodeThreadPoolSubmittedTotal
	ch <- c.nodeConnections
}
