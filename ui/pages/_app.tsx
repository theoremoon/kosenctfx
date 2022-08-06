import { defaultFetcher } from "lib/api";
import type { AppProps } from "next/app";
import React, { ReactElement, ReactNode } from "react";
import { SWRConfig } from "swr";
import { isStaticMode } from "lib/static";
import AppView from "theme/app";
import { CTF } from "lib/api/ctf";
import useAccount from "lib/api/account";
import { NextPage } from "next";

export type NextPageWithLayout = NextPage & {
  getLayout?: (page: ReactElement) => ReactNode;
};

type AppPropsWithLayout = AppProps & {
  Component: NextPageWithLayout;
};

const App = ({ Component, pageProps }: AppPropsWithLayout) => {
  let ctf: CTF | undefined = pageProps.ctf;
  const { data: account } = useAccount(null);

  if (!ctf) {
    ctf = {
      is_open: false,
      is_over: false,
      is_running: false,
      start_at: 0,
      end_at: 0,
      register_open: false,
    };
    // return <Loading />;
  }

  const siteName = "zer0pts CTF 2022";
  const canShowTasks =
    ctf.is_open && (ctf.is_over || (ctf.is_running && account));

  const leftMenuItems = [
    { item: { href: "/tasks", innerText: "TASKS" }, available: canShowTasks },
    { item: { href: "/ranking", innerText: "RANKING" }, available: true },
  ].flatMap((x) => (x.available ? [x.item] : []));

  const rightMenuItems = [
    {
      item: { href: "/admin", innerText: "ADMIN" },
      available: account && account.is_admin,
    },
    { item: { href: "/profile", innerText: "PROFILE" }, available: account },
    { item: { href: "/login", innerText: "LOGIN" }, available: !account },
    { item: { href: "/register", innerText: "REGISTER" }, available: !account },
    { item: { href: "/logout", innerText: "LOGOUT" }, available: account },
  ].flatMap((x) => (x.available && !isStaticMode ? [x.item] : []));

  if (Component.getLayout !== undefined) {
    return <>{Component.getLayout(<Component />)}</>;
  }

  return (
    <AppView
      Component={Component}
      pageProps={pageProps}
      siteName={siteName}
      leftMenuItems={leftMenuItems}
      rightMenuItems={rightMenuItems}
    />
  );
};

const AppWrapper = (appProps: AppPropsWithLayout) => {
  return (
    <SWRConfig
      value={{
        fetcher: defaultFetcher,
        revalidateOnMount: !isStaticMode,
        revalidateOnFocus: !isStaticMode,
        revalidateOnReconnect: !isStaticMode,
        refreshInterval: isStaticMode ? 0 : 30000,
      }}
    >
      <App {...appProps} />
    </SWRConfig>
  );
};

export default AppWrapper;
