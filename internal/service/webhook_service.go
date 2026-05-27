package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	mockscontext "github.com/nicopozo/mockserver/internal/context"
	"github.com/nicopozo/mockserver/internal/model"
)

const defaultWebhookTimeout = 5 * time.Second

//go:generate mockgen -destination=../utils/test/mocks/mock_service_webhook.go -package=mocks -source=./webhook_service.go

// WebhookService defines the interface for firing webhooks.
type WebhookService interface {
	Fire(
		ctx context.Context,
		webhook model.WebhookConfig,
		variables []*model.Variable,
		onResult func(model.WebhookResult),
	)
}

type webhookService struct{}

// NewWebhookService creates a new WebhookService instance.
func NewWebhookService() WebhookService {
	return &webhookService{}
}

func (s *webhookService) Fire(
	ctx context.Context,
	webhook model.WebhookConfig,
	variables []*model.Variable,
	onResult func(model.WebhookResult),
) {
	//nolint:gosec
	go s.fireAsync(ctx, webhook, variables, onResult)
}

func (s *webhookService) fireAsync(
	ctx context.Context,
	webhook model.WebhookConfig,
	variables []*model.Variable,
	onResult func(model.WebhookResult),
) {
	logger := mockscontext.Logger(ctx)
	start := time.Now()

	// Apply delay before sending the webhook, if configured.
	if webhook.Delay > 0 {
		logger.Debug(s, map[string]string{"delay_ms": fmt.Sprintf("%d", webhook.Delay)}, "delaying webhook execution")
		time.Sleep(time.Duration(webhook.Delay) * time.Millisecond)
	}

	// Use background context to prevent cancellation when parent request finishes.
	webhookCtx := context.Background()
	timeout := s.resolveTimeout(webhook.Timeout)

	webhookCtx, cancel := context.WithTimeout(webhookCtx, timeout)
	defer cancel()

	url := s.applyVariables(webhook.URL, variables)
	body := s.applyVariables(webhook.Body, variables)
	headers := s.resolveHeaders(webhook.Headers, variables)

	//nolint:contextcheck
	req, err := s.buildRequest(webhookCtx, webhook.Method, url, body)
	if err != nil {
		logger.Error(s, nil, err, "error creating webhook request")
		notifyResult(onResult, url, webhook.Method, 0, time.Since(start), err.Error(), "")

		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(s, nil, err, "error sending webhook")
		notifyResult(onResult, url, webhook.Method, 0, time.Since(start), err.Error(), "")

		return
	}

	defer resp.Body.Close()

	respBodyBytes, _ := io.ReadAll(resp.Body)
	respBody := string(respBodyBytes)

	notifyResult(onResult, url, webhook.Method, resp.StatusCode, time.Since(start), "", respBody)

	logger.Info(s, map[string]string{
		"webhook_url": url,
		"status":      fmt.Sprintf("%d", resp.StatusCode),
	}, "webhook sent successfully")
}

func (s *webhookService) resolveTimeout(timeout *int) time.Duration {
	t := defaultWebhookTimeout
	if timeout != nil && *timeout > 0 {
		t = time.Duration(*timeout) * time.Millisecond
	}

	return t
}

func (s *webhookService) resolveHeaders(
	headers map[string]string,
	variables []*model.Variable,
) map[string]string {
	result := make(map[string]string, len(headers))

	for k, v := range headers {
		result[k] = s.applyVariables(v, variables)
	}

	return result
}

func notifyResult(
	onResult func(model.WebhookResult),
	url, method string,
	statusCode int,
	duration time.Duration,
	errMsg string,
	responseBody string,
) {
	if onResult == nil {
		return
	}

	result := model.WebhookResult{
		URL:          url,
		Method:       method,
		StatusCode:   statusCode,
		DurationMs:   duration.Milliseconds(),
		ResponseBody: responseBody,
	}

	if errMsg != "" {
		result.Error = errMsg
	}

	onResult(result)
}

func (s *webhookService) buildRequest(
	ctx context.Context,
	method, url, body string,
) (*http.Request, error) {
	if body != "" {
		//nolint:wrapcheck
		return http.NewRequestWithContext(ctx, method, url, bytes.NewBufferString(body))
	}

	//nolint:wrapcheck
	return http.NewRequestWithContext(ctx, method, url, nil)
}

func (s *webhookService) applyVariables(input string, variables []*model.Variable) string {
	for _, variable := range variables {
		input = strings.ReplaceAll(input, fmt.Sprintf("{%s}", variable.Name), variable.Value)
	}

	return input
}
