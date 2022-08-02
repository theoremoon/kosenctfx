import type { AppProps as NextAppProps } from "next/app";
import { MenuItem } from "./menu";

export type AppProps = Pick<NextAppProps, "Component" | "pageProps"> & {
  siteName: string;
  leftMenuItems: MenuItem[];
  rightMenuItems: MenuItem[];
};
