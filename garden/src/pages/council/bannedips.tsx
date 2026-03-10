import { createSignal, onMount, Show, For } from "solid-js";
import { api } from "../../api";
import { auth } from "../../store/auth";
import StaffGuard from "../../components/StaffGuard";
import Modal from "../../components/Modal";

interface IPBan {
  ID: number;
  IP: string;
  Reason: string;
  CreatedAt: string;
}

interface PaginatedResponse {
  items: IPBan[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

export default function BannedIPs() {
  const [bans, setBans] = createSignal<IPBan[]>([]);
  const [page, setPage] = createSignal(1);
  const [totalPages, setTotalPages] = createSignal(1);
  const [total, setTotal] = createSignal(0);
  const [loading, setLoading] = createSignal(true);
  const [confirmBan, setConfirmBan] = createSignal<IPBan | null>(null);

  async function loadBans(p = 1) {
    setLoading(true);
    const response = await api<PaginatedResponse>(`/council/bannedips?page=${p}`, {
      token: auth.token(),
    });
    if (response.ok) {
      setBans(response.data.items ?? []);
      setPage(response.data.page);
      setTotalPages(response.data.total_pages);
      setTotal(response.data.total);
    }
    setLoading(false);
  }

  async function liftBan(ban: IPBan) {
    const response = await api(`/council/bannedips/${ban.ID}`, {
      method: "DELETE",
      token: auth.token(),
    });
    if (response.ok) {
      setConfirmBan(null);
      loadBans(page());
    }
  }

  onMount(() => loadBans());

  function formatDate(date: string) {
    return new Date(date).toLocaleDateString();
  }

  return (
    <StaffGuard>
      <section>
        <h2 class="page-title">Banned IPs</h2>

        <div class="council-grid council-grid-bannedips">
          <div class="council-grid-header">
            <span>IP Address</span>
            <span>Reason</span>
            <span>Banned At</span>
            <span></span>
          </div>
          <Show when={!loading()} fallback={
            <div class="council-grid-empty">Loading...</div>
          }>
            <Show when={bans().length} fallback={
              <div class="council-grid-empty">No banned IPs.</div>
            }>
              <For each={bans()}>
                {(ban: IPBan) => (
                  <div class="council-grid-row">
                    <span class="council-ip">{ban.IP}</span>
                    <span>{ban.Reason}</span>
                    <span>{formatDate(ban.CreatedAt)}</span>
                    <span>
                      <button type="button" class="council-detail-action-btn council-action-unban" onClick={() => setConfirmBan(ban)}>Lift</button>
                    </span>
                  </div>
                )}
              </For>
            </Show>
          </Show>
        </div>

        <Show when={totalPages() > 1}>
          <div class="council-pagination">
            <button class="council-page-btn" disabled={page() <= 1} onClick={() => loadBans(page() - 1)}>Prev</button>
            <span class="council-page-info">Page {page()} of {totalPages()} ({total()} bans)</span>
            <button class="council-page-btn" disabled={page() >= totalPages()} onClick={() => loadBans(page() + 1)}>Next</button>
          </div>
        </Show>
        <Show when={confirmBan()}>
          {(ban: () => IPBan) => (
            <Modal title="Lift IP Ban" onClose={() => setConfirmBan(null)}>
              <p style={{ "font-size": "12px", margin: "0 0 12px" }}>Lift ban on <strong>{ban().IP}</strong>?</p>
              <div class="modal-actions">
                <button type="button" class="council-detail-action-btn council-action-unban" onClick={() => liftBan(ban())}>Lift Ban</button>
                <button type="button" class="council-detail-action-btn" onClick={() => setConfirmBan(null)}>Cancel</button>
              </div>
            </Modal>
          )}
        </Show>
      </section>
    </StaffGuard>
  );
}