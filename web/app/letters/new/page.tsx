"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { useAuth } from "@/contexts/auth-context";
import { createLetter } from "@/lib/api/letters";
import { useI18n } from "@/lib/i18n/i18n-context";

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

      toast.success(`${letter.formatted_letter_number} ${t.createLetter}`);
      router.push("/letters");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create letter");
    } finally {
      setLoading(false);
    }
  }

  return (
    <ProtectedRoute allowedRoles={["superuser", "editor"]}>
      <AppShell>
        <div className="max-w-2xl space-y-6">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              {t.newLetter}
            </h1>
            <p className="text-muted-foreground">
              {t.lettersDescription}
            </p>
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
        </div>
      </AppShell>
    </ProtectedRoute>
  );
}