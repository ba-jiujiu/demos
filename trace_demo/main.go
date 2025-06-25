package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/contrib/otelfiber"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation" // 引入新包
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

func newTracerProvider(serviceName, jaegerURL string) (*sdktrace.TracerProvider, error) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerURL)))
	if err != nil {
		return nil, err
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	return tp, nil
}

// startServiceA 启动监听 8080 端口的服务
func startServiceA(wg *sync.WaitGroup, tp trace.TracerProvider, prop propagation.TextMapPropagator) {
	defer wg.Done()

	tracer := tp.Tracer("service-a-tracer")
	app := fiber.New()
	app.Use(logger.New())
	app.Use(otelfiber.Middleware(
		otelfiber.WithTracerProvider(tp),
		otelfiber.WithPropagators(prop), // 使用共享的Propagator
	))

	app.Get("/call", func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		_ctx, childSpan := tracer.Start(ctx, "call-service-b-api")
		defer childSpan.End()

		client := http.Client{
			Transport: otelhttp.NewTransport(
				http.DefaultTransport,
				otelhttp.WithTracerProvider(tp),
				otelhttp.WithPropagators(prop), // 使用共享的Propagator
			),
		}

		req, _ := http.NewRequestWithContext(_ctx, "GET", "http://localhost:8081/api", nil)

		res, err := client.Do(req)
		if err != nil {
			childSpan.RecordError(err)
			childSpan.SetAttributes(attribute.String("error", err.Error()))
			return c.Status(500).SendString("Failed to call Service B")
		}
		defer res.Body.Close()

		childSpan.SetAttributes(attribute.Int("http.response.status_code", res.StatusCode))
		return c.SendString("Called Service B, trace sent!")
	})

	log.Println("Service A [fiber-service-a] running on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Printf("Service A failed to start: %v", err)
	}
}

// startServiceB 启动监听 8081 端口的服务
func startServiceB(wg *sync.WaitGroup, tp trace.TracerProvider, prop propagation.TextMapPropagator) {
	defer wg.Done()

	tracer := tp.Tracer("service-b-tracer")
	app := fiber.New()
	app.Use(logger.New())
	app.Use(otelfiber.Middleware(
		otelfiber.WithTracerProvider(tp),
		otelfiber.WithPropagators(prop), // 使用共享的Propagator
	))

	app.Get("/api", func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		_ctx, childSpan := tracer.Start(ctx, "call-service-c-api")
		defer childSpan.End()

		client := http.Client{
			Transport: otelhttp.NewTransport(
				http.DefaultTransport,
				otelhttp.WithTracerProvider(tp),
				otelhttp.WithPropagators(prop), // 使用共享的Propagator
			),
		}

		req, _ := http.NewRequestWithContext(_ctx, "GET", "http://localhost:8082/srv", nil)

		res, err := client.Do(req)
		if err != nil {
			childSpan.RecordError(err)
			childSpan.SetAttributes(attribute.String("error", err.Error()))
			return c.Status(500).SendString("Failed to call Service C")
		}
		defer res.Body.Close()

		childSpan.SetAttributes(attribute.Int("http.response.status_code", res.StatusCode))
		return c.SendString("Called Service C, trace sent!")
	})
	log.Println("Service B [fiber-service-b] running on :8081")
	if err := app.Listen(":8081"); err != nil {
		log.Printf("Service B failed to start: %v", err)
	}
}

// startServiceC 启动监听 8082 端口的服务
func startServiceC(wg *sync.WaitGroup, tp trace.TracerProvider, prop propagation.TextMapPropagator) {
	defer wg.Done()

	tracer := tp.Tracer("service-c-tracer")
	app := fiber.New()
	app.Use(logger.New())
	app.Use(otelfiber.Middleware(
		otelfiber.WithTracerProvider(tp),
		otelfiber.WithPropagators(prop), // 使用共享的Propagator
	))

	app.Get("/srv", func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		_, dbSpan := tracer.Start(ctx, "database-query")
		time.Sleep(25 * time.Millisecond)
		dbSpan.End()

		return c.SendString("Hello from Service c!")
	})
	log.Println("Service c [fiber-service-c] running on :8082")
	if err := app.Listen(":8082"); err != nil {
		log.Printf("Service c failed to start: %v", err)
	}
}

func main() {
	jaegerURL := "http://localhost:14268/api/traces"

	tpA, err := newTracerProvider("fiber-service-a", jaegerURL)
	if err != nil {
		log.Fatal(err)
	}
	tpB, err := newTracerProvider("fiber-service-b", jaegerURL)
	if err != nil {
		log.Fatal(err)
	}
	tpC, err := newTracerProvider("fiber-service-c", jaegerURL)
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个所有服务共享的Propagator
	propagator1 := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	propagator2 := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	propagator3 := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

	defer func() {
		log.Println("Shutting down tracer providers...")
		if err := tpA.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider A: %v", err)
		}
		if err := tpB.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider B: %v", err)
		}
		if err := tpC.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider C: %v", err)
		}
		log.Println("Tracer providers shut down.")
	}()

	var wg sync.WaitGroup
	wg.Add(3)

	// 将propagator实例传给每个服务
	go startServiceA(&wg, tpA, propagator1)
	go startServiceB(&wg, tpB, propagator2)
	go startServiceC(&wg, tpC, propagator3)

	log.Println("Servers started. Press Ctrl+C to exit.")
	log.Println("Send request to http://localhost:8080/call to generate a trace.")

	wg.Wait()
	log.Println("All services have stopped.")
}
