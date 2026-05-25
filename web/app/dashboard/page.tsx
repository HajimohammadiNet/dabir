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
  totalLetters: number | null;
  usersCount: number | null;
  lastNumber: string | null;
};

export default function DashboardPage() {
  const { token, user } = useAuth();
  const { t } = useI18n();

  const [stats, setStats] = useState<DashboardStats>({
    totalLetters: null,
    usersCount: null,
    lastNumber: null,
  });
  const [loading, setLoading] = useState(false);

  const loadStats = useCallback(async () => {
    if (!token || !user) return;

    setLoading(true);

    try {
      const lettersResult = await listLetters(token, {
        page: 1,
        page_size: 1,
        sort_by: "created_at",
        sort_order: "desc",
      });

      let usersTotal: number | null = null;

      if (user.role === "superuser") {
        const usersResult = await listUsers(token, {
          page: 1,
          page_size: 1,
        });

        usersTotal = usersResult.total;
      }

      setStats({
        totalLetters: lettersResult.total,
        usersCount: usersTotal,
        lastNumber:
          lettersResult.items.length > 0
            ? lettersResult.items[0].formatted_letter_number
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

          <div className="grid gap-4 md:grid-cols-3">
            <Card>
              <CardHeader>
                <CardTitle>{t.totalLetters}</CardTitle>
              </CardHeader>
              <CardContent className="text-3xl font-bold" dir="ltr">
                {loading ? "..." : stats.totalLetters ?? 0}
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

            <Card>
              <CardHeader>
                <CardTitle>{t.lastNumber}</CardTitle>
              </CardHeader>
              <CardContent className="text-3xl font-bold">
                {loading ? (
                  "..."
                ) : stats.lastNumber ? (
                  <LetterNumberText value={stats.lastNumber} />
                ) : (
                  "-"
                )}
              </CardContent>
            </Card>
          </div>
        </div>
      </AppShell>
    </ProtectedRoute>
  );
}