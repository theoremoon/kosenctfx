import { api, makeSWRResponse, ssrApi } from "lib/api";
import { isStaticMode } from "lib/static";
import useSWR from "swr";

export interface SeriesEntry {
  teamname: string;
  score: number;
  pos: number;
  time: number;
}

const useSeries = (teams: string[], staticValue: SeriesEntry[][]) => {
  return isStaticMode
    ? makeSWRResponse(staticValue)
    : useSWR<SeriesEntry[][]>(teams.join(""), (url) => {
      return api
        .post<SeriesEntry[][]>("/series", {
          teams: teams,
        })
        .then((r) => r.data);
    });
};

export const fetchSeries = (teams: string[]) =>
  ssrApi
    .post<SeriesEntry[][]>("/series", {
      teams: teams,
    })
    .then((r) => r.data);
export default useSeries;
