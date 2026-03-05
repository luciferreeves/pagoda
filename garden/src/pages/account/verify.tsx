import { createSignal, onMount, Show } from "solid-js";
import { A, useSearchParams } from "@solidjs/router";
import { auth } from "../../store/auth";

export default function Verify() {
  const [searchParams] = useSearchParams();
  const [message, setMessage] = createSignal("");
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);

  onMount(async () => {
    const token = searchParams.token;
    if (!token) return;

    setLoading(true);
    const result = await auth.verify(token as string, "activation");
    setLoading(false);

    if (result) {
      setError(result);
    } else {
      setMessage("Your email has been verified successfully. You can now log in.");
    }
  });

  return (
    <section>
      <h2 class="page-title">Verify Account</h2>
      <Show when={loading()}>
        <p>Verifying your account...</p>
      </Show>
      <Show when={error()}>
        <div class="form-error">{error()}</div>
        <p class="form-footer">
          <A href="/account/reactivate">Request a new verification email</A>
        </p>
      </Show>
      <Show when={message()}>
        <div class="form-success">{message()}</div>
        <p class="form-footer">
          <A href="/login">Log In</A>
        </p>
      </Show>
      <Show when={!searchParams.token && !loading()}>
        <p>Please check your email for a verification link.</p>
        <p class="form-footer">
          Didn't receive an email? <A href="/account/reactivate">Request a new one</A>
        </p>
      </Show>
    </section>
  );
}