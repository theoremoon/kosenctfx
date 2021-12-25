import {
  Box,
  HStack,
  LinkBox,
  LinkOverlay,
  Progress,
  Stack,
  Text,
  useInterval,
} from "@chakra-ui/react";
import { IconProp } from "@fortawesome/fontawesome-svg-core";
import {
  faAddressCard,
  faFlag,
  faFlagUsa,
  faHome,
  faSign,
  faSignInAlt,
  faSignOutAlt,
  faTrophy,
  faUsers,
  faWrench,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import NextLink from "next/link";
import React, { useEffect, useState } from "react";
import useAccount from "../lib/api/account";
import useCTF, { CTF } from "../lib/api/ctf";
import Loading from "./loading";
import Divider from "./divider";

interface MenuLinkProps {
  label: string;
  href: string;
  icon: IconProp;
}

const MenuLink = ({ label, href, icon }: MenuLinkProps) => {
  return (
    <NextLink href={href}>
      <div className="px-4 cursor-pointer mb-2 opacity-40 hover:opacity-100 overflow-clip">
        <div className="inline-block w-6">
          <FontAwesomeIcon icon={icon} />
        </div>
        <p className="text-xl inline-block ml-2 text-white">{label}</p>
      </div>
    </NextLink>
  );
};

interface MenuProps {
  ctf: CTF;
}

const Menu = ({ ctf }: MenuProps) => {
  const { data: account } = useAccount();
  const [progress, setProgress] = useState(0);
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

    setProgress(((ctf.end_at - now) / (ctf.end_at - ctf.start_at)) * 100);

    if (now < ctf.start_at) {
      setCountdown(calcCountdown(now, ctf.start_at));
    } else if (now < ctf.end_at) {
      setCountdown(calcCountdown(now, ctf.end_at));
    }
  };

  useEffect(calcProgress, []);
  useInterval(calcProgress, 1000);

  return (
    <div className="border-r border-white-600 border-opacity-50 h-full">
      <div className="px-4">
        {!ctf.is_open && <Text fontSize="sm">CTF is closed now</Text>}
        {!ctf.is_open && now < ctf.start_at && (
          <Text fontSize="sm">CTF will start in {countdown}</Text>
        )}
        {ctf.is_running && (
          <>
            <Text>CTF is now running!</Text>
            <Text fontSize="sm">{countdown} remains</Text>
            <Progress size="xs" value={progress} colorScheme="pink" />
          </>
        )}
        {ctf.is_over && <Text>CTF is over. Thanks for playing!</Text>}
      </div>

      <MenuLink label="Top" href="/" icon={faHome} />

      {ctf && ((ctf.is_running && account) || ctf.is_over) && (
        <MenuLink label="Tasks" href="/tasks" icon={faFlag} />
      )}

      <MenuLink label="Ranking" href="/ranking" icon={faTrophy} />
      <Divider />

      {account ? (
        <>
          <MenuLink label="Profile" href="/profile" icon={faAddressCard} />
          <Divider />
          <MenuLink label="Logout" href="/logout" icon={faSignOutAlt} />
          {account.is_admin && (
            <>
              <Divider />
              <MenuLink label="Admin" href="/admin" icon={faWrench} />
              <MenuLink label="Tasks" href="/admin/tasks" icon={faFlagUsa} />
              <MenuLink label="Teams" href="/admin/teams" icon={faUsers} />
              <MenuLink label="Config" href="/admin/config" icon={faWrench} />
            </>
          )}
        </>
      ) : (
        <>
          <MenuLink label="Login" href="/login" icon={faSignInAlt} />
          <MenuLink label="Register" href="/register" icon={faSign} />
        </>
      )}
    </div>
  );
};

const MenuDefault = () => {
  const { data: ctf } = useCTF();
  if (!ctf) {
    return <Loading />;
  }
  return <Menu ctf={ctf} />;
};

export default MenuDefault;
