import { createSignal } from "solid-js";
import { A, useNavigate } from "@solidjs/router";
import { auth } from "../store/auth";

export default function Register() {
  const navigate = useNavigate();
  const [username, setUsername] = createSignal("");
  const [email, setEmail] = createSignal("");
  const [displayName, setDisplayName] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [error, setError] = createSignal("");
  const [submitting, setSubmitting] = createSignal(false);

  async function handleSubmit(event: Event) {
    event.preventDefault();
    setError("");
    setSubmitting(true);

    const result = await auth.register(username(), email(), password(), displayName());
    setSubmitting(false);

    if (result) {
      setError(result);
    } else {
      navigate("/account/verify");
    }
  }

  return (
    <section>
      <h2 class="page-title">Register</h2>
      {error() && <div class="form-error">{error()}</div>}
      <form class="auth-form" onSubmit={handleSubmit}>
        <label class="form-field">
          <span>Username</span>
          <input type="text" value={username()} onInput={(e) => setUsername(e.currentTarget.value)} required />
        </label>
        <label class="form-field">
          <span>Email</span>
          <input type="email" value={email()} onInput={(e) => setEmail(e.currentTarget.value)} required />
        </label>
        <label class="form-field">
          <span>Display Name</span>
          <input type="text" value={displayName()} onInput={(e) => setDisplayName(e.currentTarget.value)} />
        </label>
        <label class="form-field">
          <span>Password</span>
          <input type="password" value={password()} onInput={(e) => setPassword(e.currentTarget.value)} required />
        </label>
        <button type="submit" class="form-button" disabled={submitting()}>
          {submitting() ? "Registering..." : "Register"}
        </button>
      </form>
      <p class="form-footer">
        Already have an account? <A href="/login">Log In</A>
      </p>
    </section>
  );
}