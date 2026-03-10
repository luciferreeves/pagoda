import { type JSX, Show, onMount, onCleanup, createEffect } from "solid-js";
import { For } from "solid-js/web";
import { A, useLocation } from "@solidjs/router";
import Sidebar from "./Sidebar";
import NavSection from "./NavSection";
import { auth } from "../store/auth";
import { stats } from "../store/stats";
import { UserRole } from "../types/roles";
import type { CitizenSummary } from "../types/stats";


interface LayoutProps {
  children: JSX.Element;
}

export default function Layout(props: LayoutProps) {
  let heartbeatInterval: ReturnType<typeof setInterval> | undefined;
  const location = useLocation();

  onMount(() => {
    auth.initialize();
    stats.load();

    function handleExternalLinks(event: MouseEvent) {
      const anchor = (event.target as HTMLElement).closest("a");
      if (!anchor || !anchor.href) return;
      try {
        const url = new URL(anchor.href);
        if (url.hostname !== window.location.hostname) {
          anchor.setAttribute("target", "_blank");
          anchor.setAttribute("rel", "noopener noreferrer");
        }
      } catch {}
    }
    document.addEventListener("click", handleExternalLinks);
    onCleanup(() => document.removeEventListener("click", handleExternalLinks));
  });

  createEffect(() => {
    location.pathname;
    stats.load();
  });

  createEffect(() => {
    clearInterval(heartbeatInterval);
    if (auth.user()) {
      auth.heartbeat();
      heartbeatInterval = setInterval(() => auth.heartbeat(), 2 * 60 * 1000);
    }
  });

  onCleanup(() => clearInterval(heartbeatInterval));

  return (
    <>
      <header class="site-header">
        <h1>Pagoda</h1>
        <p>A community for the small web</p>
      </header>

      <nav class="top-nav">
        <A href="/" class="top-nav-link" data-accent="cyan" activeClass="active" end>Home</A>
        <A href="/districts" class="top-nav-link" data-accent="green" activeClass="active">Districts</A>
        <A href="/forums" class="top-nav-link" data-accent="pink" activeClass="active">Forums</A>
        <A href="/tavern" class="top-nav-link" data-accent="yellow" activeClass="active">Tavern</A>
        <A href="/bazaar" class="top-nav-link" data-accent="purple" activeClass="active">Bazaar</A>
      </nav>

      <div class="search-bar">
        <input type="text" placeholder="Search citizens, posts, districts..." />
      </div>

      <div class="site-main">
        <Sidebar>
          <NavSection title="Account" accent="pink">
            <ul>
              <Show when={auth.user()} fallback={
                <>
                  <li><A href="/login">Log In</A></li>
                  <li><A href="/register">Register</A></li>
                </>
              }>
                {((user) => (
                  <>
                    <li><A href={`/p/${user.username}`}>My Domain</A></li>
                    <li><A href="/letters">Letters</A></li>
                    <li><A href="/account">Account</A></li>
                    <li><A href="/account/settings">Settings</A></li>
                    <li><button type="button" class="sidebar-logout" onClick={() => auth.logout()}>Log Out</button></li>
                  </>
                ))(auth.user()!)}
              </Show>
            </ul>
          </NavSection>
          <NavSection title="Community" accent="cyan">
            <ul>
              <li><A href="/citizens">Citizens</A></li>
              <li><A href="/clubs">Clubs</A></li>
              <li><A href="/interests">Interests</A></li>
            </ul>
          </NavSection>
          <NavSection title="Explore" accent="yellow">
            <ul>
              <li><A href="/online">Who's Online</A></li>
              <li><A href="/random">Random Domain</A></li>
            </ul>
          </NavSection>
          <NavSection title="Services" accent="green">
            <ul>
              <li><A href="/caravan">Caravan</A></li>
              <li><A href="/ledger">Ledger</A></li>
              <li><A href="/census">Census</A></li>
              <li><A href="/pamphlet">Pamphlet</A></li>
            </ul>
          </NavSection>
          <Show when={auth.user()?.role === UserRole.Owner || auth.user()?.role === UserRole.Admin || auth.user()?.role === UserRole.Moderator}>
            <NavSection title="Council" accent="red">
              <ul>
                <Show when={auth.user()?.role === UserRole.Owner || auth.user()?.role === UserRole.Admin}>
                  <li><A href="/council/announcements">Announcements</A></li>
                  <li><A href="/council/audit-log">Audit Log</A></li>
                  <li><A href="/council/bannedips">Banned IPs</A></li>
                </Show>
                <li><A href="/council/bazaar">Bazaar</A></li>
                <li><A href="/council/districts">Districts</A></li>
                <li><A href="/council/forums">Forums</A></li>
                <li><A href="/council/reports">Reports</A></li>
                <li><A href="/council/users">Users</A></li>
              </ul>
            </NavSection>
          </Show>
        </Sidebar>

        <main class="content">
          {props.children}
        </main>

        <Sidebar>
          <NavSection title="Statistics" accent="yellow">
            <ul>
              <li>Citizens: {stats.data()?.citizens ?? "—"}</li>
              <li>Online: {stats.data()?.online ?? "—"}</li>
              <li>Posts Today: —</li>
            </ul>
          </NavSection>
          <NavSection title="New Citizens" accent="cyan">
            <ul>
              <Show when={stats.data()?.newest_citizens?.length} fallback={
                <li class="placeholder">Be the first to join!</li>
              }>
                <For each={stats.data()?.newest_citizens}>
                  {(citizen: CitizenSummary) => (
                    <li class="citizen-item">
                      <A href={`/p/${citizen.username}`}>
                        <img src={citizen.avatar_url} alt="" class="citizen-avatar" />
                        {citizen.display_name}
                      </A>
                    </li>
                  )}
                </For>
              </Show>
            </ul>
          </NavSection>
          <NavSection title="Who's Online" accent="green">
            <ul>
              <Show when={stats.data()?.online_citizens?.length} fallback={
                <li class="placeholder">No one online.</li>
              }>
                <For each={stats.data()?.online_citizens}>
                  {(citizen: CitizenSummary) => (
                    <li class="citizen-item">
                      <A href={`/p/${citizen.username}`}>
                        <img src={citizen.avatar_url} alt="" class="citizen-avatar" />
                        {citizen.display_name}
                      </A>
                    </li>
                  )}
                </For>
              </Show>
            </ul>
          </NavSection>
          <NavSection title="Birthdays" accent="pink">
            <ul>
              <li class="placeholder">No birthdays today.</li>
            </ul>
          </NavSection>
          <NavSection title="Top Posters" accent="yellow">
            <ul>
              <li class="placeholder">No posts yet.</li>
            </ul>
          </NavSection>
          <NavSection title="Staff Online" accent="purple">
            <ul>
              <li class="placeholder">No staff online.</li>
            </ul>
          </NavSection>
        </Sidebar>
      </div>

      <footer class="site-footer">
        <nav class="footer-links">
          <A href="/about">About</A>
          <A href="/rules">Rules</A>
          <A href="/privacy">Privacy</A>
          <A href="/terms">Terms</A>
          <A href="/contact">Contact</A>
        </nav>
        <p>&copy; {new Date().getFullYear()} Pagoda. Brought to you by <a href="https://shi.foo" target="_blank" rel="noopener noreferrer">shi.foo</a>. Powered by <a href="https://nekoweb.org" target="_blank" rel="noopener noreferrer">Nekoweb</a>.</p>
      </footer>
    </>
  );
}
