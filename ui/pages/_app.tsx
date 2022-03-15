import {
  Box,
  Center,
  ChakraProvider,
  Container,
  extendTheme,
  withDefaultVariant,
} from "@chakra-ui/react";
import Menu from "../components/menu";
import { defaultFetcher } from "lib/api";
import type { AppProps } from "next/app";
import React from "react";
import { SWRConfig } from "swr";
import { isStaticMode } from "lib/static";

const theme = extendTheme(
  {
    initialColorMode: "dark",
    useSystemColorMode: false,
    colors: {
      gray: {
        "50": "#EFF0F6",
        "100": "#D1D4E5",
        "200": "#ffffff",
        "300": "#969DC4",
        "400": "#7982B4",
        "500": "#5B67A4",
        "600": "#495283",
        "700": "#373E62",
        "800": "#0b1933",
        "900": "#121521",
      },
    },
    styles: {
      global: {
        "html, body": {
          minHeight: "100vh",
          height: "100%",
        },
        "#__next": {
          minHeight: "100%",
        },
        input: {
          textAlign: "center",
        },
      },
    },
    components: {
      Link: {
        baseStyle: {
          "&:focus": {
            boxShadow: "none",
          },
        },
      },
      Code: {
        withDefaultVariant: "solid",
      },
    },
  },

  withDefaultVariant({
    variant: "flushed",
    components: ["Input"],
  })
);

const App = ({ Component, pageProps }: AppProps) => {
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
      <ChakraProvider theme={theme}>
        <Box className="h-full">
          <Menu {...pageProps} />
          <Container maxW="container.xl" minH="100vh">
            <Component {...pageProps} />
          </Container>
          <Center color="#ffffff66" mt={10} mb={1}>
            powered by kosenctfx
          </Center>
        </Box>
      </ChakraProvider>
    </SWRConfig>
  );
};

export default App;
