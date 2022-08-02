import {
  Box,
  Button,
  Flex,
  FormControl,
  FormLabel,
  Input,
  VStack,
} from "@chakra-ui/react";
import { PasswordResetRequestProps } from "props/passwordResetRequest";

const ResetRequest = ({ register, onSubmit }: PasswordResetRequestProps) => {
  return (
    <Box w="sm" mx="auto" mt="10">
      <form onSubmit={onSubmit}>
        <VStack>
          <FormControl>
            <FormLabel htmlFor="email">email</FormLabel>
            <Input
              id="email"
              variant="flushed"
              {...register("email", { required: true })}
            ></Input>
          </FormControl>
          <FormControl>
            <Flex w="100%" direction="row-reverse">
              <Button type="submit">Send Email</Button>
            </Flex>
          </FormControl>
        </VStack>
      </form>
    </Box>
  );
};

export default ResetRequest;
