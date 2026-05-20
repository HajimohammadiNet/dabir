export type AuditLog = {
  id: string;
  actor_user_id?: string | null;

  action: string;
  entity_type: string;
  entity_id?: string | null;

  old_value?: unknown;
  new_value?: unknown;

  ip_address?: string | null;
  user_agent?: string | null;

  created_at: string;
};

export type ListAuditLogsResponse = {
  items: AuditLog[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
};