import { createSignal, onMount, Show, For } from "solid-js";
import { A } from "@solidjs/router";
import { api } from "../../api";
import { auth } from "../../store/auth";
import type { District } from "../../types/district";
import { districtImage, districtIconClass } from "../../utils/districts";

export default function Districts() {
  const [districts, setDistricts] = createSignal<District[]>([]);
  const [loading, setLoading] = createSignal(true);

  onMount(async () => {
    const response = await api<District[]>("/districts");
    if (response.ok) {
      setDistricts(response.data);
    }
    setLoading(false);
  });

  return (
    <section>
      <h1 class="page-title">Districts</h1>
      <p class="district-intro">
        Districts are themed directories of websites from across the small web.
        Browse sites by category, or <Show when={auth.user()} fallback={<span>log in to submit your own</span>}><A href="/districts/submit">submit your own</A></Show> to be listed.
      </p>

      <Show when={!loading()} fallback={<p class="loading-text">Loading districts...</p>}>
        <div class="district-grid">
          <For each={districts()}>
            {(district) => (
              <A
                href={`/districts/${district.slug}`}
                class="district-card"
                style={{
                  "background-color": district.background,
                  "border-color": district.foreground,
                }}
              >
                <div class="district-card-info">
                  <h2 class="district-card-name" style={{ color: district.foreground }}>{district.name}</h2>
                  <p class="district-card-desc" style={{ color: district.detail }}>{district.description}</p>
                  <span class="district-card-count" style={{ color: district.detail }}>
                    {district.site_count} {district.site_count === 1 ? "site" : "sites"}
                  </span>
                </div>
                <div class="district-card-icon">
                  <img src={districtImage(district.slug)} alt={district.name} class={districtIconClass(district.slug)} />
                </div>
              </A>
            )}
          </For>
        </div>
      </Show>
    </section>
  );
}