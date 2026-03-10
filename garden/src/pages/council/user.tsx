import { createSignal, onMount, onCleanup, Show, For } from "solid-js";
import { useParams, A } from "@solidjs/router";
import { api, uploadFile } from "../../api";
import { auth } from "../../store/auth";
import { UserRole, ROLE_LABELS } from "../../types/roles";
import type { AdminUser } from "../../types/admin";
import Modal from "../../components/Modal";
import Editor from "../../components/Editor";
import StaffGuard from "../../components/StaffGuard";
import DatePicker from "../../components/DatePicker";
import MonthDayPicker from "../../components/MonthDayPicker";
import MiniEditor from "../../components/MiniEditor";

export default function CouncilUser() {
  const params = useParams();
  const [user, setUser] = createSignal<AdminUser | null>(null);
  const [error, setError] = createSignal("");
  const [actionError, setActionError] = createSignal("");
  const [editing, setEditing] = createSignal<Record<string, string>>({});
  const [modal, setModal] = createSignal<"warn" | "disable" | "ban" | "jade" | null>(null);
  const [warnTitle, setWarnTitle] = createSignal("");
  const [warnBody, setWarnBody] = createSignal("");
  const [disableReason, setDisableReason] = createSignal("");
  const [disableUntil, setDisableUntil] = createSignal("");
  const [banReason, setBanReason] = createSignal("");
  const [jadeAmount, setJadeAmount] = createSignal("");
  const [editingBio, setEditingBio] = createSignal(false);
  const [bioHtml, setBioHtml] = createSignal("");
  const [editingSignature, setEditingSignature] = createSignal(false);
  const [signatureHtml, setSignatureHtml] = createSignal("");
  const [editingBirthday, setEditingBirthday] = createSignal(false);
  const [birthdayValue, setBirthdayValue] = createSignal("");
  const [signatureImage, setSignatureImage] = createSignal<File | null>(null);
  const [uploadingImage, setUploadingImage] = createSignal(false);
  const [modalError, setModalError] = createSignal("");
  const [submitting, setSubmitting] = createSignal(false);
  const [existingImageUrl, setExistingImageUrl] = createSignal("");

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
    setEditing((prev: Record<string, string>) => ({ ...prev, [field]: value }));
  }

  function cancelEdit(field: string) {
    setEditing((prev: Record<string, string>) => {
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

  async function saveField(field: string, value: string) {
    const target = user();
    if (!target) return;
    setActionError("");

    const response = await api<AdminUser>(`/council/users/${target.username}`, {
      method: "PATCH",
      token: auth.token(),
      body: { [field]: value },
    });

    if (response.ok) {
      setUser(response.data);
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

  function extractImageUrl(html: string): string {
    const match = html?.match(/<img\s+[^>]*src="([^"]+)"/);
    return match ? match[1] : "";
  }

  async function submitWarn() {
    const target = user();
    if (!target) return;
    setModalError("");
    setSubmitting(true);

    const response = await api<AdminUser>(`/council/users/${target.username}/warn`, {
      method: "POST",
      token: auth.token(),
      body: { title: warnTitle(), message: warnBody() },
    });

    setSubmitting(false);
    if (response.ok) {
      setUser(response.data);
      setModal(null);
      setWarnTitle("");
      setWarnBody("");
    } else {
      setModalError((response.data as unknown as { error: string }).error);
    }
  }

  async function submitDisable() {
    const target = user();
    if (!target) return;
    setModalError("");
    setSubmitting(true);

    const body: Record<string, string> = { reason: disableReason() };
    if (disableUntil()) body.disabled_until = new Date(disableUntil()).toISOString();

    const response = await api<AdminUser>(`/council/users/${target.username}/disable`, {
      method: "POST",
      token: auth.token(),
      body,
    });

    setSubmitting(false);
    if (response.ok) {
      setUser(response.data);
      setModal(null);
      setDisableReason("");
      setDisableUntil("");
    } else {
      setModalError((response.data as unknown as { error: string }).error);
    }
  }

  async function submitBan() {
    const target = user();
    if (!target) return;
    setModalError("");
    setSubmitting(true);

    const response = await api<AdminUser>(`/council/users/${target.username}/ban`, {
      method: "POST",
      token: auth.token(),
      body: { reason: banReason() },
    });

    setSubmitting(false);
    if (response.ok) {
      setUser(response.data);
      setModal(null);
      setBanReason("");
    } else {
      setModalError((response.data as unknown as { error: string }).error);
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

  async function submitGiftJade() {
    const target = user();
    if (!target) return;
    setModalError("");

    const amount = parseInt(jadeAmount());
    if (isNaN(amount) || amount < 1) {
      setModalError("Enter a valid jade amount (minimum 1).");
      return;
    }

    setSubmitting(true);
    const response = await api<AdminUser>(`/council/users/${target.username}`, {
      method: "PATCH",
      token: auth.token(),
      body: { jade: target.jade + amount },
    });

    setSubmitting(false);
    if (response.ok) {
      setUser(response.data);
      setModal(null);
      setJadeAmount("");
    } else {
      setModalError((response.data as unknown as { error: string }).error);
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
              onInput={(e) => setEditing((prev: Record<string, string>) => ({ ...prev, [props.field]: e.currentTarget.value }))}
              onKeyDown={(e) => { if (e.key === "Enter") saveEdit(props.field); if (e.key === "Escape") cancelEdit(props.field); }}
            />
            <button type="button" class="council-detail-edit-btn council-action-save" onClick={() => saveEdit(props.field)}>Save</button>
            <button type="button" class="council-detail-edit-btn" onClick={() => cancelEdit(props.field)}>Cancel</button>
          </span>
        </div>
      </Show>
    );
  }

  function RoleField(props: { role: string }) {
    const [roleOpen, setRoleOpen] = createSignal(false);
    let roleRef: HTMLDivElement | undefined;

    const availableRoles = () => {
      const roles: UserRole[] = [UserRole.Member, UserRole.Moderator];
      if (auth.user()?.role === UserRole.Owner) roles.push(UserRole.Admin);
      return roles;
    };

    onMount(() => {
      function handleClickOutside(e: MouseEvent) {
        if (roleRef && !roleRef.contains(e.target as Node)) setRoleOpen(false);
      }
      document.addEventListener("mousedown", handleClickOutside);
      onCleanup(() => document.removeEventListener("mousedown", handleClickOutside));
    });

    function pickRole(role: string) {
      setEditing((prev: Record<string, string>) => ({ ...prev, role }));
      setRoleOpen(false);
    }

    return (
      <Show when={editing()["role"] !== undefined} fallback={
        <div class="council-detail-row">
          <span class="council-detail-label">Role</span>
          <span class="council-detail-value">
            <span class={`council-role council-role-${props.role}`}>{ROLE_LABELS[props.role] || props.role}</span>
            <Show when={canChangeRole()}>
              <button type="button" class="council-detail-edit-trigger" onClick={() => startEdit("role", props.role)}>Edit</button>
            </Show>
          </span>
        </div>
      }>
        <div class="council-detail-row">
          <span class="council-detail-label">Role</span>
          <span class="council-detail-editable">
            <div class="council-audit-dropdown" ref={roleRef}>
              <button type="button" class="council-audit-dropdown-trigger" onClick={() => setRoleOpen(!roleOpen())}>
                {ROLE_LABELS[editing()["role"]] || editing()["role"]}
              </button>
              <Show when={roleOpen()}>
                <div class="council-audit-dropdown-menu">
                  <For each={availableRoles()}>
                    {(role: UserRole) => (
                      <button type="button" class="council-audit-dropdown-item" classList={{ "council-audit-dropdown-item-selected": editing()["role"] === role }} onClick={() => pickRole(role)}>{ROLE_LABELS[role]}</button>
                    )}
                  </For>
                </div>
              </Show>
            </div>
            <button type="button" class="council-detail-edit-btn council-action-save" onClick={() => saveRole(editing()["role"])}>Save</button>
            <button type="button" class="council-detail-edit-btn" onClick={() => cancelEdit("role")}>Cancel</button>
          </span>
        </div>
      </Show>
    );
  }

  return (
    <StaffGuard>
    <section>
      <A href="/council/users" class="council-detail-back">&larr; Back to Users</A>

      <Show when={error()}>
        <div class="form-error">{error()}</div>
      </Show>

      <Show when={user()}>
        {(target: () => AdminUser) => (
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
                    <Show when={editingBio()} fallback={
                      <div class="council-detail-row">
                        <span class="council-detail-label">Bio</span>
                        <span class="council-detail-value">
                          <Show when={target().bio} fallback={<span>—</span>}>
                            <span class="council-detail-html" innerHTML={target().bio} />
                          </Show>
                          <Show when={isAdmin()}>
                            <button type="button" class="council-detail-edit-trigger" onClick={() => setEditingBio(true)}>Edit</button>
                          </Show>
                        </span>
                      </div>
                    }>
                      <div class="council-detail-row council-detail-row-editor">
                        <span class="council-detail-label">Bio</span>
                        <div class="council-detail-editor-wrap">
                          <MiniEditor onHtml={setBioHtml} initialHtml={target().bio} />
                          <div class="council-detail-editor-actions">
                            <button type="button" class="council-detail-edit-btn council-action-save" onClick={() => { saveField("bio", bioHtml()); setEditingBio(false); }}>Save</button>
                            <button type="button" class="council-detail-edit-btn" onClick={() => setEditingBio(false)}>Cancel</button>
                          </div>
                        </div>
                      </div>
                    </Show>
                    <Show when={editing()["website"] !== undefined} fallback={
                      <div class="council-detail-row">
                        <span class="council-detail-label">Website</span>
                        <span class="council-detail-value">
                          <Show when={target().website} fallback={<span>—</span>}>
                            <a href={target().website} target="_blank" rel="noopener noreferrer">{target().website}</a>
                          </Show>
                          <Show when={isAdmin()}>
                            <button type="button" class="council-detail-edit-trigger" onClick={() => startEdit("website", target().website)}>Edit</button>
                          </Show>
                        </span>
                      </div>
                    }>
                      <div class="council-detail-row">
                        <span class="council-detail-label">Website</span>
                        <span class="council-detail-editable">
                          <input
                            type="text"
                            class="council-detail-edit-input"
                            value={editing()["website"]}
                            onInput={(event) => setEditing((prev: Record<string, string>) => ({ ...prev, website: event.currentTarget.value }))}
                            onKeyDown={(e) => { if (e.key === "Enter") saveEdit("website"); if (e.key === "Escape") cancelEdit("website"); }}
                          />
                          <button type="button" class="council-detail-edit-btn council-action-save" onClick={() => saveEdit("website")}>Save</button>
                          <button type="button" class="council-detail-edit-btn" onClick={() => cancelEdit("website")}>Cancel</button>
                        </span>
                      </div>
                    </Show>
                    <EditableField label="Location" field="location" value={target().location} />
                    <EditableField label="Pronouns" field="pronouns" value={target().pronouns} />
                    <Show when={editingBirthday()} fallback={
                      <div class="council-detail-row">
                        <span class="council-detail-label">Birthday</span>
                        <span class="council-detail-value">
                          <span>{target().birthday ? (() => { const raw = target().birthday!; const parts = raw.includes("T") ? raw.split("T")[0].split("-") : raw.split("-"); const month = parseInt(parts.length === 3 ? parts[1] : parts[0]) - 1; const day = parseInt(parts.length === 3 ? parts[2] : parts[1]); const months = ["January","February","March","April","May","June","July","August","September","October","November","December"]; return `${months[month]} ${day}`; })() : "—"}</span>
                          <Show when={isAdmin()}>
                            <button type="button" class="council-detail-edit-trigger" onClick={() => setEditingBirthday(true)}>Edit</button>
                          </Show>
                        </span>
                      </div>
                    }>
                      <div class="council-detail-row">
                        <span class="council-detail-label">Birthday</span>
                        <span class="council-detail-editable">
                          <MonthDayPicker
                            value={(() => { const raw = target().birthday; if (!raw) return ""; const parts = raw.includes("T") ? raw.split("T")[0].split("-") : raw.split("-"); return parts.length === 3 ? `${parts[1]}-${parts[2]}` : raw; })()}
                            onChange={setBirthdayValue}
                          />
                          <button type="button" class="council-detail-edit-btn council-action-save" onClick={() => { saveField("birthday", birthdayValue()); setEditingBirthday(false); }}>Save</button>
                          <button type="button" class="council-detail-edit-btn" onClick={() => setEditingBirthday(false)}>Cancel</button>
                        </span>
                      </div>
                    </Show>
                    <Show when={editingSignature()} fallback={
                      <div class="council-detail-row">
                        <span class="council-detail-label">Signature</span>
                        <span class="council-detail-value">
                          <Show when={target().signature} fallback={<span>—</span>}>
                            <span class="council-detail-html" innerHTML={target().signature} />
                          </Show>
                          <Show when={isAdmin()}>
                            <button type="button" class="council-detail-edit-trigger" onClick={() => { setExistingImageUrl(extractImageUrl(target().signature)); setEditingSignature(true); }}>Edit</button>
                          </Show>
                        </span>
                      </div>
                    }>
                      <div class="council-detail-row council-detail-row-editor">
                        <span class="council-detail-label">Signature</span>
                        <div class="council-detail-editor-wrap">
                          <MiniEditor onHtml={setSignatureHtml} initialHtml={target().signature} />
                          <div class="council-detail-signature-image">
                            <Show when={signatureImage()}>
                              <img src={URL.createObjectURL(signatureImage()!)} alt="" class="council-detail-image-preview" />
                            </Show>
                            <Show when={!signatureImage() && existingImageUrl()}>
                              <img src={existingImageUrl()} alt="" class="council-detail-image-preview" />
                            </Show>
                            <label class="council-detail-file-label">
                              <span>{signatureImage() ? signatureImage()!.name : existingImageUrl() ? "Replace image" : "Attach image (optional)"}</span>
                              <input
                                type="file"
                                accept="image/*"
                                class="council-detail-file-input"
                                onChange={(e) => setSignatureImage(e.currentTarget.files?.[0] ?? null)}
                              />
                            </label>
                            <Show when={signatureImage() || existingImageUrl()}>
                              <button type="button" class="council-detail-edit-btn" onClick={() => { setSignatureImage(null); setExistingImageUrl(""); }}>Remove</button>
                            </Show>
                          </div>
                          <div class="council-detail-editor-actions">
                            <button type="button" class="council-detail-edit-btn council-action-save" disabled={uploadingImage()} onClick={async () => {
                              let html = signatureHtml();
                              const file = signatureImage();
                              if (file) {
                                setUploadingImage(true);
                                const result = await uploadFile<{ url: string }>("/council/upload", file, auth.token()!);
                                setUploadingImage(false);
                                if (!result.ok) {
                                  setActionError((result.data as unknown as { error: string }).error);
                                  return;
                                }
                                html += `<img src="${result.data.url}" alt="" class="signature-image" />`;
                              } else if (existingImageUrl()) {
                                html += `<img src="${existingImageUrl()}" alt="" class="signature-image" />`;
                              }
                              saveField("signature", html);
                              setEditingSignature(false);
                              setSignatureImage(null);
                              setExistingImageUrl("");
                            }}>{uploadingImage() ? "Uploading..." : "Save"}</button>
                            <button type="button" class="council-detail-edit-btn" onClick={() => { setEditingSignature(false); setSignatureImage(null); setExistingImageUrl(""); }}>Cancel</button>
                          </div>
                        </div>
                      </div>
                    </Show>
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
                    <Show when={target().account_banned}>
                      <div class="council-detail-row">
                        <span class="council-detail-label">Banned At</span>
                        <span>{formatDate(target().banned_at)}</span>
                      </div>
                      <Show when={target().banned_reason}>
                        <div class="council-detail-row">
                          <span class="council-detail-label">Ban Reason</span>
                          <span class="council-detail-html" innerHTML={target().banned_reason} />
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
                          <span class="council-detail-html" innerHTML={target().disabled_reason} />
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
                    <span class={`council-role council-role-${target().role}`}>{ROLE_LABELS[target().role] || target().role}</span>
                  </div>
                </div>

                <Show when={canModerate()}>
                  <div class="council-detail-section">
                    <div class="council-detail-section-header">Actions</div>
                    <div class="council-detail-actions">
                      <button type="button" class="council-detail-action-btn council-action-warn" onClick={() => { setModalError(""); setModal("warn"); }}>
                        Warn
                      </button>
                      <Show when={!target().account_disabled} fallback={
                        <button type="button" class="council-detail-action-btn council-action-enable" onClick={submitEnable}>
                          Enable
                        </button>
                      }>
                        <button type="button" class="council-detail-action-btn council-action-disable" onClick={() => { setModalError(""); setModal("disable"); }}>
                          Disable
                        </button>
                      </Show>
                      <Show when={!target().account_banned} fallback={
                        <button type="button" class="council-detail-action-btn council-action-unban" onClick={submitUnban}>
                          Unban
                        </button>
                      }>
                        <button type="button" class="council-detail-action-btn council-action-ban" onClick={() => { setModalError(""); setModal("ban"); }}>
                          Ban
                        </button>
                      </Show>
                    </div>
                  </div>
                </Show>

                <Show when={isAdmin()}>
                  <div class="council-detail-section">
                    <div class="council-detail-section-header">Gifts</div>
                    <div class="council-detail-actions">
                      <button type="button" class="council-detail-action-btn council-action-jade" onClick={() => { setModalError(""); setModal("jade"); }}>
                        Gift Jade
                      </button>
                      <button type="button" class="council-detail-action-btn council-action-items" disabled>
                        Gift Items
                      </button>
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
                <Show when={modalError()}>
                  <div class="form-error">{modalError()}</div>
                </Show>
                <div class="modal-actions">
                  <button type="button" class="council-detail-action-btn council-action-warn" disabled={submitting()} onClick={submitWarn}>{submitting() ? "Sending..." : "Send Warning"}</button>
                  <button type="button" class="council-detail-action-btn" disabled={submitting()} onClick={() => setModal(null)}>Cancel</button>
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
                  <DatePicker value={disableUntil()} onChange={setDisableUntil} />
                </div>
                <Show when={modalError()}>
                  <div class="form-error">{modalError()}</div>
                </Show>
                <div class="modal-actions">
                  <button type="button" class="council-detail-action-btn council-action-disable" disabled={submitting()} onClick={submitDisable}>{submitting() ? "Disabling..." : "Disable"}</button>
                  <button type="button" class="council-detail-action-btn" disabled={submitting()} onClick={() => setModal(null)}>Cancel</button>
                </div>
              </Modal>
            </Show>

            <Show when={modal() === "ban"}>
              <Modal title="Ban Account" onClose={() => setModal(null)}>
                <div class="modal-field">
                  <label class="modal-label">Reason</label>
                  <Editor onHtml={setBanReason} />
                </div>
                <Show when={modalError()}>
                  <div class="form-error">{modalError()}</div>
                </Show>
                <div class="modal-actions">
                  <button type="button" class="council-detail-action-btn council-action-ban" disabled={submitting()} onClick={submitBan}>{submitting() ? "Banning..." : "Ban"}</button>
                  <button type="button" class="council-detail-action-btn" disabled={submitting()} onClick={() => setModal(null)}>Cancel</button>
                </div>
              </Modal>
            </Show>

            <Show when={modal() === "jade"}>
              <Modal title="Gift Jade" onClose={() => setModal(null)}>
                <div class="modal-field">
                  <label class="modal-label">Amount</label>
                  <input
                    type="number"
                    class="modal-input"
                    placeholder="Enter jade amount..."
                    min="1"
                    value={jadeAmount()}
                    onInput={(e) => setJadeAmount(e.currentTarget.value)}
                  />
                </div>
                <Show when={modalError()}>
                  <div class="form-error">{modalError()}</div>
                </Show>
                <div class="modal-actions">
                  <button type="button" class="council-detail-action-btn council-action-jade" disabled={submitting()} onClick={submitGiftJade}>{submitting() ? "Gifting..." : "Gift Jade"}</button>
                  <button type="button" class="council-detail-action-btn" disabled={submitting()} onClick={() => setModal(null)}>Cancel</button>
                </div>
              </Modal>
            </Show>
          </>
        )}
      </Show>
    </section>
    </StaffGuard>
  );
}