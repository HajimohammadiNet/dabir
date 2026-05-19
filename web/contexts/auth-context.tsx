"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import { useRouter } from "next/navigation";

import { getMe } from "@/lib/api/auth";
import {
  clearAccessToken,
  getAccessToken,
  setAccessToken,
} from "@/lib/auth/token";
import type { AuthUser } from "@/types/auth";

type AuthContextValue = {
  user: AuthUser | null;
  token: string | null;
  loading: boolean;
  isAuthenticated: boolean;
  setSession: (token: string, user?: AuthUser) => void;
  logout: () => void;
  refreshMe: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const router = useRouter();

  const [token, setToken] = useState<string | null>(null);
  const [user, setUser] = useState<AuthUser | null>(null);
  const [loading, setLoading] = useState(true);

  const refreshMe = useCallback(async () => {
    const storedToken = getAccessToken();

    if (!storedToken) {
      setToken(null);
      setUser(null);
      setLoading(false);
      return;
    }

    setToken(storedToken);

    try {
      const me = await getMe(storedToken);
      setUser(me);
    } catch {
      clearAccessToken();
      setToken(null);
      setUser(null);
      router.push("/login");
    } finally {
      setLoading(false);
    }
  }, [router]);

  const setSession = useCallback((nextToken: string, nextUser?: AuthUser) => {
    setAccessToken(nextToken);
    setToken(nextToken);

    if (nextUser) {
      setUser(nextUser);
    }
  }, []);

  const logout = useCallback(() => {
    clearAccessToken();
    setToken(null);
    setUser(null);
    router.push("/login");
  }, [router]);

  useEffect(() => {
    const loadSession = async () => {
      await refreshMe();
    };

    void loadSession();
  }, [refreshMe]);

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      token,
      loading,
      isAuthenticated: Boolean(token && user),
      setSession,
      logout,
      refreshMe,
    }),
    [user, token, loading, setSession, logout, refreshMe]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error("useAuth must be used inside AuthProvider");
  }

  return context;
}