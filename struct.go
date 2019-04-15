package main

// import "encoding/json"

// // import "time"

type AdminCoresName struct {
	Status map[string]struct {
		Name string `json:"name"`
	} `json:"status"`
}

// type MetricsJetty struct {
// 	Metrics struct {
// 		SolrJetty map[string]json.RawMessage `json:"solr.jetty"`
// 	} `json:"metrics"`
// }

// type Metrics struct {
// 	Metrics struct {
// 		SolrJetty map[string]json.RawMessage `json:"solr.jetty"`
// 	} `json:"metrics"`
// }

// type SimpleCount struct {
// 	Count int `json:"count"`
// }
