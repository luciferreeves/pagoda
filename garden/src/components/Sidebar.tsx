import type { JSX } from "solid-js";

interface SidebarProps {
  children: JSX.Element;
}

export default function Sidebar(props: SidebarProps) {
  return <aside class="sidebar">{props.children}</aside>;
}