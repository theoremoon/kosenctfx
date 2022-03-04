import { makeSWRResponse, ssrFetcher } from "lib/api";
import { isStaticMode } from "lib/static";
import useSWR from "swr";

export interface Account {
  teamname: string;
  team_id: number;
  country: string;
  is_admin: boolean;
}

const useAccount = (staticValue: Account|null) => {
  return (isStaticMode) ? (
    makeSWRResponse(staticValue)
  ) : (
    useSWR<Account | null>("/account")
  )
}
export const fetchAccount = () => ssrFetcher<Account|null>("/account");
export default useAccount;
