import Menu from "./components/menu";
import React from "react";
import { AppProps } from "props/app";
import { ToastContext } from "lib/useMessage";
import { toast, ToastContainer } from "react-toastify";
import Link from "next/link";
import { pink } from "lib/color";
import Head from "next/head";
import "react-toastify/dist/ReactToastify.css";

const toastProvider = {
  info: (msg: string) => {
    toast.info(msg, {
      autoClose: 2000,
      closeOnClick: true,
    });
  },
  error: (msg: string) => {
    toast.error(msg, {
      autoClose: 2000,
      closeOnClick: true,
    });
  },
};

const App = ({
  Component,
  pageProps,
  siteName,
  leftMenuItems,
  rightMenuItems,
}: AppProps) => {
  return (
    <ToastContext.Provider value={toastProvider}>
      <div style={{ margin: "0 auto", maxWidth: "1920px" }}>
        <header>
          <h1
            style={{
              fontSize: "4rem",
              fontWeight: "normal",
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              margin: "0.67em 0",
            }}
          >
            <Link href="/" passHref>
              <a
                style={{
                  display: "block",
                  borderBottom: `2px solid ${pink}`,
                  textDecoration: "none",
                  color: "inherit",
                  fontFamily: "Parisienne, cursive",
                }}
              >
                CakeCTF
              </a>
            </Link>
          </h1>
          <Menu
            siteName={siteName}
            leftMenuItems={leftMenuItems}
            rightMenuItems={rightMenuItems}
          />
        </header>
      </div>
      <main
        style={{
          height: "100%",
          margin: "0 auto",
          marginBottom: "20px",
          marginTop: "50px",
          maxWidth: "1280px",
        }}
      >
        <Component {...pageProps} />
      </main>
      <div
        style={{
          marginTop: "20px",
          marginBottom: "4px",
          color: "#00000066",
          textAlign: "center",
          position: "sticky",
          top: "100vh",
        }}
      >
        powered by kosenctfx
      </div>
      <ToastContainer />
    </ToastContext.Provider>
  );
};

const AppWrapper = (props: AppProps) => {
  return (
    <>
      <Head>
        <link
          href="https://fonts.googleapis.com/css2?family=Parisienne&display=swap"
          rel="stylesheet"
        />
      </Head>
      <App {...props} />
    </>
  );
};

export default AppWrapper;
