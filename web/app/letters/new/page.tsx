"use client";

import { FormEvent, useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { useAuth } from "@/contexts/auth-context";
import {
  createLetter,
  getLetterNumberSuggestion,
  listLetters,
  type LetterNumberSuggestion,
} from "@/lib/api/letters";

import { getPublicSettings } from "@/lib/api/settings";
import { useI18n } from "@/lib/i18n/i18n-context";
import type { Letter } from "@/types/letter";
import type { NumberingMode } from "@/types/settings";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { JalaliDatePicker } from "@/components/common/jalali-date-picker";
import { LetterNumberText } from "@/components/common/letter-number-text";

export default function NewLetterPage() {
  const router = useRouter();
  const { token } = useAuth();
  const { t } = useI18n();

  const [numberingMode, setNumberingMode] =
    useState<NumberingMode>("fixed_prefix");
  const [lastLetterNumber, setLastLetterNumber] = useState<string | null>(null);
  const [patternSuggestion, setPatternSuggestion] =
  useState<LetterNumberSuggestion | null>(null);
  const [suggestionLoading, setSuggestionLoading] = useState(false);

  const [displayLetterNumber, setDisplayLetterNumber] = useState("");
  const [title, setTitle] = useState("");
  const [letterDate, setLetterDate] = useState("");
  const [sender, setSender] = useState("");
  const [receiver, setReceiver] = useState("");
  const [description, setDescription] = useState("");
  const [loading, setLoading] = useState(false);

  const [createdLetter, setCreatedLetter] = useState<Letter | null>(null);
  const [resultDialogOpen, setResultDialogOpen] = useState(false);
  const suggestedLetterNumber =
    patternSuggestion?.suggested_number ||
    suggestNextLetterNumber(lastLetterNumber);

  const effectiveLastNumber =
    patternSuggestion?.last_number || lastLetterNumber;

  const suggestedLetterNumberPlaceholder =
    formatLetterNumberPlaceholder(suggestedLetterNumber);

  const loadPageData = useCallback(async () => {
    try {
      const settings = await getPublicSettings();
      setNumberingMode(settings.letter_config.numbering_mode);

      if (token) {
        const lettersResult = await listLetters(token, {
            page: 1,
            page_size: 1,
            sort_by: "created_at",
            sort_order: "desc",
        });

        if (lettersResult.items.length > 0) {
          setLastLetterNumber(lettersResult.items[0].formatted_letter_number);
        } else {
          setLastLetterNumber(null);
        }
      }
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to load letter settings"
      );
    }
  }, [token]);

  useEffect(() => {
    const timeoutID = window.setTimeout(() => {
      void loadPageData();
    }, 0);

    return () => {
      window.clearTimeout(timeoutID);
    };
  }, [loadPageData]);

    useEffect(() => {
    const timeoutID = window.setTimeout(async () => {
        if (!token || numberingMode !== "manual") {
        setPatternSuggestion(null);
        setSuggestionLoading(false);
        return;
        }

        const prefix = extractSuggestionPrefix(displayLetterNumber);

        setSuggestionLoading(true);

        try {
        const result = await getLetterNumberSuggestion(token, prefix);
        setPatternSuggestion(result);
        } catch {
        setPatternSuggestion(null);
        } finally {
        setSuggestionLoading(false);
        }
    }, 400);

    return () => {
        window.clearTimeout(timeoutID);
    };
    }, [token, numberingMode, displayLetterNumber]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!token) return;

    if (numberingMode === "manual" && displayLetterNumber.trim() === "") {
      toast.error("Letter number is required in manual numbering mode");
      return;
    }

    setLoading(true);

    try {
      const letter = await createLetter(token, {
        display_letter_number:
          numberingMode === "manual" ? displayLetterNumber.trim() : null,
        title,
        letter_date: letterDate,
        sender,
        receiver,
        description: description || null,
      });

      setCreatedLetter(letter);
      setResultDialogOpen(true);
      setLastLetterNumber(letter.formatted_letter_number);

      toast.success(`${letter.formatted_letter_number} ثبت شد`);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create letter");
    } finally {
      setLoading(false);
    }
  }

  function resetForm() {
    setDisplayLetterNumber("");
    setTitle("");
    setLetterDate("");
    setSender("");
    setReceiver("");
    setDescription("");
    setCreatedLetter(null);
    setResultDialogOpen(false);
  }

  return (
    <ProtectedRoute allowedRoles={["superuser", "editor"]}>
      <AppShell>
        <div className="max-w-2xl space-y-6">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              {t.newLetter}
            </h1>
            <p className="text-muted-foreground">{t.lettersDescription}</p>
          </div>

          {numberingMode === "manual" ? (
            <Card>
              <CardHeader>
                <CardTitle>{t.manualNumbering}</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="rounded-lg border bg-muted/30 p-4">
                  <div className="text-sm text-muted-foreground">
                    {t.lastLetterNumber}
                  </div>
                  <div className="mt-1 text-2xl font-bold" dir="ltr">
                    {lastLetterNumber ? (
                        <LetterNumberText value={lastLetterNumber} />
                    ) : (
                        "-"
                    )}
                  </div>
                </div>

                <div className="space-y-2">
                    <Label htmlFor="display_letter_number">
                        {t.displayLetterNumber}
                    </Label>

                    <Input
                        id="display_letter_number"
                        value={displayLetterNumber}
                        onChange={(event) =>
                            setDisplayLetterNumber(event.target.value)
                        }
                        required
                        dir="ltr"
                        placeholder={suggestedLetterNumberPlaceholder}
                        className="text-left"
                        style={{
                            direction: "ltr",
                            unicodeBidi: "plaintext",
                        }}
                    />

                    <div className="space-y-2 rounded-md border bg-muted/30 p-3 text-xs text-muted-foreground">
                        <div className="flex flex-wrap items-center justify-between gap-2">
                            <span>
                            {displayLetterNumber.trim()
                                ? "آخرین شماره مشابه:"
                                : "آخرین شماره کلی:"}
                            </span>

                            <span className="font-medium">
                            {suggestionLoading ? (
                                "..."
                            ) : effectiveLastNumber ? (
                                <LetterNumberText value={effectiveLastNumber} />
                            ) : (
                                "-"
                            )}
                            </span>
                        </div>

                        <div className="flex flex-wrap items-center justify-between gap-2">
                            <span>شماره پیشنهادی بعدی:</span>

                            <span className="font-medium">
                            {suggestionLoading ? (
                                "..."
                            ) : suggestedLetterNumber ? (
                                <LetterNumberText value={suggestedLetterNumber} />
                            ) : (
                                "-"
                            )}
                            </span>
                        </div>

                        {suggestedLetterNumber ? (
                            <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={() => setDisplayLetterNumber(suggestedLetterNumber)}
                            >
                            استفاده از پیشنهاد
                            </Button>
                        ) : null}
                        </div>
                </div>
              </CardContent>
            </Card>
          ) : null}

          <Card>
            <CardHeader>
              <CardTitle>{t.letterInformation}</CardTitle>
            </CardHeader>

            <CardContent>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="title">{t.letterTitle}</Label>
                  <Input
                    id="title"
                    value={title}
                    onChange={(event) => setTitle(event.target.value)}
                    required
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="letter_date">{t.letterDate}</Label>
                  <JalaliDatePicker
                    id="letter_date"
                    value={letterDate}
                    onChange={setLetterDate}
                    required
                  />
                </div>

                <div className="grid gap-4 md:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="sender">{t.sender}</Label>
                    <Input
                      id="sender"
                      value={sender}
                      onChange={(event) => setSender(event.target.value)}
                      required
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="receiver">{t.receiver}</Label>
                    <Input
                      id="receiver"
                      value={receiver}
                      onChange={(event) => setReceiver(event.target.value)}
                      required
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="description">{t.description}</Label>
                  <Textarea
                    id="description"
                    value={description}
                    onChange={(event) => setDescription(event.target.value)}
                  />
                </div>

                <div className="flex gap-2">
                  <Button type="submit" disabled={loading}>
                    {loading ? t.creatingLetter : t.createLetter}
                  </Button>

                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => router.push("/letters")}
                  >
                    {t.commonCancel}
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>

          <Dialog
            open={resultDialogOpen}
            onOpenChange={(open) => {
              setResultDialogOpen(open);
            }}
          >
            <DialogContent className="max-w-2xl">
              <DialogHeader>
                <DialogTitle>نامه با موفقیت ثبت شد</DialogTitle>
                <DialogDescription>
                  شماره نامه و اطلاعات ثبت‌شده در ادامه نمایش داده شده است.
                </DialogDescription>
              </DialogHeader>

              {createdLetter ? (
                <div className="space-y-6">
                  <div className="rounded-xl border bg-muted/30 p-6 text-center">
                    <div className="text-sm text-muted-foreground">
                      شماره نامه
                    </div>
                    <div
                      className="mt-2 text-4xl font-bold tracking-wide"
                      dir="ltr"
                    >
                      <LetterNumberText value={createdLetter.formatted_letter_number} />
                    </div>
                  </div>

                  <div className="grid gap-3 text-sm">
                    <InfoRow label={t.letterTitle} value={createdLetter.title} />
                    <InfoRow
                      label={t.letterDate}
                      value={createdLetter.letter_date_jalali}
                      forceLtr
                    />
                    <InfoRow label={t.sender} value={createdLetter.sender} />
                    <InfoRow label={t.receiver} value={createdLetter.receiver} />
                    <InfoRow
                      label={t.registrar}
                      value={createdLetter.registrar_name}
                    />
                    <InfoRow
                      label={t.commonStatus}
                      value={
                        createdLetter.is_deleted
                          ? t.commonDeleted
                          : t.commonActive
                      }
                      badge
                    />
                    {createdLetter.description ? (
                      <InfoRow
                        label={t.description}
                        value={createdLetter.description}
                      />
                    ) : null}
                  </div>

                  <div className="flex flex-wrap gap-2">
                    <Button onClick={() => router.push("/letters")}>
                      رفتن به لیست نامه‌ها
                    </Button>

                    <Button variant="outline" onClick={resetForm}>
                      ثبت نامه جدید
                    </Button>
                  </div>
                </div>
              ) : null}
            </DialogContent>
          </Dialog>
        </div>
      </AppShell>
    </ProtectedRoute>
  );
}

