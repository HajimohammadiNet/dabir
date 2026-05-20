"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { useAuth } from "@/contexts/auth-context";
import { createUser } from "@/lib/api/users";
import type { Role } from "@/types/auth";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export default function NewUserPage() {
  const router = useRouter();
  const { token } = useAuth();

  const [username, setUsername] = useState("");
  const [fullName, setFullName] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState<Role>("readonly");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!token) return;

    setLoading(true);

    try {
      await createUser(token, {
        username,
        full_name: fullName,
        password,
        role,
      });

      toast.success("User created");
      router.push("/users");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create user");
    } finally {
      setLoading(false);
    }
  }

  return (
    <ProtectedRoute allowedRoles={["superuser"]}>
      <AppShell>
        <div className="max-w-2xl space-y-6">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">New User</h1>
            <p className="text-muted-foreground">
              Create a new user and assign an access role.
            </p>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>User Information</CardTitle>
            </CardHeader>

            <CardContent>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="username">Username</Label>
                  <Input
                    id="username"
                    value={username}
                    onChange={(event) => setUsername(event.target.value)}
                    required
                    minLength={3}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="full_name">Full Name</Label>
                  <Input
                    id="full_name"
                    value={fullName}
                    onChange={(event) => setFullName(event.target.value)}
                    required
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="password">Password</Label>
                  <Input
                    id="password"
                    type="password"
                    value={password}
                    onChange={(event) => setPassword(event.target.value)}
                    required
                    minLength={8}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="role">Role</Label>
                  <select
                    id="role"
                    value={role}
                    onChange={(event) => setRole(event.target.value as Role)}
                    className="w-full rounded-md border bg-background px-3 py-2 text-sm"
                  >
                    <option value="readonly">Readonly</option>
                    <option value="editor">Editor</option>
                    <option value="superuser">Superuser</option>
                  </select>
                </div>

                <div className="flex gap-2">
                  <Button type="submit" disabled={loading}>
                    {loading ? "Creating..." : "Create User"}
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