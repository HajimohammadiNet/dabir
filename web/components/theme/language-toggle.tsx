"use client";

import { Button } from "@/components/ui/button";
import { useI18n } from "@/lib/i18n/i18n-context";

export function LanguageToggle() {
  const { locale, setLocale } = useI18n();

  return (
    <Button
      variant="outline"
      size="sm"
      onClick={() => setLocale(locale === "fa" ? "en" : "fa")}
    >
      {locale === "fa" ? "EN" : "FA"}
    </Button>
  );
}