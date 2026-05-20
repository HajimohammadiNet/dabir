"use client";

import { FormEvent, useState } from "react";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { useAuth } from "@/contexts/auth-context";
import { changePassword } from "@/lib/api/auth";
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

export default function ChangePasswordPage() {
  const { token, logout } = useAuth();
  const { t } = useI18n();

  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!token) return;

    setError("");

    if (newPassword.length < 8) {
      setError("New password must be at least 8 characters.");
      return;
    }

    if (newPassword !== confirmPassword) {
      setError("New password and confirmation do not match.");
      return;
    }

    setLoading(true);

    try {
      await changePassword(token, {
        current_password: currentPassword,
        new_password: newPassword,
      });

      toast.success("Password changed successfully. Please login again.");

      logout();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to change password");
    } finally {
      setLoading(false);
    }
  }

  return (
    <ProtectedRoute>
      <AppShell>
        <div className="max-w-xl space-y-6">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              {t.changePassword}
            </h1>
            <p className="text-muted-foreground">
              Update your account password.
            </p>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>{t.changePassword}</CardTitle>
              <CardDescription>
                After changing your password, you will be logged out.
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
                  <Label htmlFor="current_password">
                    {t.currentPassword}
                  </Label>
                  <Input
                    id="current_password"
                    type="password"
                    value={currentPassword}
                    onChange={(event) =>
                      setCurrentPassword(event.target.value)
                    }
                    required
                  />
                </div>

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

                <Button type="submit" disabled={loading}>
                  {loading ? "Changing..." : t.changePassword}
                </Button>
              </form>
            </CardContent>
          </Card>
        </div>
      </AppShell>
    </ProtectedRoute>
  );
}