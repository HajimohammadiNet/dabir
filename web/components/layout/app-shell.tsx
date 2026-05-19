"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

import { useAuth } from "@/contexts/auth-context";
import { Button } from "@/components/ui/button";
import type { Role } from "@/types/auth";

type NavItem = {
  href: string;
  label: string;
  roles: Role[];
};

const navItems: NavItem[] = [
  {
    href: "/dashboard",
    label: "Dashboard",
    roles: ["superuser", "editor", "readonly"],
  },
  {
    href: "/letters",
    label: "Letters",
    roles: ["superuser", "editor", "readonly"],
  },
  {
    href: "/letters/import",
    label: "Import",
    roles: ["superuser"],
  },
  {
    href: "/users",
    label: "Users",
    roles: ["superuser"],
  },
  {
    href: "/audit-logs",
    label: "Audit Logs",
    roles: ["superuser"],
  },
  {
    href: "/settings",
    label: "Settings",
    roles: ["superuser"],
  },
];

export function AppShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const { user, logout } = useAuth();

  const visibleNavItems = navItems.filter((item) =>
    user ? item.roles.includes(user.role) : false
  );

  return (
    <div className="min-h-screen bg-muted/40">
      <header className="border-b bg-background">
        <div className="h-16 px-6 flex items-center justify-between">
          <Link href="/dashboard" className="font-semibold text-lg">
            Dabir
          </Link>

          <div className="flex items-center gap-4">
            {user ? (
              <div className="text-right">
                <div className="text-sm font-medium">{user.full_name}</div>
                <div className="text-xs text-muted-foreground">
                  {user.username} · {user.role}
                </div>
              </div>
            ) : null}

            <Button variant="outline" onClick={logout}>
              Logout
            </Button>
          </div>
        </div>
      </header>

      <div className="flex">
        <aside className="w-64 border-r min-h-[calc(100vh-4rem)] bg-background p-4">
          <nav className="space-y-1">
            {visibleNavItems.map((item) => {
              const active =
                pathname === item.href ||
                (item.href !== "/dashboard" && pathname.startsWith(item.href));

              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={[
                    "block rounded-md px-3 py-2 text-sm transition-colors",
                    active
                      ? "bg-muted font-medium"
                      : "hover:bg-muted text-muted-foreground hover:text-foreground",
                  ].join(" ")}
                >
                  {item.label}
                </Link>
              );
            })}
          </nav>
        </aside>

        <main className="flex-1 p-6">{children}</main>
      </div>
    </div>
  );
}