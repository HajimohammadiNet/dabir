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
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

export default function ImportLettersPage() {
  const { token } = useAuth();
  const { t } = useI18n();

  const [file, setFile] = useState<File | null>(null);
  const [importJob, setImportJob] = useState<ImportJob | null>(null);
  const [previewLoading, setPreviewLoading] = useState(false);
  const [commitLoading, setCommitLoading] = useState(false);

  const canCommit =
    importJob &&
    importJob.status === "previewed" &&
    importJob.invalid_rows === 0 &&
    importJob.valid_rows > 0;

  async function handlePreview(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!token) return;

    if (!file) {
      toast.error("Please select an Excel file");
      return;
    }

    setPreviewLoading(true);
    setImportJob(null);

    try {
      const result = await previewLettersImport(token, file);
      setImportJob(result);

      if (result.invalid_rows > 0) {
        toast.warning("Preview completed with errors");
      } else {
        toast.success("Preview completed successfully");
      }
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Preview failed");
    } finally {
      setPreviewLoading(false);
    }
  }

  async function handleCommit() {
    if (!token || !importJob) return;

    setCommitLoading(true);

    try {
      const result = await commitLettersImport(token, importJob.id);

      toast.success(
        `Imported ${result.imported_rows} rows. Next number: ${result.next_letter_number}`
      );

      setImportJob({
        ...importJob,
        status: "committed",
        committed_at: new Date().toISOString(),
      });
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Commit failed");
    } finally {
      setCommitLoading(false);
    }
  }

  return (
    <ProtectedRoute allowedRoles={["superuser"]}>
      <AppShell>
        <div className="space-y-6">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              {t.importLetters}
            </h1>
            <p className="text-muted-foreground">
              {t.importLettersDescription}
            </p>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>{t.uploadExcelFile}</CardTitle>
              <CardDescription>{t.uploadExcelDescription}</CardDescription>
            </CardHeader>

            <CardContent>
              <form onSubmit={handlePreview} className="space-y-4">
                <Input
                  type="file"
                  accept=".xlsx"
                  onChange={(event) => {
                    setFile(event.target.files?.[0] || null);
                    setImportJob(null);
                  }}
                />

                <Button type="submit" disabled={previewLoading}>
                  {previewLoading ? t.previewing : t.previewImport}
                </Button>
              </form>
            </CardContent>
          </Card>

          {importJob ? (
            <>
              <Card>
                <CardHeader>
                  <CardTitle>{t.importSummary}</CardTitle>
                  <CardDescription>{importJob.file_name}</CardDescription>
                </CardHeader>

                <CardContent>
                  <div className="grid gap-4 md:grid-cols-5">
                    <SummaryItem label={t.commonStatus}>
                      <Badge
                        variant={
                          importJob.status === "committed"
                            ? "secondary"
                            : "outline"
                        }
                      >
                        {importJob.status}
                      </Badge>
                    </SummaryItem>

                    <SummaryItem label={t.totalRows}>
                      {importJob.total_rows}
                    </SummaryItem>

                    <SummaryItem label={t.validRows}>
                      {importJob.valid_rows}
                    </SummaryItem>

                    <SummaryItem label={t.invalidRows}>
                      {importJob.invalid_rows}
                    </SummaryItem>

                    <SummaryItem label={t.maxNumber}>
                      {importJob.max_letter_number || "-"}
                    </SummaryItem>
                  </div>

                  <div className="mt-6">
                    <Button
                      onClick={handleCommit}
                      disabled={!canCommit || commitLoading}
                    >
                      {commitLoading ? t.committing : t.commitImport}
                    </Button>

                    {!canCommit && importJob.status === "previewed" ? (
                      <p className="mt-2 text-sm text-muted-foreground">
                        Import can be committed only when there are no invalid
                        rows.
                      </p>
                    ) : null}
                  </div>
                </CardContent>
              </Card>

              <DetectedColumnsCard importJob={importJob} />

              <ImportErrorsCard importJob={importJob} />

              <PreviewRowsCard importJob={importJob} />
            </>
          ) : null}
        </div>
      </AppShell>
    </ProtectedRoute>
  );
}

function SummaryItem({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <div className="rounded-lg border bg-background p-4">
      <div className="text-sm text-muted-foreground">{label}</div>
      <div className="mt-1 text-2xl font-semibold">{children}</div>
    </div>
  );
}

function DetectedColumnsCard({ importJob }: { importJob: ImportJob }) {
  const { t } = useI18n();

  const columns = importJob.detected_columns || {};
  const entries = Object.entries(columns);

  if (entries.length === 0) {
    return null;
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t.detectedColumns}</CardTitle>
      </CardHeader>

      <CardContent>
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t.field}</TableHead>
                <TableHead>{t.detectedColumn}</TableHead>
              </TableRow>
            </TableHeader>

            <TableBody>
              {entries.map(([field, column]) => (
                <TableRow key={field}>
                  <TableCell className="font-medium">{field}</TableCell>
                  <TableCell>{column}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  );
}

function ImportErrorsCard({ importJob }: { importJob: ImportJob }) {
  const { t } = useI18n();

  const errors = importJob.errors || [];

  if (errors.length === 0) {
    return (
      <Alert>
        <AlertTitle>{t.noErrorsFound}</AlertTitle>
        <AlertDescription>{t.importReady}</AlertDescription>
      </Alert>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t.importErrors}</CardTitle>
        <CardDescription>{t.importErrorsDescription}</CardDescription>
      </CardHeader>

      <CardContent>
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t.excelRow}</TableHead>
                <TableHead>{t.field}</TableHead>
                <TableHead>{t.description}</TableHead>
              </TableRow>
            </TableHeader>

            <TableBody>
              {errors.map((error, index) => (
                <TableRow key={`${error.row}-${error.field}-${index}`}>
                  <TableCell>{error.row}</TableCell>
                  <TableCell>{error.field}</TableCell>
                  <TableCell>{error.message}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  );
}

function PreviewRowsCard({ importJob }: { importJob: ImportJob }) {
  const { t } = useI18n();

  const rows = importJob.preview_data || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t.previewRows}</CardTitle>
        <CardDescription>{t.previewRowsDescription}</CardDescription>
      </CardHeader>

      <CardContent>
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t.excelRow}</TableHead>
                <TableHead>{t.number}</TableHead>
                <TableHead>{t.letterTitle}</TableHead>
                <TableHead>{t.letterDate}</TableHead>
                <TableHead>{t.sender}</TableHead>
                <TableHead>{t.receiver}</TableHead>
              </TableRow>
            </TableHeader>

            <TableBody>
              {rows.length === 0 ? (
                <TableRow>
                  <TableCell
                    colSpan={6}
                    className="text-center text-muted-foreground"
                  >
                    {t.noValidRowsFound}
                  </TableCell>
                </TableRow>
              ) : (
                rows.map((row) => (
                  <TableRow key={`${row.row_number}-${row.letter_number}`}>
                    <TableCell>{row.row_number}</TableCell>
                    <TableCell>{row.letter_number}</TableCell>
                    <TableCell>{row.title}</TableCell>
                    <TableCell>{row.letter_date}</TableCell>
                    <TableCell>{row.sender}</TableCell>
                    <TableCell>{row.receiver}</TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  );
}