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
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faBars } from "@fortawesome/free-solid-svg-icons";
import NextLink from "next/link";
import { MenuProps } from "props/menu";

type responsiveMenuWrapperProps = Pick<
  Omit<MenuProps, "siteName">,
  "leftMenuItems" | "rightMenuItems"
> & { siteName: React.ReactNode };

const ResponsiveMenuWrapper = ({
  siteName,
  leftMenuItems,
  rightMenuItems,
}: responsiveMenuWrapperProps) => {
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

const Menu = ({ siteName, leftMenuItems, rightMenuItems }: MenuProps) => {
  return (
    <Box w="100%" borderBottom="1px solid #4491cf">
      <ResponsiveMenuWrapper
        siteName={
          <NextLink href="/" passHref>
            <Link fontSize="xl" p={1} mr={4}>
              {siteName}
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
