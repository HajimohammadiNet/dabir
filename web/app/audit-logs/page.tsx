"use client";

import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { JsonViewer } from "@/components/common/json-viewer";
import { useAuth } from "@/contexts/auth-context";
import { listAuditLogs } from "@/lib/api/audit";
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
  "letter.created",
  "letter.updated",
  "letter.deleted",
  "letters.import_previewed",
  "letters.import_committed",
];

const entityTypeOptions = ["", "setup", "auth", "user", "letter", "import_job"];

export default function AuditLogsPage() {
  const { token } = useAuth();

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
            <h1 className="text-3xl font-bold tracking-tight">Audit Logs</h1>
            <p className="text-muted-foreground">
              Review important system activities and changes.
            </p>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Filters</CardTitle>
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
                  <label className="text-sm font-medium">Action</label>
                  <select
                    value={action}
                    onChange={(event) => setAction(event.target.value)}
                    className="w-full rounded-md border bg-background px-3 py-2 text-sm"
                  >
                    {actionOptions.map((item) => (
                      <option key={item || "all"} value={item}>
                        {item || "All actions"}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium">Entity Type</label>
                  <select
                    value={entityType}
                    onChange={(event) => setEntityType(event.target.value)}
                    className="w-full rounded-md border bg-background px-3 py-2 text-sm"
                  >
                    {entityTypeOptions.map((item) => (
                      <option key={item || "all"} value={item}>
                        {item || "All entities"}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium">Actor User ID</label>
                  <Input
                    value={actorUserID}
                    onChange={(event) => setActorUserID(event.target.value)}
                    placeholder="UUID"
                  />
                </div>

                <div className="flex items-end">
                  <Button type="submit" disabled={loading} className="w-full">
                    {loading ? "Loading..." : "Apply Filters"}
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Audit Logs</CardTitle>
            </CardHeader>

            <CardContent>
              <div className="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Time</TableHead>
                      <TableHead>Action</TableHead>
                      <TableHead>Entity</TableHead>
                      <TableHead>Actor</TableHead>
                      <TableHead>IP</TableHead>
                      <TableHead className="text-right">Details</TableHead>
                    </TableRow>
                  </TableHeader>

                  <TableBody>
                    {logs.length === 0 ? (
                      <TableRow>
                        <TableCell
                          colSpan={6}
                          className="text-center text-muted-foreground"
                        >
                          {loading ? "Loading..." : "No audit logs found."}
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
                              <div className="text-xs text-muted-foreground">
                                {log.entity_id}
                              </div>
                            ) : null}
                          </TableCell>

                          <TableCell className="text-xs">
                            {log.actor_user_id || "-"}
                          </TableCell>

                          <TableCell>{log.ip_address || "-"}</TableCell>

                          <TableCell className="text-right">
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => setSelectedLog(log)}
                            >
                              View
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
                <DialogTitle>Audit Log Details</DialogTitle>
              </DialogHeader>

              {selectedLog ? (
                <div className="space-y-4">
                  <div className="grid gap-4 md:grid-cols-2 text-sm">
                    <InfoItem label="ID" value={selectedLog.id} />
                    <InfoItem label="Action" value={selectedLog.action} />
                    <InfoItem
                      label="Entity Type"
                      value={selectedLog.entity_type}
                    />
                    <InfoItem
                      label="Entity ID"
                      value={selectedLog.entity_id || "-"}
                    />
                    <InfoItem
                      label="Actor User ID"
                      value={selectedLog.actor_user_id || "-"}
                    />
                    <InfoItem
                      label="IP Address"
                      value={selectedLog.ip_address || "-"}
                    />
                    <InfoItem
                      label="Created At"
                      value={new Date(selectedLog.created_at).toLocaleString()}
                    />
                    <InfoItem
                      label="User Agent"
                      value={selectedLog.user_agent || "-"}
                    />
                  </div>

                  <div className="grid gap-4 md:grid-cols-2">
                    <div>
                      <h3 className="mb-2 text-sm font-medium">Old Value</h3>
                      <JsonViewer value={selectedLog.old_value} />
                    </div>

                    <div>
                      <h3 className="mb-2 text-sm font-medium">New Value</h3>
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
      <div className="break-all font-medium">{value}</div>
    </div>
  );
}