const parentpath = (path: string) => {
  const p2 = path.endsWith("/") ? path.slice(0, -1) : path;
  const parts = p2.split("/");
  if (parts.length === 1) {
    return p2;
  }
  return parts.slice(0, -1).join("/");
};
export default parentpath;
