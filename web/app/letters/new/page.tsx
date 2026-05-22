"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { useAuth } from "@/contexts/auth-context";
import { createLetter } from "@/lib/api/letters";
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
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";

export default function NewLetterPage() {
  const router = useRouter();
  const { token } = useAuth();
  const { t } = useI18n();

  const [title, setTitle] = useState("");
  const [letterDate, setLetterDate] = useState("");
  const [sender, setSender] = useState("");
  const [receiver, setReceiver] = useState("");
  const [description, setDescription] = useState("");
  const [loading, setLoading] = useState(false);

  const [createdLetter, setCreatedLetter] = useState<Letter | null>(null);
  const [resultDialogOpen, setResultDialogOpen] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!token) return;

    setLoading(true);

    try {
      const letter = await createLetter(token, {
        title,
        letter_date: letterDate,
        sender,
        receiver,
        description: description || null,
      });

      setCreatedLetter(letter);
      setResultDialogOpen(true);

      toast.success(`${letter.formatted_letter_number} ثبت شد`);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create letter");
    } finally {
      setLoading(false);
    }
  }

  function resetForm() {
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
                  <Input
                    id="letter_date"
                    type="date"
                    value={letterDate}
                    onChange={(event) => setLetterDate(event.target.value)}
                    required
                    dir="ltr"
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
                  شماره نامه تولیدشده و اطلاعات ثبت‌شده در ادامه نمایش داده شده است.
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
                      {createdLetter.formatted_letter_number}
                    </div>
                  </div>

                  <div className="grid gap-3 text-sm">
                    <InfoRow label={t.letterTitle} value={createdLetter.title} />
                    <InfoRow
                      label={t.letterDate}
                      value={createdLetter.letter_date}
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