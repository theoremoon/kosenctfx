import { Box, Flex, Link, Spacer, Menu as ChakraMenu, MenuButton, MenuList, MenuItem, IconButton } from '@chakra-ui/react';
import { faBars } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import useAccount, { Account } from 'lib/api/account';
import useCTF, { CTF } from 'lib/api/ctf';
import { isStaticMode } from 'lib/static';
import NextLink from 'next/link';
import Loading from './loading';

interface MenuProps {
  ctf: CTF;
  account: Account;
}

const Menu = ({ ctf: ctfDefault, account: accountDefault, ...props }: MenuProps) => {
    const { data: account } = useAccount(accountDefault);
    const { data: ctf } = useCTF(ctfDefault);

    if (ctf === undefined || account === undefined) {
      return <Loading />;
    }
    const canShowTasks = ctf.is_open && (ctf.is_over || (ctf.is_running && account));

    return (
      <Box w="100%" borderBottom="1px solid #4491cf">
        <Flex maxW="container.xl" w="100%" mx="auto" p={2} justify="space-between" display={{ base: 'none', md: 'flex' }}>
            <Flex>
                <NextLink href='/' passHref><Link fontSize="xl" p={1} mr={4}>zer0pts CTF 2022</Link></NextLink>
                {canShowTasks && (<NextLink href='/tasks' passHref><Link fontSize="xl" p={1} mr={4}>TASKS</Link></NextLink>)}
                <NextLink href='/ranking' passHref><Link fontSize="xl" p={1} mr={4}>RANKING</Link></NextLink>
            </Flex>
            {!isStaticMode && (
            <>
              <Spacer />
              <Flex>
                  {account ? (
                    <>
                    {account.is_admin && (
                      <NextLink href='/admin' passHref><Link fontSize="xl" p={1} mr={4}>ADMIN</Link></NextLink>
                    )}
                      <NextLink href='/profile' passHref><Link fontSize="xl" p={1} mr={4}>PROFILE</Link></NextLink>
                      <NextLink href='/logout' passHref><Link fontSize="xl" p={1} mr={4}>LOGOUT</Link></NextLink>
                    </>
                  ) : (
                    <>
                      <NextLink href='/login' passHref><Link fontSize="xl" p={1} mr={4}>LOGIN</Link></NextLink>
                      <NextLink href='/register' passHref><Link fontSize="xl" p={1} mr={4}>REGISTER</Link></NextLink>
                    </>
                  )}
              </Flex>
            </>
            )}
        </Flex>
        <Flex p={2} w="100%" justify="space-between" display={{ base: 'flex', md: 'none' }}>
            <NextLink href='/' passHref><Link fontSize="xl" p={1} mr={4}>zer0pts CTF 2022</Link></NextLink>
            <Spacer />
            <ChakraMenu>
                <MenuButton as={IconButton} icon={<FontAwesomeIcon icon={faBars} />} />
                <MenuList>
                {canShowTasks && (<NextLink href='/tasks' passHref><MenuItem><Link fontSize="xl" p={1} mr={4}>TASKS</Link></MenuItem></NextLink>)}
                <NextLink href='/ranking' passHref><MenuItem><Link fontSize="xl" p={1} mr={4}>RANKING</Link></MenuItem></NextLink>
                {account ? (
                  <>
                    <NextLink href='/profile' passHref><MenuItem><Link fontSize="xl" p={1} mr={4}>PROFILE</Link></MenuItem></NextLink>
                    <NextLink href='/logout' passHref><MenuItem><Link fontSize="xl" p={1} mr={4}>LOGOUT</Link></MenuItem></NextLink>
                  </>
                ) : (
                  <>
                    <NextLink href='/login' passHref><MenuItem><Link fontSize="xl" p={1} mr={4}>LOGIN</Link></MenuItem></NextLink>
                    <NextLink href='/register' passHref><MenuItem><Link fontSize="xl" p={1} mr={4}>REGISTER</Link></MenuItem></NextLink>
                  </>
                )}
                </MenuList>
            </ChakraMenu>
        </Flex>
      </Box>
    );
};

export default Menu;
