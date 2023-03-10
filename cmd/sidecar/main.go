/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	goflag "flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/spf13/pflag"
	"k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"

	"k8s.io/dns/pkg/sidecar"
	"k8s.io/dns/pkg/version"
	"k8s.io/klog/v2"
)

const (
	defaultProbeInterval = 5 * time.Second
)

func main() {
	options := sidecar.NewOptions()
	configureFlags(options, pflag.CommandLine)

	pflag.CommandLine.SetNormalizeFunc(flag.WordSepNormalizeFunc)
	logs.AddFlags(pflag.CommandLine)

	flag.InitFlags()
	// Convinces goflags that we have called Parse() to avoid noisy logs.
	// OSS Issue: kubernetes/kubernetes#17162.
	goflag.CommandLine.Parse([]string{})
	logs.InitLogs()
	defer logs.FlushLogs()

	klog.Infof("Version v%s", version.VERSION)

	version.PrintAndExitIfRequested()

	server := sidecar.NewServer()
	server.Run(options)
}

type probeOptions []sidecar.DNSProbeOption

func (po *probeOptions) String() string {
	return fmt.Sprintf("%+v", *po)
}

func (po *probeOptions) Set(value string) error {
	splits := strings.Split(value, ",")
	if !(3 <= len(splits) && len(splits) <= 5) {
		return fmt.Errorf("invalid format to --probe")
	}

	option := sidecar.DNSProbeOption{
		Label:    splits[0],
		Server:   splits[1],
		Name:     splits[2],
		Interval: defaultProbeInterval,
		Type:     dns.TypeANY,
	}

	const labelRegexp = "^[a-zA-Z0-9_]+$"
	if !regexp.MustCompile(labelRegexp).MatchString(option.Label) {
		return fmt.Errorf("label must be of format %v", labelRegexp)
	}

	if !strings.Contains(option.Server, ":") {
		option.Server = option.Server + ":53"
	}

	if !strings.HasSuffix(option.Name, ".") {
		// dns package requires a fully qualified (e.g. terminal '.') name
		option.Name = option.Name + "."
	}

	if len(splits) >= 4 {
		if interval, err := strconv.Atoi(splits[3]); err == nil {
			option.Interval = time.Duration(interval) * time.Second
		} else {
			return err
		}
	}

	if len(splits) >= 5 {
		switch splits[4] {
		case "A":
			option.Type = dns.TypeA
		case "AAAA":
			option.Type = dns.TypeAAAA
		case "ANY":
			option.Type = dns.TypeANY
		case "SRV":
			option.Type = dns.TypeSRV
		default:
			return fmt.Errorf("invalid type for DNS: %v", splits[4])
		}
	}

	*po = append(*po, option)

	return nil
}

func (po *probeOptions) Type() string {
	return "string"
}

var _ pflag.Value = (*probeOptions)(nil)

func configureFlags(opt *sidecar.Options, flagSet *pflag.FlagSet) {
	flagSet.StringVar(
		&opt.DnsMasqAddr, "dnsmasq-addr", opt.DnsMasqAddr,
		"address that the dnsmasq server is listening on")
	flagSet.IntVar(
		&opt.DnsMasqPort, "dnsmasq-port", opt.DnsMasqPort,
		"port that the dnsmasq server is listening on")
	flagSet.IntVar(
		&opt.DnsMasqPollIntervalMs, "dnsmasq-poll-interval-ms", opt.DnsMasqPollIntervalMs,
		"interval with which to poll dnsmasq for stats")
	flagSet.Var(
		(*probeOptions)(&opt.Probes), "probe",
		"probe the given DNS server with the DNS name and export probe"+
			" metrics and healthcheck URI. Specified as"+
			" <label>,<server>,<dns name>[,<interval_seconds>][,<type>]."+
			" Healthcheck url will be exported under /healthcheck/<label>."+
			" interval_seconds is optional."+
			" This option may be specified multiple times to check multiple servers."+
			" <type> is one of ANY, A, AAAA, SRV."+
			" Example: 'mydns,127.0.0.1:53,example.com,10,A'.")
	flagSet.StringVar(
		&opt.PrometheusAddr, "prometheus-addr", opt.PrometheusAddr,
		"http address to bind metrics server to")
	flagSet.IntVar(
		&opt.PrometheusPort, "prometheus-port", opt.PrometheusPort,
		"http port to use to export prometheus metrics")
	flagSet.StringVar(
		&opt.PrometheusPath, "prometheus-path", opt.PrometheusPath,
		"http path used to export metrics")
	flagSet.StringVar(
		&opt.PrometheusNamespace, "prometheus-namespace", opt.PrometheusNamespace,
		"prometheus metric namespace")
}
