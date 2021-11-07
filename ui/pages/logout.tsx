import Loading from "components/loading";
import { api } from "lib/api";
import { useRouter } from "next/router";
import React, { useEffect } from "react";
import useAccount from "../lib/api/account";

const Logout = () => {
  const router = useRouter();
  const { mutate } = useAccount();
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

export default Logout;
