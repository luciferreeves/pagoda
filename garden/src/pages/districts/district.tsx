import { createSignal, onMount, Show, For } from "solid-js";
import { useParams, useSearchParams, A } from "@solidjs/router";
import { api } from "../../api";
import type { District, DistrictSite } from "../../types/district";
import type { PaginatedResponse } from "../../types/admin";
import { districtImage, districtIconClass } from "../../utils/districts";
import { formatDate } from "../../utils/format";
import Pagination from "../../components/Pagination";

export default function DistrictDetail() {
  const params = useParams();
  const [searchParams, setSearchParams] = useSearchParams();
  const [district, setDistrict] = createSignal<District | null>(null);
  const [sites, setSites] = createSignal<DistrictSite[]>([]);
  const [total, setTotal] = createSignal(0);
  const [page, setPage] = createSignal(1);
  const [totalPages, setTotalPages] = createSignal(0);
  const [loading, setLoading] = createSignal(true);
  const [tagInput, setTagInput] = createSignal("");
  const [searchInput, setSearchInput] = createSignal("");

  onMount(async () => {
    const districtResponse = await api<District[]>("/districts");
    if (districtResponse.ok) {
      const found = districtResponse.data.find((entry) => entry.slug === params.slug);
      if (found) setDistrict(found);
    }

    const initialPage = parseInt(searchParams.page as string) || 1;
    const initialTag = (searchParams.tag as string) || "";
    const initialSearch = (searchParams.search as string) || "";
    setTagInput(initialTag);
    setSearchInput(initialSearch);
    loadSites(initialPage, initialTag, initialSearch);
  });

  async function loadSites(pageNumber = 1, tag = "", search = "") {
    setLoading(true);
    const queryParams = new URLSearchParams({
      page: String(pageNumber),
      per_page: "20",
      district: params.slug || "",
    });
    if (tag) queryParams.set("tag", tag);
    if (search) queryParams.set("search", search);

    const response = await api<PaginatedResponse<DistrictSite>>(`/districts/sites?${queryParams}`);
    if (response.ok) {
      setSites(response.data.items);
      setTotal(response.data.total);
      setPage(response.data.page);
      setTotalPages(response.data.total_pages);
    }
    setLoading(false);
  }

  function handleFilter(event: Event) {
    event.preventDefault();
    setSearchParams({ page: "1", tag: tagInput(), search: searchInput() });
    loadSites(1, tagInput(), searchInput());
  }

  function goToPage(pageNumber: number) {
    setSearchParams({ page: String(pageNumber), tag: tagInput(), search: searchInput() });
    loadSites(pageNumber, tagInput(), searchInput());
  }

  return (
    <section>
      <div class="district-header">
        <Show when={district()}>
          {(currentDistrict) => (
            <>
              <div class="district-header-info">
                <A href="/districts" class="district-back">Districts</A>
                <h1 class="page-title" style={{ color: currentDistrict().foreground }}>{currentDistrict().name}</h1>
                <p class="district-description">{currentDistrict().description}</p>
              </div>
              <div class="district-header-icon">
                <img src={districtImage(currentDistrict().slug)} alt={currentDistrict().name} class={districtIconClass(currentDistrict().slug)} />
              </div>
            </>
          )}
        </Show>
      </div>

      <form class="district-filters" onSubmit={handleFilter}>
        <input
          type="text"
          placeholder="Search sites..."
          value={searchInput()}
          onInput={(e) => setSearchInput(e.currentTarget.value)}
        />
        <input
          type="text"
          placeholder="Filter by tag..."
          value={tagInput()}
          onInput={(e) => setTagInput(e.currentTarget.value)}
        />
        <button type="submit" class="form-button">Filter</button>
      </form>

      <Show when={!loading()} fallback={<p class="loading-text">Loading sites...</p>}>
        <Show when={sites().length > 0} fallback={<p class="empty-text">No sites found in this district.</p>}>
          <div class="site-grid">
            <For each={sites()}>
              {(site) => (
                <a href={site.url} target="_blank" rel="noopener noreferrer" class="site-card">
                  <div class="site-card-thumbnail">
                    <Show when={site.thumbnail_url} fallback={<div class="site-card-placeholder" />}>
                      <img src={site.thumbnail_url} alt={site.title} />
                    </Show>
                  </div>
                  <div class="site-card-info">
                    <h3 class="site-card-title">{site.title}</h3>
                    <p class="site-card-url">{site.url}</p>
                    <Show when={site.tags.length > 0}>
                      <div class="site-card-tags">
                        <For each={site.tags}>
                          {(tag) => <span class="site-tag">{tag}</span>}
                        </For>
                      </div>
                    </Show>
                    <div class="site-card-meta">
                      <span>by {site.submitter.display_name}</span>
                      <span>{formatDate(site.created_at)}</span>
                    </div>
                  </div>
                </a>
              )}
            </For>
          </div>

          <Pagination
            page={page()}
            totalPages={totalPages()}
            total={total()}
            label="sites"
            onPage={goToPage}
          />
        </Show>
      </Show>
    </section>
  );
}