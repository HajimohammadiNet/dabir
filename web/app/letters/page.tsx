"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { useAuth } from "@/contexts/auth-context";
import { listLetters } from "@/lib/api/letters";
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

  const canCreate = user?.role === "superuser" || user?.role === "editor";

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

  return (
    <ProtectedRoute>
      <AppShell>
        <div className="space-y-6">
          <div className="flex items-start justify-between gap-4">
            <div>
              <h1 className="text-3xl font-bold tracking-tight">
                {t.letters}
              </h1>
              <p className="text-muted-foreground">
                {t.lettersDescription}
              </p>
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
                    </TableRow>
                  </TableHeader>

                  <TableBody>
                    {letters.length === 0 ? (
                      <TableRow>
                        <TableCell
                          colSpan={7}
                          className="text-center text-muted-foreground"
                        >
                          {loading ? t.commonLoading : t.noLettersFound}
                        </TableCell>
                      </TableRow>
                    ) : (
                      letters.map((letter) => (
                        <TableRow key={letter.id}>
                          <TableCell className="font-medium">
                            {letter.formatted_letter_number}
                          </TableCell>
                          <TableCell>{letter.title}</TableCell>
                          <TableCell dir="ltr">{letter.letter_date_jalali}</TableCell>
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
                        </TableRow>
                      ))
                    )}
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>
        </div>
      </AppShell>
    </ProtectedRoute>
  );
}