import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  Link,
  Text,
  VStack,
} from "@chakra-ui/react";
import { api } from "lib/api";
import useMessage from "lib/useMessage";
import NextLink from "next/link";
import { useRouter } from "next/router";
import { SubmitHandler, useForm } from "react-hook-form";
import Right from "../components/right";
import useAccount from "../lib/api/account";

type LoginParams = {
  teamname: string;
  password: string;
};

const Login = () => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginParams>();
  const router = useRouter();
  const { mutate } = useAccount(null);
  const { message, error } = useMessage();
  const onSubmit: SubmitHandler<LoginParams> = async (values) => {
    try {
      const res = await api.post("/login", {
        teamname: values.teamname,
        password: values.password,
      });
      message(res);

      mutate();
      router.push("/");
    } catch (e) {
      error(e);
    }
  };
  return (
    <Box w="sm" mx="auto" mt="10">
      <form onSubmit={handleSubmit(onSubmit)}>
        <VStack>
          <FormControl isInvalid={errors.teamname !== undefined}>
            <FormLabel htmlFor="teamname">teamname</FormLabel>
            <Input
              id="teamname"
              variant="flushed"
              {...register("teamname", { required: true })}
            ></Input>
          </FormControl>
          <FormControl isInvalid={errors.password !== undefined}>
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
