export type ItemKind = "news" | "recipe" | "deal" | "shop";

export function isItemKind(value: string | null): value is ItemKind {
  return (
    value === "news" ||
    value === "recipe" ||
    value === "deal" ||
    value === "shop"
  );
}

export function kindBadgeLabel(kind: ItemKind): string {
  switch (kind) {
    case "news":
      return "NEWS";
    case "recipe":
      return "RECIPE";
    case "deal":
      return "DEAL";
    case "shop":
      return "SHOP";
    default:
      return "";
  }
}

export function kindTitleLabel(kind: ItemKind): string {
  switch (kind) {
    case "news":
      return "ニュース";
    case "recipe":
      return "レシピ";
    case "deal":
      return "セール";
    case "shop":
      return "店舗";
    default:
      return "";
  }
}

export function cardImage(url: string | null): string {
  if (url && url.trim() !== "") {
    return url;
  }

  return "https://images.unsplash.com/photo-1509042239860-f550ce710b93?auto=format&fit=crop&w=1200&q=80";
}

export function hasRef(url: string | null): boolean {
  return !!url && url.trim() !== "";
}
