import { createSignal } from "solid-js";
import { api } from "../api";
import type { Stats } from "../types/stats";

const [data, setData] = createSignal<Stats | null>(null);

async function load() {
  const token = localStorage.getItem("token");
  const response = await api<Stats>("/stats/", { token });
  if (response.ok) {
    setData(response.data);
  }
}

export const stats = { data, load };