import { createSignal, onMount, onCleanup, Show, For } from "solid-js";
import { useNavigate, useSearchParams } from "@solidjs/router";
import { council } from "../../store/council";
import type { AuditLogEntry } from "../../types/admin";
import { AUDIT_ACTION_LABELS, AUDIT_TARGET_LABELS } from "../../types/admin";
import { formatDateTime } from "../../utils/format";
import Pagination from "../../components/Pagination";
import StaffGuard from "../../components/StaffGuard";

export default function CouncilAuditLog() {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const [actionFilter, setActionFilter] = createSignal(council.auditAction());
  const [targetFilter, setTargetFilter] = createSignal(council.auditTargetType());
  const [actionOpen, setActionOpen] = createSignal(false);
  const [targetOpen, setTargetOpen] = createSignal(false);
  let actionRef: HTMLDivElement | undefined;
  let targetRef: HTMLDivElement | undefined;

  onMount(() => {
    const p = parseInt(searchParams.page as string) || council.auditPage();
    council.loadAuditLogs(p);
    function handleClickOutside(e: MouseEvent) {
      if (actionRef && !actionRef.contains(e.target as Node)) setActionOpen(false);
      if (targetRef && !targetRef.contains(e.target as Node)) setTargetOpen(false);
    }
    document.addEventListener("mousedown", handleClickOutside);
    onCleanup(() => document.removeEventListener("mousedown", handleClickOutside));
  });

  function pickAction(value: string) {
    setActionFilter(value);
    setActionOpen(false);
    council.setAuditAction(value);
    council.setAuditTargetType(targetFilter());
    setSearchParams({ page: "1" });
    council.loadAuditLogs(1);
  }

  function pickTarget(value: string) {
    setTargetFilter(value);
    setTargetOpen(false);
    council.setAuditTargetType(value);
    council.setAuditAction(actionFilter());
    setSearchParams({ page: "1" });
    council.loadAuditLogs(1);
  }

  function clearFilters() {
    setActionFilter("");
    setTargetFilter("");
    council.setAuditAction("");
    council.setAuditTargetType("");
    setSearchParams({ page: "1" });
    council.loadAuditLogs(1);
  }

  function goToPage(p: number) {
    setSearchParams({ page: String(p) });
    council.loadAuditLogs(p);
  }

  function actionLabel(action: string) {
    return AUDIT_ACTION_LABELS[action] || action;
  }

  function targetLabel(type: string) {
    return AUDIT_TARGET_LABELS[type] || type;
  }

  return (
    <StaffGuard>
    <section>
      <h2 class="page-title">Audit Log</h2>

      <div class="council-audit-filters">
        <div class="council-audit-dropdown" ref={actionRef}>
          <button type="button" class="council-audit-dropdown-trigger" onClick={() => { setTargetOpen(false); setActionOpen(!actionOpen()); }}>
            {actionFilter() ? actionLabel(actionFilter()) : "All Actions"}
          </button>
          <Show when={actionOpen()}>
            <div class="council-audit-dropdown-menu">
              <button type="button" class="council-audit-dropdown-item" classList={{ "council-audit-dropdown-item-selected": !actionFilter() }} onClick={() => pickAction("")}>All Actions</button>
              <For each={Object.entries(AUDIT_ACTION_LABELS)}>
                {([key, label]: [string, string]) => (
                  <button type="button" class="council-audit-dropdown-item" classList={{ "council-audit-dropdown-item-selected": actionFilter() === key }} onClick={() => pickAction(key)}>{label}</button>
                )}
              </For>
            </div>
          </Show>
        </div>
        <div class="council-audit-dropdown" ref={targetRef}>
          <button type="button" class="council-audit-dropdown-trigger" onClick={() => { setActionOpen(false); setTargetOpen(!targetOpen()); }}>
            {targetFilter() ? targetLabel(targetFilter()) : "All Targets"}
          </button>
          <Show when={targetOpen()}>
            <div class="council-audit-dropdown-menu">
              <button type="button" class="council-audit-dropdown-item" classList={{ "council-audit-dropdown-item-selected": !targetFilter() }} onClick={() => pickTarget("")}>All Targets</button>
              <For each={Object.entries(AUDIT_TARGET_LABELS)}>
                {([key, label]: [string, string]) => (
                  <button type="button" class="council-audit-dropdown-item" classList={{ "council-audit-dropdown-item-selected": targetFilter() === key }} onClick={() => pickTarget(key)}>{label}</button>
                )}
              </For>
            </div>
          </Show>
        </div>
        <Show when={actionFilter() || targetFilter()}>
          <button type="button" class="council-audit-clear-btn" onClick={clearFilters}>Clear</button>
        </Show>
      </div>

      <div class="council-grid council-grid-audit">
        <div class="council-grid-header">
          <span>Date</span>
          <span>Action</span>
          <span>Actor</span>
          <span>Target</span>
          <span>Summary</span>
        </div>
        <Show when={!council.auditLoading()} fallback={
          <div class="council-grid-empty">Loading...</div>
        }>
          <Show when={council.auditLogs().length} fallback={
            <div class="council-grid-empty">No audit logs found.</div>
          }>
            <For each={council.auditLogs()}>
              {(entry: AuditLogEntry) => (
                <div class="council-grid-row" onClick={() => navigate(`/council/auditlog/${entry.system_ref}`)}>
                  <span class="council-audit-date">{formatDateTime(entry.created_at)}</span>
                  <span class="council-audit-action">{actionLabel(entry.action)}</span>
                  <span>{entry.actor}</span>
                  <span class="council-audit-target">
                    <span class="council-audit-target-type">{targetLabel(entry.target_type)}</span>
                    {entry.target_ref}
                  </span>
                  <span class="council-audit-summary">{entry.summary}</span>
                </div>
              )}
            </For>
          </Show>
        </Show>
      </div>

      <Pagination page={council.auditPage()} totalPages={council.auditTotalPages()} total={council.auditTotal()} label="entries" onPage={goToPage} />
    </section>
    </StaffGuard>
  );
}