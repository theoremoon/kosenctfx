import { api } from "lib/api";
import { fetchCTF } from "lib/api/ctf";
import { AllPageProps } from "lib/pages";
import useMessage from "lib/useMessage";
import { GetStaticProps } from "next";
import { useRouter } from "next/router";
import { RegisterParams } from "props/register";
import { useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import RegisterView from "theme/register";

const Register = () => {
  const [country, setCountry] = useState("");
  const { register, handleSubmit } = useForm<RegisterParams>();
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

  return RegisterView({
    register,
    onSubmit: handleSubmit(onSubmit),
    country,
    setCountry,
  });
};

export const getStaticProps: GetStaticProps<AllPageProps> = async () => {
  const ctf = await fetchCTF();
  return {
    props: {
      ctf: ctf,
    },
  };
};

export default Register;
