import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  VStack,
} from "@chakra-ui/react";
import CountrySelector from "components/countryselector";
import { api } from "lib/api";
import useMessage from "lib/useMessage";
import { useRouter } from "next/router";
import { useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import Right from "../components/right";

type RegisterParams = {
  email: string;
  teamname: string;
  password: string;
};

const Register = () => {
  const [country, setCountry] = useState("");
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterParams>();
  const router = useRouter();
  const { message, error } = useMessage();
  const onSubmit: SubmitHandler<RegisterParams> = async (values) => {
    try {
      const res = await api.post("/register", {
        email: values.email,
        teamname: values.teamname,
        password: values.password,
        country: country,
      });
      message(res);

      router.push("/login");
    } catch (e) {
      error(e);
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
              type="email"
              variant="flushed"
              autoComplete="email"
              {...register("email", { required: true })}
            ></Input>
          </FormControl>
          <FormControl isInvalid={errors.teamname !== undefined}>
            <FormLabel htmlFor="teamname">teamname</FormLabel>
            <Input
              id="teamname"
              variant="flushed"
              autoComplete="username"
              {...register("teamname", { required: true })}
            ></Input>
          </FormControl>
          <FormControl isInvalid={errors.password !== undefined}>
            <FormLabel htmlFor="password">password</FormLabel>
            <Input
              id="password"
              type="password"
              autoComplete="new-password"
              {...register("password", { required: true })}
            ></Input>
          </FormControl>
          <FormControl>
            <FormLabel htmlFor="country">country</FormLabel>
            <CountrySelector
              id="country"
              value={country}
              onChange={(e) => setCountry(e.target.value)}
            />
          </FormControl>

          <FormControl>
            <Right>
              <Button type="submit">Register</Button>
            </Right>
          </FormControl>
        </VStack>
      </form>
    </Box>
  );
};

export default Register;
