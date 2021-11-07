import {
  Box,
  Divider,
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

interface MenuLinkProps {
  label: string;
  href: string;
  icon: IconProp;
}

const MenuLink = ({ label, href, icon }: MenuLinkProps) => {
  return (
    <LinkBox
      css={{
        padding: "0.75em",
        paddingLeft: "1em",
        filter: "brightness(0.7)",
      }}
      _hover={{
        filter: "brightness(1)",
        textShadow: "0 0 1px #ffffff",
      }}
    >
      <NextLink href={href} passHref>
        <LinkOverlay>
          <HStack>
            <FontAwesomeIcon icon={icon} />
            <Text fontSize="xl">{label}</Text>
          </HStack>
        </LinkOverlay>
      </NextLink>
    </LinkBox>
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
    <Box
      sx={{
        minWidth: "150px",
        minHeight: "100vh",
        height: "100%",
        borderRightColor: "#EDFDFD33",
        borderRightStyle: "solid",
        borderRightWidth: "1px",
      }}
    >
      <Stack>
        <Box maxW="150px">
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
        </Box>

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
              </>
            )}
          </>
        ) : (
          <>
            <MenuLink label="Login" href="/login" icon={faSignInAlt} />
            <MenuLink label="Register" href="/register" icon={faSign} />
          </>
        )}
      </Stack>
    </Box>
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
