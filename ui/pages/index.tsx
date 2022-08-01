import { fetchAccount } from "lib/api/account";
import { AllPageProps } from "lib/pages";
import { isStaticMode } from "lib/static";
import type { GetStaticProps, NextPage } from "next";
import useCTF, { CTF, fetchCTF } from "../lib/api/ctf";
import { useEffect, useState } from "react";
import { useInterval } from "usehooks-ts";
import IndexView from "theme/index";

const useCountdown = (ctf: CTF): string => {
  const [now, setNow] = useState(0);
  const [countdown, setCountdown] = useState("");

  const calcCountdown = (current: number, to: number) => {
    const d = to - current;
    const days = ("" + Math.floor(d / (60 * 60 * 24))).padStart(2, "0");
    const hours = ("" + Math.floor((d % (60 * 60 * 24)) / (60 * 60))).padStart(
      2,
      "0"
    );
    const minutes = ("" + Math.floor((d % (60 * 60)) / 60)).padStart(2, "0");
    const seconds = ("" + Math.floor(d % 60)).padStart(2, "0");
    return days + "d " + hours + ":" + minutes + ":" + seconds;
  };

  const calcProgress = () => {
    setNow(Date.now().valueOf() / 1000);

    if (now < ctf.start_at) {
      setCountdown(calcCountdown(now, ctf.start_at));
    } else if (now < ctf.end_at) {
      setCountdown(calcCountdown(now, ctf.end_at));
    }
  };

  useEffect(calcProgress, []);
  useInterval(calcProgress, 1000);

  if (!ctf.is_open) {
    if (now < ctf.start_at) {
      return `CTF will start in ${countdown}`;
    }
    return "CTF is closed now";
  } else if (ctf.is_running) {
    return `CTF is now running! ${countdown} remains`;
  } else {
    return "CTF is over. Thanks for playing!";
  }
};

type IndexPageProps = AllPageProps;
const Index: NextPage<IndexPageProps> = ({
  ctf: fallbackCTF,
}: IndexPageProps) => {
  const { data: ctf } = useCTF(fallbackCTF);
  const countdown = useCountdown(ctf || fallbackCTF);

  return IndexView({ ctf: ctf || fallbackCTF, status: countdown });
};

export const getStaticProps: GetStaticProps<IndexPageProps> = async () => {
  const ctf = await fetchCTF();
  let account = null;
  try {
    account = isStaticMode ? null : await fetchAccount();
  } catch {
    // do nothing
  }

  return {
    props: {
      ctf: ctf,
      account: account,
    },
  };
};

export default Index;
