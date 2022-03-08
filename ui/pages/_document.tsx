import NextDocument, { Head, Html, Main, NextScript } from "next/document";
import { ColorModeScript } from "@chakra-ui/react";

type Props = Record<string, never>;

class Document extends NextDocument<Props> {
  render() {
    return (
      <Html>
        <Head>
          <title>zer0pts CTF 2022</title>
          <link rel="preconnect" href="https://fonts.googleapis.com" />
          <link rel="preconnect" href="https://fonts.gstatic.com" />
          <link
            href="https://fonts.googleapis.com/css2?family=Noto+Sans+Mono:wght@500&display=swap"
            rel="stylesheet"
          />

          <meta property="og:title" content="zer0pts CTF 2022"></meta>
          <meta property="og:site_name" content="zer0pts CTF 2022"></meta>
          <meta
            property="og:url"
            content="https://2022.ctf.bsidesahmedabad.in/"
          ></meta>
          <meta
            property="og:description"
            content="zer0pts CTF 2021 organized by zer0pts"
          ></meta>
          <meta property="og:type" content="website"></meta>
          <meta
            property="og:image"
            content="https://www.zer0pts.com/assets/zer0pts_wb.png"
          ></meta>
        </Head>
        <body>
          <ColorModeScript initialColorMode="dark" />
          <Main />
          <NextScript />
        </body>
      </Html>
    );
  }
}

export default Document;
