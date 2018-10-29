package main

import "encoding/json"

type AdminCoresStatus struct {
	Status map[string]struct {
		Index struct {
			SizeInBytes int64 `json:"sizeInBytes"`
			NumDocs     int   `json:"numDocs"`
			MaxDoc      int   `json:"maxDoc"`
			DeletedDocs int   `json:"deletedDocs"`
		} `json:"index"`
	} `json:"status"`
}

type MBeansData struct {
	Headers    ResponseHeader    `json:"responseHeader"`
	SolrMbeans []json.RawMessage `json:"solr-mbeans"`
}

type ResponseHeader struct {
	QTime  int `json:"QTime"`
	Status int `json:"status"`
}

type Core struct {
	Class string `json:"class"`
	Stats struct {
		DeletedDocs int `json:"deletedDocs"`
		MaxDoc      int `json:"maxDoc"`
		NumDocs     int `json:"numDocs"`
	} `json:"stats"`
}

type QueryHandler struct {
	Class string `json:"class"`
	Stats struct {
		One5minRateReqsPerSecond     *float64 `json:"15minRateReqsPerSecond,omitempty"`
		FiveMinRateReqsPerSecond     *float64 `json:"5MinRateReqsPerSecond,omitempty"`
		One5minRateRequestsPerSecond float64  `json:"15minRateRequestsPerSecond,omitempty"`
		FiveminRateRequestsPerSecond float64  `json:"5minRateRequestsPerSecond,omitempty"`
		Seven5thPcRequestTime        float64  `json:"75thPcRequestTime"`
		Nine5thPcRequestTime         float64  `json:"95thPcRequestTime"`
		Nine9thPcRequestTime         float64  `json:"99thPcRequestTime"`
		Nine99thPcRequestTime        float64  `json:"999thPcRequestTime"`
		AvgRequestsPerSecond         float64  `json:"avgRequestsPerSecond"`
		AvgTimePerRequest            float64  `json:"avgTimePerRequest"`
		Errors                       int      `json:"errors"`
		HandlerStart                 int      `json:"handlerStart"`
		MedianRequestTime            float64  `json:"medianRequestTime"`
		Requests                     int      `json:"requests"`
		Timeouts                     int      `json:"timeouts"`
		TotalTime                    float64  `json:"totalTime"`
	} `json:"stats"`
}

type UpdateHandler struct {
	Class string `json:"class"`
	Stats struct {
		Adds                     int    `json:"adds"`
		AutocommitMaxDocs        int    `json:"autocommit maxDocs"`
		AutocommitMaxTime        string `json:"autocommit maxTime"`
		Autocommits              int    `json:"autocommits"`
		Commits                  int    `json:"commits"`
		CumulativeAdds           int    `json:"cumulative_adds"`
		CumulativeDeletesByID    int    `json:"cumulative_deletesById"`
		CumulativeDeletesByQuery int    `json:"cumulative_deletesByQuery"`
		CumulativeErrors         int    `json:"cumulative_errors"`
		DeletesByID              int    `json:"deletesById"`
		DeletesByQuery           int    `json:"deletesByQuery"`
		DocsPending              int    `json:"docsPending"`
		Errors                   int    `json:"errors"`
		ExpungeDeletes           int    `json:"expungeDeletes"`
		Optimizes                int    `json:"optimizes"`
		Rollbacks                int    `json:"rollbacks"`
		SoftAutocommits          int    `json:"soft autocommits"`
	} `json:"stats"`
}

type Cache struct {
	Class string `json:"class"`
	Stats struct {
		CumulativeEvictions int         `json:"cumulative_evictions"`
		CumulativeHitratio  json.Number `json:"cumulative_hitratio,Number"`
		CumulativeHits      int         `json:"cumulative_hits"`
		CumulativeInserts   int         `json:"cumulative_inserts"`
		CumulativeLookups   int         `json:"cumulative_lookups"`
		Evictions           int         `json:"evictions"`
		Hitratio            json.Number `json:"hitratio,Number"`
		Hits                int         `json:"hits"`
		Inserts             int         `json:"inserts"`
		Lookups             int         `json:"lookups"`
		Size                int         `json:"size"`
		WarmupTime          int         `json:"warmupTime"`
	} `json:"stats"`
}

