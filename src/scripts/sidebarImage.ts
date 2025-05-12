document.addEventListener("DOMContentLoaded", (): void => {
  const currentMinute: number = Math.floor(Date.now() / 60000);
  const sidebarImagePathPrefix: string = "/images/internal/sidebar";

  type ImageRule = {
    divisor: number;
    filename: string;
  };

  const imageRules: ImageRule[] = [
    { divisor: 7, filename: "7.gif" },
    { divisor: 5, filename: "5.gif" },
    { divisor: 3, filename: "3.gif" },
    { divisor: 2, filename: "2.gif" },
  ];

  const matchedRule: ImageRule | undefined = imageRules.find(
    (rule: ImageRule): boolean => currentMinute % rule.divisor === 0
  );

  const imageName: string = matchedRule ? matchedRule.filename : "default.gif";
  const sidebarImage: string = `${sidebarImagePathPrefix}/${imageName}`;

  const imageElement: HTMLImageElement | null = document.getElementById(
    "sidebar-image"
  ) as HTMLImageElement;
  if (imageElement) {
    imageElement.src = sidebarImage;
  }
});
