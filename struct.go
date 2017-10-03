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
		One5minRateReqsPerSecond *float64 `json:"15minRateReqsPerSecond,omitempty"`
		FiveMinRateReqsPerSecond *float64 `json:"5MinRateReqsPerSecond,omitempty"`
		One5minRateRequestsPerSecond float64 `json:"15minRateRequestsPerSecond,omitempty"`
		FiveminRateRequestsPerSecond float64 `json:"5minRateRequestsPerSecond,omitempty"`
		Seven5thPcRequestTime    float64 `json:"75thPcRequestTime"`
		Nine5thPcRequestTime     float64 `json:"95thPcRequestTime"`
		Nine9thPcRequestTime     float64 `json:"99thPcRequestTime"`
		Nine99thPcRequestTime    float64 `json:"999thPcRequestTime"`
		AvgRequestsPerSecond     float64 `json:"avgRequestsPerSecond"`
		AvgTimePerRequest        float64 `json:"avgTimePerRequest"`
		Errors                   int     `json:"errors"`
		HandlerStart             int     `json:"handlerStart"`
		MedianRequestTime        float64 `json:"medianRequestTime"`
		Requests                 int     `json:"requests"`
		Timeouts                 int     `json:"timeouts"`
		TotalTime                float64 `json:"totalTime"`
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
		CumulativeEvictions int     `json:"cumulative_evictions"`
		CumulativeHitratio  float64 `json:"cumulative_hitratio"`
		CumulativeHits      int     `json:"cumulative_hits"`
		CumulativeInserts   int     `json:"cumulative_inserts"`
		CumulativeLookups   int     `json:"cumulative_lookups"`
		Evictions           int     `json:"evictions"`
		Hitratio            float64 `json:"hitratio"`
		Hits                int     `json:"hits"`
		Inserts             int     `json:"inserts"`
		Lookups             int     `json:"lookups"`
		Size                int     `json:"size"`
		WarmupTime          int     `json:"warmupTime"`
	} `json:"stats"`
}

type CacheSolrV4 struct {
	Class string `json:"class"`
	Stats struct {
		CumulativeEvictions int    `json:"cumulative_evictions"`
		CumulativeHitratio  string `json:"cumulative_hitratio"`
		CumulativeHits      int    `json:"cumulative_hits"`
		CumulativeInserts   int    `json:"cumulative_inserts"`
		CumulativeLookups   int    `json:"cumulative_lookups"`
		Evictions           int    `json:"evictions"`
		Hitratio            string `json:"hitratio"`
		Hits                int    `json:"hits"`
		Inserts             int    `json:"inserts"`
		Lookups             int    `json:"lookups"`
		Size                int    `json:"size"`
		WarmupTime          int    `json:"warmupTime"`
	} `json:"stats"`
}
