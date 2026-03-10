import { createSignal, Show } from "solid-js";
import { A, useNavigate } from "@solidjs/router";
import { auth } from "../store/auth";

const EMAIL_NOT_VERIFIED = "Your email address has not been verified. Please check your inbox.";

export default function Login() {
  const navigate = useNavigate();
  const [username, setUsername] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [error, setError] = createSignal("");
  const [submitting, setSubmitting] = createSignal(false);

  async function handleSubmit(event: Event) {
    event.preventDefault();
    setError("");
    setSubmitting(true);

    const result = await auth.login(username(), password());
    setSubmitting(false);

    if (result) {
      setError(result);
    } else {
      navigate("/");
    }
  }

  return (
    <section>
      <h2 class="page-title">Log In</h2>
      <Show when={error()}>
        <div class="form-error">
          {error()}
          <Show when={error() === EMAIL_NOT_VERIFIED}>
            {" "}<A href="/account/reactivate">Resend verification email</A>
          </Show>
        </div>
      </Show>
      <form class="auth-form" onSubmit={handleSubmit}>
        <label class="form-field">
          <span>Username</span>
          <input type="text" value={username()} onInput={(e) => setUsername(e.currentTarget.value)} required />
        </label>
        <label class="form-field">
          <span>Password</span>
          <input type="password" value={password()} onInput={(e) => setPassword(e.currentTarget.value)} required />
        </label>
        <button type="submit" class="form-button" disabled={submitting()}>
          {submitting() ? "Logging in..." : "Log In"}
        </button>
      </form>
      <p class="form-footer">
        Don't have an account? <A href="/register">Register</A>
      </p>
    </section>
  );
}