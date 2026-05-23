"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

import { useAuth } from "@/contexts/auth-context";
import { Button } from "@/components/ui/button";
import { ThemeToggle } from "@/components/theme/theme-toggle";
import { LanguageToggle } from "@/components/theme/language-toggle";
import { useI18n } from "@/lib/i18n/i18n-context";
import { dictionaries } from "@/lib/i18n/dictionaries";
import type { Role } from "@/types/auth";

type NavItem = {
  href: string;
  labelKey: keyof typeof dictionaries.en;
  roles: Role[];
};

const navItems: NavItem[] = [
  {
    href: "/dashboard",
    labelKey: "dashboard",
    roles: ["superuser", "editor", "readonly"],
  },
  {
    href: "/profile/change-password",
    labelKey: "changePassword",
    roles: ["superuser", "editor", "readonly"],
  },
  {
    href: "/letters",
    labelKey: "letters",
    roles: ["superuser", "editor", "readonly"],
  },
  {
    href: "/letters/import",
    labelKey: "import",
    roles: ["superuser"],
  },
  {
    href: "/users",
    labelKey: "users",
    roles: ["superuser"],
  },
  {
    href: "/audit-logs",
    labelKey: "auditLogs",
    roles: ["superuser"],
  },
  {
    href: "/settings",
    labelKey: "settings",
    roles: ["superuser"],
  },
];

export function AppShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const { user, logout } = useAuth();
  const { t, direction } = useI18n();

  const visibleNavItems = navItems.filter((item) =>
    user ? item.roles.includes(user.role) : false
  );

  return (
    <div
      className="min-h-screen overflow-x-hidden bg-muted/40"
      dir={direction}
    >
      <header className="sticky top-0 z-40 border-b bg-background/95 backdrop-blur">
        <div className="flex h-16 min-w-0 items-center justify-between gap-4 px-4 md:px-6">
          <Link
            href="/dashboard"
            className="shrink-0 text-lg font-semibold"
          >
            {t.appName}
          </Link>

          <div className="flex min-w-0 items-center gap-2 md:gap-3">
            {user ? (
              <div className="hidden min-w-0 text-end sm:block">
                <div className="truncate text-sm font-medium">
                  {user.full_name}
                </div>
                <div className="truncate text-xs text-muted-foreground">
                  {user.username} · {user.role}
                </div>
              </div>
            ) : null}

            <LanguageToggle />
            <ThemeToggle />

            <Button variant="outline" size="sm" onClick={logout}>
              {t.logout}
            </Button>
          </div>
        </div>
      </header>

      <div className="flex min-w-0">
        <aside className="hidden min-h-[calc(100vh-4rem)] w-64 shrink-0 border-e bg-background p-4 md:block">
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
                      ? "bg-muted font-medium text-foreground"
                      : "text-muted-foreground hover:bg-muted hover:text-foreground",
                  ].join(" ")}
                >
                  {t[item.labelKey]}
                </Link>
              );
            })}
          </nav>
        </aside>

        <main className="min-w-0 flex-1 overflow-x-hidden p-4 md:p-6">
          <div className="mx-auto min-w-0 max-w-full">{children}</div>
        </main>
      </div>

      <nav className="fixed inset-x-0 bottom-0 z-40 border-t bg-background md:hidden">
        <div className="flex overflow-x-auto px-2 py-2">
          {visibleNavItems.map((item) => {
            const active =
              pathname === item.href ||
              (item.href !== "/dashboard" && pathname.startsWith(item.href));

            return (
              <Link
                key={item.href}
                href={item.href}
                className={[
                  "shrink-0 rounded-md px-3 py-2 text-sm transition-colors",
                  active
                    ? "bg-muted font-medium text-foreground"
                    : "text-muted-foreground hover:bg-muted hover:text-foreground",
                ].join(" ")}
              >
                {t[item.labelKey]}
              </Link>
            );
          })}
        </div>
      </nav>
    </div>
  );
}