import dayjs from "dayjs";
const customParseFormat = require("dayjs/plugin/customParseFormat");
dayjs.extend(customParseFormat);

export function dateFormat(ts) {
  return dayjs(ts * 1000).format("YYYY-MM-DD HH:mm:ss Z");
}

export function parseDateString(s) {
  const d = dayjs(s, "YYYY-MM-DD HH:mm:ss Z", true);
  if (!d.isValid()) {
    return null;
  }
  return d.unix();
}
