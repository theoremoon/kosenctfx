import {
    Box,
    Center,
    ChakraProvider,
    Container,
    extendTheme,
    withDefaultVariant,
} from "@chakra-ui/react";
import Menu from "../components/Menu";
import Loading from "../components/loading";
import { defaultFetcher } from "lib/api";
import type { AppProps } from "next/app";
import React, { useEffect } from "react";
import { SWRConfig } from "swr";
import useAccount from "lib/api/account";
import { fetchCTF } from "lib/api/ctf";
import { isStaticMode } from "lib/static";

const theme = extendTheme(
    {
        initialColorMode: "dark",
        useSystemColorMode: false,
        "colors": {
            "gray": {
                "50": "#EFF0F6",
                "100": "#D1D4E5",
                "200": "#B4B9D5",
                "300": "#969DC4",
                "400": "#7982B4",
                "500": "#5B67A4",
                "600": "#495283",
                "700": "#373E62",
                "800": "#0b1933",
                "900": "#121521"
            }
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
                    '&:focus': {
                        boxShadow: 'none',
                    },
                },
            }
        },
    },

    withDefaultVariant({
        variant: "flushed",
        components: ["Input"],
    }),
);

const App = ({ Component, pageProps }: AppProps) => {
    return (
        <SWRConfig value={{
            fetcher: defaultFetcher,
            revalidateOnMount: !isStaticMode,
            revalidateOnFocus: !isStaticMode,
            revalidateOnReconnect: !isStaticMode,
            refreshInterval: (isStaticMode) ? 0 : 30000,
        }}>
            <ChakraProvider theme={theme}>
                <Box className="h-full">
                    <Menu {...pageProps} />
                    <Container maxW="container.xl" minH="100vh">
                        <Component {...pageProps} />
                    </Container>
                    <Center color="#ffffff66">powered by kosenctfx</Center>
                </Box>
            </ChakraProvider>
        </SWRConfig>
    );
};

export default App;
