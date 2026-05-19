"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { clearAccessToken } from "@/lib/auth/token";
import { Button } from "@/components/ui/button";

const navItems = [
  { href: "/dashboard", label: "Dashboard" },
  { href: "/letters", label: "Letters" },
  { href: "/letters/import", label: "Import" },
  { href: "/users", label: "Users" },
  { href: "/audit-logs", label: "Audit Logs" },
  { href: "/settings", label: "Settings" },
];

export function AppShell({ children }: { children: React.ReactNode }) {
  const router = useRouter();

  function logout() {
    clearAccessToken();
    router.push("/login");
  }

  return (
    <div className="min-h-screen bg-muted/40">
      <header className="border-b bg-background">
        <div className="h-16 px-6 flex items-center justify-between">
          <Link href="/dashboard" className="font-semibold text-lg">
            Dabir
          </Link>

          <Button variant="outline" onClick={logout}>
            Logout
          </Button>
        </div>
      </header>

      <div className="flex">
        <aside className="w-64 border-r min-h-[calc(100vh-4rem)] bg-background p-4">
          <nav className="space-y-1">
            {navItems.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className="block rounded-md px-3 py-2 text-sm hover:bg-muted"
              >
                {item.label}
              </Link>
            ))}
          </nav>
        </aside>

        <main className="flex-1 p-6">{children}</main>
      </div>
    </div>
  );
}