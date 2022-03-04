import {
  Box,
  Flex,
  Link,
  Spacer,
  Menu as ChakraMenu,
  MenuButton,
  MenuList,
  MenuItem,
  IconButton,
} from "@chakra-ui/react";
import { faBars } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import useAccount, { Account } from "lib/api/account";
import useCTF, { CTF } from "lib/api/ctf";
import { isStaticMode } from "lib/static";
import NextLink from "next/link";
import React from "react";
import Loading from "./loading";

type MenuItem = {
  href: string;
  innerText: string;
};

interface ResponsiveMenuWrapperProps {
  siteName: React.ReactNode;
  leftMenuItems: MenuItem[];
  rightMenuItems: MenuItem[];
}

const ResponsiveMenuWrapper = ({
  siteName,
  leftMenuItems,
  rightMenuItems,
}: ResponsiveMenuWrapperProps) => {
  return (
    <>
      <Flex
        maxW="container.xl"
        w="100%"
        mx="auto"
        p={2}
        justify="space-between"
        display={{ base: "none", md: "flex" }}
      >
        {siteName}
        {leftMenuItems.map((item) => (
          <NextLink href={item.href} passHref key={item.href}>
            <Link fontSize="xl" p={1} mr={4}>
              {item.innerText}
            </Link>
          </NextLink>
        ))}
        <Spacer />
        <Flex>
          {rightMenuItems.map((item) => (
            <NextLink href={item.href} passHref key={item.href}>
              <Link fontSize="xl" p={1} mr={4}>
                {item.innerText}
              </Link>
            </NextLink>
          ))}
        </Flex>
      </Flex>
      <Flex
        p={2}
        w="100%"
        justify="space-between"
        display={{ base: "flex", md: "none" }}
      >
        {siteName}
        <Spacer />
        <ChakraMenu>
          <MenuButton
            as={IconButton}
            icon={<FontAwesomeIcon icon={faBars} />}
          />
          <MenuList>
            {leftMenuItems.map((item) => (
              <MenuItem key={item.href}>
                <NextLink href={item.href} passHref>
                  <Link fontSize="xl" p={1} mr={4}>
                    {item.innerText}
                  </Link>
                </NextLink>
              </MenuItem>
            ))}
            {rightMenuItems.map((item) => (
              <MenuItem key={item.href}>
                <NextLink href={item.href} passHref>
                  <Link fontSize="xl" p={1} mr={4}>
                    {item.innerText}
                  </Link>
                </NextLink>
              </MenuItem>
            ))}
          </MenuList>
        </ChakraMenu>
      </Flex>
    </>
  );
};

interface MenuProps {
  ctf: CTF;
  account: Account;
}

const Menu = ({
  ctf: ctfDefault,
  account: accountDefault,
  ...props
}: MenuProps) => {
  const { data: account } = useAccount(accountDefault);
  const { data: ctf } = useCTF(ctfDefault);

  if (ctf === undefined || account === undefined) {
    return <Loading />;
  }
  const canShowTasks =
    ctf.is_open && (ctf.is_over || (ctf.is_running && account));

  const leftMenuItems = [
    { item: { href: "/task", innerText: "TASKS" }, available: canShowTasks },
    { item: { href: "/ranking", innerText: "RANKING" }, available: true },
  ].flatMap((x) => (x.available ? [x.item] : []));

  const rightMenuItems = [
    {
      item: { href: "/admin", innerText: "ADMIN" },
      available: account && account.is_admin,
    },
    { item: { href: "/profile", innerText: "PROFILE" }, available: account },
    { item: { href: "/login", innerText: "LOGIN" }, available: !account },
    { item: { href: "/register", innerText: "REGISTER" }, available: !account },
    { item: { href: "/logout", innerText: "LOGOUT" }, available: account },
  ].flatMap((x) => (x.available && !isStaticMode ? [x.item] : []));

  return (
    <Box w="100%" borderBottom="1px solid #4491cf">
      <ResponsiveMenuWrapper
        siteName={
          <NextLink href="/" passHref>
            <Link fontSize="xl" p={1} mr={4}>
              zer0pts CTF 2022
            </Link>
          </NextLink>
        }
        leftMenuItems={leftMenuItems}
        rightMenuItems={rightMenuItems}
      />
    </Box>
  );
};

export default Menu;
