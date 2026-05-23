"use client";

import { useCallback, useEffect, useState } from "react";
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

export default function LettersPage() {
  const { token, user } = useAuth();
  const { t } = useI18n();

  const [letters, setLetters] = useState<Letter[]>([]);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(false);

  const [selectedLetter, setSelectedLetter] = useState<Letter | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Letter | null>(null);
  const [deleteLoading, setDeleteLoading] = useState(false);

  const canCreate = user?.role === "superuser" || user?.role === "editor";
  const canDelete = user?.role === "superuser" || user?.role === "editor";

  const loadLetters = useCallback(async () => {
    if (!token) return;

    setLoading(true);

    try {
      const result = await listLetters(token, {
        page: 1,
        page_size: 20,
        search,
      });

      setLetters(result.items);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to load letters");
    } finally {
      setLoading(false);
    }
  }, [token, search]);

  useEffect(() => {
    const load = async () => {
      await loadLetters();
    };

    void load();
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

  return (
    <ProtectedRoute>
      <AppShell>
        <div className="space-y-6">
          <div className="flex items-start justify-between gap-4">
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

          <Card>
            <CardHeader>
              <CardTitle>{t.commonSearch}</CardTitle>
            </CardHeader>
            <CardContent>
              <form
                className="flex gap-2"
                onSubmit={(event) => {
                  event.preventDefault();
                  void loadLetters();
                }}
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

          <Card>
            <CardHeader>
              <CardTitle>{t.letters}</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>{t.number}</TableHead>
                      <TableHead>{t.letterTitle}</TableHead>
                      <TableHead>{t.letterDate}</TableHead>
                      <TableHead>{t.sender}</TableHead>
                      <TableHead>{t.receiver}</TableHead>
                      <TableHead>{t.registrar}</TableHead>
                      <TableHead>{t.commonStatus}</TableHead>
                      <TableHead className="text-end">
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
                          <TableCell className="font-medium" dir="ltr">
                            {letter.formatted_letter_number}
                          </TableCell>

                          <TableCell>{letter.title}</TableCell>

                          <TableCell dir="ltr">
                            {letter.letter_date_jalali}
                          </TableCell>

                          <TableCell>{letter.sender}</TableCell>

                          <TableCell>{letter.receiver}</TableCell>

                          <TableCell>{letter.registrar_name}</TableCell>

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
              <InfoRow
                label={t.description}
                value={letter.description || "-"}
              />
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
        <div className="max-w-md whitespace-pre-wrap break-words text-end font-medium" dir={forceLtr ? "ltr" : "auto"}>
          {value}
        </div>
      )}
    </div>
  );
}