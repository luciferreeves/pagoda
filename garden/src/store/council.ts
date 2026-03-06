import { createSignal } from "solid-js";
import { api } from "../api";
import { auth } from "./auth";
import type { AdminUser, PaginatedResponse } from "../types/admin";

const [users, setUsers] = createSignal<AdminUser[]>([]);
const [total, setTotal] = createSignal(0);
const [page, setPage] = createSignal(1);
const [totalPages, setTotalPages] = createSignal(0);
const [loading, setLoading] = createSignal(false);
const [search, setSearch] = createSignal("");

async function loadUsers(p = 1, q = "") {
  setLoading(true);
  const params = new URLSearchParams({ page: String(p), per_page: "20" });
  if (q) params.set("search", q);

  const response = await api<PaginatedResponse<AdminUser>>(`/council/users?${params}`, {
    token: auth.token(),
  });

  if (response.ok) {
    setUsers(response.data.items);
    setTotal(response.data.total);
    setPage(response.data.page);
    setTotalPages(response.data.total_pages);
  }
  setLoading(false);
}

export const council = {
  users,
  total,
  page,
  totalPages,
  loading,
  search,
  setSearch,
  loadUsers,
};