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

  const letterConfig = settings?.letter_config;

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
                label={t.numberingMode}
                value={
                  letterConfig
                    ? formatNumberingMode(letterConfig.numbering_mode, t)
                    : "-"
                }
              />

              {letterConfig?.numbering_mode === "fixed_prefix" ? (
                <>
                  <InfoRow
                    label={t.numberPrefix}
                    value={letterConfig.number_prefix || "-"}
                  />

                  <InfoRow
                    label={t.numberPadding}
                    value={String(letterConfig.number_padding)}
                  />
                </>
              ) : null}

              {letterConfig?.numbering_mode === "jalali_yearly" ? (
                <>
                  <InfoRow
                    label={t.yearlySerialPadding}
                    value={String(letterConfig.yearly_serial_padding)}
                  />

                  <InfoRow
                    label={t.yearlySeparator}
                    value={letterConfig.yearly_separator || "-"}
                  />

                  <InfoRow
                    label={t.yearSource}
                    value={
                      letterConfig.year_source === "created_at"
                        ? t.yearSourceCreatedAt
                        : t.yearSourceLetterDate
                    }
                  />

                  <InfoRow
                    label={t.numberPrefix}
                    value={String(letterConfig.yearly_prefix_digits)}
                  />
                </>
              ) : null}

              {letterConfig ? (
                <div className="rounded-lg border bg-muted/30 p-4">
                  <div className="text-sm text-muted-foreground">
                    {t.exampleFormattedNumber}
                  </div>
                  <div className="mt-2">
                    <Badge variant="secondary" className="text-base" dir="ltr">
                      {formatExample(letterConfig)}
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
    <div className="flex items-center justify-between gap-4 border-b py-3 last:border-b-0">
      <div className="text-sm text-muted-foreground">{label}</div>
      <div className="font-medium text-end" dir="auto">
        {value}
      </div>
    </div>
  );
}

function formatNumberingMode(
  mode: PublicSettings["letter_config"]["numbering_mode"],
  t: ReturnType<typeof useI18n>["t"]
) {
  if (mode === "jalali_yearly") {
    return t.jalaliYearlyNumbering;
  }

  return t.fixedPrefixNumbering;
}

function formatExample(config: PublicSettings["letter_config"]) {
  if (config.numbering_mode === "jalali_yearly") {
    const separator = config.yearly_separator || "-";
    const serialPadding = config.yearly_serial_padding || 4;

    return `405${separator}${"1".padStart(serialPadding, "0")}`;
  }

  const prefix = config.number_prefix || "DABIR";
  const padding = config.number_padding || 6;

  return `${prefix}-${"1".padStart(padding, "0")}`;
}