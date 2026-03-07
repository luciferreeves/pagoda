declare module "solid-js" {
  namespace JSX {
    interface IntrinsicElements {
      "emoji-picker": HTMLAttributes<HTMLElement> & {
        class?: string;
        "on:emoji-click"?: (event: CustomEvent) => void;
      };
    }
  }
}