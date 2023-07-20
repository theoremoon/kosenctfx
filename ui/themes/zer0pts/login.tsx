import {
  Box,
  FormControl,
  FormLabel,
  Input,
  Button,
  Text,
  VStack,
} from "@chakra-ui/react";
import { Link } from "@chakra-ui/next-js";
import Right from "./components/right";
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
            <Link href="/passwordreset_request">
              Forgot your password? You can reset your password here.
            </Link>
          </Text>
        </VStack>
      </form>
    </Box>
  );
};

export default Login;
