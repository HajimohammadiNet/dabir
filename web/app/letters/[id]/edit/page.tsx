"use client";

import { FormEvent, useCallback, useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { LetterNumberText } from "@/components/common/letter-number-text";
import { useAuth } from "@/contexts/auth-context";
import { getLetter, updateLetter } from "@/lib/api/letters";
import { useI18n } from "@/lib/i18n/i18n-context";
import type { Letter } from "@/types/letter";

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
import { JalaliDatePicker } from "@/components/common/jalali-date-picker";

export default function EditLetterPage() {
  const router = useRouter();
  const params = useParams<{ id: string }>();
  const { token } = useAuth();
  const { t } = useI18n();

  const letterID = params.id;

  const [letter, setLetter] = useState<Letter | null>(null);

  const [displayLetterNumber, setDisplayLetterNumber] = useState("");
  const [title, setTitle] = useState("");
  const [letterDate, setLetterDate] = useState("");
  const [sender, setSender] = useState("");
  const [receiver, setReceiver] = useState("");
  const [description, setDescription] = useState("");

  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);

  const loadLetter = useCallback(async () => {
    if (!token || !letterID) return;

    setLoading(true);

    try {
      const result = await getLetter(token, letterID);

      setLetter(result);
      setDisplayLetterNumber(result.formatted_letter_number || "");
      setTitle(result.title || "");
      setLetterDate(result.letter_date_jalali || "");
      setSender(result.sender || "");
      setReceiver(result.receiver || "");
      setDescription(result.description || "");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to load letter");
    } finally {
      setLoading(false);
    }
  }, [token, letterID]);

  useEffect(() => {
    const timeoutID = window.setTimeout(() => {
      void loadLetter();
    }, 0);

    return () => {
      window.clearTimeout(timeoutID);
    };
  }, [loadLetter]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!token || !letterID) return;

    setSaving(true);

    try {
      const updated = await updateLetter(token, letterID, {
        display_letter_number: displayLetterNumber.trim(),
        title,
        letter_date: letterDate,
        sender,
        receiver,
        description: description.trim() || null,
      });

      toast.success("نامه با موفقیت ویرایش شد");
      setLetter(updated);

      router.push("/letters");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update letter");
    } finally {
      setSaving(false);
    }
  }

  return (
    <ProtectedRoute allowedRoles={["superuser", "editor"]}>
      <AppShell>
        <div className="max-w-2xl space-y-6">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              ویرایش نامه
            </h1>
            <p className="text-muted-foreground">
              اطلاعات نامه را ویرایش و ذخیره کنید.
            </p>
          </div>

          {letter ? (
            <Card>
              <CardHeader>
                <CardTitle>شماره فعلی نامه</CardTitle>
              </CardHeader>
              <CardContent>
                <LetterNumberText
                  value={letter.formatted_letter_number}
                  className="text-3xl font-bold"
                />
              </CardContent>
            </Card>
          ) : null}

          <Card>
            <CardHeader>
              <CardTitle>{t.letterInformation}</CardTitle>
            </CardHeader>

            <CardContent>
              {loading ? (
                <div className="text-sm text-muted-foreground">
                  {t.commonLoading}
                </div>
              ) : (
                <form onSubmit={handleSubmit} className="space-y-4">
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
                      className="text-left"
                      style={{
                        direction: "ltr",
                        unicodeBidi: "plaintext",
                      }}
                    />
                  </div>

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

                  <div className="flex flex-wrap gap-2">
                    <Button type="submit" disabled={saving}>
                      {saving ? t.commonLoading : "ذخیره تغییرات"}
                    </Button>

                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => router.push("/letters")}
                      disabled={saving}
                    >
                      {t.commonCancel}
                    </Button>
                  </div>
                </form>
              )}
            </CardContent>
          </Card>
        </div>
      </AppShell>
    </ProtectedRoute>
  );
}