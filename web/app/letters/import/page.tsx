"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { useAuth } from "@/contexts/auth-context";
import { deleteLetter, listLetters } from "@/lib/api/letters";
import { useI18n } from "@/lib/i18n/i18n-context";
import type { Letter } from "@/types/letter";

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
  DialogDescription,
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

const DEFAULT_PAGE_SIZE = 20;

export default function LettersPage() {
  const { token, user } = useAuth();
  const { t } = useI18n();

  const [letters, setLetters] = useState<Letter[]>([]);
  const [search, setSearch] = useState("");
  const [appliedSearch, setAppliedSearch] = useState("");
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);

  const [selectedLetter, setSelectedLetter] = useState<Letter | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Letter | null>(null);
  const [deleteLoading, setDeleteLoading] = useState(false);

  const canCreate = user?.role === "superuser" || user?.role === "editor";
  const canDelete = user?.role === "superuser" || user?.role === "editor";

  const totalPages = useMemo(() => {
    return Math.max(1, Math.ceil(total / pageSize));
  }, [total, pageSize]);

  const loadLetters = useCallback(async () => {
    if (!token) return;

    setLoading(true);

    try {
      const result = await listLetters(token, {
        page,
        page_size: pageSize,
        search: appliedSearch,
      });

      setLetters(result.items);
      setTotal(result.total);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to load letters");
    } finally {
      setLoading(false);
    }
  }, [token, page, pageSize, appliedSearch]);

  useEffect(() => {
    const timeoutID = window.setTimeout(() => {
        void loadLetters();
    }, 0);

    return () => {
        window.clearTimeout(timeoutID);
    };
    }, [loadLetters]);

  async function handleDeleteLetter() {
    if (!token || !deleteTarget) return;

    setDeleteLoading(true);

    try {
      await deleteLetter(token, deleteTarget.id);

      toast.success(t.letterDeleted);
      setDeleteTarget(null);
      setSelectedLetter(null);

      await loadLetters();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete letter");
    } finally {
      setDeleteLoading(false);
    }
  }

  function handleSearchSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setPage(1);
    setAppliedSearch(search.trim());
  }

  function handlePageSizeChange(value: string) {
    setPage(1);
    setPageSize(Number(value));
  }

  return (
    <ProtectedRoute>
      <AppShell>
        <div className="min-w-0 space-y-6">
          <div className="flex flex-wrap items-start justify-between gap-4">
            <div>
              <h1 className="text-3xl font-bold tracking-tight">
                {t.letters}
              </h1>
              <p className="text-muted-foreground">{t.lettersDescription}</p>
            </div>

            {canCreate ? (
              <Link href="/letters/new">
                <Button>{t.newLetter}</Button>
              </Link>
            ) : null}
          </div>

          <Card className="max-w-full">
            <CardHeader>
              <CardTitle>{t.commonSearch}</CardTitle>
            </CardHeader>
            <CardContent>
              <form
                className="grid gap-3 md:grid-cols-[1fr_auto]"
                onSubmit={handleSearchSubmit}
              >
                <Input
                  placeholder={t.searchLettersPlaceholder}
                  value={search}
                  onChange={(event) => setSearch(event.target.value)}
                />
                <Button type="submit" disabled={loading}>
                  {loading ? t.commonLoading : t.commonSearch}
                </Button>
              </form>
            </CardContent>
          </Card>

          <Card className="max-w-full overflow-hidden">
            <CardHeader className="flex flex-row items-center justify-between gap-4">
              <CardTitle>{t.letters}</CardTitle>

              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <span dir="ltr">
                  {total.toLocaleString()} {t.letters}
                </span>

                <select
                  value={pageSize}
                  onChange={(event) => handlePageSizeChange(event.target.value)}
                  className="rounded-md border bg-background px-2 py-1 text-sm"
                >
                  <option value="10">10</option>
                  <option value="20">20</option>
                  <option value="50">50</option>
                  <option value="100">100</option>
                </select>
              </div>
            </CardHeader>

            <CardContent className="min-w-0">
              <div className="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead className="w-[130px] whitespace-nowrap">
                        {t.number}
                      </TableHead>
                      <TableHead>{t.letterTitle}</TableHead>
                      <TableHead className="w-[120px] whitespace-nowrap">
                        {t.letterDate}
                      </TableHead>
                      <TableHead className="hidden lg:table-cell">
                        {t.sender}
                      </TableHead>
                      <TableHead className="hidden xl:table-cell">
                        {t.receiver}
                      </TableHead>
                      <TableHead className="hidden 2xl:table-cell">
                        {t.registrar}
                      </TableHead>
                      <TableHead className="w-[90px] whitespace-nowrap">
                        {t.commonStatus}
                      </TableHead>
                      <TableHead className="w-[150px] text-end whitespace-nowrap">
                        {t.commonActions}
                      </TableHead>
                    </TableRow>
                  </TableHeader>

                  <TableBody>
                    {letters.length === 0 ? (
                      <TableRow>
                        <TableCell
                          colSpan={8}
                          className="text-center text-muted-foreground"
                        >
                          {loading ? t.commonLoading : t.noLettersFound}
                        </TableCell>
                      </TableRow>
                    ) : (
                      letters.map((letter) => (
                        <TableRow
                          key={letter.id}
                          className="cursor-pointer"
                          onClick={() => setSelectedLetter(letter)}
                        >
                          <TableCell
                            className="font-medium whitespace-nowrap"
                            dir="ltr"
                          >
                            {letter.formatted_letter_number}
                          </TableCell>

                          <TableCell className="max-w-[420px]">
                            <div className="line-clamp-2 font-medium">
                              {letter.title}
                            </div>
                            <div className="mt-1 text-xs text-muted-foreground lg:hidden">
                              {letter.sender} ← {letter.receiver}
                            </div>
                          </TableCell>

                          <TableCell className="whitespace-nowrap" dir="ltr">
                            {letter.letter_date_jalali}
                          </TableCell>

                          <TableCell className="hidden max-w-[260px] lg:table-cell">
                            <div className="line-clamp-2">{letter.sender}</div>
                          </TableCell>

                          <TableCell className="hidden max-w-[260px] xl:table-cell">
                            <div className="line-clamp-2">{letter.receiver}</div>
                          </TableCell>

                          <TableCell className="hidden 2xl:table-cell">
                            {letter.registrar_name}
                          </TableCell>

                          <TableCell>
                            {letter.is_deleted ? (
                              <Badge variant="destructive">
                                {t.commonDeleted}
                              </Badge>
                            ) : (
                              <Badge variant="secondary">
                                {t.commonActive}
                              </Badge>
                            )}
                          </TableCell>

                          <TableCell
                            className="text-end"
                            onClick={(event) => event.stopPropagation()}
                          >
                            <div className="flex justify-end gap-2">
                              <Button
                                variant="outline"
                                size="sm"
                                onClick={() => setSelectedLetter(letter)}
                              >
                                {t.commonView}
                              </Button>

                              {canDelete ? (
                                <Button
                                  variant="destructive"
                                  size="sm"
                                  onClick={() => setDeleteTarget(letter)}
                                >
                                  {t.commonDelete}
                                </Button>
                              ) : null}
                            </div>
                          </TableCell>
                        </TableRow>
                      ))
                    )}
                  </TableBody>
                </Table>
              </div>

              <div className="mt-4 flex flex-wrap items-center justify-between gap-3">
                <div className="text-sm text-muted-foreground" dir="ltr">
                  Page {page} of {totalPages}
                </div>

                <div className="flex items-center gap-2">
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    disabled={page <= 1 || loading}
                    onClick={() => setPage((current) => Math.max(1, current - 1))}
                  >
                    قبلی
                  </Button>

                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    disabled={page >= totalPages || loading}
                    onClick={() =>
                      setPage((current) => Math.min(totalPages, current + 1))
                    }
                  >
                    بعدی
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>

          <LetterPreviewDialog
            letter={selectedLetter}
            canDelete={canDelete}
            onClose={() => setSelectedLetter(null)}
            onDelete={(letter) => setDeleteTarget(letter)}
          />

          <DeleteLetterDialog
            letter={deleteTarget}
            loading={deleteLoading}
            onClose={() => setDeleteTarget(null)}
            onConfirm={() => void handleDeleteLetter()}
          />
        </div>
      </AppShell>
    </ProtectedRoute>
  );
}

