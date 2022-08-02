import { defaultFetcher } from "lib/api";
import type { AppProps } from "next/app";
import React from "react";
import { SWRConfig } from "swr";
import { isStaticMode } from "lib/static";
import AppView from "theme/app";
import { CTF } from "lib/api/ctf";
import Loading from "components/loading";
import useAccount from "lib/api/account";

const App = ({ Component, pageProps }: AppProps) => {
  const ctf: CTF | undefined = pageProps.ctf;
  const { data: account } = useAccount(pageProps.account || null);

  if (!ctf) {
    return <Loading />;
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
      <AppView
        Component={Component}
        pageProps={pageProps}
        siteName={siteName}
        leftMenuItems={leftMenuItems}
        rightMenuItems={rightMenuItems}
      />
    </SWRConfig>
  );
};

export default App;
