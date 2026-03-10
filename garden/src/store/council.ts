import { createSignal } from "solid-js";
import { api } from "../api";
import { auth } from "./auth";
import type { AdminUser, AuditLogEntry, PaginatedResponse } from "../types/admin";

const [users, setUsers] = createSignal<AdminUser[]>([]);
const [total, setTotal] = createSignal(0);
const [page, setPage] = createSignal(1);
const [totalPages, setTotalPages] = createSignal(0);
const [loading, setLoading] = createSignal(false);
const [search, setSearch] = createSignal("");
const [sortField, setSortField] = createSignal("created_at");
const [sortOrder, setSortOrder] = createSignal<"asc" | "desc">("desc");

async function loadUsers(p = 1, q = "") {
  setLoading(true);
  const params = new URLSearchParams({
    page: String(p),
    per_page: "20",
    sort: sortField(),
    order: sortOrder(),
  });
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

function toggleSort(field: string) {
  if (sortField() === field) {
    setSortOrder(sortOrder() === "asc" ? "desc" : "asc");
  } else {
    setSortField(field);
    setSortOrder("desc");
  }
  loadUsers(1, search());
}

const [auditLogs, setAuditLogs] = createSignal<AuditLogEntry[]>([]);
const [auditTotal, setAuditTotal] = createSignal(0);
const [auditPage, setAuditPage] = createSignal(1);
const [auditTotalPages, setAuditTotalPages] = createSignal(0);
const [auditLoading, setAuditLoading] = createSignal(false);
const [auditAction, setAuditAction] = createSignal("");
const [auditTargetType, setAuditTargetType] = createSignal("");

async function loadAuditLogs(p = 1) {
  setAuditLoading(true);
  const params = new URLSearchParams({
    page: String(p),
    per_page: "20",
  });
  if (auditAction()) params.set("action", auditAction());
  if (auditTargetType()) params.set("target_type", auditTargetType());

  const response = await api<PaginatedResponse<AuditLogEntry>>(`/council/audit?${params}`, {
    token: auth.token(),
  });

  if (response.ok) {
    setAuditLogs(response.data.items);
    setAuditTotal(response.data.total);
    setAuditPage(response.data.page);
    setAuditTotalPages(response.data.total_pages);
  }
  setAuditLoading(false);
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
  sortField,
  sortOrder,
  toggleSort,
  auditLogs,
  auditTotal,
  auditPage,
  auditTotalPages,
  auditLoading,
  auditAction,
  setAuditAction,
  auditTargetType,
  setAuditTargetType,
  loadAuditLogs,
};