"use client";

import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { LetterNumberText } from "@/components/common/letter-number-text";
import { useAuth } from "@/contexts/auth-context";
import { listLetters } from "@/lib/api/letters";
import { listUsers } from "@/lib/api/users";
import { useI18n } from "@/lib/i18n/i18n-context";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

type DashboardStats = {
  totalIncomingLetters: number | null;
  totalOutgoingLetters: number | null;
  usersCount: number | null;
  latestIncomingNumber: string | null;
  latestOutgoingNumber: string | null;
};

export default function DashboardPage() {
  const { token, user } = useAuth();
  const { t } = useI18n();

  const [stats, setStats] = useState<DashboardStats>({
    totalIncomingLetters: null,
    totalOutgoingLetters: null,
    usersCount: null,
    latestIncomingNumber: null,
    latestOutgoingNumber: null,
  });
  const [loading, setLoading] = useState(false);

  const loadStats = useCallback(async () => {
    if (!token || !user) return;

    setLoading(true);

    try {
      const [incomingResult, outgoingResult] = await Promise.all([
        listLetters(token, {
          direction: "incoming",
          page: 1,
          page_size: 1,
          sort_by: "created_at",
          sort_order: "desc",
        }),
        listLetters(token, {
          direction: "outgoing",
          page: 1,
          page_size: 1,
          sort_by: "created_at",
          sort_order: "desc",
        }),
      ]);

      let usersTotal: number | null = null;

      if (user.role === "superuser") {
        const usersResult = await listUsers(token, {
          page: 1,
          page_size: 1,
        });

        usersTotal = usersResult.total;
      }

      setStats({
        totalIncomingLetters: incomingResult.total,
        totalOutgoingLetters: outgoingResult.total,
        usersCount: usersTotal,
        latestIncomingNumber:
          incomingResult.items.length > 0
            ? incomingResult.items[0].formatted_letter_number
            : null,
        latestOutgoingNumber:
          outgoingResult.items.length > 0
            ? outgoingResult.items[0].formatted_letter_number
            : null,
      });
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to load dashboard stats"
      );
    } finally {
      setLoading(false);
    }
  }, [token, user]);

  useEffect(() => {
    const timeoutID = window.setTimeout(() => {
      void loadStats();
    }, 0);

    return () => {
      window.clearTimeout(timeoutID);
    };
  }, [loadStats]);

  return (
    <ProtectedRoute>
      <AppShell>
        <div className="space-y-6">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              {t.dashboard}
            </h1>
            <p className="text-muted-foreground">{t.dashboardDescription}</p>
          </div>

          <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-5">
            <Card>
              <CardHeader>
                <CardTitle>{t.totalIncomingLetters}</CardTitle>
              </CardHeader>
              <CardContent className="text-3xl font-bold" dir="ltr">
                {loading ? "..." : stats.totalIncomingLetters ?? 0}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>{t.totalOutgoingLetters}</CardTitle>
              </CardHeader>
              <CardContent className="text-3xl font-bold" dir="ltr">
                {loading ? "..." : stats.totalOutgoingLetters ?? 0}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>{t.latestIncomingNumber}</CardTitle>
              </CardHeader>
              <CardContent className="text-3xl font-bold">
                {loading ? (
                  "..."
                ) : stats.latestIncomingNumber ? (
                  <LetterNumberText value={stats.latestIncomingNumber} />
                ) : (
                  "-"
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>{t.latestOutgoingNumber}</CardTitle>
              </CardHeader>
              <CardContent className="text-3xl font-bold">
                {loading ? (
                  "..."
                ) : stats.latestOutgoingNumber ? (
                  <LetterNumberText value={stats.latestOutgoingNumber} />
                ) : (
                  "-"
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>{t.usersCount}</CardTitle>
              </CardHeader>
              <CardContent className="text-3xl font-bold" dir="ltr">
                {loading
                  ? "..."
                  : stats.usersCount !== null
                    ? stats.usersCount
                    : "-"}
              </CardContent>
            </Card>
          </div>
        </div>
      </AppShell>
    </ProtectedRoute>
  );
}
