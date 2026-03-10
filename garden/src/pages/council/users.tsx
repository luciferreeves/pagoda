import { createSignal, onMount, Show, For } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { council } from "../../store/council";
import type { AdminUser } from "../../types/admin";
import StaffGuard from "../../components/StaffGuard";

export default function CouncilUsers() {
  const navigate = useNavigate();
  const [searchInput, setSearchInput] = createSignal("");

  onMount(() => council.loadUsers());

  function handleSearch(event: Event) {
    event.preventDefault();
    council.setSearch(searchInput());
    council.loadUsers(1, searchInput());
  }

  function goToPage(p: number) {
    council.loadUsers(p, council.search());
  }

  function statusBadge(user: AdminUser) {
    if (user.account_banned) return "banned";
    if (user.account_disabled) return "disabled";
    if (!user.email_verified) return "unverified";
    return "active";
  }

  function formatDate(date: string) {
    return new Date(date).toLocaleDateString();
  }

  function sortIndicator(field: string) {
    if (council.sortField() !== field) return "";
    return council.sortOrder() === "asc" ? " \u25B2" : " \u25BC";
  }

  return (
    <StaffGuard>
    <section>
      <h2 class="page-title">Users</h2>

      <form class="council-search" onSubmit={handleSearch}>
        <input
          type="text"
          placeholder="Search by username, display name, or email..."
          value={searchInput()}
          onInput={(e) => setSearchInput(e.currentTarget.value)}
        />
        <button type="submit" class="form-button">Search</button>
      </form>

      <div class="council-grid">
        <div class="council-grid-header">
          <span class="council-sortable" onClick={() => council.toggleSort("display_name")}>User{sortIndicator("display_name")}</span>
          <span class="council-sortable" onClick={() => council.toggleSort("email")}>Email{sortIndicator("email")}</span>
          <span class="council-sortable" onClick={() => council.toggleSort("role")}>Role{sortIndicator("role")}</span>
          <span>Status</span>
          <span class="council-sortable" onClick={() => council.toggleSort("created_at")}>Joined{sortIndicator("created_at")}</span>
        </div>
        <Show when={!council.loading()} fallback={
          <div class="council-grid-empty">Loading...</div>
        }>
          <Show when={council.users().length} fallback={
            <div class="council-grid-empty">No users found.</div>
          }>
            <For each={council.users()}>
              {(user: AdminUser) => (
                <div class="council-grid-row" onClick={() => navigate(`/council/users/${user.username}`)}>
                  <span class="council-user-cell">
                    <img src={user.avatar_url} alt="" class="council-avatar" />
                    <span>
                      <span class="council-display-name">{user.display_name}</span>
                      <span class="council-username">@{user.username}</span>
                    </span>
                  </span>
                  <span>{user.email}</span>
                  <span>
                    <span class={`council-role council-role-${user.role}`}>{user.role}</span>
                  </span>
                  <span>
                    <span class={`council-status council-status-${statusBadge(user)}`}>
                      {statusBadge(user)}
                    </span>
                  </span>
                  <span>{formatDate(user.created_at)}</span>
                </div>
              )}
            </For>
          </Show>
        </Show>
      </div>

      <Show when={council.totalPages() > 1}>
        <div class="council-pagination">
          <button
            class="council-page-btn"
            disabled={council.page() <= 1}
            onClick={() => goToPage(council.page() - 1)}
          >
            Prev
          </button>
          <span class="council-page-info">
            Page {council.page()} of {council.totalPages()} ({council.total()} users)
          </span>
          <button
            class="council-page-btn"
            disabled={council.page() >= council.totalPages()}
            onClick={() => goToPage(council.page() + 1)}
          >
            Next
          </button>
        </div>
      </Show>
    </section>
    </StaffGuard>
  );
}