import Loading from "components/loading";
import { api } from "lib/api";
import { fetchCTF } from "lib/api/ctf";
import { AllPageProps } from "lib/pages";
import { isStaticMode, revalidateInterval } from "lib/static";
import { GetStaticProps } from "next";
import { useRouter } from "next/router";
import React, { useEffect } from "react";
import useAccount from "../lib/api/account";

const Logout = () => {
  const router = useRouter();
  const { mutate } = useAccount(null);
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

export const getStaticProps: GetStaticProps<AllPageProps> = async () => {
  const ctf = await fetchCTF();
  return {
    props: {
      ctf: ctf,
    },
    revalidate: isStaticMode ? false : revalidateInterval,
  };
};

export default Logout;
