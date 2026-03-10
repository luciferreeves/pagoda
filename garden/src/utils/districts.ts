const DISTRICT_IMAGES: Record<string, string> = {
  arcadia: "/images/districts/arcadia.png",
  arles: "/images/districts/arles.png",
  hollywood: "/images/districts/hollywood.png",
  oxford: "/images/districts/oxford.png",
  petsburg: "/images/districts/petsburg.png",
  purgatory: "/images/districts/purgatory.png",
  "silicon-valley": "/images/districts/silicon.png",
  "silver-lake": "/images/districts/silver.png",
  "stratford-upon-avon": "/images/districts/stratford.png",
  tokyo: "/images/districts/tokyo.png",
};

export function districtImage(slug: string): string {
  return DISTRICT_IMAGES[slug] || "";
}

const OUTLINED_SLUGS = new Set([
  "arcadia",
  "arles",
  "silver-lake",
  "silicon-valley",
  "stratford-upon-avon",
  "purgatory",
  "tokyo",
]);

export function districtIconClass(slug: string): string {
  return OUTLINED_SLUGS.has(slug) ? "icon-outlined" : "";
}
