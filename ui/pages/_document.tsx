import NextDocument, { Head, Html, Main, NextScript } from "next/document";
import { bgColor, bgSubColor, pink, white } from "../lib/color";

type Props = {};
class Document extends NextDocument<Props> {
  render() {
    return (
      <Html
        style={{
          backgroundColor: bgColor,
          backgroundImage: `radial-gradient(${bgSubColor} 1px, transparent 1px)`,
          backgroundSize: `10px 10px`,
          color: "#eaf1f1",
        }}
      >
        <Head>
          <title>BSides Ahmedabad CTF 2021</title>
          <link rel="preconnect" href="https://fonts.googleapis.com" />
          <link rel="preconnect" href="https://fonts.gstatic.com" />
          <link
            href="https://fonts.googleapis.com/css2?family=Noto+Sans+Mono:wght@500&display=swap"
            rel="stylesheet"
          />

          <meta property="og:title" content="BSides Ahmedabad CTF 2021"></meta>
          <meta
            property="og:site_name"
            content="BSides Ahmedabad CTF 2021"
          ></meta>
          <meta
            property="og:url"
            content="https://score.bsidesahmedabad.in/"
          ></meta>
          <meta
            property="og:description"
            content="BSides Ahmedabad CTF 2021 organized by zer0pts"
          ></meta>
          <meta property="og:type" content="website"></meta>
          <meta
            property="og:image"
            content="https://www.zer0pts.com/assets/zer0pts_wb.png"
          ></meta>
        </Head>
        <body>
          <Main />
          <NextScript />
        </body>
      </Html>
    );
  }
}

export default Document;
