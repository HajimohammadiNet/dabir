"use client";

import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { JsonViewer } from "@/components/common/json-viewer";
import { useAuth } from "@/contexts/auth-context";
import { listAuditLogs } from "@/lib/api/audit";
import { useI18n } from "@/lib/i18n/i18n-context";
import type { AuditLog } from "@/types/audit";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

const actionOptions = [
  "",
  "setup.initialized",
  "auth.login_success",
  "auth.login_failed",
  "user.created",
  "user.updated",
  "user.activated",
  "user.deactivated",
  "user.password_changed",
  "user.password_reset",
  "letter.created",
  "letter.updated",
  "letter.deleted",
  "letters.import_previewed",
  "letters.import_committed",
];

const entityTypeOptions = ["", "setup", "auth", "user", "letter", "import_job"];

export default function AuditLogsPage() {
  const { token } = useAuth();
  const { t } = useI18n();

  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [action, setAction] = useState("");
  const [entityType, setEntityType] = useState("");
  const [actorUserID, setActorUserID] = useState("");
  const [loading, setLoading] = useState(false);
  const [selectedLog, setSelectedLog] = useState<AuditLog | null>(null);

  const loadLogs = useCallback(async () => {
    if (!token) return;

    setLoading(true);

    try {
      const result = await listAuditLogs(token, {
        page: 1,
        page_size: 50,
        action,
        entity_type: entityType,
        actor_user_id: actorUserID,
      });

      setLogs(result.items);
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to load audit logs"
      );
    } finally {
      setLoading(false);
    }
  }, [token, action, entityType, actorUserID]);

  useEffect(() => {
    const load = async () => {
      await loadLogs();
    };

    void load();
  }, [loadLogs]);

  return (
    <ProtectedRoute allowedRoles={["superuser"]}>
      <AppShell>
        <div className="space-y-6">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              {t.auditLogs}
            </h1>
            <p className="text-muted-foreground">
              {t.auditLogsDescription}
            </p>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>{t.filters}</CardTitle>
            </CardHeader>

            <CardContent>
              <form
                className="grid gap-4 md:grid-cols-4"
                onSubmit={(event) => {
                  event.preventDefault();
                  void loadLogs();
                }}
              >
                <div className="space-y-2">
                  <label className="text-sm font-medium">{t.action}</label>
                  <select
                    value={action}
                    onChange={(event) => setAction(event.target.value)}
                    className="w-full rounded-md border bg-background px-3 py-2 text-sm"
                  >
                    {actionOptions.map((item) => (
                      <option key={item || "all"} value={item}>
                        {item || t.allActions}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium">
                    {t.entityType}
                  </label>
                  <select
                    value={entityType}
                    onChange={(event) => setEntityType(event.target.value)}
                    className="w-full rounded-md border bg-background px-3 py-2 text-sm"
                  >
                    {entityTypeOptions.map((item) => (
                      <option key={item || "all"} value={item}>
                        {item || t.allEntities}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium">
                    {t.actorUserId}
                  </label>
                  <Input
                    value={actorUserID}
                    onChange={(event) => setActorUserID(event.target.value)}
                    placeholder="UUID"
                    dir="ltr"
                  />
                </div>

                <div className="flex items-end">
                  <Button type="submit" disabled={loading} className="w-full">
                    {loading ? t.commonLoading : t.applyFilters}
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>{t.auditLogs}</CardTitle>
            </CardHeader>

            <CardContent>
              <div className="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>{t.createdAt}</TableHead>
                      <TableHead>{t.action}</TableHead>
                      <TableHead>{t.entityType}</TableHead>
                      <TableHead>{t.actorUserId}</TableHead>
                      <TableHead>{t.ipAddress}</TableHead>
                      <TableHead className="text-end">
                        {t.commonView}
                      </TableHead>
                    </TableRow>
                  </TableHeader>

                  <TableBody>
                    {logs.length === 0 ? (
                      <TableRow>
                        <TableCell
                          colSpan={6}
                          className="text-center text-muted-foreground"
                        >
                          {loading ? t.commonLoading : t.noAuditLogsFound}
                        </TableCell>
                      </TableRow>
                    ) : (
                      logs.map((log) => (
                        <TableRow key={log.id}>
                          <TableCell>
                            {new Date(log.created_at).toLocaleString()}
                          </TableCell>

                          <TableCell>
                            <Badge variant="outline">{log.action}</Badge>
                          </TableCell>

                          <TableCell>
                            <div className="font-medium">
                              {log.entity_type}
                            </div>
                            {log.entity_id ? (
                              <div className="text-xs text-muted-foreground break-all">
                                {log.entity_id}
                              </div>
                            ) : null}
                          </TableCell>

                          <TableCell className="text-xs break-all">
                            {log.actor_user_id || "-"}
                          </TableCell>

                          <TableCell dir="ltr">{log.ip_address || "-"}</TableCell>

                          <TableCell className="text-end">
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => setSelectedLog(log)}
                            >
                              {t.commonView}
                            </Button>
                          </TableCell>
                        </TableRow>
                      ))
                    )}
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>

          <Dialog
            open={Boolean(selectedLog)}
            onOpenChange={(open) => {
              if (!open) setSelectedLog(null);
            }}
          >
            <DialogContent className="max-w-4xl">
              <DialogHeader>
                <DialogTitle>{t.auditLogDetails}</DialogTitle>
              </DialogHeader>

              {selectedLog ? (
                <div className="space-y-4">
                  <div className="grid gap-4 md:grid-cols-2 text-sm">
                    <InfoItem label="ID" value={selectedLog.id} />
                    <InfoItem label={t.action} value={selectedLog.action} />
                    <InfoItem
                      label={t.entityType}
                      value={selectedLog.entity_type}
                    />
                    <InfoItem
                      label={t.entityId}
                      value={selectedLog.entity_id || "-"}
                    />
                    <InfoItem
                      label={t.actorUserId}
                      value={selectedLog.actor_user_id || "-"}
                    />
                    <InfoItem
                      label={t.ipAddress}
                      value={selectedLog.ip_address || "-"}
                    />
                    <InfoItem
                      label={t.createdAt}
                      value={new Date(selectedLog.created_at).toLocaleString()}
                    />
                    <InfoItem
                      label={t.userAgent}
                      value={selectedLog.user_agent || "-"}
                    />
                  </div>

                  <div className="grid gap-4 md:grid-cols-2">
                    <div>
                      <h3 className="mb-2 text-sm font-medium">
                        {t.oldValue}
                      </h3>
                      <JsonViewer value={selectedLog.old_value} />
                    </div>

                    <div>
                      <h3 className="mb-2 text-sm font-medium">
                        {t.newValue}
                      </h3>
                      <JsonViewer value={selectedLog.new_value} />
                    </div>
                  </div>
                </div>
              ) : null}
            </DialogContent>
          </Dialog>
        </div>
      </AppShell>
    </ProtectedRoute>
  );
}

function InfoItem({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <div className="text-muted-foreground">{label}</div>
      <div className="break-all font-medium" dir="auto">
        {value}
      </div>
    </div>
  );
}