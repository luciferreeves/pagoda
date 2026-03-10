import { type JSX, Show, createEffect } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { auth } from "../store/auth";
import { UserRole } from "../types/roles";

interface StaffGuardProps {
  children: JSX.Element;
}

function isStaff(role?: string) {
  return role === UserRole.Owner || role === UserRole.Admin || role === UserRole.Moderator;
}

export default function StaffGuard(props: StaffGuardProps) {
  const navigate = useNavigate();

  createEffect(() => {
    if (!auth.loading() && (!auth.user() || !isStaff(auth.user()?.role))) {
      navigate("/", { replace: true });
    }
  });

  return (
    <Show when={!auth.loading() && auth.user() && isStaff(auth.user()?.role)}>
      {props.children}
    </Show>
  );
}