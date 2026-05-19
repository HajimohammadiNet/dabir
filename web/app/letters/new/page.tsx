"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { useAuth } from "@/contexts/auth-context";
import { createLetter } from "@/lib/api/letters";

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

      toast.success(`Letter ${letter.formatted_letter_number} created`);
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
            <h1 className="text-3xl font-bold tracking-tight">New Letter</h1>
            <p className="text-muted-foreground">
              Register a new letter and generate a new number.
            </p>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Letter Information</CardTitle>
            </CardHeader>

            <CardContent>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="title">Title</Label>
                  <Input
                    id="title"
                    value={title}
                    onChange={(event) => setTitle(event.target.value)}
                    required
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="letter_date">Letter Date</Label>
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
                    <Label htmlFor="sender">Sender</Label>
                    <Input
                      id="sender"
                      value={sender}
                      onChange={(event) => setSender(event.target.value)}
                      required
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="receiver">Receiver</Label>
                    <Input
                      id="receiver"
                      value={receiver}
                      onChange={(event) => setReceiver(event.target.value)}
                      required
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="description">Description</Label>
                  <Textarea
                    id="description"
                    value={description}
                    onChange={(event) => setDescription(event.target.value)}
                  />
                </div>

                <div className="flex gap-2">
                  <Button type="submit" disabled={loading}>
                    {loading ? "Creating..." : "Create Letter"}
                  </Button>

                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => router.push("/letters")}
                  >
                    Cancel
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