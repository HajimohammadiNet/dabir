"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

import { getSetupStatus } from "@/lib/api/setup";

export default function HomePage() {
  const router = useRouter();

  useEffect(() => {
    async function checkSetup() {
      try {
        const status = await getSetupStatus();

        if (status.setup_needed) {
          router.replace("/setup");
          return;
        }

        router.replace("/login");
      } catch {
        router.replace("/login");
      }
    }

    void checkSetup();
  }, [router]);

  return (
    <main className="min-h-screen flex items-center justify-center text-sm text-muted-foreground">
      Loading...
    </main>
  );
}