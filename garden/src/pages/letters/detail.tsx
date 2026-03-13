import { createSignal, onMount, Show, For } from "solid-js";
import { useParams, useNavigate } from "@solidjs/router";
import { api } from "../../api";
import { auth } from "../../store/auth";
import { extractError } from "../../utils/api";
import { formatDateTime } from "../../utils/format";
import type { LetterDetail, LetterMessage } from "../../types/letter";
import Modal from "../../components/Modal";

export default function LetterDetailPage() {
  const params = useParams();
  const navigate = useNavigate();

  const [letter, setLetter] = createSignal<LetterDetail | null>(null);
  const [loading, setLoading] = createSignal(true);
  const [error, setError] = createSignal("");

  const [replyBody, setReplyBody] = createSignal("");
  const [replyError, setReplyError] = createSignal("");
  const [replying, setReplying] = createSignal(false);

  const [editingRef, setEditingRef] = createSignal("");
  const [editBody, setEditBody] = createSignal("");
  const [editError, setEditError] = createSignal("");
  const [saving, setSaving] = createSignal(false);

  const [deleteTarget, setDeleteTarget] = createSignal<LetterMessage | null>(null);
  const [deleting, setDeleting] = createSignal(false);

  const [renaming, setRenaming] = createSignal(false);
  const [renameTitle, setRenameTitle] = createSignal("");
  const [renameError, setRenameError] = createSignal("");
  const [renameSaving, setRenameSaving] = createSignal(false);

  const [leaving, setLeaving] = createSignal(false);
  const [leaveError, setLeaveError] = createSignal("");

  let messagesEndRef: HTMLDivElement | undefined;

  onMount(() => {
    if (!auth.token()) {
      navigate("/login");
      return;
    }
    loadLetter();
  });

  async function loadLetter() {
    setLoading(true);
    const response = await api<LetterDetail>(`/letters/${params.ref}`, {
      token: auth.token(),
    });

    if (response.ok) {
      setLetter(response.data);
      setTimeout(() => messagesEndRef?.scrollIntoView({ behavior: "smooth" }), 100);
    } else {
      setError(extractError(response.data));
    }
    setLoading(false);
  }

  async function sendReply() {
    if (!replyBody().trim()) return;
    setReplying(true);
    setReplyError("");

    const response = await api<LetterMessage>(`/letters/${params.ref}/messages`, {
      method: "POST",
      token: auth.token(),
      body: { body: replyBody().trim(), attachment_refs: [] },
    });

    if (response.ok) {
      setReplyBody("");
      loadLetter();
    } else {
      setReplyError(extractError(response.data));
    }
    setReplying(false);
  }

  function startEdit(message: LetterMessage) {
    setEditingRef(message.ref);
    setEditBody(message.body.replace(/<[^>]*>/g, ""));
    setEditError("");
  }

  async function submitEdit() {
    if (!editBody().trim()) return;
    setSaving(true);
    setEditError("");

    const response = await api<LetterMessage>(`/letters/${params.ref}/messages/${editingRef()}`, {
      method: "PATCH",
      token: auth.token(),
      body: { body: editBody().trim() },
    });

    if (response.ok) {
      setEditingRef("");
      loadLetter();
    } else {
      setEditError(extractError(response.data));
    }
    setSaving(false);
  }

  async function confirmDelete() {
    const target = deleteTarget();
    if (!target) return;
    setDeleting(true);

    const response = await api(`/letters/${params.ref}/messages/${target.ref}`, {
      method: "DELETE",
      token: auth.token(),
    });

    if (response.ok) {
      setDeleteTarget(null);
      loadLetter();
    }
    setDeleting(false);
  }

  function openRename() {
    setRenameTitle(letter()?.title || "");
    setRenameError("");
    setRenaming(true);
  }

  async function submitRename() {
    setRenameSaving(true);
    setRenameError("");

    const response = await api(`/letters/${params.ref}`, {
      method: "PATCH",
      token: auth.token(),
      body: { title: renameTitle().trim() },
    });

    if (response.ok) {
      setRenaming(false);
      loadLetter();
    } else {
      setRenameError(extractError(response.data));
    }
    setRenameSaving(false);
  }

  async function leaveLetter() {
    setLeaving(true);
    setLeaveError("");

    const response = await api(`/letters/${params.ref}/leave`, {
      method: "POST",
      token: auth.token(),
    });

    if (response.ok) {
      navigate("/letters");
    } else {
      setLeaveError(extractError(response.data));
    }
    setLeaving(false);
  }

  function isOwnMessage(message: LetterMessage) {
    return message.sender?.username === auth.user()?.username;
  }

  function handleReplyKeyDown(e: KeyboardEvent) {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      sendReply();
    }
  }

  return (
    <section>
      <Show when={!loading()} fallback={<div class="letters-empty">Loading...</div>}>
        <Show when={!error()} fallback={<div class="form-error">{error()}</div>}>
          <Show when={letter()}>
            {(l) => (
              <>
                <div class="letter-detail-header">
                  <button type="button" class="letter-back" onClick={() => navigate("/letters")}>&larr; Back</button>
                  <h2 class="letter-detail-title">{l().title}</h2>
                  <div class="letter-detail-actions">
                    <Show when={!l().is_system}>
                      <button type="button" class="action-btn" onClick={openRename}>Rename</button>
                      <button type="button" class="action-btn deny" onClick={() => setLeaveError("confirm")}>Leave</button>
                    </Show>
                  </div>
                </div>

                <div class="letter-participants">
                  <For each={l().participants}>
                    {(p) => (
                      <span class="letter-participant">
                        <img src={p.avatar_url} alt="" class="letter-participant-avatar" />
                        {p.display_name}
                      </span>
                    )}
                  </For>
                </div>

                <div class="letter-messages">
                  <For each={l().messages}>
                    {(message) => (
                      <div class="letter-message" classList={{ "letter-message-own": isOwnMessage(message), "letter-message-deleted": message.deleted }}>
                        <Show when={!message.deleted} fallback={
                          <div class="letter-message-body letter-message-deleted-text">This message was deleted.</div>
                        }>
                          <div class="letter-message-sender">
                            <Show when={message.sender}>
                              <img src={message.sender!.avatar_url} alt="" class="letter-message-avatar" />
                              <span class="letter-message-name">{message.sender!.display_name}</span>
                            </Show>
                            <span class="letter-message-time">
                              {formatDateTime(message.created_at)}
                              <Show when={message.edited_at}> (edited)</Show>
                            </span>
                          </div>

                          <Show when={editingRef() === message.ref} fallback={
                            <div class="letter-message-body" innerHTML={message.body} />
                          }>
                            <div class="letter-message-edit">
                              <textarea
                                value={editBody()}
                                onInput={(e) => setEditBody(e.currentTarget.value)}
                                rows={3}
                              />
                              <Show when={editError()}>
                                <div class="form-error">{editError()}</div>
                              </Show>
                              <div class="letter-message-edit-actions">
                                <button type="button" class="action-btn approve" onClick={submitEdit} disabled={saving()}>Save</button>
                                <button type="button" class="action-btn" onClick={() => setEditingRef("")}>Cancel</button>
                              </div>
                            </div>
                          </Show>

                          <Show when={message.attachments.length > 0}>
                            <div class="letter-attachments">
                              <For each={message.attachments}>
                                {(att) => (
                                  <a href={att.url} target="_blank" rel="noopener noreferrer" class="letter-attachment">
                                    <Show when={att.category === "image"} fallback={
                                      <span class="letter-attachment-file">{att.file_name}</span>
                                    }>
                                      <img src={att.url} alt={att.file_name} class="letter-attachment-image" />
                                    </Show>
                                  </a>
                                )}
                              </For>
                            </div>
                          </Show>

                          <Show when={isOwnMessage(message) && editingRef() !== message.ref}>
                            <div class="letter-message-actions">
                              <button type="button" class="letter-msg-action" onClick={() => startEdit(message)}>Edit</button>
                              <button type="button" class="letter-msg-action letter-msg-action-delete" onClick={() => setDeleteTarget(message)}>Delete</button>
                            </div>
                          </Show>
                        </Show>
                      </div>
                    )}
                  </For>
                  <div ref={messagesEndRef} />
                </div>

                <Show when={!l().is_system}>
                  <div class="letter-reply">
                    <textarea
                      placeholder="Write a reply..."
                      value={replyBody()}
                      onInput={(e) => setReplyBody(e.currentTarget.value)}
                      onKeyDown={handleReplyKeyDown}
                      rows={3}
                    />
                    <Show when={replyError()}>
                      <div class="form-error">{replyError()}</div>
                    </Show>
                    <button type="button" class="form-button" onClick={sendReply} disabled={replying()}>
                      {replying() ? "Sending..." : "Send"}
                    </button>
                  </div>
                </Show>
              </>
            )}
          </Show>
        </Show>
      </Show>

      <Show when={deleteTarget()}>
        <Modal title="Delete Message" onClose={() => setDeleteTarget(null)}>
          <p>Are you sure you want to delete this message?</p>
          <div class="modal-actions">
            <button type="button" class="form-button" onClick={confirmDelete} disabled={deleting()}>
              {deleting() ? "Deleting..." : "Delete"}
            </button>
            <button type="button" class="form-button secondary" onClick={() => setDeleteTarget(null)}>Cancel</button>
          </div>
        </Modal>
      </Show>

      <Show when={renaming()}>
        <Modal title="Rename Conversation" onClose={() => setRenaming(false)}>
          <div class="form-field">
            <label>Title</label>
            <input
              type="text"
              value={renameTitle()}
              onInput={(e) => setRenameTitle(e.currentTarget.value)}
              maxLength={200}
            />
          </div>
          <Show when={renameError()}>
            <div class="form-error">{renameError()}</div>
          </Show>
          <div class="modal-actions">
            <button type="button" class="form-button" onClick={submitRename} disabled={renameSaving()}>
              {renameSaving() ? "Saving..." : "Save"}
            </button>
            <button type="button" class="form-button secondary" onClick={() => setRenaming(false)}>Cancel</button>
          </div>
        </Modal>
      </Show>

      <Show when={leaveError() === "confirm"}>
        <Modal title="Leave Conversation" onClose={() => setLeaveError("")}>
          <p>Are you sure you want to leave this conversation? You will no longer receive messages.</p>
          <Show when={leaveError() && leaveError() !== "confirm"}>
            <div class="form-error">{leaveError()}</div>
          </Show>
          <div class="modal-actions">
            <button type="button" class="form-button" onClick={leaveLetter} disabled={leaving()}>
              {leaving() ? "Leaving..." : "Leave"}
            </button>
            <button type="button" class="form-button secondary" onClick={() => setLeaveError("")}>Cancel</button>
          </div>
        </Modal>
      </Show>
    </section>
  );
}