import { A } from "@solidjs/router";

export default function NotFound() {
  return (
    <section>
      <h1>404</h1>
      <p>This page does not exist. <A href="/">Return home</A>.</p>
    </section>
  );
}