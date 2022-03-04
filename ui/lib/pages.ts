import { CTF } from "./api/ctf";
import { Account } from "./api/account";

export interface AllPageProps {
  ctf: CTF;
  account: Account | null;
}
