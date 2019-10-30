export function ellipsify (text: string, head: number = 6, tail: number = 6): string {
  return `${text.slice(0, head)}...${text.slice(-tail)}`;
}