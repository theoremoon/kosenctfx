import {
  Box,
  Container,
  extendTheme,
  Flex,
  withDefaultColorScheme,
  withDefaultVariant,
} from "@chakra-ui/react";
import { defaultFetcher } from "lib/api";
import type { AppProps } from "next/app";
import React from "react";
import { SWRConfig } from "swr";
import Menu from "../components/menu";
import { bgColor, bgSubColor, pink, white } from "../lib/color";
import "tailwindcss/tailwind.css";
import "../style/global.scss";

const App = ({ Component, pageProps }: AppProps) => {
  return (
    <SWRConfig value={{ fetcher: defaultFetcher }}>
      <div className="min-h-full h-full flex flex-row">
        <div className="w-44 break-normal">
          <Menu />
        </div>
        <div className="container h-full">
          <div className="w-min md:w-2/3 mx-auto h-full">
            <div className="h-full flex flex-col">
              <div className="flex-1">
                <Component {...pageProps} />
              </div>
              <div className="opacity-20 text-center">powered by kosenctfx</div>
            </div>
          </div>
        </div>
      </div>
    </SWRConfig>
  );
};
export default App;
