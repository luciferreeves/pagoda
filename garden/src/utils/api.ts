export function extractError(data: unknown): string {
  return (data as { error: string }).error || "An error occurred.";
}