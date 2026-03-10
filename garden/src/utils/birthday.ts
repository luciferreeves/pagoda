const MONTHS = [
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December",
];

function parseParts(raw: string): string[] {
  const base = raw.includes("T") ? raw.split("T")[0] : raw;
  return base.split("-");
}

export function formatBirthday(raw: string | null): string {
  if (!raw) return "\u2014";
  const parts = parseParts(raw);
  const month = parseInt(parts.length === 3 ? parts[1] : parts[0]) - 1;
  const day = parseInt(parts.length === 3 ? parts[2] : parts[1]);
  return `${MONTHS[month]} ${day}`;
}

export function birthdayToMonthDay(raw: string | null): string {
  if (!raw) return "";
  const parts = parseParts(raw);
  return parts.length === 3 ? `${parts[1]}-${parts[2]}` : raw;
}