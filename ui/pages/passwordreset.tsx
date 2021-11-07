import {
  Box,
  Button,
  Flex,
  FormControl,
  FormLabel,
  Input,
  useToast,
  VStack,
} from "@chakra-ui/react";
import { api, ErrorResponse } from "lib/api";
import { useRouter } from "next/router";
import { SubmitHandler, useForm } from "react-hook-form";

type ResetParams = {
  token: string;
  password: string;
};

const PasswordReset = () => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ResetParams>();
  const router = useRouter();
  const toast = useToast();
  const onSubmit: SubmitHandler<ResetParams> = async (values) => {
    try {
      await api.post("/passwordreset", {
        token: values.token,
        new_password: values.password,
      });

      router.push("/login");
    } catch (e) {
      const message = (e as ErrorResponse).response?.data.message;
      if (message) {
        toast({
          description: message,
          status: "error",
          duration: 2000,
          isClosable: true,
        });
      }
    }
  };
  return (
    <Box w="sm" mx="auto" mt="10">
      <form onSubmit={handleSubmit(onSubmit)}>
        <VStack>
          <FormControl isInvalid={errors.token !== undefined}>
            <FormLabel htmlFor="token">token</FormLabel>
            <Input
              id="token"
              variant="flushed"
              {...register("token", { required: true })}
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
            <Flex w="100%" direction="row-reverse">
              <Button type="submit">Reset Password</Button>
            </Flex>
          </FormControl>
        </VStack>
      </form>
    </Box>
  );
};

export default PasswordReset;
