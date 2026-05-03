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
  assertions: Assertion[];
}

export interface Response {
  description: string;
  body: string;
  content_type: string;
  http_status: number;
  delay: number;
  scene?: string;
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
}

export interface PaginatedMocks {
  results: Mock[];
  paging: Paging;
}