type JVMStatus struct {
	Metrics struct {
		JVM struct {
			GCConcurrentMarkSweepCount int64 `json:"gc.ConcurrentMarkSweep.count"`
			GCConcurrentMarkSweepTime  int64 `json:"gc.ConcurrentMarkSweep.time"`
			GCParNewCount              int64 `json:"gc.ParNew.count"`
			GCParNewTime               int64 `json:"gc.ParNew.time"`

			MemoryHeapCommitted int64   `json:"memory.heap.committed"`
			MemoryHeapInit      int64   `json:"memory.heap.init"`
			MemoryHeapMax       int64   `json:"memory.heap.max"`
			MemoryHeapUsage     float64 `json:"memory.heap.usage"`
			MemoryHeapUsed      int64   `json:"memory.heap.used"`

			MemoryNonHeapCommitted int64   `json:"memory.non-heap.committed"`
			MemoryNonHeapInit      int64   `json:"memory.non-heap.init"`
			MemoryNonHeapMax       int64   `json:"memory.non-heap.max"`
			MemoryNonHeapUsage     float64 `json:"memory.non-heap.usage"`
			MemoryNonHeapUsed      int64   `json:"memory.non-heap.used"`

			MemoryTotalCommitted int64 `json:"memory.total.committed"`
			MemoryTotalInit      int64 `json:"memory.total.init"`
			MemoryTotalMax       int64 `json:"memory.total.max"`
			MemoryTotalUsed      int64 `json:"memory.total.used"`

			OSAvailableProcessors        int64   `json:"os.availableProcessors"`
			OSCommittedVirtualMemorySize int64   `json:"os.committedVirtualMemorySize"`
			OSFreePhysicalMemorySize     int64   `json:"os.freePhysicalMemorySize"`
			OSFreeSwapSpaceSize          int64   `json:"os.freeSwapSpaceSize"`
			OSMaxFileDescriptorCount     int64   `json:"os.maxFileDescriptorCount"`
			OSOpenFileDescriptorCount    int64   `json:"os.openFileDescriptorCount"`
			OSProcessCPUTime             int64   `json:"os.processCpuTime"`
			OSSystemLoadAverage          float64 `json:"os.systemLoadAverage"`
			OSTotalPhysicalMemorySize    int64   `json:"os.totalPhysicalMemorySize"`
			OSTotalSwapSapceSize         int64   `json:"os.totalSwapSpaceSize"`

			ThreadsBlockedCount      int64 `json:"threads.blocked.count"`
			ThreadsDaemonCount       int64 `json:"threads.daemon.count"`
			ThreadsDeadlockCount     int64 `json:"threads.deadlock.count"`
			ThreadsNewCount          int64 `json:"threads.new.count"`
			ThreadsRunnableCount     int64 `json:"threads.runnable.count"`
			ThreadsTerminatedCount   int64 `json:"threads.terminated.count"`
			ThreadsTimedWaitingCount int64 `json:"threads.timed_waiting.count"`
			ThreadsWaitingCount      int64 `json:"threads.waiting.count"`
		} `json:"solr.jvm"`
	} `json:"metrics"`
}

