import { createSignal, onMount, Show, For } from "solid-js";
import { useParams, A } from "@solidjs/router";
import { api } from "../../api";
import { auth } from "../../store/auth";
import { UserRole } from "../../types/roles";
import type { AdminUser } from "../../types/admin";
import Modal from "../../components/Modal";
import Editor from "../../components/Editor";

export default function CouncilUser() {
  const params = useParams();
  const [user, setUser] = createSignal<AdminUser | null>(null);
  const [error, setError] = createSignal("");
  const [actionError, setActionError] = createSignal("");
  const [editing, setEditing] = createSignal<Record<string, string>>({});
  const [modal, setModal] = createSignal<"warn" | "disable" | "ban" | null>(null);
  const [warnTitle, setWarnTitle] = createSignal("");
  const [warnBody, setWarnBody] = createSignal("");
  const [disableReason, setDisableReason] = createSignal("");
  const [disableUntil, setDisableUntil] = createSignal("");
  const [banReason, setBanReason] = createSignal("");

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
    const target = user();
    if (!me || !target || target.username === me.username) return false;
    if (target.role === UserRole.Owner) return false;
    if (me.role === UserRole.Owner) return true;
    if (me.role === UserRole.Admin && target.role !== UserRole.Admin) return true;
    return false;
  }

  function canModerate() {
    const me = auth.user();
    const target = user();
    if (!me || !target || target.username === me.username) return false;
    if (target.role === UserRole.Owner) return false;
    if (me.role === UserRole.Owner) return true;
    if (target.role === UserRole.Admin) return false;
    return true;
  }

  function isAdmin() {
    const me = auth.user();
    return me?.role === UserRole.Owner || me?.role === UserRole.Admin;
  }

  function statusBadge() {
    const target = user();
    if (!target) return "";
    if (target.account_banned) return "banned";
    if (target.account_disabled) return "disabled";
    if (!target.email_verified) return "unverified";
    return "active";
  }

  function formatDate(date: string | null) {
    if (!date) return "\u2014";
    return new Date(date).toLocaleDateString();
  }

  function startEdit(field: string, value: string) {
    setEditing((prev) => ({ ...prev, [field]: value }));
  }

  function cancelEdit(field: string) {
    setEditing((prev) => {
      const next = { ...prev };
      delete next[field];
      return next;
    });
  }

  async function saveEdit(field: string) {
    const target = user();
    if (!target) return;
    setActionError("");

    const response = await api<AdminUser>(`/council/users/${target.username}`, {
      method: "PATCH",
      token: auth.token(),
      body: { [field]: editing()[field] },
    });

    if (response.ok) {
      setUser(response.data);
      cancelEdit(field);
    } else {
      setActionError((response.data as unknown as { error: string }).error);
    }
  }

  async function saveRole(role: string) {
    const target = user();
    if (!target) return;
    setActionError("");

    const response = await api<AdminUser>(`/council/users/${target.username}/role`, {
      method: "POST",
      token: auth.token(),
      body: { role },
    });

    if (response.ok) {
      setUser(response.data);
      cancelEdit("role");
    } else {
      setActionError((response.data as unknown as { error: string }).error);
    }
  }

  async function submitWarn() {
    const target = user();
    if (!target) return;
    setActionError("");

    const response = await api<AdminUser>(`/council/users/${target.username}/warn`, {
      method: "POST",
      token: auth.token(),
      body: { title: warnTitle(), message: warnBody() },
    });

    if (response.ok) {
      setUser(response.data);
      setModal(null);
      setWarnTitle("");
      setWarnBody("");
    } else {
      setActionError((response.data as unknown as { error: string }).error);
    }
  }

  async function submitDisable() {
    const target = user();
    if (!target) return;
    setActionError("");

    const body: Record<string, string> = { reason: disableReason() };
    if (disableUntil()) body.disabled_until = new Date(disableUntil()).toISOString();

    const response = await api<AdminUser>(`/council/users/${target.username}/disable`, {
      method: "POST",
      token: auth.token(),
      body,
    });

    if (response.ok) {
      setUser(response.data);
      setModal(null);
      setDisableReason("");
      setDisableUntil("");
    } else {
      setActionError((response.data as unknown as { error: string }).error);
    }
  }

  async function submitBan() {
    const target = user();
    if (!target) return;
    setActionError("");

    const response = await api<AdminUser>(`/council/users/${target.username}/ban`, {
      method: "POST",
      token: auth.token(),
      body: { reason: banReason() },
    });

    if (response.ok) {
      setUser(response.data);
      setModal(null);
      setBanReason("");
    } else {
      setActionError((response.data as unknown as { error: string }).error);
    }
  }

  async function submitUnban() {
    const target = user();
    if (!target) return;
    setActionError("");

    const response = await api<AdminUser>(`/council/users/${target.username}/unban`, {
      method: "POST",
      token: auth.token(),
    });

    if (response.ok) {
      setUser(response.data);
    } else {
      setActionError((response.data as unknown as { error: string }).error);
    }
  }

  async function submitEnable() {
    const target = user();
    if (!target) return;
    setActionError("");

    const response = await api<AdminUser>(`/council/users/${target.username}/enable`, {
      method: "POST",
      token: auth.token(),
    });

    if (response.ok) {
      setUser(response.data);
    } else {
      setActionError((response.data as unknown as { error: string }).error);
    }
  }

  function EditableField(props: { label: string; field: string; value: string }) {
    return (
      <Show when={editing()[props.field] !== undefined} fallback={
        <div class="council-detail-row">
          <span class="council-detail-label">{props.label}</span>
          <span class="council-detail-value">
            <span>{props.value || "\u2014"}</span>
            <Show when={isAdmin()}>
              <button type="button" class="council-detail-edit-trigger" onClick={() => startEdit(props.field, props.value)}>Edit</button>
            </Show>
          </span>
        </div>
      }>
        <div class="council-detail-row">
          <span class="council-detail-label">{props.label}</span>
          <span class="council-detail-editable">
            <input
              type="text"
              class="council-detail-edit-input"
              value={editing()[props.field]}
              onInput={(e) => setEditing((prev) => ({ ...prev, [props.field]: e.currentTarget.value }))}
            />
            <button type="button" class="council-detail-edit-btn council-action-save" onClick={() => saveEdit(props.field)}>Save</button>
            <button type="button" class="council-detail-edit-btn" onClick={() => cancelEdit(props.field)}>Cancel</button>
          </span>
        </div>
      </Show>
    );
  }

  function RoleField(props: { role: string }) {
    const availableRoles = () => {
      const roles: UserRole[] = [UserRole.Member, UserRole.Moderator];
      if (auth.user()?.role === UserRole.Owner) roles.push(UserRole.Admin);
      return roles;
    };

    return (
      <Show when={editing()["role"] !== undefined} fallback={
        <div class="council-detail-row">
          <span class="council-detail-label">Role</span>
          <span class="council-detail-value">
            <span class={`council-role council-role-${props.role}`}>{props.role}</span>
            <Show when={canChangeRole()}>
              <button type="button" class="council-detail-edit-trigger" onClick={() => startEdit("role", props.role)}>Edit</button>
            </Show>
          </span>
        </div>
      }>
        <div class="council-detail-row">
          <span class="council-detail-label">Role</span>
          <span class="council-detail-editable">
            <select
              class="council-detail-role-select"
              aria-label="Change role"
              title="Change role"
              value={editing()["role"]}
              onChange={(e) => setEditing((prev) => ({ ...prev, role: e.currentTarget.value }))}
            >
              <For each={availableRoles()}>
                {(role) => <option value={role}>{role}</option>}
              </For>
            </select>
            <button type="button" class="council-detail-edit-btn council-action-save" onClick={() => saveRole(editing()["role"])}>Save</button>
            <button type="button" class="council-detail-edit-btn" onClick={() => cancelEdit("role")}>Cancel</button>
          </span>
        </div>
      </Show>
    );
  }

  return (
    <section>
      <A href="/council/users" class="council-detail-back">&larr; Back to Users</A>

      <Show when={error()}>
        <div class="form-error">{error()}</div>
      </Show>

      <Show when={user()}>
        {(target) => (
          <>
            <Show when={actionError()}>
              <div class="form-error">{actionError()}</div>
            </Show>

            <div class="council-detail-layout">
              <div class="council-detail-left">
                <div class="council-detail-section">
                  <div class="council-detail-section-header">Profile</div>
                  <div class="council-detail-table">
                    <EditableField label="Username" field="username" value={target().username} />
                    <EditableField label="Display Name" field="display_name" value={target().display_name} />
                    <EditableField label="Email" field="email" value={target().email} />
                    <EditableField label="Bio" field="bio" value={target().bio} />
                    <EditableField label="Website" field="website" value={target().website} />
                    <EditableField label="Location" field="location" value={target().location} />
                    <EditableField label="Pronouns" field="pronouns" value={target().pronouns} />
                    <EditableField label="Birthday" field="birthday" value={target().birthday?.split("T")[0] ?? ""} />
                    <EditableField label="Signature" field="signature" value={target().signature} />
                  </div>
                </div>

                <div class="council-detail-section">
                  <div class="council-detail-section-header">Account</div>
                  <div class="council-detail-table">
                    <RoleField role={target().role} />
                    <div class="council-detail-row">
                      <span class="council-detail-label">Verified</span>
                      <span>{target().email_verified ? "Yes" : "No"}</span>
                    </div>
                    <div class="council-detail-row">
                      <span class="council-detail-label">Jade</span>
                      <span>{target().jade}</span>
                    </div>
                    <div class="council-detail-row">
                      <span class="council-detail-label">Honor</span>
                      <span>{target().honor}</span>
                    </div>
                    <div class="council-detail-row">
                      <span class="council-detail-label">Warnings</span>
                      <span>{target().warning_count}</span>
                    </div>
                    <div class="council-detail-row">
                      <span class="council-detail-label">Joined</span>
                      <span>{formatDate(target().created_at)}</span>
                    </div>
                    <div class="council-detail-row">
                      <span class="council-detail-label">Last Seen</span>
                      <span>{formatDate(target().last_seen_at)}</span>
                    </div>
                    <div class="council-detail-row">
                      <span class="council-detail-label">IP</span>
                      <span>{target().registration_ip || "\u2014"}</span>
                    </div>
                    <Show when={target().account_banned}>
                      <div class="council-detail-row">
                        <span class="council-detail-label">Banned At</span>
                        <span>{formatDate(target().banned_at)}</span>
                      </div>
                      <Show when={target().banned_reason}>
                        <div class="council-detail-row">
                          <span class="council-detail-label">Ban Reason</span>
                          <span>{target().banned_reason}</span>
                        </div>
                      </Show>
                    </Show>
                    <Show when={target().account_disabled}>
                      <div class="council-detail-row">
                        <span class="council-detail-label">Disabled At</span>
                        <span>{formatDate(target().disabled_at)}</span>
                      </div>
                      <Show when={target().disabled_until}>
                        <div class="council-detail-row">
                          <span class="council-detail-label">Until</span>
                          <span>{formatDate(target().disabled_until)}</span>
                        </div>
                      </Show>
                      <Show when={target().disabled_reason}>
                        <div class="council-detail-row">
                          <span class="council-detail-label">Reason</span>
                          <span>{target().disabled_reason}</span>
                        </div>
                      </Show>
                    </Show>
                  </div>
                </div>
              </div>

              <div class="council-detail-right">
                <div class="council-detail-section">
                  <div class="council-detail-section-header">{target().display_name}</div>
                  <div class="council-detail-card">
                    <img src={target().avatar_url} alt="" class="council-detail-card-avatar" />
                    <span class="council-detail-card-username">@{target().username}</span>
                    <span class={`council-status council-status-${statusBadge()}`}>{statusBadge()}</span>
                    <span class={`council-role council-role-${target().role}`}>{target().role}</span>
                  </div>
                </div>

                <Show when={canModerate()}>
                  <div class="council-detail-section">
                    <div class="council-detail-section-header">Actions</div>
                    <div class="council-detail-actions">
                      <button type="button" class="council-detail-action-btn council-action-warn" onClick={() => setModal("warn")}>
                        Warn
                      </button>
                      <Show when={!target().account_disabled} fallback={
                        <button type="button" class="council-detail-action-btn council-action-enable" onClick={submitEnable}>
                          Enable
                        </button>
                      }>
                        <button type="button" class="council-detail-action-btn council-action-disable" onClick={() => setModal("disable")}>
                          Disable
                        </button>
                      </Show>
                      <Show when={!target().account_banned} fallback={
                        <button type="button" class="council-detail-action-btn council-action-unban" onClick={submitUnban}>
                          Unban
                        </button>
                      }>
                        <button type="button" class="council-detail-action-btn council-action-ban" onClick={() => setModal("ban")}>
                          Ban
                        </button>
                      </Show>
                    </div>
                  </div>
                </Show>
              </div>
            </div>

            <Show when={modal() === "warn"}>
              <Modal title="Issue Warning" onClose={() => setModal(null)}>
                <div class="modal-field">
                  <label class="modal-label">Title</label>
                  <input
                    type="text"
                    class="modal-input"
                    placeholder="Warning title..."
                    value={warnTitle()}
                    onInput={(e) => setWarnTitle(e.currentTarget.value)}
                  />
                </div>
                <div class="modal-field">
                  <label class="modal-label">Message</label>
                  <Editor onHtml={setWarnBody} />
                </div>
                <div class="modal-actions">
                  <button type="button" class="council-detail-action-btn council-action-warn" onClick={submitWarn}>Send Warning</button>
                  <button type="button" class="council-detail-action-btn" onClick={() => setModal(null)}>Cancel</button>
                </div>
              </Modal>
            </Show>

            <Show when={modal() === "disable"}>
              <Modal title="Disable Account" onClose={() => setModal(null)}>
                <div class="modal-field">
                  <label class="modal-label">Reason</label>
                  <Editor onHtml={setDisableReason} />
                </div>
                <div class="modal-field">
                  <label class="modal-label">Disabled Until (optional)</label>
                  <input
                    type="date"
                    class="modal-input"
                    value={disableUntil()}
                    onInput={(e) => setDisableUntil(e.currentTarget.value)}
                  />
                </div>
                <div class="modal-actions">
                  <button type="button" class="council-detail-action-btn council-action-disable" onClick={submitDisable}>Disable</button>
                  <button type="button" class="council-detail-action-btn" onClick={() => setModal(null)}>Cancel</button>
                </div>
              </Modal>
            </Show>

            <Show when={modal() === "ban"}>
              <Modal title="Ban Account" onClose={() => setModal(null)}>
                <div class="modal-field">
                  <label class="modal-label">Reason</label>
                  <Editor onHtml={setBanReason} />
                </div>
                <div class="modal-actions">
                  <button type="button" class="council-detail-action-btn council-action-ban" onClick={submitBan}>Ban</button>
                  <button type="button" class="council-detail-action-btn" onClick={() => setModal(null)}>Cancel</button>
                </div>
              </Modal>
            </Show>
          </>
        )}
      </Show>
    </section>
  );
}