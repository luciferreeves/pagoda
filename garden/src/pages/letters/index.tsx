import { createSignal, onMount, Show, For } from "solid-js";
import { A, useNavigate, useSearchParams } from "@solidjs/router";
import { api } from "../../api";
import { auth } from "../../store/auth";
import { formatDateTime } from "../../utils/format";
import type { Letter } from "../../types/letter";
import type { PaginatedResponse } from "../../types/admin";
import Pagination from "../../components/Pagination";

export default function Letters() {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();

  const [letters, setLetters] = createSignal<Letter[]>([]);
  const [total, setTotal] = createSignal(0);
  const [page, setPage] = createSignal(1);
  const [totalPages, setTotalPages] = createSignal(0);
  const [loading, setLoading] = createSignal(true);

  onMount(() => {
    if (!auth.token()) {
      navigate("/login");
      return;
    }
    const p = parseInt(searchParams.page as string) || 1;
    loadLetters(p);
  });

  async function loadLetters(pageNumber = 1) {
    setLoading(true);
    const response = await api<PaginatedResponse<Letter>>(`/letters?page=${pageNumber}&per_page=20`, {
      token: auth.token(),
    });

    if (response.ok) {
      setLetters(response.data.items);
      setTotal(response.data.total);
      setPage(response.data.page);
      setTotalPages(response.data.total_pages);
    }
    setLoading(false);
  }

  function goToPage(p: number) {
    setSearchParams({ page: String(p) });
    loadLetters(p);
  }

  function previewBody(message?: Letter["last_message"]) {
    if (!message) return "No messages yet";
    if (message.deleted) return "Message deleted";
    const text = message.body.replace(/<[^>]*>/g, "");
    return text.length > 80 ? text.slice(0, 80) + "..." : text;
  }

  return (
    <section>
      <h2 class="page-title">Letters</h2>

      <div class="letters-list">
        <Show when={!loading()} fallback={
          <div class="letters-empty">Loading...</div>
        }>
          <Show when={letters().length} fallback={
            <div class="letters-empty">No letters yet.</div>
          }>
            <For each={letters()}>
              {(letter) => (
                <A href={`/letters/${letter.ref}`} class="letter-item" classList={{ "letter-unread": letter.unread }}>
                  <div class="letter-avatars">
                    <For each={letter.participants.slice(0, 3)}>
                      {(p) => (
                        <img src={p.avatar_url} alt="" class="letter-avatar" />
                      )}
                    </For>
                    <Show when={letter.participants.length > 3}>
                      <span class="letter-avatar-more">+{letter.participants.length - 3}</span>
                    </Show>
                  </div>
                  <div class="letter-content">
                    <div class="letter-top">
                      <span class="letter-title">
                        <Show when={letter.is_system}>
                          <span class="letter-system-badge">System</span>
                        </Show>
                        {letter.title}
                      </span>
                      <span class="letter-time">{formatDateTime(letter.updated_at)}</span>
                    </div>
                    <div class="letter-preview">{previewBody(letter.last_message)}</div>
                  </div>
                  <Show when={letter.unread}>
                    <span class="letter-unread-dot" />
                  </Show>
                </A>
              )}
            </For>
          </Show>
        </Show>
      </div>

      <Pagination page={page()} totalPages={totalPages()} total={total()} label="letters" onPage={goToPage} />
    </section>
  );
}