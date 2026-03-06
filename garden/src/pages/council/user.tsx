import { createSignal, onMount, Show } from "solid-js";
import { useParams, A } from "@solidjs/router";
import { api } from "../../api";
import { auth } from "../../store/auth";
import { UserRole } from "../../types/roles";
import type { AdminUser } from "../../types/admin";

export default function CouncilUser() {
  const params = useParams();
  const [user, setUser] = createSignal<AdminUser | null>(null);
  const [error, setError] = createSignal("");
  const [actionError, setActionError] = createSignal("");
  const [reason, setReason] = createSignal("");

  onMount(async () => {
    const response = await api<AdminUser>(`/council/users/${params.username}`, {
      token: auth.token(),
    });
    if (response.ok) {
      setUser(response.data);
    } else {
      setError("User not found.");
    }
  });

  function canChangeRole() {
    const me = auth.user();
    const u = user();
    if (!me || !u || u.username === me.username) return false;
    if (u.role === UserRole.Owner) return false;
    if (me.role === UserRole.Owner) return true;
    if (me.role === UserRole.Admin && u.role !== UserRole.Admin) return true;
    return false;
  }

  function canModerate() {
    const me = auth.user();
    const u = user();
    if (!me || !u || u.username === me.username) return false;
    if (u.role === UserRole.Owner) return false;
    if (me.role === UserRole.Owner) return true;
    if (u.role === UserRole.Admin) return false;
    return true;
  }

  function statusBadge() {
    const u = user();
    if (!u) return "";
    if (u.account_banned) return "banned";
    if (u.account_disabled) return "disabled";
    if (!u.email_verified) return "unverified";
    return "active";
  }

  function formatDate(date: string | null) {
    if (!date) return "\u2014";
    return new Date(date).toLocaleDateString();
  }

  async function handleAction(action: string) {
    const u = user();
    if (!u) return;
    setActionError("");

    const body = action === "ban" || action === "disable" ? { reason: reason() } : undefined;

    const response = await api<AdminUser | { error: string }>(`/council/users/${u.username}/${action}`, {
      method: "POST",
      token: auth.token(),
      body,
    });

    if (response.ok) {
      setUser(response.data as AdminUser);
      setReason("");
    } else {
      setActionError((response.data as { error: string }).error);
    }
  }

  async function handleRoleChange(role: string) {
    const u = user();
    if (!u) return;
    setActionError("");

    const response = await api<AdminUser | { error: string }>(`/council/users/${u.username}/role`, {
      method: "POST",
      token: auth.token(),
      body: { role },
    });

    if (response.ok) {
      setUser(response.data as AdminUser);
    } else {
      setActionError((response.data as { error: string }).error);
    }
  }

  return (
    <section>
      <Show when={error()}>
        <div class="form-error">{error()}</div>
      </Show>

      <Show when={user()}>
        {(u) => (
          <>
            <div class="council-detail-header">
              <img src={u().avatar_url} alt="" class="council-detail-avatar" />
              <div>
                <h2 class="page-title">{u().display_name}</h2>
                <span class="council-detail-handle">@{u().username}</span>
              </div>
            </div>

            <Show when={actionError()}>
              <div class="form-error">{actionError()}</div>
            </Show>

            <div class="council-detail-section">
              <div class="council-detail-section-header">Details</div>
              <div class="council-detail-table">
                <div class="council-detail-row">
                  <span class="council-detail-label">Email</span>
                  <span>{u().email}</span>
                </div>
                <div class="council-detail-row">
                  <span class="council-detail-label">Role</span>
                  <span>
                    <Show when={canChangeRole()} fallback={
                      <span class={`council-role council-role-${u().role}`}>{u().role}</span>
                    }>
                      <select
                        class="council-detail-role-select"
                        title="Change role"
                        value={u().role}
                        onChange={(e) => handleRoleChange(e.currentTarget.value)}
                      >
                        <option value="member">member</option>
                        <option value="moderator">moderator</option>
                        <Show when={auth.user()?.role === UserRole.Owner}>
                          <option value="admin">admin</option>
                        </Show>
                      </select>
                    </Show>
                  </span>
                </div>
                <div class="council-detail-row">
                  <span class="council-detail-label">Status</span>
                  <span class={`council-status council-status-${statusBadge()}`}>
                    {statusBadge()}
                  </span>
                </div>
                <div class="council-detail-row">
                  <span class="council-detail-label">Verified</span>
                  <span>{u().email_verified ? "Yes" : "No"}</span>
                </div>
                <div class="council-detail-row">
                  <span class="council-detail-label">Joined</span>
                  <span>{formatDate(u().created_at)}</span>
                </div>
                <div class="council-detail-row">
                  <span class="council-detail-label">Last Seen</span>
                  <span>{formatDate(u().last_seen_at)}</span>
                </div>
                <Show when={u().account_banned}>
                  <div class="council-detail-row">
                    <span class="council-detail-label">Banned At</span>
                    <span>{formatDate(u().banned_at)}</span>
                  </div>
                  <Show when={u().banned_reason}>
                    <div class="council-detail-row">
                      <span class="council-detail-label">Ban Reason</span>
                      <span>{u().banned_reason}</span>
                    </div>
                  </Show>
                </Show>
                <Show when={u().account_disabled}>
                  <div class="council-detail-row">
                    <span class="council-detail-label">Disabled At</span>
                    <span>{formatDate(u().disabled_at)}</span>
                  </div>
                  <Show when={u().disabled_reason}>
                    <div class="council-detail-row">
                      <span class="council-detail-label">Disable Reason</span>
                      <span>{u().disabled_reason}</span>
                    </div>
                  </Show>
                </Show>
              </div>
            </div>

            <Show when={canModerate()}>
              <div class="council-detail-section">
                <div class="council-detail-section-header">Actions</div>
                <div class="council-detail-actions">
                  <Show when={!u().account_banned} fallback={
                    <button type="button" class="council-detail-action-btn council-action-unban" onClick={() => handleAction("unban")}>
                      Unban
                    </button>
                  }>
                    <div class="council-detail-action-row">
                      <input
                        type="text"
                        class="council-detail-reason"
                        placeholder="Reason for ban..."
                        value={reason()}
                        onInput={(e) => setReason(e.currentTarget.value)}
                      />
                      <button type="button" class="council-detail-action-btn council-action-ban" onClick={() => handleAction("ban")}>
                        Ban
                      </button>
                    </div>
                  </Show>

                  <Show when={!u().account_disabled} fallback={
                    <button type="button" class="council-detail-action-btn council-action-enable" onClick={() => handleAction("enable")}>
                      Enable
                    </button>
                  }>
                    <div class="council-detail-action-row">
                      <input
                        type="text"
                        class="council-detail-reason"
                        placeholder="Reason for disable..."
                        value={reason()}
                        onInput={(e) => setReason(e.currentTarget.value)}
                      />
                      <button type="button" class="council-detail-action-btn council-action-disable" onClick={() => handleAction("disable")}>
                        Disable
                      </button>
                    </div>
                  </Show>
                </div>
              </div>
            </Show>

            <A href="/council/users" class="council-detail-back">&larr; Back to Users</A>
          </>
        )}
      </Show>
    </section>
  );
}