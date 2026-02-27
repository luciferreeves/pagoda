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

      <div class="site-main">
        <Sidebar>
          <NavSection title="Navigation">
            <li><A href="/">Home</A></li>
            <li><A href="/districts">Districts</A></li>
            <li><A href="/forums">Forums</A></li>
            <li><A href="/chat">Chat</A></li>
          </NavSection>
          <NavSection title="Resources" variant="alternate">
            <li><A href="/bazaar">Bazaar</A></li>
          </NavSection>
        </Sidebar>

        <main class="content">
          {props.children}
        </main>

        <Sidebar>
          <NavSection title="Account" variant="alternate">
            <li><A href="/login">Log In</A></li>
            <li><A href="/register">Register</A></li>
          </NavSection>
        </Sidebar>
      </div>

      <footer class="site-footer">
        <p>&copy; {new Date().getFullYear()} Pagoda. Brought to you by <a href="https://shi.foo" target="_blank" rel="noopener noreferrer">shi.foo</a>.</p>
      </footer>
    </>
  );
}