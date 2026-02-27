import type { JSX } from "solid-js";

interface NavSectionProps {
  title: string;
  accent?: "cyan" | "green" | "pink" | "purple" | "yellow";
  children: JSX.Element;
}

export default function NavSection(props: NavSectionProps) {
  return (
    <section class="nav-section" data-accent={props.accent || "purple"}>
      <div class="nav-section-header">{props.title}</div>
      <div class="nav-section-body">
        {props.children}
      </div>
    </section>
  );
}