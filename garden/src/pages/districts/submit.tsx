import { createSignal, onMount, Show, For } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { api } from "../../api";
import { auth } from "../../store/auth";
import { extractError } from "../../utils/api";
import type { District, SiteRequest } from "../../types/district";
import { districtImage, districtIconClass } from "../../utils/districts";

export default function SubmitSite() {
  const navigate = useNavigate();
  const [districts, setDistricts] = createSignal<District[]>([]);
  const [selectedDistrict, setSelectedDistrict] = createSignal("");
  const [title, setTitle] = createSignal("");
  const [url, setUrl] = createSignal("");
  const [description, setDescription] = createSignal("");
  const [tagInput, setTagInput] = createSignal("");
  const [tags, setTags] = createSignal<string[]>([]);
  const [error, setError] = createSignal("");
  const [submitting, setSubmitting] = createSignal(false);

  onMount(async () => {
    const response = await api<District[]>("/districts");
    if (response.ok) {
      setDistricts(response.data);
    }
  });

  function addTag() {
    const tag = tagInput().trim().toLowerCase();
    if (tag && tags().length < 5 && !tags().includes(tag)) {
      setTags([...tags(), tag]);
      setTagInput("");
    }
  }

  function removeTag(tag: string) {
    setTags(tags().filter((existingTag) => existingTag !== tag));
  }

  function handleTagKeyDown(event: KeyboardEvent) {
    if (event.key === "Enter") {
      event.preventDefault();
      addTag();
    }
  }

  async function handleSubmit(event: Event) {
    event.preventDefault();
    setError("");
    setSubmitting(true);

    const response = await api<SiteRequest>("/districts/sites", {
      method: "POST",
      token: auth.token(),
      body: {
        district: selectedDistrict(),
        title: title(),
        url: url(),
        description: description(),
        tags: tags(),
      },
    });

    if (response.ok) {
      navigate("/districts");
    } else {
      setError(extractError(response.data));
    }
    setSubmitting(false);
  }

  return (
    <section>
      <h1 class="page-title">Submit a Site</h1>
      <p class="district-intro">
        Submit a website to be listed in a district. Your submission will be reviewed by staff before it appears.
      </p>

      <Show when={error()}>
        <div class="form-error">{error()}</div>
      </Show>

      <form class="submit-site-form" onSubmit={handleSubmit}>
        <div class="form-field">
          <label>District</label>
          <div class="district-select-grid">
            <For each={districts()}>
              {(district) => (
                <button
                  type="button"
                  class={`district-select-option ${selectedDistrict() === district.slug ? "selected" : ""}`}
                  style={{
                    "border-color": selectedDistrict() === district.slug ? district.foreground : district.background,
                    "background-color": district.background,
                  }}
                  onClick={() => setSelectedDistrict(district.slug)}
                >
                  <img src={districtImage(district.slug)} alt={district.name} class={`district-select-icon ${districtIconClass(district.slug)}`} />
                  <span style={{ color: district.foreground }}>{district.name}</span>
                </button>
              )}
            </For>
          </div>
        </div>

        <div class="form-field">
          <label>Site Title</label>
          <input
            type="text"
            value={title()}
            onInput={(e) => setTitle(e.currentTarget.value)}
            placeholder="My Cool Website"
            maxLength={200}
          />
        </div>

        <div class="form-field">
          <label>URL</label>
          <input
            type="url"
            value={url()}
            onInput={(e) => setUrl(e.currentTarget.value)}
            placeholder="https://example.nekoweb.org"
          />
        </div>

        <div class="form-field">
          <label>Description</label>
          <textarea
            value={description()}
            onInput={(e) => setDescription(e.currentTarget.value)}
            placeholder="A short description of the site..."
            maxLength={1000}
            rows={3}
          />
        </div>

        <div class="form-field">
          <label>Tags (up to 5)</label>
          <div class="tag-input-row">
            <input
              type="text"
              value={tagInput()}
              onInput={(e) => setTagInput(e.currentTarget.value)}
              onKeyDown={handleTagKeyDown}
              placeholder="Add a tag..."
              maxLength={50}
              disabled={tags().length >= 5}
            />
            <button type="button" class="form-button" onClick={addTag} disabled={tags().length >= 5}>Add</button>
          </div>
          <Show when={tags().length > 0}>
            <div class="tag-list">
              <For each={tags()}>
                {(tag) => (
                  <span class="site-tag removable" onClick={() => removeTag(tag)}>
                    {tag} &times;
                  </span>
                )}
              </For>
            </div>
          </Show>
        </div>

        <button type="submit" class="form-button" disabled={submitting() || !selectedDistrict()}>
          {submitting() ? "Submitting..." : "Submit Site"}
        </button>
      </form>
    </section>
  );
}