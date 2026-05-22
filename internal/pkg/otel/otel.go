package otel

import (
	"bookstore/internal/pkg/config"
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func InitOtel(ctx context.Context, cfg *config.Otel, srvName string, srvVersion string, env string) (func(context.Context) error, error) {
	// 创建资源
	res, err := initResource(ctx, srvName, srvVersion, env)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// 创建TracerProvider
	tracerProvider, err := initTracerProvider(ctx, cfg, res)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer provider: %w", err)
	}

	// 创建MeterProvider
	meterProvider, err := initMeterProvider(ctx, cfg, res)
	if err != nil {
		if tracerErr := tracerProvider.Shutdown(ctx); tracerErr != nil {
			return nil, errors.Join(
				fmt.Errorf("failed to create meter provider: %w", err),
				fmt.Errorf("failed to shutdown tracer provider: %w", tracerErr),
			)
		}
		return nil, fmt.Errorf("failed to create meter provider: %w", err)
	}
	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)

	return func(ctx context.Context) error {
		errs := make([]error, 0)
		tracerErr := tracerProvider.Shutdown(ctx)
		if tracerErr != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown tracer provider: %w", tracerErr))
		}
		meterErr := meterProvider.Shutdown(ctx)
		if meterErr != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown meter provider: %w", meterErr))
		}
		if len(errs) > 0 {
			return errors.Join(errs...)
		}
		return nil
	}, nil
}

// 创建资源
func initResource(ctx context.Context, srvName string, srvVersion string, env string) (*resource.Resource, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(srvName),
			semconv.ServiceVersionKey.String(srvVersion),
			attribute.String("deployment.environment", env),
		),
		resource.WithProcess(), // 添加进程信息（PID、可执行文件路径等）
		resource.WithOS(),      // 添加操作系统信息
		resource.WithHost(),    // 添加主机名
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	return resource.Merge(resource.Default(), res)
}

// 创建TracerProvider
func initTracerProvider(ctx context.Context, cfg *config.Otel, res *resource.Resource) (*trace.TracerProvider, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.OtelEndpoint),
		otlptracegrpc.WithTimeout(cfg.ExportTimeout),
	}
	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	provider := trace.NewTracerProvider(
		// 设置导出器
		trace.WithBatcher(exporter,
			trace.WithBatchTimeout(cfg.BatchTimeout),
		),
		// 设置资源
		trace.WithResource(res),
		// 设置采样器
		trace.WithSampler(trace.TraceIDRatioBased(cfg.TraceSampleRate)),
	)
	return provider, nil
}

// 创建MeterProvider
func initMeterProvider(ctx context.Context, cfg *config.Otel, res *resource.Resource) (*metric.MeterProvider, error) {
	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.OtelEndpoint),
		otlpmetricgrpc.WithTimeout(cfg.ExportTimeout),
	}
	if cfg.Insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP metric exporter: %w", err)
	}

	provider := metric.NewMeterProvider(
		// 设置导出器
		metric.WithReader(metric.NewPeriodicReader(exporter,
			metric.WithInterval(cfg.ExportInterval),
		)),
		// 设置资源
		metric.WithResource(res),
	)
	return provider, nil
}
