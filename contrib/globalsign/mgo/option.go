// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package mgo

import (
	"context"
	"math"

	"gopkg.in/DataDog/dd-trace-go.v1/internal"
	"gopkg.in/DataDog/dd-trace-go.v1/internal/namingschema"
)

const defaultServiceName = "mongodb"

type mongoConfig struct {
	ctx           context.Context
	serviceName   string
	spanName      string
	analyticsRate float64
}

func newConfig() *mongoConfig {
	rate := math.NaN()
	if internal.BoolEnv("DD_TRACE_GIN_ANALYTICS_ENABLED", false) {
		rate = 1.0
	}
	return &mongoConfig{
		serviceName: namingschema.NewDefaultServiceName(
			defaultServiceName,
			namingschema.WithOverrideV0(defaultServiceName),
		).GetName(),
		spanName: namingschema.NewMongoDBOutboundOp().GetName(),
		ctx:      context.Background(),
		// analyticsRate: globalconfig.AnalyticsRate(),
		analyticsRate: rate,
	}
}

// DialOption represents an option that can be passed to Dial
type DialOption func(*mongoConfig)

// WithServiceName sets the service name for a given MongoDB context.
func WithServiceName(name string) DialOption {
	return func(cfg *mongoConfig) {
		cfg.serviceName = name
	}
}

// WithContext sets the context.
func WithContext(ctx context.Context) DialOption {
	return func(cfg *mongoConfig) {
		cfg.ctx = ctx
	}
}

// WithAnalytics enables Trace Analytics for all started spans.
func WithAnalytics(on bool) DialOption {
	return func(cfg *mongoConfig) {
		if on {
			cfg.analyticsRate = 1.0
		} else {
			cfg.analyticsRate = math.NaN()
		}
	}
}

// WithAnalyticsRate sets the sampling rate for Trace Analytics events
// correlated to started spans.
func WithAnalyticsRate(rate float64) DialOption {
	return func(cfg *mongoConfig) {
		if rate >= 0.0 && rate <= 1.0 {
			cfg.analyticsRate = rate
		} else {
			cfg.analyticsRate = math.NaN()
		}
	}
}
