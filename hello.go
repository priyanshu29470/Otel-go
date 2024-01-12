package main

import (
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	tracer         = otel.Tracer("hello")
	meter          = otel.Meter("hello")
	requestCounter metric.Int64Counter
	errorCounter   metric.Int64Counter
)

func init() {
	var err error
	requestCounter, err = meter.Int64Counter("http.requests",
		metric.WithDescription("The number of http requests made"),
		metric.WithUnit("{request}"))
	errorCounter, err = meter.Int64Counter("http.errors",
		metric.WithDescription("The number of errors"),
		metric.WithUnit("{error}"))
	if err != nil {
		panic(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "success-span")
	defer span.End()

	requestStatusAttr := attribute.Int("request.status", 200)
	requestTypeAttr := attribute.String("request.type", "success")
	span.SetAttributes(requestStatusAttr)
	requestCounter.Add(ctx, 1, metric.WithAttributes(requestStatusAttr, requestTypeAttr))

	resp := "http status: 200"
	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}

func getError(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "error-span")
	defer span.End()

	requestStatusAttr := attribute.Int("request.status", 404)
	requestTypeAttr := attribute.String("request.type", "error")
	span.SetAttributes(requestStatusAttr)
	errorCounter.Add(ctx, 1, metric.WithAttributes(requestStatusAttr, requestTypeAttr))

	resp := "http status: 404"
	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}
