import { api } from "lib/api";
import useMessage from "lib/useMessage";
import { useRouter } from "next/router";
import { SubmitHandler, useForm } from "react-hook-form";
import useAccount from "lib/api/account";
import LoginView from "theme/login";
import { fetchCTF } from "lib/api/ctf";
import { AllPageProps } from "lib/pages";
import { GetStaticProps } from "next";
import { LoginParams } from "props/login";
import { isStaticMode, revalidateInterval } from "lib/static";

const Login = () => {
  const router = useRouter();
  const { mutate } = useAccount(null);
  const { message, error } = useMessage();

  const { register, handleSubmit } = useForm<LoginParams>();
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
      console.log(e);
      error(e);
    }
  };

  return LoginView({ register, onSubmit: handleSubmit(onSubmit) });
};

export const getStaticProps: GetStaticProps<AllPageProps> = async () => {
  const ctf = await fetchCTF();
  return {
    props: {
      ctf: ctf,
    },
    revalidate: isStaticMode ? false : revalidateInterval,
  };
};

export default Login;
