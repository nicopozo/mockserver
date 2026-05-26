export interface Assertion {
  fail_on_error: boolean;
  type: string;
  variable_name?: string;
  value?: string;
  min?: number;
  max?: number;
}

export interface Variable {
  type: string;
  name: string;
  key: string;
  min?: number;
  max?: number;
  decimals?: number;
  assertions: Assertion[];
}

export interface WebhookConfig {
  url: string;
  method: string;
  headers: Record<string, string>;
  body: string;
  enabled: boolean;
  timeout?: number;
}

export interface Response {
  description: string;
  body: string;
  content_type: string;
  http_status: number;
  delay: number;
  scene?: string;
  webhook?: WebhookConfig;
}

export interface Mock {
  key?: string;
  group: string;
  name: string;
  path: string;
  strategy: string;
  method: string;
  status: 'enabled' | 'disabled';
  responses: Response[];
  variables: Variable[];
}

export interface Paging {
  total: number;
  limit: number;
  offset: number;
  last_id?: string;
}

export interface PaginatedMocks {
  results: Mock[];
  paging: Paging;
}

export interface WebhookResult {
  url: string;
  method: string;
  status_code: number;
  duration_ms: number;
  error?: string;
  response_body?: string;
}

export interface LogEntry {
  id: string;
  timestamp: string;
  method: string;
  url: string;
  request_body: string;
  request_headers: Record<string, string>;
  query_params: Record<string, string>;
  response_status: number;
  response_body: string;
  assertion_errors?: string[];
  webhook_results?: WebhookResult[];
}

export interface LogList {
  results: LogEntry[];
  paging: Paging;
}
