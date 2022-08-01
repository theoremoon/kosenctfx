import { api } from "lib/api";
import useMessage from "lib/useMessage";
import { useRouter } from "next/router";
import { SubmitHandler, useForm } from "react-hook-form";
import useAccount from "../lib/api/account";
import LoginView from "theme/login";

type LoginParams = {
  teamname: string;
  password: string;
};

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

export default Login;
