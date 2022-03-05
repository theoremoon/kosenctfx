import Loading from "components/loading";
import { api } from "lib/api";
import { isStaticMode } from "lib/static";
import { GetStaticProps } from "next";
import { useRouter } from "next/router";
import React, { useEffect } from "react";
import useAccount, { Account, fetchAccount } from "../lib/api/account";

interface LogoutProps {
  account: Account | null;
}

const Logout = ({ account }: LogoutProps) => {
  const router = useRouter();
  const { mutate } = useAccount(account);
  useEffect(() => {
    const f = async () => {
      await api.post("/logout");
      mutate();
      router.push("/");
    };
    f();
  }, []);
  return <Loading />;
};

export const getStaticProps: GetStaticProps<LogoutProps> = async () => {
  const account = isStaticMode ? null : await fetchAccount();
  return {
    props: {
      account: account,
    },
  };
};

export default Logout;
