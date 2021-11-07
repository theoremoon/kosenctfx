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

type ResetRequestParams = {
  email: string;
};

const ResetRequest = () => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ResetRequestParams>();
  const router = useRouter();
  const toast = useToast();
  const onSubmit: SubmitHandler<ResetRequestParams> = async (values) => {
    try {
      await api.post("/passwordreset-request", {
        email: values.email,
      });
      router.push("/passwordreset");
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
          <FormControl isInvalid={errors.email !== undefined}>
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
