import { createSignal } from "solid-js";
import { api } from "../api";

interface User {
  id: number;
  username: string;
  email: string;
  display_name: string;
  bio: string;
  birthday: string | null;
  avatar_url: string;
  blinkie_url: string;
  website: string;
  location: string;
  pronouns: string;
  signature: string;
  role: string;
  created_at: string;
}

interface AuthResponse {
  token: string;
  user: User;
}

interface MessageResponse {
  message: string;
}

interface ErrorResponse {
  error: string;
}

const [user, setUser] = createSignal<User | null>(null);
const [token, setToken] = createSignal<string | null>(localStorage.getItem("token"));
const [loading, setLoading] = createSignal(true);

async function initialize() {
  const stored = token();
  if (!stored) {
    setLoading(false);
    return;
  }

  try {
    const response = await api<User>("/auth/me", { token: stored });
    if (response.ok) {
      setUser(response.data);
    } else {
      localStorage.removeItem("token");
      setToken(null);
    }
  } catch {
    localStorage.removeItem("token");
    setToken(null);
  }
  setLoading(false);
}

async function login(username: string, password: string): Promise<string | null> {
  const response = await api<AuthResponse | ErrorResponse>("/auth/login", {
    method: "POST",
    body: { username, password },
  });

  if (response.ok) {
    const data = response.data as AuthResponse;
    localStorage.setItem("token", data.token);
    setToken(data.token);
    setUser(data.user);
    return null;
  }

  return (response.data as ErrorResponse).error;
}

async function register(username: string, email: string, password: string, displayName: string): Promise<string | null> {
  const response = await api<MessageResponse | ErrorResponse>("/auth/register", {
    method: "POST",
    body: { username, email, password, display_name: displayName },
  });

  if (response.ok) return null;
  return (response.data as ErrorResponse).error;
}

async function verify(verificationToken: string, type: string): Promise<string | null> {
  const response = await api<MessageResponse | ErrorResponse>("/auth/verify", {
    method: "POST",
    body: { token: verificationToken, type },
  });

  if (response.ok) return null;
  return (response.data as ErrorResponse).error;
}

async function reactivate(email: string): Promise<string | null> {
  const response = await api<MessageResponse | ErrorResponse>("/auth/reactivate", {
    method: "POST",
    body: { email },
  });

  if (response.ok) return null;
  return (response.data as ErrorResponse).error;
}

async function logout() {
  const stored = token();
  if (stored) {
    await api("/auth/logout", { method: "POST", token: stored });
  }
  localStorage.removeItem("token");
  setToken(null);
  setUser(null);
}

export const auth = { user, token, loading, initialize, login, register, verify, reactivate, logout };