type JVMStatusV6 struct {
	Metrics struct {
		JVM struct {
			GCConcurrentMarkSweepCount struct {
				Value int64 `json:"value"`
			} `json:"gc.ConcurrentMarkSweep.count"`
			GCConcurrentMarkSweepTime struct {
				Value int64 `json:"value"`
			} `json:"gc.ConcurrentMarkSweep.time"`
			GCParNewCount struct {
				Value int64 `json:"value"`
			} `json:"gc.ParNew.count"`
			GCParNewTime struct {
				Value int64 `json:"value"`
			} `json:"gc.ParNew.time"`

			MemoryHeapCommitted struct {
				Value int64 `json:"value"`
			} `json:"memory.heap.committed"`
			MemoryHeapInit struct {
				Value int64 `json:"value"`
			} `json:"memory.heap.init"`
			MemoryHeapMax struct {
				Value int64 `json:"value"`
			} `json:"memory.heap.max"`
			MemoryHeapUsage struct {
				Value float64 `json:"value"`
			} `json:"memory.heap.usage"`
			MemoryHeapUsed struct {
				Value int64 `json:"value"`
			} `json:"memory.heap.used"`

			MemoryNonHeapCommitted struct {
				Value int64 `json:"value"`
			} `json:"memory.non-heap.committed"`
			MemoryNonHeapInit struct {
				Value int64 `json:"value"`
			} `json:"memory.non-heap.init"`
			MemoryNonHeapMax struct {
				Value int64 `json:"value"`
			} `json:"memory.non-heap.max"`
			MemoryNonHeapUsage struct {
				Value float64 `json:"value"`
			} `json:"memory.non-heap.usage"`
			MemoryNonHeapUsed struct {
				Value int64 `json:"value"`
			} `json:"memory.non-heap.used"`

			MemoryTotalCommitted struct {
				Value int64 `json:"value"`
			} `json:"memory.total.committed"`
			MemoryTotalInit struct {
				Value int64 `json:"value"`
			} `json:"memory.total.init"`
			MemoryTotalMax struct {
				Value int64 `json:"value"`
			} `json:"memory.total.max"`
			MemoryTotalUsed struct {
				Value int64 `json:"value"`
			} `json:"memory.total.used"`

			OSAvailableProcessors struct {
				Value int64 `json:"value"`
			} `json:"os.availableProcessors"`
			OSCommittedVirtualMemorySize struct {
				Value int64 `json:"value"`
			} `json:"os.committedVirtualMemorySize"`
			OSFreePhysicalMemorySize struct {
				Value int64 `json:"value"`
			} `json:"os.freePhysicalMemorySize"`
			OSFreeSwapSpaceSize struct {
				Value int64 `json:"value"`
			} `json:"os.freeSwapSpaceSize"`
			OSMaxFileDescriptorCount struct {
				Value int64 `json:"value"`
			} `json:"os.maxFileDescriptorCount"`
			OSOpenFileDescriptorCount struct {
				Value int64 `json:"value"`
			} `json:"os.openFileDescriptorCount"`
			OSProcessCPUTime struct {
				Value int64 `json:"value"`
			} `json:"os.processCpuTime"`
			OSSystemLoadAverage struct {
				Value float64 `json:"value"`
			} `json:"os.systemLoadAverage"`
			OSTotalPhysicalMemorySize struct {
				Value int64 `json:"value"`
			} `json:"os.totalPhysicalMemorySize"`
			OSTotalSwapSapceSize struct {
				Value int64 `json:"value"`
			} `json:"os.totalSwapSpaceSize"`

			ThreadsBlockedCount struct {
				Value int64 `json:"value"`
			} `json:"threads.blocked.count"`
			ThreadsDaemonCount struct {
				Value int64 `json:"value"`
			} `json:"threads.daemon.count"`
			ThreadsDeadlockCount struct {
				Value int64 `json:"value"`
			} `json:"threads.deadlock.count"`
			ThreadsNewCount struct {
				Value int64 `json:"value"`
			} `json:"threads.new.count"`
			ThreadsRunnableCount struct {
				Value int64 `json:"value"`
			} `json:"threads.runnable.count"`
			ThreadsTerminatedCount struct {
				Value int64 `json:"value"`
			} `json:"threads.terminated.count"`
			ThreadsTimedWaitingCount struct {
				Value int64 `json:"value"`
			} `json:"threads.timed_waiting.count"`
			ThreadsWaitingCount struct {
				Value int64 `json:"value"`
			} `json:"threads.waiting.count"`
		} `json:"solr.jvm"`
	} `json:"metrics"`
}

type InfoSystem struct {
	Lucene struct {
		SolrVersion string `json:"solr-spec-version"`
	} `json:"lucene"`
}
