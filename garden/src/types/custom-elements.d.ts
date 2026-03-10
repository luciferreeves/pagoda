import "solid-js";

declare module "solid-js" {
  namespace JSX {
    interface IntrinsicElements {
      "emoji-picker": JSX.HTMLAttributes<HTMLElement> & {
        class?: string;
        "on:emoji-click"?: (event: CustomEvent) => void;
      };
    }
  }
}