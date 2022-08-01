import { api } from "lib/api";
import useMessage from "lib/useMessage";
import { useRouter } from "next/router";
import { SubmitHandler, useForm } from "react-hook-form";
import useAccount, { fetchAccount } from "lib/api/account";
import LoginView from "theme/login";
import { isStaticMode } from "lib/static";
import { fetchCTF } from "lib/api/ctf";
import { AllPageProps } from "lib/pages";
import { GetStaticProps } from "next";
import { LoginParams } from "props/login";

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
      error(e);
    }
  };

  return LoginView({ register, onSubmit: handleSubmit(onSubmit) });
};

export const getStaticProps: GetStaticProps<AllPageProps> = async () => {
  const account = isStaticMode ? null : await fetchAccount().catch(() => null);
  const ctf = await fetchCTF();
  return {
    props: {
      account: account,
      ctf: ctf,
    },
  };
};

export default Login;
