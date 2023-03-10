package server

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) RunMetrics() error {
	errCh := make(chan error)

	metricsServer := echo.New()
	metricsServer.HidePort = true
	metricsServer.HideBanner = true

	metricsServer.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	go func() {
		s.logger.Infof(
			"Metrics available URL: %s, ServiceName: %s",
			s.cfg.Metrics.URL,
			s.cfg.Metrics.ServiceName,
		)
		errCh <- metricsServer.Start(s.cfg.Metrics.URL)
	}()

	select {
	case <-s.ctx.Done():
		c := make(chan bool)
		go func() {
			defer close(c)
			errCh <- metricsServer.Shutdown(s.ctx)
		}()
		select {
		case <-c:
		case <-time.After(5 * time.Second):
		}
		s.logger.Info("Metrics Exited Properly")
		return nil
	case err := <-errCh:
		return err
	}
}
