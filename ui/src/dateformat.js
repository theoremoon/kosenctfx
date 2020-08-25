import dayjs from "dayjs";

export function dateFormat(ts) {
  return dayjs(ts * 1000).format("YYYY-MM-DD HH:mm:ss Z");
}
