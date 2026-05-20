"use client";

import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";

import { dictionaries, type Locale } from "./dictionaries";

type Dictionary = (typeof dictionaries)[Locale];

type I18nContextValue = {
  locale: Locale;
  direction: "rtl" | "ltr";
  t: Dictionary;
  setLocale: (locale: Locale) => void;
};

const I18N_STORAGE_KEY = "dabir_locale";

const I18nContext = createContext<I18nContextValue | null>(null);

function getInitialLocale(): Locale {
  if (typeof window === "undefined") {
    return "fa";
  }

  const stored = window.localStorage.getItem(I18N_STORAGE_KEY);

  if (stored === "fa" || stored === "en") {
    return stored;
  }

  return "fa";
}

export function I18nProvider({ children }: { children: ReactNode }) {
  const [locale, setLocaleState] = useState<Locale>(() => getInitialLocale());

  const direction: "rtl" | "ltr" = locale === "fa" ? "rtl" : "ltr";

  function setLocale(nextLocale: Locale) {
    setLocaleState(nextLocale);

    if (typeof window !== "undefined") {
      window.localStorage.setItem(I18N_STORAGE_KEY, nextLocale);
    }
  }

  useEffect(() => {
    document.documentElement.lang = locale;
    document.documentElement.dir = direction;
  }, [locale, direction]);

  const value = useMemo<I18nContextValue>(
    () => ({
      locale,
      direction,
      t: dictionaries[locale] as Dictionary,
      setLocale,
    }),
    [locale, direction]
  );

  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>;
}

export function useI18n() {
  const context = useContext(I18nContext);

  if (!context) {
    throw new Error("useI18n must be used inside I18nProvider");
  }

  return context;
}