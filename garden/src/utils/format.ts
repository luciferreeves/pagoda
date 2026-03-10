export function formatDate(date: string | null): string {
  if (!date) return "\u2014";
  return new Date(date).toLocaleDateString();
}

export function formatDateTime(date: string | null): string {
  if (!date) return "\u2014";
  const d = new Date(date);
  return d.toLocaleDateString() + " " + d.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
}

export function formatDateTimeFull(date: string | null): string {
  if (!date) return "\u2014";
  const d = new Date(date);
  return d.toLocaleDateString() + " " + d.toLocaleTimeString();
}