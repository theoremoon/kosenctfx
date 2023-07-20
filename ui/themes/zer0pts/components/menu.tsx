import {
  Box,
  Flex,
  Spacer,
  Menu as ChakraMenu,
  MenuButton,
  MenuList,
  MenuItem,
  IconButton,
} from "@chakra-ui/react";
import { Link } from "@chakra-ui/next-js";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faBars } from "@fortawesome/free-solid-svg-icons";
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
          <Link href={item.href} key={item.href} fontSize="xl" p={1} mr={4}>
            {item.innerText}
          </Link>
        ))}
        <Spacer />
        <Flex>
          {rightMenuItems.map((item) => (
            <Link href={item.href} key={item.href} fontSize="xl" p={1} mr={4}>
              {item.innerText}
            </Link>
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
                <Link
                  href={item.href}
                  key={item.href}
                  fontSize="xl"
                  p={1}
                  mr={4}
                >
                  {item.innerText}
                </Link>
              </MenuItem>
            ))}
            {rightMenuItems.map((item) => (
              <MenuItem key={item.href}>
                <Link
                  href={item.href}
                  key={item.href}
                  fontSize="xl"
                  p={1}
                  mr={4}
                >
                  {item.innerText}
                </Link>
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
          <Link href="/" fontSize="xl" p={1} mr={4}>
            {siteName}
          </Link>
        }
        leftMenuItems={leftMenuItems}
        rightMenuItems={rightMenuItems}
      />
    </Box>
  );
};

export default Menu;
