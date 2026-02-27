import type { JSX } from "solid-js";
import { A } from "@solidjs/router";
import Sidebar from "./Sidebar";
import NavSection from "./NavSection";

interface LayoutProps {
  children: JSX.Element;
}

export default function Layout(props: LayoutProps) {
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
        <A href="/chat" class="top-nav-link" data-accent="yellow" activeClass="active">Chat</A>
        <A href="/bazaar" class="top-nav-link" data-accent="purple" activeClass="active">Bazaar</A>
      </nav>

      <div class="search-bar">
        <input type="text" placeholder="Search members, posts, districts..." />
      </div>

      <div class="site-main">
        <Sidebar>
          <NavSection title="Account" accent="pink">
            <ul>
              <li><A href="/login">Log In</A></li>
              <li><A href="/register">Register</A></li>
            </ul>
          </NavSection>
          <NavSection title="Community" accent="cyan">
            <ul>
              <li><A href="/members">Members</A></li>
              <li><A href="/online">Who's Online</A></li>
              <li><A href="/buttons">Button Wall</A></li>
              <li><A href="/random">Random Member</A></li>
            </ul>
          </NavSection>
          <NavSection title="Services" accent="green">
            <ul>
              <li><A href="/webring">Webring</A></li>
              <li><A href="/guestbook">Guestbook</A></li>
              <li><A href="/hitcounter">Hit Counter</A></li>
            </ul>
          </NavSection>
        </Sidebar>

        <main class="content">
          {props.children}
        </main>

        <Sidebar>
          <NavSection title="Statistics" accent="yellow">
            <ul>
              <li>Members: —</li>
              <li>Online: —</li>
              <li>Posts Today: —</li>
              <li>Newest: —</li>
            </ul>
          </NavSection>
          <NavSection title="New Members" accent="cyan">
            <ul>
              <li class="placeholder">Be the first to join!</li>
            </ul>
          </NavSection>
          <NavSection title="Birthdays" accent="green">
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
        <p>&copy; {new Date().getFullYear()} Pagoda. Brought to you by <a href="https://shi.foo" target="_blank" rel="noopener noreferrer">shi.foo</a>.</p>
      </footer>
    </>
  );
}
