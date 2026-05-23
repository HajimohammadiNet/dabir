"use client";

import { FormEvent, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

import { getSetupStatus, initializeSetup } from "@/lib/api/setup";
import { useI18n } from "@/lib/i18n/i18n-context";
import type { NumberingMode, YearSource } from "@/types/setup";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";

export default function SetupPage() {
  const router = useRouter();
  const { t } = useI18n();

  const [checking, setChecking] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  const [organizationName, setOrganizationName] = useState("Dabir");

  const [username, setUsername] = useState("admin");
  const [fullName, setFullName] = useState("System Administrator");
  const [password, setPassword] = useState("");

  const [numberingMode, setNumberingMode] =
    useState<NumberingMode>("fixed_prefix");

  const [numberPrefix, setNumberPrefix] = useState("DABIR");
  const [numberPadding, setNumberPadding] = useState(6);

  const [yearlyPrefixDigits, setYearlyPrefixDigits] = useState(3);
  const [yearlySerialPadding, setYearlySerialPadding] = useState(4);
  const [yearlySeparator, setYearlySeparator] = useState("-");
  const [yearSource, setYearSource] = useState<YearSource>("letter_date");

  const [error, setError] = useState("");

  useEffect(() => {
    async function checkSetup() {
      try {
        const status = await getSetupStatus();

        if (!status.setup_needed) {
          router.replace("/login");
          return;
        }
      } catch {
        setError("Could not check setup status. Please make sure API is running.");
      } finally {
        setChecking(false);
      }
    }

    void checkSetup();
  }, [router]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setError("");
    setSubmitting(true);

    try {
      await initializeSetup({
        organization_name: organizationName,
        superuser: {
          username,
          full_name: fullName,
          password,
        },
        letter_config: {
          numbering_mode: numberingMode,

          number_prefix: numberPrefix,
          number_padding: numberPadding,

          yearly_prefix_digits: yearlyPrefixDigits,
          yearly_serial_padding: yearlySerialPadding,
          yearly_separator: yearlySeparator,
          year_source: yearSource,
        },
      });

      toast.success("Application initialized successfully");
      router.push("/login");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Setup failed");
    } finally {
      setSubmitting(false);
    }
  }

  function fixedPrefixExample() {
    const prefix = numberPrefix || "DABIR";
    const padding = numberPadding || 6;
    return `${prefix}-${"1".padStart(padding, "0")}`;
  }

  function jalaliYearlyExample() {
    const yearSuffix = "405";
    const padding = yearlySerialPadding || 4;
    const separator = yearlySeparator || "-";
    return `${yearSuffix}${separator}${"1".padStart(padding, "0")}`;
  }

  if (checking) {
    return (
      <main className="min-h-screen flex items-center justify-center text-sm text-muted-foreground">
        {t.commonLoading}
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-muted/40 flex items-center justify-center p-4">
      <Card className="w-full max-w-3xl">
        <CardHeader>
          <CardTitle>{t.setupDabir}</CardTitle>
          <CardDescription>{t.setupDescription}</CardDescription>
        </CardHeader>

        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-8">
            {error ? (
              <Alert variant="destructive">
                <AlertTitle>{t.setup}</AlertTitle>
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            ) : null}

            <section className="space-y-4">
              <div>
                <h2 className="font-semibold">{t.organization}</h2>
                <p className="text-sm text-muted-foreground">
                  {t.organizationSettingsDescription}
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="organization_name">{t.organizationName}</Label>
                <Input
                  id="organization_name"
                  value={organizationName}
                  onChange={(event) => setOrganizationName(event.target.value)}
                  required
                />
              </div>
            </section>

            <section className="space-y-4">
              <div>
                <h2 className="font-semibold">{t.superuser}</h2>
                <p className="text-sm text-muted-foreground">
                  {t.setupDescription}
                </p>
              </div>

              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="username">{t.username}</Label>
                  <Input
                    id="username"
                    value={username}
                    onChange={(event) => setUsername(event.target.value)}
                    required
                    minLength={3}
                    dir="ltr"
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="full_name">{t.fullName}</Label>
                  <Input
                    id="full_name"
                    value={fullName}
                    onChange={(event) => setFullName(event.target.value)}
                    required
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="password">{t.password}</Label>
                <Input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                  required
                  minLength={8}
                  dir="ltr"
                />
              </div>
            </section>

            <section className="space-y-4">
              <div>
                <h2 className="font-semibold">{t.letterNumbering}</h2>
                <p className="text-sm text-muted-foreground">
                  {t.letterNumberingDescription}
                </p>
              </div>

              <div className="grid gap-4 md:grid-cols-2">
                <button
                  type="button"
                  onClick={() => setNumberingMode("fixed_prefix")}
                  className={[
                    "rounded-xl border p-4 text-start transition-colors",
                    numberingMode === "fixed_prefix"
                      ? "border-primary bg-primary/10"
                      : "hover:bg-muted",
                  ].join(" ")}
                >
                  <div className="flex items-center justify-between gap-2">
                    <div className="font-semibold">
                      {t.fixedPrefixNumbering}
                    </div>
                    {numberingMode === "fixed_prefix" ? (
                      <Badge>{t.commonActive}</Badge>
                    ) : null}
                  </div>
                  <p className="mt-2 text-sm text-muted-foreground">
                    {t.fixedPrefixNumberingDescription}
                  </p>
                  <div className="mt-3 font-mono text-lg" dir="ltr">
                    {fixedPrefixExample()}
                  </div>
                </button>

                <button
                  type="button"
                  onClick={() => setNumberingMode("jalali_yearly")}
                  className={[
                    "rounded-xl border p-4 text-start transition-colors",
                    numberingMode === "jalali_yearly"
                      ? "border-primary bg-primary/10"
                      : "hover:bg-muted",
                  ].join(" ")}
                >
                  <div className="flex items-center justify-between gap-2">
                    <div className="font-semibold">
                      {t.jalaliYearlyNumbering}
                    </div>
                    {numberingMode === "jalali_yearly" ? (
                      <Badge>{t.commonActive}</Badge>
                    ) : null}
                  </div>
                  <p className="mt-2 text-sm text-muted-foreground">
                    {t.jalaliYearlyNumberingDescription}
                  </p>
                  <div className="mt-3 font-mono text-lg" dir="ltr">
                    {jalaliYearlyExample()}
                  </div>
                </button>
              </div>

              {numberingMode === "fixed_prefix" ? (
                <div className="grid gap-4 md:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="number_prefix">{t.numberPrefix}</Label>
                    <Input
                      id="number_prefix"
                      value={numberPrefix}
                      onChange={(event) => setNumberPrefix(event.target.value)}
                      required
                      dir="ltr"
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="number_padding">{t.numberPadding}</Label>
                    <Input
                      id="number_padding"
                      type="number"
                      min={1}
                      max={12}
                      value={numberPadding}
                      onChange={(event) =>
                        setNumberPadding(Number(event.target.value))
                      }
                      required
                      dir="ltr"
                    />
                  </div>
                </div>
              ) : (
                <div className="grid gap-4 md:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="yearly_serial_padding">
                      {t.yearlySerialPadding}
                    </Label>
                    <Input
                      id="yearly_serial_padding"
                      type="number"
                      min={1}
                      max={12}
                      value={yearlySerialPadding}
                      onChange={(event) =>
                        setYearlySerialPadding(Number(event.target.value))
                      }
                      required
                      dir="ltr"
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="yearly_separator">
                      {t.yearlySeparator}
                    </Label>
                    <Input
                      id="yearly_separator"
                      value={yearlySeparator}
                      onChange={(event) =>
                        setYearlySeparator(event.target.value)
                      }
                      required
                      dir="ltr"
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="year_source">{t.yearSource}</Label>
                    <select
                      id="year_source"
                      value={yearSource}
                      onChange={(event) =>
                        setYearSource(event.target.value as YearSource)
                      }
                      className="w-full rounded-md border bg-background px-3 py-2 text-sm"
                    >
                      <option value="letter_date">
                        {t.yearSourceLetterDate}
                      </option>
                      <option value="created_at">{t.yearSourceCreatedAt}</option>
                    </select>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="yearly_prefix_digits">
                      {t.numberPrefix}
                    </Label>
                    <Input
                      id="yearly_prefix_digits"
                      type="number"
                      min={1}
                      max={4}
                      value={yearlyPrefixDigits}
                      onChange={(event) =>
                        setYearlyPrefixDigits(Number(event.target.value))
                      }
                      required
                      dir="ltr"
                    />
                  </div>
                </div>
              )}

              <div className="rounded-md border bg-muted/30 p-3 text-sm">
                {t.example}:{" "}
                <span className="font-medium" dir="ltr">
                  {numberingMode === "fixed_prefix"
                    ? fixedPrefixExample()
                    : jalaliYearlyExample()}
                </span>
              </div>
            </section>

            <Button type="submit" className="w-full" disabled={submitting}>
              {submitting ? t.initializing : t.initializeApplication}
            </Button>
          </form>
        </CardContent>
      </Card>
    </main>
  );
}