function LetterPreviewDialog({
  letter,
  canDelete,
  onClose,
  onDelete,
}: {
  letter: Letter | null;
  canDelete: boolean;
  onClose: () => void;
  onDelete: (letter: Letter) => void;
}) {
  const { t } = useI18n();

  return (
    <Dialog open={Boolean(letter)} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>{t.letterDetails}</DialogTitle>
          <DialogDescription>
            {letter?.formatted_letter_number || ""}
          </DialogDescription>
        </DialogHeader>

        {letter ? (
          <div className="space-y-6">
            <div className="rounded-xl border bg-muted/30 p-6 text-center">
              <div className="text-sm text-muted-foreground">{t.number}</div>
              <div className="mt-2 text-4xl font-bold tracking-wide" dir="ltr">
                {letter.formatted_letter_number}
              </div>
            </div>

            <div className="grid gap-3 text-sm">
              <InfoRow label={t.letterTitle} value={letter.title} />
              <InfoRow
                label={t.letterDate}
                value={letter.letter_date_jalali}
                forceLtr
              />
              <InfoRow label={t.sender} value={letter.sender} />
              <InfoRow label={t.receiver} value={letter.receiver} />
              <InfoRow label={t.registrar} value={letter.registrar_name} />
              <InfoRow
                label={t.commonStatus}
                value={letter.is_deleted ? t.commonDeleted : t.commonActive}
                badge
                badgeVariant={letter.is_deleted ? "destructive" : "secondary"}
              />
              <InfoRow label={t.description} value={letter.description || "-"} />
              <InfoRow
                label={t.createdAt}
                value={new Date(letter.created_at).toLocaleString()}
                forceLtr
              />
            </div>

            <div className="flex flex-wrap gap-2">
              <Button variant="outline" onClick={onClose}>
                {t.commonCancel}
              </Button>

              {canDelete && !letter.is_deleted ? (
                <Button
                  variant="destructive"
                  onClick={() => {
                    onClose();
                    onDelete(letter);
                  }}
                >
                  {t.deleteLetter}
                </Button>
              ) : null}
            </div>
          </div>
        ) : null}
      </DialogContent>
    </Dialog>
  );
}