function InfoRow({
  label,
  value,
  forceLtr,
  badge,
}: {
  label: string;
  value: string;
  forceLtr?: boolean;
  badge?: boolean;
}) {
  return (
    <div className="flex items-start justify-between gap-4 border-b pb-2 last:border-b-0">
      <div className="text-muted-foreground">{label}</div>

      {badge ? (
        <Badge variant="secondary">{value}</Badge>
      ) : (
        <div className="font-medium text-end" dir={forceLtr ? "ltr" : "auto"}>
          {value}
        </div>
      )}
    </div>
  );
}

function extractSuggestionPrefix(value: string) {
  const trimmed = normalizeDigitsForSuggestion(value.trim());

  if (!trimmed) {
    return "";
  }

  // اگر کاربر خودش prefix وارد کرده باشد، مثل:
  // 405-ق-
  // ۴۰۵-ق-
  // HR-2026-
  if (/[-_/\\.\s]$/.test(trimmed)) {
    return trimmed;
  }

  // اگر آخر ورودی عدد دارد، عدد انتهایی را حذف کن:
  // 405-ق-001 -> 405-ق-
  // ۴۰۵-ق-۰۰۱ -> 405-ق-
  // HR-2026-0042 -> HR-2026-
  return trimmed.replace(/[0-9]+$/, "");
}

function suggestNextLetterNumber(lastNumber: string | null) {
  if (!lastNumber) {
    return "405-158";
  }

  const value = lastNumber.trim();
  if (!value) {
    return "405-158";
  }

  const match = value.match(/(\d+)(?!.*\d)/);
  if (!match || match.index === undefined) {
    return value;
  }

  const numericPart = match[1];
  const nextNumber = String(Number(numericPart) + 1).padStart(
    numericPart.length,
    "0"
  );

  return (
    value.slice(0, match.index) +
    nextNumber +
    value.slice(match.index + numericPart.length)
  );
}

function formatLetterNumberPlaceholder(value: string) {
  if (!value) {
    return "";
  }

  // LRI + PDI forces the placeholder to keep mixed Persian/Latin order.
  return `\u2066${value}\u2069`;
}

function normalizeDigitsForSuggestion(value: string) {
  return value
    .replace(/[۰-۹]/g, (digit) => String("۰۱۲۳۴۵۶۷۸۹".indexOf(digit)))
    .replace(/[٠-٩]/g, (digit) => String("٠١٢٣٤٥٦٧٨٩".indexOf(digit)));
}