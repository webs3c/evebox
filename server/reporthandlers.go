/* Copyright (c) 2014-2015 Jason Ish
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED ``AS IS'' AND ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package server

import (
	"github.com/jasonish/evebox/core"
	"net/http"
	"strconv"
)

func ReportDnsRequestRrnames(appContext AppContext, r *http.Request) interface{} {

	options := core.ReportOptions{}

	if r.Method == http.MethodPost {
		var requestBody struct {
			TimeRange   string `json:"timeRange"`
			Size        int64  `json:"size"`
			QueryString string `json:"queryString"`
		}
		DecodeRequestBody(r, &requestBody)
		options.TimeRange = requestBody.TimeRange
		options.Size = requestBody.Size
		options.QueryString = requestBody.QueryString
	} else {
		options.TimeRange = r.FormValue("timeRange")
		options.Size, _ = strconv.ParseInt(r.FormValue("size"), 10, 64)
		options.QueryString = r.FormValue("queryString")
	}

	data, err := appContext.ReportService.ReportDnsRequestRrnames(options)
	if err != nil {
		return err
	}

	return map[string]interface{}{
		"data": data,
	}

}

func ReportAlertAggs(appContext AppContext, r *http.Request) interface{} {
	options := core.ReportOptions{}

	agg := r.FormValue("agg")
	options.TimeRange = r.FormValue("timeRange")
	options.Size, _ = strconv.ParseInt(r.FormValue("size"), 10, 64)
	options.AddressFilter = r.FormValue("addressFilter")
	options.QueryString = r.FormValue("queryString")

	response, err := appContext.ReportService.ReportAlertAggs(agg, options)
	if err != nil {
		return err
	}
	return response
}

func ReportHistogram(appContext AppContext, r *http.Request) interface{} {
	options := core.ReportOptions{}

	options.TimeRange = r.FormValue("timeRange")
	options.AddressFilter = r.FormValue("addressFilter")
	options.QueryString = r.FormValue("queryString")
	options.SensorFilter = r.FormValue("sensorFilter")
	options.EventType = r.FormValue("eventType")
	options.DnsType = r.FormValue("dnsType")

	interval := r.FormValue("interval")

	response, err := appContext.ReportService.ReportHistogram(interval, options)
	if err != nil {
		return err
	}
	return response
}