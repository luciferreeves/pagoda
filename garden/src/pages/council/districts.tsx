import { createSignal, onMount, Show, For } from "solid-js";
import { useSearchParams } from "@solidjs/router";
import { api } from "../../api";
import { auth } from "../../store/auth";
import { extractError } from "../../utils/api";
import { useClickOutside } from "../../utils/clickOutside";
import type { SiteRequest, AdminSite } from "../../types/district";
import type { PaginatedResponse } from "../../types/admin";
import Pagination from "../../components/Pagination";
import StaffGuard from "../../components/StaffGuard";
import Modal from "../../components/Modal";

const STATUS_LABELS: Record<string, string> = {
  "": "Pending & Hold",
  pending: "Pending",
  hold: "On Hold",
  approved: "Approved",
  denied: "Denied",
};

export default function CouncilDistricts() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [tab, setTab] = createSignal<"requests" | "sites">("requests");

  const [requests, setRequests] = createSignal<SiteRequest[]>([]);
  const [reqTotal, setReqTotal] = createSignal(0);
  const [reqPage, setReqPage] = createSignal(1);
  const [reqTotalPages, setReqTotalPages] = createSignal(0);
  const [reqLoading, setReqLoading] = createSignal(false);
  const [reqStatus, setReqStatus] = createSignal("");
  const [statusOpen, setStatusOpen] = createSignal(false);
  let statusRef: HTMLDivElement | undefined;
  useClickOutside(() => statusRef, setStatusOpen);

  const [adminSites, setAdminSites] = createSignal<AdminSite[]>([]);
  const [siteTotal, setSiteTotal] = createSignal(0);
  const [sitePage, setSitePage] = createSignal(1);
  const [siteTotalPages, setSiteTotalPages] = createSignal(0);
  const [siteLoading, setSiteLoading] = createSignal(false);
  const [siteSearch, setSiteSearch] = createSignal("");

  const [reviewTarget, setReviewTarget] = createSignal<SiteRequest | null>(null);
  const [reviewAction, setReviewAction] = createSignal("");
  const [reviewError, setReviewError] = createSignal("");
  const [reviewing, setReviewing] = createSignal(false);

  const [editTarget, setEditTarget] = createSignal<AdminSite | null>(null);
  const [editTitle, setEditTitle] = createSignal("");
  const [editDescription, setEditDescription] = createSignal("");
  const [editError, setEditError] = createSignal("");
  const [editing, setEditing] = createSignal(false);

  onMount(() => {
    const initialTab = searchParams.tab as string;
    if (initialTab === "sites") setTab("sites");
    loadRequests();
  });

  async function loadRequests(pageNumber = 1) {
    setReqLoading(true);
    const queryParams = new URLSearchParams({ page: String(pageNumber), per_page: "20" });
    if (reqStatus()) queryParams.set("status", reqStatus());

    const response = await api<PaginatedResponse<SiteRequest>>(`/council/districts/requests?${queryParams}`, {
      token: auth.token(),
    });

    if (response.ok) {
      setRequests(response.data.items);
      setReqTotal(response.data.total);
      setReqPage(response.data.page);
      setReqTotalPages(response.data.total_pages);
    }
    setReqLoading(false);
  }

  async function loadSites(pageNumber = 1) {
    setSiteLoading(true);
    const queryParams = new URLSearchParams({ page: String(pageNumber), per_page: "20" });
    if (siteSearch()) queryParams.set("search", siteSearch());

    const response = await api<PaginatedResponse<AdminSite>>(`/council/districts/sites?${queryParams}`, {
      token: auth.token(),
    });

    if (response.ok) {
      setAdminSites(response.data.items);
      setSiteTotal(response.data.total);
      setSitePage(response.data.page);
      setSiteTotalPages(response.data.total_pages);
    }
    setSiteLoading(false);
  }

  function switchTab(selectedTab: "requests" | "sites") {
    setTab(selectedTab);
    setSearchParams({ tab: selectedTab });
    if (selectedTab === "requests") loadRequests();
    else loadSites();
  }

  function pickStatus(value: string) {
    setReqStatus(value);
    setStatusOpen(false);
    loadRequests(1);
  }

  function openReview(site: SiteRequest, action: string) {
    setReviewTarget(site);
    setReviewAction(action);
    setReviewError("");
  }

  async function submitReview() {
    const target = reviewTarget();
    if (!target) return;
    setReviewing(true);
    setReviewError("");

    const response = await api<SiteRequest>(`/council/districts/sites/${target.ref}/review`, {
      method: "POST",
      token: auth.token(),
      body: { status: reviewAction() },
    });

    if (response.ok) {
      setReviewTarget(null);
      loadRequests(reqPage());
    } else {
      setReviewError(extractError(response.data));
    }
    setReviewing(false);
  }

  function openEdit(site: AdminSite) {
    setEditTarget(site);
    setEditTitle(site.title);
    setEditDescription(site.description);
    setEditError("");
  }

  async function submitEdit() {
    const target = editTarget();
    if (!target) return;
    setEditing(true);
    setEditError("");

    const response = await api<AdminSite>(`/council/districts/sites/${target.ref}`, {
      method: "PATCH",
      token: auth.token(),
      body: {
        title: editTitle(),
        description: editDescription(),
      },
    });

    if (response.ok) {
      setEditTarget(null);
      loadSites(sitePage());
    } else {
      setEditError(extractError(response.data));
    }
    setEditing(false);
  }

  function statusLabel(status: string) {
    return STATUS_LABELS[status] || status;
  }

  return (
    <StaffGuard>
      <section>
        <h2 class="page-title">Districts</h2>

        <div class="council-tabs">
          <button
            type="button"
            class={`council-tab ${tab() === "requests" ? "active" : ""}`}
            onClick={() => switchTab("requests")}
          >
            Requests
          </button>
          <button
            type="button"
            class={`council-tab ${tab() === "sites" ? "active" : ""}`}
            onClick={() => switchTab("sites")}
          >
            Sites
          </button>
        </div>

        <Show when={tab() === "requests"}>
          <div class="council-audit-filters">
            <div class="council-audit-dropdown" ref={statusRef}>
              <button type="button" class="council-audit-dropdown-trigger" onClick={() => setStatusOpen(!statusOpen())}>
                {statusLabel(reqStatus())}
              </button>
              <Show when={statusOpen()}>
                <div class="council-audit-dropdown-menu">
                  <For each={Object.entries(STATUS_LABELS)}>
                    {([key, label]: [string, string]) => (
                      <button
                        type="button"
                        class="council-audit-dropdown-item"
                        classList={{ "council-audit-dropdown-item-selected": reqStatus() === key }}
                        onClick={() => pickStatus(key)}
                      >
                        {label}
                      </button>
                    )}
                  </For>
                </div>
              </Show>
            </div>
            <Show when={reqStatus()}>
              <button type="button" class="council-audit-clear-btn" onClick={() => pickStatus("")}>Clear</button>
            </Show>
          </div>

          <div class="council-grid district-req-grid">
            <div class="council-grid-header">
              <span>Title</span>
              <span>District</span>
              <span>Submitter</span>
              <span>Status</span>
              <span>Actions</span>
            </div>
            <Show when={!reqLoading()} fallback={
              <div class="council-grid-empty">Loading...</div>
            }>
              <Show when={requests().length} fallback={
                <div class="council-grid-empty">No requests found.</div>
              }>
                <For each={requests()}>
                  {(request) => (
                    <div class="council-grid-row">
                      <span>
                        <a href={request.url} target="_blank" rel="noopener noreferrer">{request.title}</a>
                      </span>
                      <span>{request.district}</span>
                      <span>{request.submitter.display_name}</span>
                      <span class={`status-badge status-${request.status}`}>{statusLabel(request.status)}</span>
                      <span class="council-actions">
                        <Show when={request.status === "pending" || request.status === "hold"}>
                          <button type="button" class="action-btn approve" onClick={() => openReview(request, "approved")}>Approve</button>
                          <button type="button" class="action-btn deny" onClick={() => openReview(request, "denied")}>Deny</button>
                          <Show when={request.status === "pending"}>
                            <button type="button" class="action-btn hold" onClick={() => openReview(request, "hold")}>Hold</button>
                          </Show>
                        </Show>
                      </span>
                    </div>
                  )}
                </For>
              </Show>
            </Show>
          </div>

          <Pagination page={reqPage()} totalPages={reqTotalPages()} total={reqTotal()} label="requests" onPage={(pageNumber) => loadRequests(pageNumber)} />
        </Show>

        <Show when={tab() === "sites"}>
          <div class="council-search">
            <input
              type="text"
              placeholder="Search sites..."
              value={siteSearch()}
              onInput={(e) => setSiteSearch(e.currentTarget.value)}
              onKeyDown={(e) => { if (e.key === "Enter") loadSites(1); }}
            />
            <button type="button" class="form-button" onClick={() => loadSites(1)}>Search</button>
          </div>

          <div class="council-grid district-site-grid">
            <div class="council-grid-header">
              <span>Title</span>
              <span>District</span>
              <span>URL</span>
              <span>Actions</span>
            </div>
            <Show when={!siteLoading()} fallback={
              <div class="council-grid-empty">Loading...</div>
            }>
              <Show when={adminSites().length} fallback={
                <div class="council-grid-empty">No approved sites found.</div>
              }>
                <For each={adminSites()}>
                  {(site) => (
                    <div class="council-grid-row">
                      <span>{site.title}</span>
                      <span>{site.district}</span>
                      <span>
                        <a href={site.url} target="_blank" rel="noopener noreferrer">{site.url}</a>
                      </span>
                      <span>
                        <button type="button" class="action-btn" onClick={() => openEdit(site)}>Edit</button>
                      </span>
                    </div>
                  )}
                </For>
              </Show>
            </Show>
          </div>

          <Pagination page={sitePage()} totalPages={siteTotalPages()} total={siteTotal()} label="sites" onPage={(pageNumber) => loadSites(pageNumber)} />
        </Show>

        <Show when={reviewTarget()}>
          <Modal title={`${reviewAction() === "approved" ? "Approve" : reviewAction() === "denied" ? "Deny" : "Hold"} Site`} onClose={() => setReviewTarget(null)}>
            <p>
              Are you sure you want to {reviewAction() === "approved" ? "approve" : reviewAction() === "denied" ? "deny" : "put on hold"}{" "}
              <strong>{reviewTarget()!.title}</strong>?
            </p>
            <Show when={reviewError()}>
              <div class="form-error">{reviewError()}</div>
            </Show>
            <div class="modal-actions">
              <button type="button" class="form-button" onClick={submitReview} disabled={reviewing()}>
                {reviewing() ? "Processing..." : "Confirm"}
              </button>
              <button type="button" class="form-button secondary" onClick={() => setReviewTarget(null)}>Cancel</button>
            </div>
          </Modal>
        </Show>

        <Show when={editTarget()}>
          <Modal title="Edit Site" onClose={() => setEditTarget(null)}>
            <Show when={editError()}>
              <div class="form-error">{editError()}</div>
            </Show>
            <div class="form-field">
              <label>Title</label>
              <input type="text" value={editTitle()} onInput={(e) => setEditTitle(e.currentTarget.value)} maxLength={200} />
            </div>
            <div class="form-field">
              <label>Description</label>
              <textarea value={editDescription()} onInput={(e) => setEditDescription(e.currentTarget.value)} maxLength={1000} rows={3} />
            </div>
            <div class="modal-actions">
              <button type="button" class="form-button" onClick={submitEdit} disabled={editing()}>
                {editing() ? "Saving..." : "Save"}
              </button>
              <button type="button" class="form-button secondary" onClick={() => setEditTarget(null)}>Cancel</button>
            </div>
          </Modal>
        </Show>
      </section>
    </StaffGuard>
  );
}