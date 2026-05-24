"use client";

import { FormEvent, useState } from "react";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { useAuth } from "@/contexts/auth-context";
import {
  commitLettersImport,
  previewLettersImport,
} from "@/lib/api/imports";
import { useI18n } from "@/lib/i18n/i18n-context";
import type { ImportJob } from "@/types/import";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

export default function ImportLettersPage() {
  const { token } = useAuth();
  const { t } = useI18n();

  const [file, setFile] = useState<File | null>(null);
  const [job, setJob] = useState<ImportJob | null>(null);
  const [previewLoading, setPreviewLoading] = useState(false);
  const [commitLoading, setCommitLoading] = useState(false);
  const [committed, setCommitted] = useState(false);

  async function handlePreview(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!token) return;

    if (!file) {
      toast.error("Please select an Excel file");
      return;
    }

    setPreviewLoading(true);
    setCommitted(false);

    try {
      const result = await previewLettersImport(token, file);
      setJob(result);

      if (result.invalid_rows > 0) {
        toast.warning(
          `Preview completed with ${result.invalid_rows} invalid rows`
        );
      } else {
        toast.success("Import preview completed");
      }
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Import preview failed");
    } finally {
      setPreviewLoading(false);
    }
  }

  async function handleCommit() {
    if (!token || !job) return;

    setCommitLoading(true);

    try {
      const result = await commitLettersImport(token, job.id);

      setCommitted(true);

      toast.success(
        `Imported ${result.imported_rows} rows. Skipped ${result.skipped_rows} rows.`
      );
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Import commit failed");
    } finally {
      setCommitLoading(false);
    }
  }

  const previewRows = job?.preview_data || [];
  const errors = job?.errors || [];
  const canCommit =
    Boolean(job) &&
    job?.status === "previewed" &&
    !committed &&
    (job?.valid_rows || 0) > 0;

  return (
    <ProtectedRoute allowedRoles={["superuser"]}>
      <AppShell>
        <div className="min-w-0 space-y-6 pb-16 md:pb-0">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              {t.import}
            </h1>
            <p className="text-muted-foreground">
              Import letters from an Excel file.
            </p>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Upload Excel File</CardTitle>
              <CardDescription>
                Select an `.xlsx` file and preview rows before committing them
                to the database.
              </CardDescription>
            </CardHeader>

            <CardContent>
              <form onSubmit={handlePreview} className="space-y-4">
                <div className="space-y-2">
                  <Input
                    type="file"
                    accept=".xlsx"
                    onChange={(event) => {
                      setFile(event.target.files?.[0] || null);
                      setJob(null);
                      setCommitted(false);
                    }}
                  />
                </div>

                <div className="flex flex-wrap gap-2">
                  <Button type="submit" disabled={previewLoading || !file}>
                    {previewLoading ? "Previewing..." : "Preview Import"}
                  </Button>

                  {job ? (
                    <Button
                      type="button"
                      variant="default"
                      disabled={!canCommit || commitLoading}
                      onClick={() => void handleCommit()}
                    >
                      {commitLoading ? "Importing..." : "Commit Import"}
                    </Button>
                  ) : null}
                </div>
              </form>
            </CardContent>
          </Card>

          {job ? (
            <Card>
              <CardHeader>
                <CardTitle>Import Summary</CardTitle>
                <CardDescription>{job.file_name}</CardDescription>
              </CardHeader>

              <CardContent className="grid gap-3 md:grid-cols-4">
                <SummaryItem label="Total Rows" value={job.total_rows} />
                <SummaryItem label="Valid Rows" value={job.valid_rows} />
                <SummaryItem label="Invalid Rows" value={job.invalid_rows} />
                <SummaryItem
                  label="Status"
                  value={committed ? "committed" : job.status}
                />
              </CardContent>
            </Card>
          ) : null}

          {job?.detected_columns ? (
            <Card>
              <CardHeader>
                <CardTitle>Detected Columns</CardTitle>
              </CardHeader>

              <CardContent>
                <div className="grid gap-2 md:grid-cols-2">
                  {Object.entries(job.detected_columns).map(([key, value]) => (
                    <div
                      key={key}
                      className="flex items-center justify-between gap-4 rounded-md border p-3 text-sm"
                    >
                      <span className="text-muted-foreground">{key}</span>
                      <span className="font-medium">{value}</span>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          ) : null}

          {previewRows.length > 0 ? (
            <Card className="max-w-full overflow-hidden">
              <CardHeader>
                <CardTitle>Preview Rows</CardTitle>
                <CardDescription>
                  Showing preview rows returned by the backend.
                </CardDescription>
              </CardHeader>

              <CardContent className="min-w-0">
                <div className="rounded-md border">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead className="w-[90px] whitespace-nowrap">
                          Row
                        </TableHead>
                        <TableHead className="w-[140px] whitespace-nowrap">
                          Number
                        </TableHead>
                        <TableHead>Title</TableHead>
                        <TableHead className="w-[130px] whitespace-nowrap">
                          Date
                        </TableHead>
                        <TableHead>Sender</TableHead>
                        <TableHead>Receiver</TableHead>
                      </TableRow>
                    </TableHeader>

                    <TableBody>
                      {previewRows.map((row) => (
                        <TableRow key={row.row_number}>
                          <TableCell className="whitespace-nowrap" dir="ltr">
                            {row.row_number}
                          </TableCell>

                          <TableCell className="whitespace-nowrap" dir="ltr">
                            {row.display_letter_number || row.letter_number}
                          </TableCell>

                          <TableCell className="max-w-[420px]">
                            <div className="line-clamp-2">{row.title}</div>
                          </TableCell>

                          <TableCell className="whitespace-nowrap" dir="ltr">
                            {row.letter_date_jalali || row.letter_date}
                          </TableCell>

                          <TableCell className="max-w-[260px]">
                            <div className="line-clamp-2">{row.sender}</div>
                          </TableCell>

                          <TableCell className="max-w-[260px]">
                            <div className="line-clamp-2">{row.receiver}</div>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>
              </CardContent>
            </Card>
          ) : null}

          {errors.length > 0 ? (
            <Card>
              <CardHeader>
                <CardTitle>Validation Errors</CardTitle>
                <CardDescription>
                  These rows were rejected during preview.
                </CardDescription>
              </CardHeader>

              <CardContent>
                <div className="rounded-md border">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead className="w-[90px] whitespace-nowrap">
                          Row
                        </TableHead>
                        <TableHead className="w-[160px] whitespace-nowrap">
                          Field
                        </TableHead>
                        <TableHead>Message</TableHead>
                      </TableRow>
                    </TableHeader>

                    <TableBody>
                      {errors.map((error, index) => (
                        <TableRow key={`${error.row}-${error.field}-${index}`}>
                          <TableCell className="whitespace-nowrap" dir="ltr">
                            {error.row}
                          </TableCell>
                          <TableCell>{error.field}</TableCell>
                          <TableCell>{error.message}</TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>
              </CardContent>
            </Card>
          ) : null}
        </div>
      </AppShell>
    </ProtectedRoute>
  );
}

function SummaryItem({
  label,
  value,
}: {
  label: string;
  value: string | number;
}) {
  return (
    <div className="rounded-lg border bg-muted/30 p-4">
      <div className="text-sm text-muted-foreground">{label}</div>
      <div className="mt-2">
        <Badge variant="secondary" className="text-base">
          {value}
        </Badge>
      </div>
    </div>
  );
}