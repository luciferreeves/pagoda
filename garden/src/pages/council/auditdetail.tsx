import { createSignal, onMount, Show, For } from "solid-js";
import { useParams, A } from "@solidjs/router";
import { api } from "../../api";
import { auth } from "../../store/auth";
import type { AuditLogDetail } from "../../types/admin";
import { AUDIT_ACTION_LABELS } from "../../types/admin";
import { ROLE_LABELS } from "../../types/roles";
import { formatDateTimeFull } from "../../utils/format";
import StaffGuard from "../../components/StaffGuard";

interface FieldChange {
  field: string;
  old: string | number | boolean | null;
  new: string | number | boolean | null;
}

export default function CouncilAuditDetail() {
  const params = useParams();
  const [entry, setEntry] = createSignal<AuditLogDetail | null>(null);
  const [error, setError] = createSignal("");

  onMount(async () => {
    const response = await api<AuditLogDetail>(`/council/audit/${params.ref}`, {
      token: auth.token(),
    });
    if (response.ok) {
      setEntry(response.data);
    } else {
      setError("Audit log not found.");
    }
  });

  function actionLabel(action: string) {
    return AUDIT_ACTION_LABELS[action] || action;
  }

  function parseDetails(details: string): Record<string, unknown> | null {
    if (!details) return null;
    try {
      return JSON.parse(details);
    } catch {
      return null;
    }
  }

  function formatValue(val: string | number | boolean | null | undefined): string {
    if (val === null || val === undefined) return "—";
    if (typeof val === "boolean") return val ? "Yes" : "No";
    return String(val);
  }

  function renderDetails(action: string, data: Record<string, unknown>) {
    switch (action) {
      case "user.ban":
        return renderBanDetails(data as { reason: string; system_ref: string });
      case "user.disable":
        return renderDisableDetails(data as { reason: string; disabled_until: string | null; system_ref: string });
      case "user.role_change":
        return renderRoleChange(data as { old_role: string; new_role: string });
      case "user.warn":
        return renderWarning(data as { warning_ref: string; title: string; message: string });
      case "user.unwarn":
        return renderDeactivateWarning(data as { warning_ref: string });
      case "user.edit":
        return renderEditUser(data as { changes: FieldChange[] });
      default:
        return renderGeneric(data);
    }
  }

  function renderBanDetails(data: { reason: string; system_ref: string }) {
    return (
      <div class="council-detail-table">
        <div class="council-detail-row">
          <span class="council-detail-label">Reason</span>
          <span class="council-detail-value council-detail-html" innerHTML={data.reason} />
        </div>
        <div class="council-detail-row">
          <span class="council-detail-label">Letter Ref</span>
          <span class="council-detail-value council-audit-ref">{data.system_ref}</span>
        </div>
      </div>
    );
  }

  function renderDisableDetails(data: { reason: string; disabled_until: string | null; system_ref: string }) {
    return (
      <div class="council-detail-table">
        <div class="council-detail-row">
          <span class="council-detail-label">Reason</span>
          <span class="council-detail-value council-detail-html" innerHTML={data.reason} />
        </div>
        <Show when={data.disabled_until}>
          <div class="council-detail-row">
            <span class="council-detail-label">Until</span>
            <span class="council-detail-value">{formatDateTimeFull(data.disabled_until!)}</span>
          </div>
        </Show>
        <div class="council-detail-row">
          <span class="council-detail-label">Letter Ref</span>
          <span class="council-detail-value council-audit-ref">{data.system_ref}</span>
        </div>
      </div>
    );
  }

  function renderRoleChange(data: { old_role: string; new_role: string }) {
    return (
      <div class="council-detail-table">
        <div class="council-detail-row">
          <span class="council-detail-label">Old Role</span>
          <span class="council-detail-value">
            <span class={`council-role council-role-${data.old_role}`}>{ROLE_LABELS[data.old_role] || data.old_role}</span>
          </span>
        </div>
        <div class="council-detail-row">
          <span class="council-detail-label">New Role</span>
          <span class="council-detail-value">
            <span class={`council-role council-role-${data.new_role}`}>{ROLE_LABELS[data.new_role] || data.new_role}</span>
          </span>
        </div>
      </div>
    );
  }

  function renderWarning(data: { warning_ref: string; title: string; message: string }) {
    return (
      <div class="council-detail-table">
        <div class="council-detail-row">
          <span class="council-detail-label">Ref</span>
          <span class="council-detail-value council-audit-ref">{data.warning_ref}</span>
        </div>
        <div class="council-detail-row">
          <span class="council-detail-label">Title</span>
          <span class="council-detail-value">{data.title}</span>
        </div>
        <div class="council-detail-row">
          <span class="council-detail-label">Message</span>
          <span class="council-detail-value council-detail-html" innerHTML={data.message} />
        </div>
      </div>
    );
  }

  function renderDeactivateWarning(data: { warning_ref: string }) {
    return (
      <div class="council-detail-table">
        <div class="council-detail-row">
          <span class="council-detail-label">Warning Ref</span>
          <span class="council-detail-value council-audit-ref">{data.warning_ref}</span>
        </div>
      </div>
    );
  }

  function renderEditUser(data: { changes: FieldChange[] }) {
    return (
      <div class="council-detail-table">
        <For each={data.changes}>
          {(change: FieldChange) => (
            <div class="council-detail-row">
              <span class="council-detail-label">{change.field}</span>
              <span class="council-detail-value council-audit-change">
                <Show when={change.old !== null && change.old !== undefined}>
                  <span class="council-audit-old">{formatValue(change.old)}</span>
                  <span class="council-audit-arrow">→</span>
                </Show>
                <span class="council-audit-new">{formatValue(change.new)}</span>
              </span>
            </div>
          )}
        </For>
      </div>
    );
  }

  function renderGeneric(data: Record<string, unknown>) {
    return (
      <div class="council-detail-table">
        <For each={Object.entries(data)}>
          {([key, val]: [string, unknown]) => (
            <div class="council-detail-row">
              <span class="council-detail-label">{key}</span>
              <span class="council-detail-value">{typeof val === "object" ? JSON.stringify(val) : formatValue(val as string | number | boolean | null)}</span>
            </div>
          )}
        </For>
      </div>
    );
  }

  return (
    <StaffGuard>
    <section>
      <A href="/council/auditlog" class="council-detail-back">← Back to Audit Log</A>
      <h2 class="page-title">Audit Log Detail</h2>

      <Show when={error()}>
        <div class="form-error">{error()}</div>
      </Show>

      <Show when={entry()}>
        {(e) => {
          const details = parseDetails(e().details);
          return (
            <>
              <div class="council-detail-section">
                <div class="council-detail-section-header">Entry</div>
                <div class="council-detail-table">
                  <div class="council-detail-row">
                    <span class="council-detail-label">Ref</span>
                    <span class="council-detail-value council-audit-ref">{e().system_ref}</span>
                  </div>
                  <div class="council-detail-row">
                    <span class="council-detail-label">Date</span>
                    <span class="council-detail-value">{formatDateTimeFull(e().created_at)}</span>
                  </div>
                  <div class="council-detail-row">
                    <span class="council-detail-label">Action</span>
                    <span class="council-detail-value council-audit-action">{actionLabel(e().action)}</span>
                  </div>
                  <div class="council-detail-row">
                    <span class="council-detail-label">Actor</span>
                    <span class="council-detail-value">
                      <A href={`/council/users/${e().actor}`}>{e().actor}</A>
                    </span>
                  </div>
                  <div class="council-detail-row">
                    <span class="council-detail-label">Target</span>
                    <span class="council-detail-value">
                      <span class="council-audit-target-type">{e().target_type}</span>
                      <Show when={e().target_type === "user"} fallback={<span>{e().target_ref}</span>}>
                        <A href={`/council/users/${e().target_ref}`}>{e().target_ref}</A>
                      </Show>
                    </span>
                  </div>
                  <div class="council-detail-row">
                    <span class="council-detail-label">Summary</span>
                    <span class="council-detail-value">{e().summary}</span>
                  </div>
                </div>
              </div>

              <Show when={details}>
                <div class="council-detail-section">
                  <div class="council-detail-section-header">Details</div>
                  {renderDetails(e().action, details!)}
                </div>
              </Show>
            </>
          );
        }}
      </Show>
    </section>
    </StaffGuard>
  );
}