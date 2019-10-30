export const getHHMMSS = (iso: string): string => {
  return iso.split(/T|\./)[1];
};