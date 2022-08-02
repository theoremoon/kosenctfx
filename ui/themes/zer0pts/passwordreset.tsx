import {
  Box,
  Button,
  Flex,
  FormControl,
  FormLabel,
  Input,
  VStack,
} from "@chakra-ui/react";
import { PasswordResetProps } from "props/passwordreset";

const PasswordReset = ({ onSubmit, register }: PasswordResetProps) => {
  return (
    <Box w="sm" mx="auto" mt="10">
      <form onSubmit={onSubmit}>
        <VStack>
          <FormControl>
            <FormLabel htmlFor="token">token</FormLabel>
            <Input
              id="token"
              variant="flushed"
              {...register("token", { required: true })}
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
