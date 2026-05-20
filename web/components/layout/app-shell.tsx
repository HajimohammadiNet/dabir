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
    <div className="min-h-screen bg-muted/40" dir={direction}>
      <header className="border-b bg-background">
        <div className="h-16 px-6 flex items-center justify-between gap-4">
          <Link href="/dashboard" className="font-semibold text-lg">
            {t.appName}
          </Link>

          <div className="flex items-center gap-3">
            {user ? (
              <div className="text-end">
                <div className="text-sm font-medium">{user.full_name}</div>
                <div className="text-xs text-muted-foreground">
                  {user.username} · {user.role}
                </div>
              </div>
            ) : null}

            <LanguageToggle />
            <ThemeToggle />

            <Button variant="outline" onClick={logout}>
              {t.logout}
            </Button>
          </div>
        </div>
      </header>

      <div className="flex">
        <aside className="w-64 border-e min-h-[calc(100vh-4rem)] bg-background p-4">
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

        <main className="flex-1 p-6">{children}</main>
      </div>
    </div>
  );
}