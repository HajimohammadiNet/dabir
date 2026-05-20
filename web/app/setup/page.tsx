"use client";

import { FormEvent, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

import { getSetupStatus, initializeSetup } from "@/lib/api/setup";

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
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

export default function SetupPage() {
  const router = useRouter();

  const [checking, setChecking] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  const [organizationName, setOrganizationName] = useState("Dabir");

  const [username, setUsername] = useState("admin");
  const [fullName, setFullName] = useState("System Administrator");
  const [password, setPassword] = useState("");

  const [numberPrefix, setNumberPrefix] = useState("DABIR");
  const [numberPadding, setNumberPadding] = useState(6);

  const [error, setError] = useState("");

  useEffect(() => {
    async function checkSetup() {
      try {
        const status = await getSetupStatus();

        if (!status.setup_needed) {
          router.replace("/login");
          return;
        }
      } catch {
        setError("Could not check setup status. Please make sure API is running.");
      } finally {
        setChecking(false);
      }
    }

    void checkSetup();
  }, [router]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setError("");
    setSubmitting(true);

    try {
      await initializeSetup({
        organization_name: organizationName,
        superuser: {
          username,
          full_name: fullName,
          password,
        },
        letter_config: {
          number_prefix: numberPrefix,
          number_padding: numberPadding,
        },
      });

      toast.success("Application initialized successfully");
      router.push("/login");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Setup failed");
    } finally {
      setSubmitting(false);
    }
  }

  if (checking) {
    return (
      <main className="min-h-screen flex items-center justify-center text-sm text-muted-foreground">
        Checking setup status...
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-muted/40 flex items-center justify-center p-4">
      <Card className="w-full max-w-2xl">
        <CardHeader>
          <CardTitle>Setup Dabir</CardTitle>
          <CardDescription>
            Initialize the application and create the first superuser.
          </CardDescription>
        </CardHeader>

        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            {error ? (
              <Alert variant="destructive">
                <AlertTitle>Setup error</AlertTitle>
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            ) : null}

            <section className="space-y-4">
              <div>
                <h2 className="font-semibold">Organization</h2>
                <p className="text-sm text-muted-foreground">
                  Basic organization information.
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="organization_name">Organization Name</Label>
                <Input
                  id="organization_name"
                  value={organizationName}
                  onChange={(event) => setOrganizationName(event.target.value)}
                  required
                />
              </div>
            </section>

            <section className="space-y-4">
              <div>
                <h2 className="font-semibold">Superuser</h2>
                <p className="text-sm text-muted-foreground">
                  This user will have full access to Dabir.
                </p>
              </div>

              <div className="grid gap-4 md:grid-cols-2">
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
            </section>

            <section className="space-y-4">
              <div>
                <h2 className="font-semibold">Letter Numbering</h2>
                <p className="text-sm text-muted-foreground">
                  Configure how letter numbers are displayed.
                </p>
              </div>

              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="number_prefix">Number Prefix</Label>
                  <Input
                    id="number_prefix"
                    value={numberPrefix}
                    onChange={(event) => setNumberPrefix(event.target.value)}
                    required
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="number_padding">Number Padding</Label>
                  <Input
                    id="number_padding"
                    type="number"
                    min={1}
                    max={12}
                    value={numberPadding}
                    onChange={(event) =>
                      setNumberPadding(Number(event.target.value))
                    }
                    required
                  />
                </div>
              </div>

              <div className="rounded-md border bg-muted/30 p-3 text-sm">
                Example:{" "}
                <span className="font-medium">
                  {numberPrefix || "DABIR"}-
                  {"1".padStart(numberPadding || 6, "0")}
                </span>
              </div>
            </section>

            <Button type="submit" className="w-full" disabled={submitting}>
              {submitting ? "Initializing..." : "Initialize Application"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </main>
  );
}