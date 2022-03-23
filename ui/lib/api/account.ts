import useSWR from "swr";

export interface Account {
  teamname: string;
  team_id: number;
  country: string;
  is_admin: boolean;
}

const useAccount = () => useSWR<Account | null>("/account");
export default useAccount;
