"use client";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { useI18n } from "@/lib/i18n/i18n-context";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function DashboardPage() {
  const { t } = useI18n();

  return (
    <ProtectedRoute>
      <AppShell>
        <div className="space-y-6">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              {t.dashboard}
            </h1>
            <p className="text-muted-foreground">
              {t.dashboardDescription}
            </p>
          </div>

          <div className="grid gap-4 md:grid-cols-3">
            <Card>
              <CardHeader>
                <CardTitle>{t.totalLetters}</CardTitle>
              </CardHeader>
              <CardContent className="text-3xl font-bold">-</CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>{t.usersCount}</CardTitle>
              </CardHeader>
              <CardContent className="text-3xl font-bold">-</CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>{t.lastNumber}</CardTitle>
              </CardHeader>
              <CardContent className="text-3xl font-bold">-</CardContent>
            </Card>
          </div>
        </div>
      </AppShell>
    </ProtectedRoute>
  );
}