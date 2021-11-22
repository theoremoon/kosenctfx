import { useToast } from "@chakra-ui/react";
import { api, ErrorResponse } from "lib/api";
import { useRouter } from "next/router";
import { SubmitHandler, useForm } from "react-hook-form";
import FormWrapper from "../components/formwrapper";
import FormItem from "../components/formitem";
import Input from "../components/input";
import Label from "../components/label";
import Right from "../components/right";
import Button from "../components/button";

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
    <FormWrapper>
      <form onSubmit={handleSubmit(onSubmit)}>
        <FormItem>
          <Label htmlFor="email">email</Label>
          <Input id="email" {...register("email", { required: true })} />
        </FormItem>

        <FormItem>
          <Right>
            <Button type="submit">Send Email</Button>
          </Right>
        </FormItem>
      </form>
    </FormWrapper>
  );
};

export default ResetRequest;
