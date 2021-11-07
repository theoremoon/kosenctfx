import {
  Box,
  ChakraProvider,
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

const theme = extendTheme(
  {
    initialColorMode: "dark",
    useSystemColorMode: false,
    styles: {
      global: {
        "html, body": {
          minHeight: "100vh",
          height: "100%",
          backgroundColor: bgColor,
          backgroundImage: `radial-gradient(${bgSubColor} 1px, transparent 1px)`,
          backgroundSize: `10px 10px`,
          color: "#eaf1f1",
        },
        "#__next": {
          minHeight: "100%",
        },
        input: {
          textAlign: "center",
        },
        "button.chakra-button:hover": {
          background: pink,
          color: white,
        },
      },
    },
  },

  withDefaultVariant({
    variant: "flushed",
    components: ["Input"],
  }),
  withDefaultVariant({
    variant: "outline",
    components: ["Button"],
  }),
  withDefaultColorScheme({
    colorScheme: "pink",
  })
);

const App = ({ Component, pageProps }: AppProps) => {
  return (
    <SWRConfig value={{ fetcher: defaultFetcher }}>
      <ChakraProvider theme={theme}>
        <Flex
          sx={{
            minHeight: "100vh",
            height: "100%",
          }}
        >
          <div style={{ overflow: "hidden" }}>
            <Menu />
          </div>
          <Container maxW="container.xl">
            <Flex direction="column" h="100%">
              <Box flex="1">
                <Component {...pageProps} />
              </Box>
              <Box
                sx={{
                  opacity: "0.2",
                  textAlign: "center",
                }}
              >
                powered by kosenctfx
              </Box>
            </Flex>
          </Container>
        </Flex>
      </ChakraProvider>
    </SWRConfig>
  );
};
export default App;
