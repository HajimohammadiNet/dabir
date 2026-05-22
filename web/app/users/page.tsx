"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { useAuth } from "@/contexts/auth-context";
import { activateUser, deactivateUser, listUsers } from "@/lib/api/users";
import { useI18n } from "@/lib/i18n/i18n-context";
import type { User } from "@/types/user";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

export default function UsersPage() {
  const { token, user: currentUser } = useAuth();
  const { t } = useI18n();

  const [users, setUsers] = useState<User[]>([]);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(false);
  const [actionLoadingUserID, setActionLoadingUserID] = useState<string | null>(
    null
  );

  const loadUsers = useCallback(async () => {
    if (!token) return;

    setLoading(true);

    try {
      const result = await listUsers(token, {
        page: 1,
        page_size: 50,
        search,
      });

      setUsers(result.items);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to load users");
    } finally {
      setLoading(false);
    }
  }, [token, search]);

  useEffect(() => {
    const load = async () => {
      await loadUsers();
    };

    void load();
  }, [loadUsers]);

  async function handleToggleActive(targetUser: User) {
    if (!token) return;

    if (targetUser.id === currentUser?.id) {
      toast.error("You cannot deactivate your own user from here");
      return;
    }

    setActionLoadingUserID(targetUser.id);

    try {
      if (targetUser.is_active) {
        await deactivateUser(token, targetUser.id);
        toast.success("User deactivated");
      } else {
        await activateUser(token, targetUser.id);
        toast.success("User activated");
      }

      await loadUsers();
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to update user status"
      );
    } finally {
      setActionLoadingUserID(null);
    }
  }

  return (
    <ProtectedRoute allowedRoles={["superuser"]}>
      <AppShell>
        <div className="space-y-6">
          <div className="flex items-start justify-between gap-4">
            <div>
              <h1 className="text-3xl font-bold tracking-tight">{t.users}</h1>
              <p className="text-muted-foreground">{t.usersDescription}</p>
            </div>

            <Link href="/users/new">
              <Button>{t.newUser}</Button>
            </Link>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>{t.commonSearch}</CardTitle>
            </CardHeader>
            <CardContent>
              <form
                className="flex gap-2"
                onSubmit={(event) => {
                  event.preventDefault();
                  void loadUsers();
                }}
              >
                <Input
                  placeholder={t.searchUsersPlaceholder}
                  value={search}
                  onChange={(event) => setSearch(event.target.value)}
                />

                <Button type="submit" disabled={loading}>
                  {loading ? t.commonLoading : t.commonSearch}
                </Button>
              </form>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>{t.users}</CardTitle>
            </CardHeader>

            <CardContent>
              <div className="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>{t.username}</TableHead>
                      <TableHead>{t.fullName}</TableHead>
                      <TableHead>{t.role}</TableHead>
                      <TableHead>{t.commonStatus}</TableHead>
                      <TableHead>{t.createdAt}</TableHead>
                      <TableHead className="text-end">
                        {t.commonActions}
                      </TableHead>
                    </TableRow>
                  </TableHeader>

                  <TableBody>
                    {users.length === 0 ? (
                      <TableRow>
                        <TableCell
                          colSpan={6}
                          className="text-center text-muted-foreground"
                        >
                          {loading ? t.commonLoading : t.noUsersFound}
                        </TableCell>
                      </TableRow>
                    ) : (
                      users.map((user) => (
                        <TableRow key={user.id}>
                          <TableCell className="font-medium">
                            {user.username}
                          </TableCell>

                          <TableCell>{user.full_name}</TableCell>

                          <TableCell>
                            <Badge variant="outline">{user.role}</Badge>
                          </TableCell>

                          <TableCell>
                            {user.is_active ? (
                              <Badge variant="secondary">
                                {t.commonActive}
                              </Badge>
                            ) : (
                              <Badge variant="destructive">
                                {t.commonInactive}
                              </Badge>
                            )}
                          </TableCell>

                          <TableCell>
                            {new Date(user.created_at).toLocaleString()}
                          </TableCell>

                          <TableCell className="text-end">
                            <div className="flex justify-end gap-2">
                              <Link href={`/users/${user.id}/reset-password`}>
                                <Button variant="outline" size="sm">
                                  {t.resetPassword}
                                </Button>
                              </Link>

                              <Button
                                variant="outline"
                                size="sm"
                                disabled={actionLoadingUserID === user.id}
                                onClick={() => void handleToggleActive(user)}
                              >
                                {user.is_active ? t.deactivate : t.activate}
                              </Button>
                            </div>
                          </TableCell>
                        </TableRow>
                      ))
                    )}
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>
        </div>
      </AppShell>
    </ProtectedRoute>
  );
}