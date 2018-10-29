package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"reflect"
	"testing"
	"time"
)

type solrVersionJSON struct {
	version string
	json    io.Reader
}

var solrResponseDir = "utils/solr-responses"

func loadAllMbeans() ([]solrVersionJSON, error) {
	mbeans := []solrVersionJSON{}
	solrVersions, err := ioutil.ReadDir(solrResponseDir)
	if err != nil {
		return mbeans, fmt.Errorf("failed to list solr versions")
	}
	for _, sv := range solrVersions {
		mbeansFile := path.Join(solrResponseDir, sv.Name(), "mbeans.json")
		fi, err := os.Stat(mbeansFile)
		if err != nil {
			continue
		}
		file, err := os.Open(mbeansFile)
		if err != nil {
			return nil, fmt.Errorf("Failed to open file %s", fi.Name())
		}
		mbeans = append(mbeans, solrVersionJSON{sv.Name(), file})

	}
	return mbeans, nil
}

func Test_processMbeans(t *testing.T) {
	type args struct {
		e        *Exporter
		coreName string
		data     io.Reader
	}
	type test struct {
		name string
		args args
		want []error
	}
	hc := http.Client{}
	exporter := NewExporter("", time.Second, "", hc)
	tests := []test{}
	mbeans, err := loadAllMbeans()
	if err != nil {
		t.Errorf("enumerating solr versions: %v", err)
	}

	for _, mbean := range mbeans {
		tests = append(tests, test{name: mbean.version, args: args{exporter, "gettingstarted", mbean.json}, want: []error{}})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processMbeans(tt.args.e, tt.args.coreName, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				for _, err := range got {
					t.Errorf("processMbeans() returned error: %v", err)
				}
			}
		})
	}
}