function DeleteLetterDialog({
  letter,
  loading,
  onClose,
  onConfirm,
}: {
  letter: Letter | null;
  loading: boolean;
  onClose: () => void;
  onConfirm: () => void;
}) {
  const { t } = useI18n();

  return (
    <Dialog open={Boolean(letter)} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>{t.deleteLetterConfirmTitle}</DialogTitle>
          <DialogDescription>
            {t.deleteLetterConfirmDescription}
          </DialogDescription>
        </DialogHeader>

        {letter ? (
          <div className="space-y-4">
            <div className="rounded-lg border bg-muted/30 p-4">
              <div className="text-sm text-muted-foreground">{t.number}</div>
              <div className="mt-1 font-bold" dir="ltr">
                {letter.formatted_letter_number}
              </div>
              <div className="mt-2 text-sm">{letter.title}</div>
            </div>

            <div className="flex gap-2">
              <Button
                variant="destructive"
                disabled={loading}
                onClick={onConfirm}
              >
                {loading ? t.commonLoading : t.confirmDelete}
              </Button>

              <Button variant="outline" onClick={onClose} disabled={loading}>
                {t.commonCancel}
              </Button>
            </div>
          </div>
        ) : null}
      </DialogContent>
    </Dialog>
  );
}

function InfoRow({
  label,
  value,
  forceLtr,
  badge,
  badgeVariant = "secondary",
}: {
  label: string;
  value: string;
  forceLtr?: boolean;
  badge?: boolean;
  badgeVariant?: "default" | "secondary" | "destructive" | "outline";
}) {
  return (
    <div className="flex items-start justify-between gap-4 border-b pb-2 last:border-b-0">
      <div className="text-muted-foreground">{label}</div>

      {badge ? (
        <Badge variant={badgeVariant}>{value}</Badge>
      ) : (
        <div
          className="max-w-md whitespace-pre-wrap break-words text-end font-medium"
          dir={forceLtr ? "ltr" : "auto"}
        >
          {value}
        </div>
      )}
    </div>
  );
}