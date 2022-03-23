import CountrySelector from "components/countryselector";
import { api } from "lib/api";
import useMessage from "lib/useMessage";
import { useRouter } from "next/router";
import { useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import Right from "../components/right";
import FormWrapper from "../components/formwrapper";
import FormItem from "../components/formitem";
import Input from "../components/input";
import Label from "../components/label";
import Button from "../components/button";

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
    <FormWrapper>
      <form onSubmit={handleSubmit(onSubmit)}>
        <FormItem>
          <Label htmlFor="email">email</Label>
          <Input
            id="email"
            type="email"
            autoComplete="email"
            {...register("email", { required: true })}
          />
        </FormItem>

        <FormItem>
          <Label htmlFor="teamname">teamname</Label>
          <Input
            id="teamname"
            autoComplete="username"
            {...register("teamname", { required: true })}
          />
        </FormItem>

        <FormItem>
          <Label htmlFor="password">password</Label>
          <Input
            id="password"
            type="password"
            autoComplete="new-password"
            {...register("password", { required: true })}
          />
        </FormItem>

        <FormItem>
          <Label htmlFor="country">country</Label>
          <CountrySelector
            id="country"
            value={country}
            onChange={(e) => setCountry(e.target.value)}
          />
        </FormItem>

        <FormItem>
          <Right>
            <Button type="submit">Register</Button>
          </Right>
        </FormItem>
      </form>
    </FormWrapper>
  );
};

export default Register;
