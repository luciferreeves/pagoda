import { onMount, onCleanup, Setter } from "solid-js";

export function useClickOutside(
  getRef: () => HTMLElement | undefined,
  setOpen: Setter<boolean>,
) {
  onMount(() => {
    function handler(e: MouseEvent) {
      const ref = getRef();
      if (ref && !ref.contains(e.target as Node)) setOpen(false);
    }
    document.addEventListener("mousedown", handler);
    onCleanup(() => document.removeEventListener("mousedown", handler));
  });
}