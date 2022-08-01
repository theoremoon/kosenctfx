import {
  Box,
  FormControl,
  FormLabel,
  Input,
  Button,
  Text,
  Link,
  VStack,
} from "@chakra-ui/react";
import Right from "./components/right";
import NextLink from "next/link";
import { LoginProps } from "props/login";

const Login = ({ register, onSubmit }: LoginProps) => {
  return (
    <Box w="sm" mx="auto" mt="10">
      <form onSubmit={onSubmit}>
        <VStack>
          <FormControl>
            <FormLabel htmlFor="teamname">teamname</FormLabel>
            <Input
              id="teamname"
              variant="flushed"
              {...register("teamname", { required: true })}
            ></Input>
          </FormControl>
          <FormControl>
            <FormLabel htmlFor="password">password</FormLabel>
            <Input
              id="password"
              type="password"
              {...register("password", { required: true })}
            ></Input>
          </FormControl>
          <FormControl>
            <Right>
              <Button type="submit">Login</Button>
            </Right>
          </FormControl>
          <Text>
            <Link as={NextLink} href="/passwordreset_request">
              Forgot your password? You can reset your password here.
            </Link>
          </Text>
        </VStack>
      </form>
    </Box>
  );
};

export default Login;
