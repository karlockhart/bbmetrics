package bbmetrics

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"
)

// DownloadMeter represents a meter that downloads a file.
type DownloadMeter struct {
}

type DownloadMeasurement struct {
	StartConnect      int64
	TimeToConnect     int64
	TimeToFirstByte   int64
	dnsStartTime      int64
	TimeToDNSStart    int64
	DNSDuration       int64
	tlsStartTime      int64
	TimeToTLSStart    int64
	TLSDuration       int64
	TimeToSendHeaders int64
	TimeToSendRequest int64
	ConnError         int
	TLSError          int
	RequestTries      int
	TotalTime         int64
}

func (dm *DownloadMeasurement) GotFirstResponseByte() {
	dm.TimeToFirstByte = time.Now().UnixNano() - dm.StartConnect
}

func (dm *DownloadMeasurement) GetConn(a string) {
	dm.StartConnect = time.Now().UnixNano()
}

func (dm *DownloadMeasurement) DNSStart(d httptrace.DNSStartInfo) {
	dm.dnsStartTime = time.Now().UnixNano()
	dm.TimeToDNSStart = dm.dnsStartTime - dm.StartConnect
}

func (dm *DownloadMeasurement) DNSDone(d httptrace.DNSDoneInfo) {
	dm.DNSDuration = time.Now().UnixNano() - dm.dnsStartTime
}

func (dm *DownloadMeasurement) TLSHandshakeStart() {
	dm.tlsStartTime = time.Now().UnixNano()
	dm.TimeToTLSStart = dm.tlsStartTime - dm.StartConnect
}

func (dm *DownloadMeasurement) TLSHandshakeDone(t tls.ConnectionState, err error) {
	if err != nil {
		dm.TLSError++
	}
	dm.TLSDuration = time.Now().UnixNano() - dm.StartConnect

}

func (dm *DownloadMeasurement) WroteRequest(httptrace.WroteRequestInfo) {
	dm.RequestTries++
	dm.TimeToSendRequest = time.Now().UnixNano() - dm.StartConnect
}

func (dm *DownloadMeasurement) ConnectDone(network, a string, err error) {
	if err != nil {
		dm.ConnError++
		return
	}
	dm.TimeToConnect = time.Now().UnixNano() - dm.StartConnect
}

func (dm *DownloadMeasurement) WroteHeaders() {
	dm.TimeToSendHeaders = time.Now().UnixNano() - dm.StartConnect
}

// NewDownloadMeter returns a download meter.
func NewDownloadMeter() (dlm *DownloadMeter) {

	return
}

// Measure gets a sample of the download.
func (dlm *DownloadMeter) Measure() {
	url := "http://generictld.xyz/50MB.zip"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	var m DownloadMeasurement
	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: m.GotFirstResponseByte,
		GetConn:              m.GetConn,
		ConnectDone:          m.ConnectDone,
		DNSDone:              m.DNSDone,
		DNSStart:             m.DNSStart,
		TLSHandshakeStart:    m.TLSHandshakeStart,
		TLSHandshakeDone:     m.TLSHandshakeDone,
		WroteRequest:         m.WroteRequest,
		WroteHeaders:         m.WroteHeaders,
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	res, err := http.DefaultTransport.RoundTrip(req)
	m.TotalTime = time.Now().UnixNano() - m.StartConnect
	if err != nil {
		log.Fatal(err)
	}
	log.Println(m)
	log.Println(res.ContentLength)

}
