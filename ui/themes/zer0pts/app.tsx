import {
  Box,
  Center,
  ChakraProvider,
  ColorModeProviderProps,
  Container,
  extendTheme,
  withDefaultVariant,
} from "@chakra-ui/react";
import Menu from "./components/menu";
import React from "react";
import { AppProps } from "props/app";

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

const App = ({
  Component,
  pageProps,
  siteName,
  leftMenuItems,
  rightMenuItems,
}: AppProps) => {
  const forceDarkMode: ColorModeProviderProps["colorModeManager"] = {
    get: () => "dark",
    set: () => {
      // noop
    },
    type: "localStorage",
  };
  return (
    <>
      <ChakraProvider theme={theme} colorModeManager={forceDarkMode}>
        <Box className="h-full">
          <Menu
            siteName={siteName}
            leftMenuItems={leftMenuItems}
            rightMenuItems={rightMenuItems}
          />
          <Container maxW="container.xl" minH="100vh">
            <Component {...pageProps} />
          </Container>
          <Center color="#ffffff66" mt={10} mb={1}>
            powered by kosenctfx
          </Center>
        </Box>
      </ChakraProvider>
    </>
  );
};

export default App;
