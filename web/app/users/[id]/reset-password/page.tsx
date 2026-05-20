"use client";

import { FormEvent, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { useAuth } from "@/contexts/auth-context";
import { resetUserPassword } from "@/lib/api/users";
import { useI18n } from "@/lib/i18n/i18n-context";

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
import { Alert, AlertDescription } from "@/components/ui/alert";

export default function ResetPasswordPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const { token, user } = useAuth();
  const { t } = useI18n();

  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const userID = params.id;

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!token) return;

    setError("");

    if (newPassword.length < 8) {
      setError("Password must be at least 8 characters.");
      return;
    }

    if (newPassword !== confirmPassword) {
      setError("Password and confirmation do not match.");
      return;
    }

    setLoading(true);

    try {
      await resetUserPassword(token, userID, newPassword);

      toast.success("Password reset successfully");
      router.push("/users");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to reset password");
    } finally {
      setLoading(false);
    }
  }

  return (
    <ProtectedRoute allowedRoles={["superuser"]}>
      <AppShell>
        <div className="max-w-xl space-y-6">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              {t.resetPassword}
            </h1>
            <p className="text-muted-foreground">
              Set a new password for this user.
            </p>
          </div>

          {user?.id === userID ? (
            <Alert>
              <AlertDescription>
                You are resetting your own password. Prefer using Change
                Password from your profile.
              </AlertDescription>
            </Alert>
          ) : null}

          <Card>
            <CardHeader>
              <CardTitle>{t.resetPassword}</CardTitle>
              <CardDescription>
                The target user can login with this new password immediately.
              </CardDescription>
            </CardHeader>

            <CardContent>
              <form onSubmit={handleSubmit} className="space-y-4">
                {error ? (
                  <Alert variant="destructive">
                    <AlertDescription>{error}</AlertDescription>
                  </Alert>
                ) : null}

                <div className="space-y-2">
                  <Label htmlFor="new_password">{t.newPassword}</Label>
                  <Input
                    id="new_password"
                    type="password"
                    value={newPassword}
                    onChange={(event) => setNewPassword(event.target.value)}
                    required
                    minLength={8}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="confirm_password">
                    {t.confirmPassword}
                  </Label>
                  <Input
                    id="confirm_password"
                    type="password"
                    value={confirmPassword}
                    onChange={(event) =>
                      setConfirmPassword(event.target.value)
                    }
                    required
                    minLength={8}
                  />
                </div>

                <div className="flex gap-2">
                  <Button type="submit" disabled={loading}>
                    {loading ? "Resetting..." : t.resetPassword}
                  </Button>

                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => router.push("/users")}
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