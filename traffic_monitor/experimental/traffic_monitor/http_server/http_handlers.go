package http_server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Type int

// TODO rename these, all caps isn't recommended Go style
const (
	TR_CONFIG Type = (1 << iota)
	TR_STATE_DERIVED
	TR_STATE_SELF
	CACHE_STATS
	DS_STATS
	EVENT_LOG
	PEER_STATES
	STAT_SUMMARY
	STATS
	CONFIG_DOC
	API_CACHE_COUNT
	API_CACHE_AVAILABLE_COUNT
	API_CACHE_DOWN_COUNT
	API_VERSION
	API_TRAFFIC_OPS_URI
	API_CACHE_STATES
)

type Format int

const (
	XML Format = (1 << iota)
	JSON
)

type DataRequest struct {
	T          Type
	F          Format
	C          chan []byte
	Date       string
	Parameters map[string][]string
}

func dataRequest(w http.ResponseWriter, req *http.Request, t Type, f Format) {
	//pp: "0=[my-ats-edge-cache-0], hc=[1]",
	//dateLayout := "Thu Oct 09 20:28:36 UTC 2014"
	dateLayout := "Mon Jan 02 15:04:05 MST 2006"
	time := time.Now()
	p := make(map[string][]string)

	for key, v := range req.URL.Query() {
		for _, value := range v {
			p[key] = append(p[key], value)
		}
	}

	dr := DataRequest{
		T:          t,
		F:          f,
		C:          make(chan []byte, 1), // must be buffered, so if this is killed, the writer doesn't block forever
		Date:       time.UTC().Format(dateLayout),
		Parameters: p,
	}

	mgrReqChan <- dr
	writeResponse(w, f, dr)
}

func handleCrStates(w http.ResponseWriter, req *http.Request) {
	t := TR_STATE_DERIVED

	if req.URL.RawQuery == "raw" {
		t = TR_STATE_SELF
	}

	dataRequest(w, req, t, JSON)
}

func handleCrConfig(w http.ResponseWriter, req *http.Request) {
	dataRequest(w, req, TR_CONFIG, JSON)
}

func handleCacheStats(w http.ResponseWriter, req *http.Request) {
	dataRequest(w, req, CACHE_STATS, JSON)
}

func handleDsStats(w http.ResponseWriter, req *http.Request) {
	dataRequest(w, req, DS_STATS, JSON)
}

func handleEventLog(w http.ResponseWriter, req *http.Request) {
	dataRequest(w, req, EVENT_LOG, JSON)
}

func handlePeerStates(w http.ResponseWriter, req *http.Request) {
	dataRequest(w, req, PEER_STATES, JSON)
}

func handleStatSummary(w http.ResponseWriter, req *http.Request) {
	dataRequest(w, req, STAT_SUMMARY, JSON)
}

func handleStats(w http.ResponseWriter, req *http.Request) {
	dataRequest(w, req, STATS, JSON)
}

func handleConfigDoc(w http.ResponseWriter, req *http.Request) {
	dataRequest(w, req, CONFIG_DOC, JSON)
}

func handleRootFunc() (http.HandlerFunc, error) {
	index, err := ioutil.ReadFile("index.html")
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%s", index)
	}, nil
}

func handleApiCacheCount(w http.ResponseWriter, req *http.Request) {
	dataRequest(w, req, API_CACHE_COUNT, JSON)
}

func handleApiCacheAvailableCount(w http.ResponseWriter, req *http.Request) {
	dataRequest(w, req, API_CACHE_AVAILABLE_COUNT, JSON)
}

func handleApiCacheDownCount(w http.ResponseWriter, req *http.Request) {
	dataRequest(w, req, API_CACHE_DOWN_COUNT, JSON)
}

func handleApiVersion(w http.ResponseWriter, req *http.Request) {
	dataRequest(w, req, API_VERSION, JSON)
}

func handleApiTrafficOpsURI(w http.ResponseWriter, req *http.Request) {
	dataRequest(w, req, API_TRAFFIC_OPS_URI, JSON)
}

func handleApiCacheStates(w http.ResponseWriter, req *http.Request) {
	dataRequest(w, req, API_CACHE_STATES, JSON)
}
