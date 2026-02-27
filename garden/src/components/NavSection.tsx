import type { JSX } from "solid-js";

interface NavSectionProps {
  title: string;
  variant?: "primary" | "alternate";
  children: JSX.Element;
}

export default function NavSection(props: NavSectionProps) {
  const variant = () => props.variant ?? "primary";

  return (
    <section class={`nav-section ${variant()}`}>
      <div class="nav-section-header">
        <div class="bar-left" />
        <h3>{props.title}</h3>
        <div class="bar-right" />
      </div>
      <div class="nav-section-body">
        <ul>{props.children}</ul>
      </div>
    </section>
  );
}