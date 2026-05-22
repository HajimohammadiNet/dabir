"use client";

import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { getPublicSettings } from "@/lib/api/settings";
import { useI18n } from "@/lib/i18n/i18n-context";
import type { PublicSettings } from "@/types/settings";

import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

export default function SettingsPage() {
  const { t } = useI18n();

  const [settings, setSettings] = useState<PublicSettings | null>(null);
  const [loading, setLoading] = useState(false);

  const loadSettings = useCallback(async () => {
    setLoading(true);

    try {
      const result = await getPublicSettings();
      setSettings(result);
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to load settings"
      );
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    const load = async () => {
      await loadSettings();
    };

    void load();
  }, [loadSettings]);

  return (
    <ProtectedRoute allowedRoles={["superuser"]}>
      <AppShell>
        <div className="space-y-6">
          <div className="flex items-start justify-between gap-4">
            <div>
              <h1 className="text-3xl font-bold tracking-tight">
                {t.settings}
              </h1>
              <p className="text-muted-foreground">
                {t.settingsDescription}
              </p>
            </div>

            <Button variant="outline" onClick={() => void loadSettings()}>
              {loading ? t.commonLoading : t.commonRefresh}
            </Button>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>{t.organization}</CardTitle>
              <CardDescription>
                {t.organizationSettingsDescription}
              </CardDescription>
            </CardHeader>

            <CardContent>
              <InfoRow
                label={t.organizationName}
                value={settings?.organization_name || "-"}
              />
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>{t.letterNumbering}</CardTitle>
              <CardDescription>
                {t.letterNumberingDescription}
              </CardDescription>
            </CardHeader>

            <CardContent className="space-y-4">
              <InfoRow
                label={t.numberPrefix}
                value={settings?.letter_config.number_prefix || "-"}
              />

              <InfoRow
                label={t.numberPadding}
                value={
                  settings
                    ? String(settings.letter_config.number_padding)
                    : "-"
                }
              />

              {settings ? (
                <div className="rounded-lg border bg-muted/30 p-4">
                  <div className="text-sm text-muted-foreground">
                    {t.exampleFormattedNumber}
                  </div>
                  <div className="mt-2">
                    <Badge variant="secondary" className="text-base">
                      {formatExample(
                        settings.letter_config.number_prefix,
                        settings.letter_config.number_padding
                      )}
                    </Badge>
                  </div>
                </div>
              ) : null}
            </CardContent>
          </Card>
        </div>
      </AppShell>
    </ProtectedRoute>
  );
}

function InfoRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between border-b py-3 last:border-b-0">
      <div className="text-sm text-muted-foreground">{label}</div>
      <div className="font-medium">{value}</div>
    </div>
  );
}

function formatExample(prefix: string, padding: number) {
  const number = "1".padStart(padding || 6, "0");
  return `${prefix || "DABIR"}-${number}`;
}