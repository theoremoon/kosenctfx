import {
  Box,
  FormControl,
  FormLabel,
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
import FormWrapper from "../components/formwrapper";
import FormItem from "../components/formitem";
import Input from "../components/input";
import Label from "../components/label";
import Button from "../components/button";

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
  const { mutate } = useAccount();
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
    <FormWrapper>
      <form onSubmit={handleSubmit(onSubmit)}>
        <FormItem>
          <Label htmlFor="teamname">teamname</Label>
          <Input id="teamname" {...register("teamname", { required: true })} />
        </FormItem>
        <FormItem>
          <Label htmlFor="password">password</Label>
          <Input
            id="password"
            type="password"
            {...register("password", { required: true })}
          />
        </FormItem>
        <FormItem>
          <Right>
            <Button type="submit">Login</Button>
          </Right>
        </FormItem>
        <p>
          <Link as={NextLink} href="/passwordreset_request">
            Forgot your password? You can reset your password here.
          </Link>
        </p>
      </form>
    </FormWrapper>
  );
};

export default Login;
