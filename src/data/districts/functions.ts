import type { District, Site } from "./types";
import { districts } from "./districts";
import { arcadia } from "./arcadia";
import { arles } from "./arles";
import { hollywood } from "./hollywood";
import { oxford } from "./oxford";
import { petsburg } from "./petsburg";
import { purgatory } from "./purgatory";
import { siliconValley } from "./siliconValley";
import { silverLake } from "./silverLake";
import { stratfordUponAvon } from "./stratfordUponAvon";
import { tokyo } from "./tokyo";

export const getDistrictById = (id: string): District | undefined => {
  return districts.find((district) => district.id === id);
};

export const getSitesByDistrictId = (id: string): Site[] => {
  switch (id) {
    case "arcadia":
      return arcadia;
    case "arles":
      return arles;
    case "hollywood":
      return hollywood;
    case "oxford":
      return oxford;
    case "petsburg":
      return petsburg;
    case "purgatory":
      return purgatory;
    case "silicon-valley":
      return siliconValley;
    case "silver-lake":
      return silverLake;
    case "stratford-upon-avon":
      return stratfordUponAvon;
    case "tokyo":
      return tokyo;
    default:
      return [];
  }
};

export const getPaginatedSitesByDistrictId = (
  id: string,
  page: number,
  limit: number
): { sites: Site[]; totalPages: number; currentPage: number } => {
  const sites = getSitesByDistrictId(id);
  const startIndex = (page - 1) * limit;
  const endIndex = startIndex + limit;
  return {
    sites: sites.slice(startIndex, endIndex),
    totalPages: Math.ceil(sites.length / limit),
    currentPage: page,
  };
};
