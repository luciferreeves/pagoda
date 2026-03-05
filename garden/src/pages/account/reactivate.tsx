import { createSignal } from "solid-js";
import { A } from "@solidjs/router";
import { auth } from "../../store/auth";

export default function Reactivate() {
  const [email, setEmail] = createSignal("");
  const [message, setMessage] = createSignal("");
  const [error, setError] = createSignal("");
  const [submitting, setSubmitting] = createSignal(false);

  async function handleSubmit(event: Event) {
    event.preventDefault();
    setError("");
    setMessage("");
    setSubmitting(true);

    const result = await auth.reactivate(email());
    setSubmitting(false);

    if (result) {
      setError(result);
    } else {
      setMessage("A new verification email has been sent. Please check your inbox.");
    }
  }

  return (
    <section>
      <h2 class="page-title">Resend Verification</h2>
      {error() && <div class="form-error">{error()}</div>}
      {message() && <div class="form-success">{message()}</div>}
      <form class="auth-form" onSubmit={handleSubmit}>
        <label class="form-field">
          <span>Email</span>
          <input type="email" value={email()} onInput={(e) => setEmail(e.currentTarget.value)} required />
        </label>
        <button type="submit" class="form-button" disabled={submitting()}>
          {submitting() ? "Sending..." : "Resend Verification Email"}
        </button>
      </form>
      <p class="form-footer">
        Already verified? <A href="/login">Log In</A>
      </p>
    </section>
  );
}